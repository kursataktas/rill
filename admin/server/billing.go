package server

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/rilldata/rill/admin/billing"
	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/server/auth"
	adminv1 "github.com/rilldata/rill/proto/gen/rill/admin/v1"
	"github.com/rilldata/rill/runtime/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetBillingSubscription(ctx context.Context, req *adminv1.GetBillingSubscriptionRequest) (*adminv1.GetBillingSubscriptionResponse, error) {
	observability.AddRequestAttributes(ctx, attribute.String("args.org", req.Organization))

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Organization)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.OrganizationPermissions(ctx, org.ID).ManageOrg && !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "not allowed to read org subscriptions")
	}

	if org.BillingCustomerID == "" {
		return &adminv1.GetBillingSubscriptionResponse{Organization: organizationToDTO(org)}, nil
	}

	sub, err := s.admin.Biller.GetActiveSubscription(ctx, org.BillingCustomerID)
	if err != nil {
		if errors.Is(err, billing.ErrNotFound) {
			return &adminv1.GetBillingSubscriptionResponse{Organization: organizationToDTO(org)}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminv1.GetBillingSubscriptionResponse{
		Organization:     organizationToDTO(org),
		Subscription:     subscriptionToDTO(sub),
		BillingPortalUrl: sub.Customer.PortalURL,
	}, nil
}

func (s *Server) UpdateBillingSubscription(ctx context.Context, req *adminv1.UpdateBillingSubscriptionRequest) (*adminv1.UpdateBillingSubscriptionResponse, error) {
	observability.AddRequestAttributes(ctx, attribute.String("args.org", req.Organization))
	observability.AddRequestAttributes(ctx, attribute.String("args.plan_name", req.PlanName))

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Organization)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.OrganizationPermissions(ctx, org.ID).ManageOrg && !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "not allowed to update org billing plan")
	}

	if req.PlanName == "" {
		return nil, status.Error(codes.InvalidArgument, "plan name must be provided")
	}

	if org.BillingCustomerID == "" {
		return nil, status.Error(codes.FailedPrecondition, "billing not yet initialized for the organization")
	}

	bisc, err := s.admin.DB.FindBillingIssueByTypeForOrg(ctx, org.ID, database.BillingIssueTypeSubscriptionCancelled)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if bisc != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "plan cannot be changed on existing subscription as it was cancelled, please renew the subscription")
	}

	plan, err := s.admin.Biller.GetPlanByName(ctx, req.PlanName)
	if err != nil {
		if errors.Is(err, billing.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	// if its a trial plan, start trial only if its a new org
	if plan.Default {
		bi, err := s.admin.DB.FindBillingIssueByTypeForOrg(ctx, org.ID, database.BillingIssueTypeNeverSubscribed)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return nil, status.Errorf(codes.FailedPrecondition, "only new organizations can subscribe to the trial plan %s", plan.Name)
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
		if bi != nil {
			// check against trial orgs quota
			if org.CreatedByUserID != nil {
				u, err := s.admin.DB.FindUser(ctx, *org.CreatedByUserID)
				if err != nil {
					return nil, status.Error(codes.Internal, err.Error())
				}
				if u.QuotaTrialOrgs >= 0 && u.CurrentTrialOrgsCount >= u.QuotaTrialOrgs {
					return nil, status.Errorf(codes.FailedPrecondition, "trial orgs quota exceeded for user %s", u.Email)
				}
			}

			updatedOrg, sub, err := s.admin.StartTrial(ctx, org)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			return &adminv1.UpdateBillingSubscriptionResponse{
				Organization: organizationToDTO(updatedOrg),
				Subscription: subscriptionToDTO(sub),
			}, nil
		}
	}

	forceAccess := claims.Superuser(ctx) && req.SuperuserForceAccess

	if !plan.Public && !forceAccess {
		return nil, status.Errorf(codes.FailedPrecondition, "cannot assign a private plan %s", plan.Name)
	}

	// check for validation errors
	err = s.planChangeValidationChecks(ctx, org, forceAccess)
	if err != nil {
		return nil, err
	}

	if planDowngrade(plan, org) {
		if !forceAccess {
			return nil, status.Errorf(codes.FailedPrecondition, "plan downgrade not supported")
		}
		s.logger.Named("billing").Warn("plan downgrade request", zap.String("org_id", org.ID), zap.String("org_name", org.Name), zap.String("plan_name", plan.Name))
	}

	sub, err := s.admin.Biller.GetActiveSubscription(ctx, org.BillingCustomerID)
	if err != nil {
		if !errors.Is(err, billing.ErrNotFound) {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if sub == nil {
		// create new subscription
		sub, err = s.admin.Biller.CreateSubscription(ctx, org.BillingCustomerID, plan)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		s.logger.Named("billing").Info("new subscription created", zap.String("org_id", org.ID), zap.String("org_name", org.Name), zap.String("plan_id", sub.Plan.ID), zap.String("plan_name", sub.Plan.Name))
	} else {
		// schedule plan change
		oldPlan := sub.Plan
		if oldPlan.ID != plan.ID {
			sub, err = s.admin.Biller.ChangeSubscriptionPlan(ctx, sub.ID, plan)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			s.logger.Named("billing").Info("plan changed", zap.String("org_id", org.ID), zap.String("org_name", org.Name), zap.String("old_plan_id", oldPlan.ID), zap.String("old_plan_name", oldPlan.Name), zap.String("new_plan_id", sub.Plan.ID), zap.String("new_plan_name", sub.Plan.Name))
		}
	}

	org, err = s.updateQuotasAndHandleBillingIssues(ctx, org, sub)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminv1.UpdateBillingSubscriptionResponse{
		Organization: organizationToDTO(org),
		Subscription: subscriptionToDTO(sub),
	}, nil
}

// CancelBillingSubscription cancels the billing subscription for the organization and puts them on default plan
func (s *Server) CancelBillingSubscription(ctx context.Context, req *adminv1.CancelBillingSubscriptionRequest) (*adminv1.CancelBillingSubscriptionResponse, error) {
	observability.AddRequestAttributes(ctx, attribute.String("args.org", req.Organization))

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Organization)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.OrganizationPermissions(ctx, org.ID).ManageOrg && !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "not allowed to cancel org subscription")
	}

	endDate, err := s.admin.Biller.CancelSubscriptionsForCustomer(ctx, org.BillingCustomerID, billing.SubscriptionCancellationOptionEndOfSubscriptionTerm)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !endDate.IsZero() {
		// raise a billing issue of the subscription cancellation
		_, err = s.admin.DB.UpsertBillingIssue(ctx, &database.UpsertBillingIssueOptions{
			OrgID: org.ID,
			Type:  database.BillingIssueTypeSubscriptionCancelled,
			Metadata: database.BillingIssueMetadataSubscriptionCancelled{
				EndDate: endDate,
			},
			EventTime: time.Now(),
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// clean up any trial related billing issues if present
	err = s.admin.CleanupTrialBillingIssues(ctx, org.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logger.Named("billing").Warn("subscription cancelled", zap.String("org_id", org.ID), zap.String("org_name", org.Name))

	return &adminv1.CancelBillingSubscriptionResponse{}, nil
}

func (s *Server) RenewBillingSubscription(ctx context.Context, req *adminv1.RenewBillingSubscriptionRequest) (*adminv1.RenewBillingSubscriptionResponse, error) {
	observability.AddRequestAttributes(ctx, attribute.String("args.org", req.Organization))
	observability.AddRequestAttributes(ctx, attribute.String("args.plan_name", req.PlanName))

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Organization)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.OrganizationPermissions(ctx, org.ID).ManageOrg && !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "not allowed to renew org subscription")
	}

	if org.BillingCustomerID == "" {
		return nil, status.Error(codes.FailedPrecondition, "billing not yet initialized for the organization")
	}

	bisc, err := s.admin.DB.FindBillingIssueByTypeForOrg(ctx, org.ID, database.BillingIssueTypeSubscriptionCancelled)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, status.Errorf(codes.FailedPrecondition, "subscription not cancelled for the organization %s", org.Name)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	plan, err := s.admin.Biller.GetPlanByName(ctx, req.PlanName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if plan.Default {
		return nil, status.Errorf(codes.FailedPrecondition, "cannot renew to trial plan %s", plan.Name)
	}

	forceAccess := claims.Superuser(ctx) && req.SuperuserForceAccess

	// check for validation errors
	err = s.planChangeValidationChecks(ctx, org, forceAccess)
	if err != nil {
		return nil, err
	}

	sub, err := s.admin.Biller.GetActiveSubscription(ctx, org.BillingCustomerID)
	if err != nil {
		if !errors.Is(err, billing.ErrNotFound) {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if sub == nil {
		sub, err = s.admin.Biller.CreateSubscription(ctx, org.BillingCustomerID, plan)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else if sub.EndDate == sub.CurrentBillingCycleEndDate {
		// To make request idempotent, if subscription is still on cancellation schedule, unschedule it
		sub, err = s.admin.Biller.UnscheduleCancellation(ctx, sub.ID)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if sub.Plan.ID != plan.ID {
		// change the plan, won't happen for new subscriptions
		sub, err = s.admin.Biller.ChangeSubscriptionPlan(ctx, sub.ID, plan)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// update quotas
	org, err = s.admin.DB.UpdateOrganization(ctx, org.ID, &database.UpdateOrganizationOptions{
		Name:                                org.Name,
		DisplayName:                         org.DisplayName,
		Description:                         org.Description,
		CustomDomain:                        org.CustomDomain,
		QuotaProjects:                       valOrDefault(sub.Plan.Quotas.NumProjects, org.QuotaProjects),
		QuotaDeployments:                    valOrDefault(sub.Plan.Quotas.NumDeployments, org.QuotaDeployments),
		QuotaSlotsTotal:                     valOrDefault(sub.Plan.Quotas.NumSlotsTotal, org.QuotaSlotsTotal),
		QuotaSlotsPerDeployment:             valOrDefault(sub.Plan.Quotas.NumSlotsPerDeployment, org.QuotaSlotsPerDeployment),
		QuotaOutstandingInvites:             valOrDefault(sub.Plan.Quotas.NumOutstandingInvites, org.QuotaOutstandingInvites),
		QuotaStorageLimitBytesPerDeployment: valOrDefault(sub.Plan.Quotas.StorageLimitBytesPerDeployment, org.QuotaStorageLimitBytesPerDeployment),
		BillingCustomerID:                   org.BillingCustomerID,
		PaymentCustomerID:                   org.PaymentCustomerID,
		BillingEmail:                        org.BillingEmail,
		CreatedByUserID:                     org.CreatedByUserID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// delete the billing issue
	err = s.admin.DB.DeleteBillingIssue(ctx, bisc.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logger.Named("billing").Info("subscription renewed", zap.String("org_id", org.ID), zap.String("org_name", org.Name), zap.String("plan_id", sub.Plan.ID), zap.String("plan_name", sub.Plan.Name))

	return &adminv1.RenewBillingSubscriptionResponse{
		Organization: organizationToDTO(org),
		Subscription: subscriptionToDTO(sub),
	}, nil
}

func (s *Server) GetPaymentsPortalURL(ctx context.Context, req *adminv1.GetPaymentsPortalURLRequest) (*adminv1.GetPaymentsPortalURLResponse, error) {
	observability.AddRequestAttributes(ctx, attribute.String("args.org", req.Organization))
	observability.AddRequestAttributes(ctx, attribute.String("args.return_url", req.ReturnUrl))

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Organization)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.OrganizationPermissions(ctx, org.ID).ManageOrg && !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "not allowed to manage org billing")
	}

	if org.PaymentCustomerID == "" {
		return nil, status.Error(codes.FailedPrecondition, "payment customer not initialized yet for the organization")
	}

	url, err := s.admin.PaymentProvider.GetBillingPortalURL(ctx, org.PaymentCustomerID, req.ReturnUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminv1.GetPaymentsPortalURLResponse{Url: url}, nil
}

// SudoUpdateOrganizationBillingCustomer updates the billing customer id for an organization. May be useful if customer is initialized manually in billing system
func (s *Server) SudoUpdateOrganizationBillingCustomer(ctx context.Context, req *adminv1.SudoUpdateOrganizationBillingCustomerRequest) (*adminv1.SudoUpdateOrganizationBillingCustomerResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.Organization),
		attribute.String("args.billing_customer_id", req.BillingCustomerId),
	)

	claims := auth.GetClaims(ctx)
	if !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "only superusers can manage billing customer")
	}

	if req.BillingCustomerId == "" {
		return nil, status.Error(codes.InvalidArgument, "billing customer id is required")
	}

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Organization)
	if err != nil {
		return nil, err
	}

	opts := &database.UpdateOrganizationOptions{
		Name:                                req.Organization,
		DisplayName:                         org.DisplayName,
		Description:                         org.Description,
		CustomDomain:                        org.CustomDomain,
		QuotaProjects:                       org.QuotaProjects,
		QuotaDeployments:                    org.QuotaDeployments,
		QuotaSlotsTotal:                     org.QuotaSlotsTotal,
		QuotaSlotsPerDeployment:             org.QuotaSlotsPerDeployment,
		QuotaOutstandingInvites:             org.QuotaOutstandingInvites,
		QuotaStorageLimitBytesPerDeployment: org.QuotaStorageLimitBytesPerDeployment,
		BillingCustomerID:                   req.BillingCustomerId,
		PaymentCustomerID:                   org.PaymentCustomerID,
		BillingEmail:                        org.BillingEmail,
		CreatedByUserID:                     org.CreatedByUserID,
	}

	org, err = s.admin.DB.UpdateOrganization(ctx, org.ID, opts)
	if err != nil {
		return nil, err
	}

	// get active subscriptions if present
	sub, err := s.admin.Biller.GetActiveSubscription(ctx, org.BillingCustomerID)
	if err != nil {
		if errors.Is(err, billing.ErrNotFound) {
			return &adminv1.SudoUpdateOrganizationBillingCustomerResponse{
				Organization: organizationToDTO(org),
			}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminv1.SudoUpdateOrganizationBillingCustomerResponse{
		Organization: organizationToDTO(org),
		Subscription: subscriptionToDTO(sub),
	}, nil
}

func (s *Server) ListPublicBillingPlans(ctx context.Context, req *adminv1.ListPublicBillingPlansRequest) (*adminv1.ListPublicBillingPlansResponse, error) {
	observability.AddRequestAttributes(ctx)

	// no permissions required to list public billing plans
	plans, err := s.admin.Biller.GetPublicPlans(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var dtos []*adminv1.BillingPlan
	for _, plan := range plans {
		dtos = append(dtos, billingPlanToDTO(plan))
	}

	return &adminv1.ListPublicBillingPlansResponse{
		Plans: dtos,
	}, nil
}

func (s *Server) ListOrganizationBillingIssues(ctx context.Context, req *adminv1.ListOrganizationBillingIssuesRequest) (*adminv1.ListOrganizationBillingIssuesResponse, error) {
	observability.AddRequestAttributes(ctx, attribute.String("args.org", req.Organization))

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Organization)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "org not found")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims := auth.GetClaims(ctx)
	if !claims.OrganizationPermissions(ctx, org.ID).ReadOrg && !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "not allowed to read org billing errors")
	}

	issues, err := s.admin.DB.FindBillingIssuesForOrg(ctx, org.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var dtos []*adminv1.BillingIssue
	for _, i := range issues {
		dtos = append(dtos, &adminv1.BillingIssue{
			Organization: org.Name,
			Type:         billingIssueTypeToDTO(i.Type),
			Level:        billingIssueLevelToDTO(i.Level),
			Metadata:     billingIssueMetadataToDTO(i.Type, i.Metadata),
			EventTime:    timestamppb.New(i.EventTime),
			CreatedOn:    timestamppb.New(i.CreatedOn),
		})
	}

	return &adminv1.ListOrganizationBillingIssuesResponse{
		Issues: dtos,
	}, nil
}

func (s *Server) SudoDeleteOrganizationBillingIssue(ctx context.Context, req *adminv1.SudoDeleteOrganizationBillingIssueRequest) (*adminv1.SudoDeleteOrganizationBillingIssueResponse, error) {
	observability.AddRequestAttributes(ctx, attribute.String("args.org", req.Organization), attribute.String("args.type", req.Type.String()))

	claims := auth.GetClaims(ctx)
	if !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "only superusers can delete billing errors")
	}

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Organization)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	t, err := dtoBillingIssueTypeToDB(req.Type)
	if err != nil {
		return nil, err
	}

	err = s.admin.DB.DeleteBillingIssueByTypeForOrg(ctx, org.ID, t)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminv1.SudoDeleteOrganizationBillingIssueResponse{}, nil
}

func (s *Server) updateQuotasAndHandleBillingIssues(ctx context.Context, org *database.Organization, sub *billing.Subscription) (*database.Organization, error) {
	org, err := s.admin.DB.UpdateOrganization(ctx, org.ID, &database.UpdateOrganizationOptions{
		Name:                                org.Name,
		DisplayName:                         org.DisplayName,
		Description:                         org.Description,
		CustomDomain:                        org.CustomDomain,
		QuotaProjects:                       valOrDefault(sub.Plan.Quotas.NumProjects, org.QuotaProjects),
		QuotaDeployments:                    valOrDefault(sub.Plan.Quotas.NumDeployments, org.QuotaDeployments),
		QuotaSlotsTotal:                     valOrDefault(sub.Plan.Quotas.NumSlotsTotal, org.QuotaSlotsTotal),
		QuotaSlotsPerDeployment:             valOrDefault(sub.Plan.Quotas.NumSlotsPerDeployment, org.QuotaSlotsPerDeployment),
		QuotaOutstandingInvites:             valOrDefault(sub.Plan.Quotas.NumOutstandingInvites, org.QuotaOutstandingInvites),
		QuotaStorageLimitBytesPerDeployment: valOrDefault(sub.Plan.Quotas.StorageLimitBytesPerDeployment, org.QuotaStorageLimitBytesPerDeployment),
		BillingCustomerID:                   org.BillingCustomerID,
		PaymentCustomerID:                   org.PaymentCustomerID,
		BillingEmail:                        org.BillingEmail,
		CreatedByUserID:                     org.CreatedByUserID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// schedule job to handle billing issues
	_, err = s.admin.Jobs.HandlePlanChangeBillingIssues(ctx, org.ID, sub.ID, sub.Plan.ID, sub.CurrentBillingCycleStartDate)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (s *Server) planChangeValidationChecks(ctx context.Context, org *database.Organization, forceAccess bool) error {
	// not a trial plan, check for a payment method and a valid billing address
	var validationErrs []string
	pc, err := s.admin.PaymentProvider.FindCustomer(ctx, org.PaymentCustomerID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if !pc.HasPaymentMethod {
		validationErrs = append(validationErrs, "no payment method found")
	}

	if !pc.HasBillableAddress {
		validationErrs = append(validationErrs, "no billing address found, click on update information to add billing address")
	}

	be, err := s.admin.DB.FindBillingIssueByTypeForOrg(ctx, org.ID, database.BillingIssueTypePaymentFailed)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			return status.Error(codes.Internal, err.Error())
		}
	}
	if be != nil {
		validationErrs = append(validationErrs, "a previous payment is due, please pay the outstanding amount")
	}

	if len(validationErrs) > 0 && !forceAccess {
		return status.Errorf(codes.FailedPrecondition, "please fix following by visiting billing portal: %s", strings.Join(validationErrs, ", "))
	}

	return nil
}

func subscriptionToDTO(sub *billing.Subscription) *adminv1.Subscription {
	return &adminv1.Subscription{
		Id:                           sub.ID,
		Plan:                         billingPlanToDTO(sub.Plan),
		StartDate:                    timestamppb.New(sub.StartDate),
		EndDate:                      timestamppb.New(sub.EndDate),
		CurrentBillingCycleStartDate: timestamppb.New(sub.CurrentBillingCycleStartDate),
		CurrentBillingCycleEndDate:   timestamppb.New(sub.CurrentBillingCycleEndDate),
		TrialEndDate:                 timestamppb.New(sub.TrialEndDate),
	}
}

func billingPlanToDTO(plan *billing.Plan) *adminv1.BillingPlan {
	return &adminv1.BillingPlan{
		Id:              plan.ID,
		Name:            plan.Name,
		DisplayName:     plan.DisplayName,
		Description:     plan.Description,
		TrialPeriodDays: uint32(plan.TrialPeriodDays),
		Default:         plan.Default,
		Quotas: &adminv1.Quotas{
			Projects:                       valOrEmptyString(plan.Quotas.NumProjects),
			Deployments:                    valOrEmptyString(plan.Quotas.NumDeployments),
			SlotsTotal:                     valOrEmptyString(plan.Quotas.NumSlotsTotal),
			SlotsPerDeployment:             valOrEmptyString(plan.Quotas.NumSlotsPerDeployment),
			OutstandingInvites:             valOrEmptyString(plan.Quotas.NumOutstandingInvites),
			StorageLimitBytesPerDeployment: val64OrEmptyString(plan.Quotas.StorageLimitBytesPerDeployment),
		},
	}
}

func billingIssueTypeToDTO(t database.BillingIssueType) adminv1.BillingIssueType {
	switch t {
	case database.BillingIssueTypeOnTrial:
		return adminv1.BillingIssueType_BILLING_ISSUE_TYPE_ON_TRIAL
	case database.BillingIssueTypeTrialEnded:
		return adminv1.BillingIssueType_BILLING_ISSUE_TYPE_TRIAL_ENDED
	case database.BillingIssueTypeNoPaymentMethod:
		return adminv1.BillingIssueType_BILLING_ISSUE_TYPE_NO_PAYMENT_METHOD
	case database.BillingIssueTypeNoBillableAddress:
		return adminv1.BillingIssueType_BILLING_ISSUE_TYPE_NO_BILLABLE_ADDRESS
	case database.BillingIssueTypePaymentFailed:
		return adminv1.BillingIssueType_BILLING_ISSUE_TYPE_PAYMENT_FAILED
	case database.BillingIssueTypeSubscriptionCancelled:
		return adminv1.BillingIssueType_BILLING_ISSUE_TYPE_SUBSCRIPTION_CANCELLED
	case database.BillingIssueTypeNeverSubscribed:
		return adminv1.BillingIssueType_BILLING_ISSUE_TYPE_NEVER_SUBSCRIBED
	default:
		return adminv1.BillingIssueType_BILLING_ISSUE_TYPE_UNSPECIFIED
	}
}

func billingIssueLevelToDTO(l database.BillingIssueLevel) adminv1.BillingIssueLevel {
	switch l {
	case database.BillingIssueLevelError:
		return adminv1.BillingIssueLevel_BILLING_ISSUE_LEVEL_ERROR
	case database.BillingIssueLevelWarning:
		return adminv1.BillingIssueLevel_BILLING_ISSUE_LEVEL_WARNING
	default:
		return adminv1.BillingIssueLevel_BILLING_ISSUE_LEVEL_UNSPECIFIED
	}
}

func dtoBillingIssueTypeToDB(t adminv1.BillingIssueType) (database.BillingIssueType, error) {
	switch t {
	case adminv1.BillingIssueType_BILLING_ISSUE_TYPE_ON_TRIAL:
		return database.BillingIssueTypeOnTrial, nil
	case adminv1.BillingIssueType_BILLING_ISSUE_TYPE_TRIAL_ENDED:
		return database.BillingIssueTypeTrialEnded, nil
	case adminv1.BillingIssueType_BILLING_ISSUE_TYPE_NO_PAYMENT_METHOD:
		return database.BillingIssueTypeNoPaymentMethod, nil
	case adminv1.BillingIssueType_BILLING_ISSUE_TYPE_NO_BILLABLE_ADDRESS:
		return database.BillingIssueTypeNoBillableAddress, nil
	case adminv1.BillingIssueType_BILLING_ISSUE_TYPE_PAYMENT_FAILED:
		return database.BillingIssueTypePaymentFailed, nil
	case adminv1.BillingIssueType_BILLING_ISSUE_TYPE_SUBSCRIPTION_CANCELLED:
		return database.BillingIssueTypeSubscriptionCancelled, nil
	case adminv1.BillingIssueType_BILLING_ISSUE_TYPE_NEVER_SUBSCRIBED:
		return database.BillingIssueTypeNeverSubscribed, nil
	default:
		return database.BillingIssueTypeUnspecified, status.Error(codes.InvalidArgument, "invalid billing error type")
	}
}

func billingIssueMetadataToDTO(t database.BillingIssueType, m database.BillingIssueMetadata) *adminv1.BillingIssueMetadata {
	switch t {
	case database.BillingIssueTypeOnTrial:
		return &adminv1.BillingIssueMetadata{
			Metadata: &adminv1.BillingIssueMetadata_OnTrial{
				OnTrial: &adminv1.BillingIssueMetadataOnTrial{
					EndDate: timestamppb.New(m.(*database.BillingIssueMetadataOnTrial).EndDate),
				},
			},
		}
	case database.BillingIssueTypeTrialEnded:
		return &adminv1.BillingIssueMetadata{
			Metadata: &adminv1.BillingIssueMetadata_TrialEnded{
				TrialEnded: &adminv1.BillingIssueMetadataTrialEnded{
					GracePeriodEndDate: timestamppb.New(m.(*database.BillingIssueMetadataTrialEnded).GracePeriodEndDate),
				},
			},
		}
	case database.BillingIssueTypeNoPaymentMethod:
		return &adminv1.BillingIssueMetadata{
			Metadata: &adminv1.BillingIssueMetadata_NoPaymentMethod{
				NoPaymentMethod: &adminv1.BillingIssueMetadataNoPaymentMethod{},
			},
		}
	case database.BillingIssueTypeNoBillableAddress:
		return &adminv1.BillingIssueMetadata{
			Metadata: &adminv1.BillingIssueMetadata_NoBillableAddress{
				NoBillableAddress: &adminv1.BillingIssueMetadataNoBillableAddress{},
			},
		}
	case database.BillingIssueTypePaymentFailed:
		paymentFailed := m.(*database.BillingIssueMetadataPaymentFailed)
		invoices := make([]*adminv1.BillingIssueMetadataPaymentFailedMeta, 0)
		for k := range paymentFailed.Invoices {
			invoices = append(invoices, &adminv1.BillingIssueMetadataPaymentFailedMeta{
				InvoiceId:     paymentFailed.Invoices[k].ID,
				InvoiceNumber: paymentFailed.Invoices[k].Number,
				InvoiceUrl:    paymentFailed.Invoices[k].URL,
				AmountDue:     paymentFailed.Invoices[k].Amount,
				Currency:      paymentFailed.Invoices[k].Currency,
				DueDate:       timestamppb.New(paymentFailed.Invoices[k].DueDate),
			})
		}
		return &adminv1.BillingIssueMetadata{
			Metadata: &adminv1.BillingIssueMetadata_PaymentFailed{
				PaymentFailed: &adminv1.BillingIssueMetadataPaymentFailed{
					Invoices: invoices,
				},
			},
		}
	case database.BillingIssueTypeSubscriptionCancelled:
		return &adminv1.BillingIssueMetadata{
			Metadata: &adminv1.BillingIssueMetadata_SubscriptionCancelled{
				SubscriptionCancelled: &adminv1.BillingIssueMetadataSubscriptionCancelled{
					EndDate: timestamppb.New(m.(*database.BillingIssueMetadataSubscriptionCancelled).EndDate),
				},
			},
		}
	case database.BillingIssueTypeNeverSubscribed:
		return &adminv1.BillingIssueMetadata{
			Metadata: &adminv1.BillingIssueMetadata_NeverSubscribed{
				NeverSubscribed: &adminv1.BillingIssueMetadataNeverSubscribed{},
			},
		}
	default:
		return &adminv1.BillingIssueMetadata{}
	}
}

func planDowngrade(newPan *billing.Plan, org *database.Organization) bool {
	// nil or negative values are considered as unlimited
	if comparableInt(newPan.Quotas.NumProjects) < comparableInt(&org.QuotaProjects) {
		return true
	}
	if comparableInt(newPan.Quotas.NumDeployments) < comparableInt(&org.QuotaDeployments) {
		return true
	}
	if comparableInt(newPan.Quotas.NumSlotsTotal) < comparableInt(&org.QuotaSlotsTotal) {
		return true
	}
	if comparableInt(newPan.Quotas.NumSlotsPerDeployment) < comparableInt(&org.QuotaSlotsPerDeployment) {
		return true
	}
	if comparableInt(newPan.Quotas.NumOutstandingInvites) < comparableInt(&org.QuotaOutstandingInvites) {
		return true
	}
	if comparableInt64(newPan.Quotas.StorageLimitBytesPerDeployment) < comparableInt64(&org.QuotaStorageLimitBytesPerDeployment) {
		return true
	}
	return false
}

func comparableInt(v *int) int {
	if v == nil || *v < 0 {
		return math.MaxInt
	}
	return *v
}

func comparableInt64(v *int64) int64 {
	if v == nil || *v < 0 {
		return math.MaxInt64
	}
	return *v
}

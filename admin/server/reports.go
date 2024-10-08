package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/admin"
	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/pkg/authtoken"
	"github.com/rilldata/rill/admin/server/auth"
	adminv1 "github.com/rilldata/rill/proto/gen/rill/admin/v1"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/drivers/slack"
	"github.com/rilldata/rill/runtime/pkg/observability"
	"github.com/rilldata/rill/runtime/pkg/pbutil"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v3"
)

func (s *Server) GetReportMeta(ctx context.Context, req *adminv1.GetReportMetaRequest) (*adminv1.GetReportMetaResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.project_id", req.ProjectId),
		attribute.String("args.branch", req.Branch),
		attribute.String("args.report", req.Report),
	)

	proj, err := s.admin.DB.FindProject(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}

	permissions := auth.GetClaims(ctx).ProjectPermissions(ctx, proj.OrganizationID, proj.ID)
	if !permissions.ReadProdStatus {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to read report meta")
	}

	if proj.ProdBranch != req.Branch {
		return nil, status.Error(codes.InvalidArgument, "branch not found")
	}

	org, err := s.admin.DB.FindOrganization(ctx, proj.OrganizationID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	opts, err := recreateReportOptionsFromSpec(req.Spec)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	annotations := parseReportAnnotations(req.Spec.Annotations)

	var externalEmails []string
	for _, email := range opts.EmailRecipients {
		_, err := s.admin.DB.FindUserByEmail(ctx, email)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				externalEmails = append(externalEmails, email)
			} else {
				return nil, status.Errorf(codes.Internal, "failed to find user by email: %s", err.Error())
			}
		}
	}

	externalUsersUrls := make(map[string]*adminv1.GetReportMetaResponse_Urls, len(externalEmails))
	if len(externalEmails) > 0 {
		// issue magic tokens for external emails
		var ownerID *string
		if annotations.AdminOwnerUserID != "" {
			ownerID = &annotations.AdminOwnerUserID
		}
		emailTokens, err := s.issueMagicTokensForExternalEmails(ctx, proj.OrganizationID, proj.ID, req.Report, ownerID, opts, externalEmails)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to issue magic auth tokens: %s", err.Error())
		}

		for email, token := range emailTokens {
			externalUsersUrls[email] = &adminv1.GetReportMetaResponse_Urls{
				OpenUrl:   s.admin.URLs.WithCustomDomain(org.CustomDomain).ReportOpenExternal(org.Name, proj.Name, req.Report, token, req.ExecutionTime.AsTime()),
				ExportUrl: s.admin.URLs.WithCustomDomain(org.CustomDomain).ReportExportExternal(org.Name, proj.Name, req.Report, token),
			}
		}
	} else {
		// cleanup all report tokens if no external emails if found
		err = s.cleanUpAllReportTokens(ctx, req.Report)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to clean up report tokens: %s", err.Error())
		}
	}

	return &adminv1.GetReportMetaResponse{
		InternalUsersUrls: &adminv1.GetReportMetaResponse_Urls{
			OpenUrl:   s.admin.URLs.WithCustomDomain(org.CustomDomain).ReportOpen(org.Name, proj.Name, req.Report, req.ExecutionTime.AsTime()),
			ExportUrl: s.admin.URLs.WithCustomDomain(org.CustomDomain).ReportExport(org.Name, proj.Name, req.Report),
			EditUrl:   s.admin.URLs.WithCustomDomain(org.CustomDomain).ReportEdit(org.Name, proj.Name, req.Report),
		},
		ExternalUsersUrls: externalUsersUrls,
	}, nil
}

func (s *Server) CreateReport(ctx context.Context, req *adminv1.CreateReportRequest) (*adminv1.CreateReportResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.organization", req.Organization),
		attribute.String("args.project", req.Project),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.Organization, req.Project)
	if err != nil {
		return nil, err
	}

	claims := auth.GetClaims(ctx)
	permissions := claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID)
	if !permissions.CreateReports {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to read project repo")
	}

	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.PermissionDenied, "only users can create reports")
	}

	if proj.ProdDeploymentID == nil {
		return nil, status.Error(codes.FailedPrecondition, "project does not have a production deployment")
	}

	depl, err := s.admin.DB.FindDeployment(ctx, *proj.ProdDeploymentID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	name, err := s.generateReportName(ctx, depl, req.Options.Title)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data, err := s.yamlForManagedReport(req.Options, claims.OwnerID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to generate report YAML: %s", err.Error())
	}

	err = s.admin.DB.UpsertVirtualFile(ctx, &database.InsertVirtualFileOptions{
		ProjectID: proj.ID,
		Branch:    proj.ProdBranch,
		Path:      virtualFilePathForManagedReport(name),
		Data:      data,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert virtual file: %s", err.Error())
	}

	err = s.admin.TriggerParserAndAwaitResource(ctx, depl, name, runtime.ResourceKindReport)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "timed out waiting for report to be created")
		}
		return nil, status.Errorf(codes.Internal, "failed to reconcile report: %s", err.Error())
	}

	return &adminv1.CreateReportResponse{
		Name: name,
	}, nil
}

func (s *Server) EditReport(ctx context.Context, req *adminv1.EditReportRequest) (*adminv1.EditReportResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.organization", req.Organization),
		attribute.String("args.project", req.Project),
		attribute.String("args.name", req.Name),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.Organization, req.Project)
	if err != nil {
		return nil, err
	}

	claims := auth.GetClaims(ctx)
	permissions := claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID)
	if !permissions.ReadProd {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to read project repo")
	}

	if proj.ProdDeploymentID == nil {
		return nil, status.Error(codes.FailedPrecondition, "project does not have a production deployment")
	}

	depl, err := s.admin.DB.FindDeployment(ctx, *proj.ProdDeploymentID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	spec, err := s.admin.LookupReport(ctx, depl, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not get report: %s", err.Error())
	}
	annotations := parseReportAnnotations(spec.Annotations)

	if !annotations.AdminManaged {
		return nil, status.Error(codes.FailedPrecondition, "can't edit report because it was not created from the UI")
	}

	isOwner := claims.OwnerType() == auth.OwnerTypeUser && annotations.AdminOwnerUserID == claims.OwnerID()
	if !permissions.ManageReports && !isOwner {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to edit report")
	}

	data, err := s.yamlForManagedReport(req.Options, annotations.AdminOwnerUserID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to generate report YAML: %s", err.Error())
	}

	err = s.admin.DB.UpsertVirtualFile(ctx, &database.InsertVirtualFileOptions{
		ProjectID: proj.ID,
		Branch:    proj.ProdBranch,
		Path:      virtualFilePathForManagedReport(req.Name),
		Data:      data,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update virtual file: %s", err.Error())
	}

	err = s.admin.TriggerParserAndAwaitResource(ctx, depl, req.Name, runtime.ResourceKindReport)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "timed out waiting for report to be updated")
		}
		return nil, status.Errorf(codes.Internal, "failed to reconcile report: %s", err.Error())
	}

	return &adminv1.EditReportResponse{}, nil
}

func (s *Server) UnsubscribeReport(ctx context.Context, req *adminv1.UnsubscribeReportRequest) (*adminv1.UnsubscribeReportResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.organization", req.Organization),
		attribute.String("args.project", req.Project),
		attribute.String("args.name", req.Name),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.Organization, req.Project)
	if err != nil {
		return nil, err
	}

	claims := auth.GetClaims(ctx)
	permissions := claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID)
	if !permissions.ReadProd {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to read project repo")
	}

	if proj.ProdDeploymentID == nil {
		return nil, status.Error(codes.FailedPrecondition, "project does not have a production deployment")
	}

	depl, err := s.admin.DB.FindDeployment(ctx, *proj.ProdDeploymentID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	spec, err := s.admin.LookupReport(ctx, depl, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not get report: %s", err.Error())
	}
	annotations := parseReportAnnotations(spec.Annotations)

	if !annotations.AdminManaged {
		return nil, status.Error(codes.FailedPrecondition, "can't edit report because it was not created from the UI")
	}

	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.PermissionDenied, "only users can unsubscribe from reports")
	}
	user, err := s.admin.DB.FindUser(ctx, claims.OwnerID())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	opts, err := recreateReportOptionsFromSpec(spec)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to recreate report options: %s", err.Error())
	}

	found := false
	for idx, email := range opts.EmailRecipients {
		if strings.EqualFold(user.Email, email) {
			opts.EmailRecipients = slices.Delete(opts.EmailRecipients, idx, idx+1)
			found = true
			break
		}
	}
	for idx, email := range opts.SlackUsers {
		if strings.EqualFold(user.Email, email) {
			opts.SlackUsers = slices.Delete(opts.SlackUsers, idx, idx+1)
			found = true
			break
		}
	}

	if !found {
		return nil, status.Error(codes.InvalidArgument, "user is not subscribed to report")
	}

	if len(opts.EmailRecipients) == 0 && len(opts.SlackUsers) == 0 && len(opts.SlackChannels) == 0 && len(opts.SlackWebhooks) == 0 {
		err = s.admin.DB.UpdateVirtualFileDeleted(ctx, proj.ID, proj.ProdBranch, virtualFilePathForManagedReport(req.Name))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update virtual file: %s", err.Error())
		}
		err = s.cleanUpAllReportTokens(ctx, req.Name)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to clean up report tokens: %s", err.Error())
		}
	} else {
		data, err := s.yamlForManagedReport(opts, annotations.AdminOwnerUserID)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to generate report YAML: %s", err.Error())
		}

		err = s.admin.DB.UpsertVirtualFile(ctx, &database.InsertVirtualFileOptions{
			ProjectID: proj.ID,
			Branch:    proj.ProdBranch,
			Path:      virtualFilePathForManagedReport(req.Name),
			Data:      data,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update virtual file: %s", err.Error())
		}
	}

	err = s.admin.TriggerParserAndAwaitResource(ctx, depl, req.Name, runtime.ResourceKindReport)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "timed out waiting for report to be updated")
		}
		return nil, status.Errorf(codes.Internal, "failed to reconcile report: %s", err.Error())
	}

	return &adminv1.UnsubscribeReportResponse{}, nil
}

func (s *Server) DeleteReport(ctx context.Context, req *adminv1.DeleteReportRequest) (*adminv1.DeleteReportResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.organization", req.Organization),
		attribute.String("args.project", req.Project),
		attribute.String("args.name", req.Name),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.Organization, req.Project)
	if err != nil {
		return nil, err
	}

	claims := auth.GetClaims(ctx)
	permissions := claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID)
	if !permissions.ReadProd {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to read project repo")
	}

	if proj.ProdDeploymentID == nil {
		return nil, status.Error(codes.FailedPrecondition, "project does not have a production deployment")
	}

	depl, err := s.admin.DB.FindDeployment(ctx, *proj.ProdDeploymentID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	spec, err := s.admin.LookupReport(ctx, depl, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not get report: %s", err.Error())
	}
	annotations := parseReportAnnotations(spec.Annotations)

	if !annotations.AdminManaged {
		return nil, status.Error(codes.FailedPrecondition, "can't edit report because it was not created from the UI")
	}

	isOwner := claims.OwnerType() == auth.OwnerTypeUser && annotations.AdminOwnerUserID == claims.OwnerID()
	if !permissions.ManageReports && !isOwner {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to edit report")
	}

	err = s.admin.DB.UpdateVirtualFileDeleted(ctx, proj.ID, proj.ProdBranch, virtualFilePathForManagedReport(req.Name))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete virtual file: %s", err.Error())
	}

	err = s.cleanUpAllReportTokens(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clean up report tokens: %s", err.Error())
	}

	err = s.admin.TriggerParserAndAwaitResource(ctx, depl, req.Name, runtime.ResourceKindReport)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "timed out waiting for report to be deleted")
		}
		return nil, status.Errorf(codes.Internal, "failed to reconcile report: %s", err.Error())
	}

	return &adminv1.DeleteReportResponse{}, nil
}

func (s *Server) TriggerReport(ctx context.Context, req *adminv1.TriggerReportRequest) (*adminv1.TriggerReportResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.organization", req.Organization),
		attribute.String("args.project", req.Project),
		attribute.String("args.name", req.Name),
	)

	proj, err := s.admin.DB.FindProjectByName(ctx, req.Organization, req.Project)
	if err != nil {
		return nil, err
	}

	claims := auth.GetClaims(ctx)
	permissions := claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID)
	if !permissions.ReadProd {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to read project repo")
	}

	if proj.ProdDeploymentID == nil {
		return nil, status.Error(codes.FailedPrecondition, "project does not have a production deployment")
	}

	depl, err := s.admin.DB.FindDeployment(ctx, *proj.ProdDeploymentID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	spec, err := s.admin.LookupReport(ctx, depl, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not get report: %s", err.Error())
	}
	annotations := parseReportAnnotations(spec.Annotations)

	isOwner := claims.OwnerType() == auth.OwnerTypeUser && annotations.AdminOwnerUserID == claims.OwnerID()
	if !permissions.ManageReports && !isOwner {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to edit report")
	}

	err = s.admin.TriggerReport(ctx, depl, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to trigger report: %s", err.Error())
	}

	return &adminv1.TriggerReportResponse{}, nil
}

func (s *Server) GenerateReportYAML(ctx context.Context, req *adminv1.GenerateReportYAMLRequest) (*adminv1.GenerateReportYAMLResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.organization", req.Organization),
		attribute.String("args.project", req.Project),
	)

	data, err := s.yamlForCommittedReport(req.Options)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to generate report YAML: %s", err.Error())
	}

	return &adminv1.GenerateReportYAMLResponse{
		Yaml: string(data),
	}, nil
}

func (s *Server) yamlForManagedReport(opts *adminv1.ReportOptions, ownerUserID string) ([]byte, error) {
	res := reportYAML{}
	res.Type = "report"
	res.Title = opts.Title
	res.Refresh.Cron = opts.RefreshCron
	res.Refresh.TimeZone = opts.RefreshTimeZone
	res.Watermark = "inherit"
	res.Intervals.Duration = opts.IntervalDuration
	res.Query.Name = opts.QueryName
	res.Query.ArgsJSON = opts.QueryArgsJson
	res.Export.Format = opts.ExportFormat.String()
	res.Export.Limit = uint(opts.ExportLimit)
	res.Notify.Email.Recipients = opts.EmailRecipients
	res.Notify.Slack.Channels = opts.SlackChannels
	res.Notify.Slack.Users = opts.SlackUsers
	res.Notify.Slack.Webhooks = opts.SlackWebhooks
	res.Annotations.AdminOwnerUserID = ownerUserID
	res.Annotations.AdminManaged = true
	res.Annotations.AdminNonce = time.Now().Format(time.RFC3339Nano)
	res.Annotations.WebOpenPath = opts.WebOpenPath
	res.Annotations.WebOpenState = opts.WebOpenState
	return yaml.Marshal(res)
}

func (s *Server) yamlForCommittedReport(opts *adminv1.ReportOptions) ([]byte, error) {
	// Format args as pretty YAML
	var args map[string]interface{}
	if opts.QueryArgsJson != "" {
		err := json.Unmarshal([]byte(opts.QueryArgsJson), &args)
		if err != nil {
			return nil, fmt.Errorf("failed to parse queryArgsJSON: %w", err)
		}
	}

	// Format export format as pretty string
	var exportFormat string
	switch opts.ExportFormat {
	case runtimev1.ExportFormat_EXPORT_FORMAT_CSV:
		exportFormat = "csv"
	case runtimev1.ExportFormat_EXPORT_FORMAT_PARQUET:
		exportFormat = "parquet"
	case runtimev1.ExportFormat_EXPORT_FORMAT_XLSX:
		exportFormat = "xlsx"
	default:
		exportFormat = opts.ExportFormat.String()
	}

	res := reportYAML{}
	res.Type = "report"
	res.Title = opts.Title
	res.Refresh.Cron = opts.RefreshCron
	res.Refresh.TimeZone = opts.RefreshTimeZone
	res.Watermark = "inherit"
	res.Intervals.Duration = opts.IntervalDuration
	res.Query.Name = opts.QueryName
	res.Query.Args = args
	res.Export.Format = exportFormat
	res.Export.Limit = uint(opts.ExportLimit)
	res.Notify.Email.Recipients = opts.EmailRecipients
	res.Notify.Slack.Channels = opts.SlackChannels
	res.Notify.Slack.Users = opts.SlackUsers
	res.Notify.Slack.Webhooks = opts.SlackWebhooks
	res.Annotations.WebOpenPath = opts.WebOpenPath
	res.Annotations.WebOpenState = opts.WebOpenState
	return yaml.Marshal(res)
}

// generateReportName generates a random report name with the title as a seed.
// Example: "My report!" -> "my-report-5b3f7e1a".
// It verifies that the name is not taken (the random component makes any collision unlikely, but we check to be sure).
func (s *Server) generateReportName(ctx context.Context, depl *database.Deployment, title string) (string, error) {
	for i := 0; i < 5; i++ {
		name := randomReportName(title)

		_, err := s.admin.LookupReport(ctx, depl, name)
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				// Success! Name isn't taken
				return name, nil
			}
			return "", fmt.Errorf("failed to check report name: %w", err)
		}
	}

	// Fail-safe in case all names we tried were taken
	return uuid.New().String(), nil
}

func (s *Server) issueMagicTokensForExternalEmails(ctx context.Context, orgID, projectID, reportName string, ownerID *string, opts *adminv1.ReportOptions, emails []string) (map[string]string, error) {
	emailSet := make(map[string]struct{}, len(emails))
	for _, email := range emails {
		emailSet[email] = struct{}{}
	}

	emailTokens := make(map[string]string, len(emails))

	// check if magic token already exists for the recipient, recipient list might have changed
	tkns, err := s.admin.DB.FindReportTokensWithSecret(ctx, reportName)
	if err != nil {
		return nil, fmt.Errorf("failed to find report tokens: %w", err)
	}

	var unusedTknIDs []string
	for _, tkn := range tkns {
		if _, ok := emailSet[tkn.RecipientEmail]; ok {
			id, err := uuid.Parse(tkn.MagicAuthTokenID)
			if err != nil {
				return nil, fmt.Errorf("failed to parse token ID: %w", err)
			}
			token, err := authtoken.FromParts(authtoken.TypeMagic, id, tkn.MagicAuthTokenSecret)
			if err != nil {
				return nil, fmt.Errorf("failed to create magic auth token from parts: %w", err)
			}
			emailTokens[tkn.RecipientEmail] = token.String()
			delete(emailSet, tkn.RecipientEmail)
		} else {
			unusedTknIDs = append(unusedTknIDs, tkn.MagicAuthTokenID)
		}
	}

	// cleanup unused tokens
	if len(unusedTknIDs) > 0 {
		err = s.admin.DB.DeleteMagicAuthTokens(ctx, unusedTknIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to delete unused report tokens: %w", err)
		}
	}

	var args map[string]interface{}
	if opts.QueryArgsJson != "" {
		err := json.Unmarshal([]byte(opts.QueryArgsJson), &args)
		if err != nil {
			return nil, fmt.Errorf("failed to parse queryArgsJSON: %w", err)
		}
	}

	mgcOpts := &admin.IssueMagicAuthTokenOptions{
		ProjectID:       projectID,
		CreatedByUserID: ownerID,
		ResourceType:    runtime.ResourceKindReport,
		ResourceName:    reportName,
		State:           opts.WebOpenState,
		Title:           opts.Title,
	}

	if ownerID != nil {
		// Get the project-level permissions for the creating user.
		orgPerms, err := s.admin.OrganizationPermissionsForUser(ctx, orgID, *ownerID)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		projectPermissions, err := s.admin.ProjectPermissionsForUser(ctx, projectID, *ownerID, orgPerms)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Generate JWT attributes based on the creating user's, but with limited project-level permissions.
		// We store these attributes with the magic token, so it can simulate the creating user (even if the creating user is later deleted or their permissions change).
		//
		// NOTE: A problem with this approach is that if we change the built-in format of JWT attributes, these will remain as they were when captured.
		// NOTE: Another problem is that if the creator is an admin, attrs["admin"] will be true. It shouldn't be a problem today, but could end up leaking some privileges in the future if we're not careful.
		attrs, err := s.jwtAttributesForUser(ctx, *ownerID, orgID, projectPermissions)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		mgcOpts.Attributes = attrs
	}

	// issue magic tokens for new external emails
	cctx, tx, err := s.admin.DB.NewTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	for email := range emailSet {
		if ownerID == nil {
			// set user attrs as per the email
			mgcOpts.Attributes = map[string]interface{}{
				"name":   "",
				"email":  email,
				"domain": email[strings.LastIndex(email, "@")+1:],
				"groups": []string{},
				"admin":  false,
			}
		}

		tkn, err := s.admin.IssueMagicAuthToken(cctx, mgcOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to issue magic auth token for email %s: %w", email, err)
		}

		emailTokens[email] = tkn.Token().String()

		_, err = s.admin.DB.InsertReportToken(cctx, &database.InsertReportTokenOptions{
			ReportName:       reportName,
			RecipientEmail:   email,
			MagicAuthTokenID: tkn.Token().ID.String(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to insert report token for email %s: %w", email, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return emailTokens, nil
}

func (s *Server) cleanUpAllReportTokens(ctx context.Context, reportName string) error {
	tkns, err := s.admin.DB.FindReportTokensWithSecret(ctx, reportName)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed to find report tokens: %w", err)
	}

	var unusedTknIDs []string
	for _, tkn := range tkns {
		unusedTknIDs = append(unusedTknIDs, tkn.MagicAuthTokenID)
	}

	if len(unusedTknIDs) > 0 {
		err = s.admin.DB.DeleteMagicAuthTokens(ctx, unusedTknIDs)
		if err != nil {
			return fmt.Errorf("failed to delete unused report tokens: %w", err)
		}
	}

	return nil
}

var reportNameToDashCharsRegexp = regexp.MustCompile(`[ _]+`)

var reportNameExcludeCharsRegexp = regexp.MustCompile(`[^a-zA-Z0-9-]+`)

func randomReportName(title string) string {
	name := reportNameToDashCharsRegexp.ReplaceAllString(title, "-")
	name = reportNameExcludeCharsRegexp.ReplaceAllString(name, "")
	name = strings.ToLower(name)
	name = strings.Trim(name, "-")
	if name == "" {
		name = uuid.New().String()
	} else {
		name = name + "-" + uuid.New().String()[0:8]
	}
	return name
}

func recreateReportOptionsFromSpec(spec *runtimev1.ReportSpec) (*adminv1.ReportOptions, error) {
	annotations := parseReportAnnotations(spec.Annotations)

	opts := &adminv1.ReportOptions{}
	opts.Title = spec.Title
	if spec.RefreshSchedule != nil && spec.RefreshSchedule.Cron != "" {
		opts.RefreshCron = spec.RefreshSchedule.Cron
		opts.RefreshTimeZone = spec.RefreshSchedule.TimeZone
	}
	opts.IntervalDuration = spec.IntervalsIsoDuration
	opts.QueryName = spec.QueryName
	opts.QueryArgsJson = spec.QueryArgsJson
	opts.ExportLimit = spec.ExportLimit
	opts.ExportFormat = spec.ExportFormat
	for _, notifier := range spec.Notifiers {
		switch notifier.Connector {
		case "email":
			opts.EmailRecipients = pbutil.ToSliceString(notifier.Properties.AsMap()["recipients"])
		case "slack":
			props, err := slack.DecodeProps(notifier.Properties.AsMap())
			if err != nil {
				return nil, err
			}
			opts.SlackUsers = props.Users
			opts.SlackChannels = props.Channels
			opts.SlackWebhooks = props.Webhooks
		default:
			return nil, fmt.Errorf("unknown notifier connector: %s", notifier.Connector)
		}
	}
	opts.WebOpenPath = annotations.WebOpenPath
	opts.WebOpenState = annotations.WebOpenState
	return opts, nil
}

// reportYAML is derived from rillv1.ReportYAML, but adapted for generating (as opposed to parsing) the report YAML.
type reportYAML struct {
	Type    string `yaml:"type"`
	Title   string `yaml:"title"`
	Refresh struct {
		Cron     string `yaml:"cron"`
		TimeZone string `yaml:"time_zone"`
	} `yaml:"refresh"`
	Watermark string `yaml:"watermark"`
	Intervals struct {
		Duration string `yaml:"duration"`
	} `yaml:"intervals"`
	Query struct {
		Name     string         `yaml:"name"`
		Args     map[string]any `yaml:"args,omitempty"`
		ArgsJSON string         `yaml:"args_json,omitempty"`
	} `yaml:"query"`
	Export struct {
		Format string `yaml:"format"`
		Limit  uint   `yaml:"limit"`
	} `yaml:"export"`
	Notify struct {
		Email struct {
			Recipients []string `yaml:"recipients"`
		} `yaml:"email"`
		Slack struct {
			Users    []string `yaml:"users"`
			Channels []string `yaml:"channels"`
			Webhooks []string `yaml:"webhooks"`
		} `yaml:"slack"`
	} `yaml:"notify"`
	Annotations reportAnnotations `yaml:"annotations,omitempty"`
}

type reportAnnotations struct {
	AdminOwnerUserID string `yaml:"admin_owner_user_id"`
	AdminManaged     bool   `yaml:"admin_managed"`
	AdminNonce       string `yaml:"admin_nonce"` // To ensure spec version gets updated on writes, to enable polling in TriggerReconcileAndAwaitReport
	WebOpenPath      string `yaml:"web_open_path"`
	WebOpenState     string `yaml:"web_open_state"`
}

func parseReportAnnotations(annotations map[string]string) reportAnnotations {
	if annotations == nil {
		return reportAnnotations{}
	}

	res := reportAnnotations{}
	res.AdminOwnerUserID = annotations["admin_owner_user_id"]
	res.AdminManaged, _ = strconv.ParseBool(annotations["admin_managed"])
	res.AdminNonce = annotations["admin_nonce"]
	res.WebOpenPath = annotations["web_open_path"]
	res.WebOpenState = annotations["web_open_state"]

	return res
}

func virtualFilePathForManagedReport(name string) string {
	return path.Join("reports", name+".yaml")
}

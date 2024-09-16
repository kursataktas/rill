package worker

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/rilldata/rill/admin"
	"github.com/rilldata/rill/admin/database"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (w *Worker) deploymentsHealthCheck(ctx context.Context) error {
	afterID := ""
	limit := 100
	seenHosts := map[string]bool{}
	expectedInstances := map[string][]string{}
	actualInstances := map[string][]string{}
	var mu sync.Mutex
	for {
		deployments, err := w.admin.DB.FindDeployments(ctx, afterID, limit)
		if err != nil {
			return fmt.Errorf("deployment health check: failed to get deployments: %w", err)
		}
		if len(deployments) == 0 {
			break
		}

		group, cctx := errgroup.WithContext(ctx)
		group.SetLimit(8)
		for _, d := range deployments {
			d := d
			if d.Status != database.DeploymentStatusOK {
				if time.Since(d.UpdatedOn) > time.Hour {
					w.logger.Error("deployment health check: failing deployment", zap.String("id", d.ID), zap.String("status", d.Status.String()), zap.Duration("duration", time.Since(d.UpdatedOn)))
				}
				continue
			}
			addExpectedInstance(expectedInstances, d)
			if seenHosts[d.RuntimeHost] {
				continue
			}
			seenHosts[d.RuntimeHost] = true
			group.Go(func() error {
				instances, ok := w.deploymentHealthCheck(cctx, d)
				if ok {
					mu.Lock()
					actualInstances[d.RuntimeHost] = instances
					mu.Unlock()
				}
				return nil
			})
		}
		if err := group.Wait(); err != nil {
			return err
		}
		if len(deployments) < limit {
			break
		}
		afterID = deployments[len(deployments)-1].ID
		// fetch again
	}

	// compare expected and actual instances
	for host, expected := range expectedInstances {
		actual, ok := actualInstances[host]
		if !ok {
			// runtime health check failed
			continue
		}
		for _, instance := range expected {
			if slices.Contains(actual, instance) {
				continue
			}
			// an expected instance is missing
			// re verify that the deployment is not deleted
			d, err := w.admin.DB.FindDeploymentByInstanceID(ctx, instance)
			if err != nil {
				if errors.Is(err, database.ErrNotFound) {
					// Deployment was deleted
					continue
				}
				w.logger.Error("deployment health check: failed to find deployment", zap.String("instance_id", instance), zap.Error(err))
				continue
			}
			annotations, err := w.annotationsForDeployment(ctx, d)
			if err != nil {
				w.logger.Error("deployment health check: failed to find deployment_annotations", zap.String("project", d.ProjectID), zap.String("deployment", d.ID), zap.Error(err))
				continue
			}
			f := []zap.Field{zap.String("host", d.RuntimeHost), zap.String("instance_id", instance)}
			for k, v := range annotations.ToMap() {
				f = append(f, zap.String(k, v))
			}
			w.logger.Error("deployment health check: missing instance on runtime", f...)
		}
	}
	return nil
}

func (w *Worker) deploymentHealthCheck(ctx context.Context, d *database.Deployment) (instances []string, runtimeOK bool) {
	client, err := w.admin.OpenRuntimeClient(d)
	if err != nil {
		w.logger.Error("deployment health check: failed to open runtime client", zap.String("host", d.RuntimeHost), zap.Error(err))
		return nil, false
	}
	defer client.Close()

	resp, err := client.Health(ctx, &runtimev1.HealthRequest{})
	if err != nil {
		if status.Code(err) != codes.Unavailable {
			w.logger.Error("deployment health check: health check call failed", zap.String("host", d.RuntimeHost), zap.Error(err))
			return nil, false
		}
		// an unavailable error could also be because the deployment got deleted
		d, dbErr := w.admin.DB.FindDeployment(ctx, d.ID)
		if dbErr != nil {
			if errors.Is(dbErr, database.ErrNotFound) {
				// Deployment was deleted
				return nil, false
			}
			w.logger.Error("deployment health check: failed to find deployment", zap.String("deployment", d.ID), zap.Error(dbErr))
			return nil, false
		}
		if d.Status == database.DeploymentStatusOK {
			w.logger.Error("deployment health check: health check call failed", zap.String("host", d.RuntimeHost), zap.Error(err))
		}
		// Deployment status changed (probably being deleted)
		return nil, false
	}

	if runtimeUnhealthy(resp) {
		f := []zap.Field{zap.String("host", d.RuntimeHost)}
		if resp.LimiterError != "" {
			f = append(f, zap.String("limiter_error", resp.LimiterError))
		}
		if resp.ConnCacheError != "" {
			f = append(f, zap.String("conn_cache_error", resp.ConnCacheError))
		}
		if resp.MetastoreError != "" {
			f = append(f, zap.String("metastore_error", resp.MetastoreError))
		}
		if resp.NetworkError != "" {
			f = append(f, zap.String("network_error", resp.NetworkError))
		}
		w.logger.Error("deployment health check: runtime is unhealthy", f...)
		return nil, false
	}
	for instanceID, health := range resp.InstancesHealth {
		instances = append(instances, instanceID)
		if !instanceUnhealthy(health) {
			continue
		}
		// In case of multiple instances on same host the runtime API will return health of all instances
		// But the deployment for instances except one will be different from the deployment passed as argument
		d := d
		if d.RuntimeInstanceID != instanceID {
			d, err = w.admin.DB.FindDeploymentByInstanceID(ctx, instanceID)
			if err != nil {
				// NOTE :: In some race conditions this may return a false alert when instance is deleted
				// but we alert on not found errors as well to handle cases when runtime reports extra instance
				w.logger.Error("deployment health check: failed to find deployment", zap.String("instance_id", instanceID), zap.Error(err))
				continue
			}
		}
		annotations, err := w.annotationsForDeployment(ctx, d)
		if err != nil {
			w.logger.Error("deployment health check: failed to find deployment_annotations", zap.String("project", d.ProjectID), zap.String("deployment", d.ID), zap.Error(err))
			continue
		}
		f := []zap.Field{zap.String("host", d.RuntimeHost), zap.String("instance_id", instanceID)}
		for k, v := range annotations.ToMap() {
			f = append(f, zap.String(k, v))
		}
		if health.OlapError != "" {
			f = append(f, zap.String("olap_error", health.OlapError))
		}
		if health.ControllerError != "" {
			f = append(f, zap.String("controller_error", health.ControllerError))
		}
		if health.RepoError != "" {
			f = append(f, zap.String("repo_error", health.RepoError))
		}
		w.logger.Error("deployment health check: runtime instance is unhealthy", f...)
	}
	return instances, true
}

func (w *Worker) annotationsForDeployment(ctx context.Context, d *database.Deployment) (*admin.DeploymentAnnotations, error) {
	proj, err := w.admin.DB.FindProject(ctx, d.ProjectID)
	if err != nil {
		return nil, err
	}
	org, err := w.admin.DB.FindOrganization(ctx, proj.OrganizationID)
	if err != nil {
		return nil, err
	}
	annotations := w.admin.NewDeploymentAnnotations(org, proj)
	return &annotations, nil
}

func runtimeUnhealthy(r *runtimev1.HealthResponse) bool {
	return r.LimiterError != "" || r.ConnCacheError != "" || r.MetastoreError != "" || r.NetworkError != ""
}

func instanceUnhealthy(i *runtimev1.InstanceHealth) bool {
	return i.OlapError != "" || i.ControllerError != "" || i.RepoError != ""
}

func addExpectedInstance(expectedInstances map[string][]string, d *database.Deployment) {
	if expectedInstances[d.RuntimeHost] == nil {
		expectedInstances[d.RuntimeHost] = []string{}
	}
	expectedInstances[d.RuntimeHost] = append(expectedInstances[d.RuntimeHost], d.RuntimeInstanceID)
}

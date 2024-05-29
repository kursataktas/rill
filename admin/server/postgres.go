package server

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	wire "github.com/jeroenrinzema/psql-wire"
	"github.com/lib/pq/oid"
	"github.com/rilldata/rill/admin/server/auth"
	runtimeauth "github.com/rilldata/rill/runtime/server/auth"
	"go.uber.org/zap"
)

func (s *Server) QueryHandler(ctx context.Context, query string) (wire.PreparedStatements, error) {
	s.logger.Debug("query", zap.String("query", query))
	if strings.Trim(query, " ") == "" {
		return wire.Prepared(wire.NewStatement(func(ctx context.Context, writer wire.DataWriter, parameters []wire.Parameter) error {
			return writer.Empty()
		})), nil
	}

	upperQuery := strings.ToUpper(query)
	if strings.HasPrefix(upperQuery, "SET") {
		return wire.Prepared(wire.NewStatement(func(ctx context.Context, writer wire.DataWriter, parameters []wire.Parameter) error {
			return writer.Complete("SET")
		}, wire.WithColumns(nil))), nil
	}

	if strings.HasPrefix(upperQuery, "BEGIN") || strings.HasPrefix(upperQuery, "COMMIT") || strings.HasPrefix(upperQuery, "ROLLBACK") {
		return wire.Prepared(wire.NewStatement(func(ctx context.Context, writer wire.DataWriter, parameters []wire.Parameter) error {
			return writer.Complete(strings.Trim(upperQuery, ";"))
		}, wire.WithColumns(nil))), nil
	}

	clientParams := wire.ClientParameters(ctx)
	// database is org.password
	tokens := strings.Split(clientParams[wire.ParamDatabase], ".")
	if len(tokens) != 2 {
		return nil, fmt.Errorf("invalid org or project")
	}
	org := tokens[0]
	project := tokens[1]

	// Find the production deployment for the project we're proxying to
	proj, err := s.admin.DB.FindProjectByName(ctx, org, project)
	if err != nil {
		return nil, fmt.Errorf("invalid org or project")
	}

	if proj.ProdDeploymentID == nil {
		return nil, fmt.Errorf("no prod deployment for project")
	}
	depl, err := s.admin.DB.FindDeployment(ctx, *proj.ProdDeploymentID)
	if err != nil {
		return nil, fmt.Errorf("no prod deployment for project")
	}

	var jwt string
	claims := auth.GetClaims(ctx)
	switch claims.OwnerType() {
	case auth.OwnerTypeAnon:
		// If the client is not authenticated with the admin service, we just proxy the contents of the password to the runtime (if any).
		password := ctx.Value(auth.PostgresPassword{}).(string)
		if len(password) >= 6 && strings.EqualFold(password[0:6], "bearer") {
			jwt = strings.TrimSpace(password[6:])
		}
	case auth.OwnerTypeUser, auth.OwnerTypeService:
		// If the client is authenticated with the admin service, we issue a new ephemeral runtime JWT.
		// The JWT should have the same permissions/configuration as one they would get by calling AdminService.GetProject.
		permissions := claims.ProjectPermissions(ctx, proj.OrganizationID, depl.ProjectID)
		if !permissions.ReadProd {
			return nil, fmt.Errorf("does not have permission to access the production deployment")
		}

		var attr map[string]any
		if claims.OwnerType() == auth.OwnerTypeUser {
			attr, err = s.jwtAttributesForUser(ctx, claims.OwnerID(), proj.OrganizationID, permissions)
			if err != nil {
				return nil, err
			}
		}

		jwt, err = s.issuer.NewToken(runtimeauth.TokenOptions{
			AudienceURL: depl.RuntimeAudience,
			Subject:     claims.OwnerID(),
			TTL:         runtimeAccessTokenDefaultTTL,
			InstancePermissions: map[string][]runtimeauth.Permission{
				depl.RuntimeInstanceID: {
					// TODO: Remove ReadProfiling and ReadRepo (may require frontend changes)
					runtimeauth.ReadObjects,
					runtimeauth.ReadMetrics,
					runtimeauth.ReadProfiling,
					runtimeauth.ReadRepo,
					runtimeauth.ReadAPI,
				},
			},
			Attributes: attr,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("runtime proxy not available for owner type %q", claims.OwnerType())
	}

	// Track usage of the deployment
	s.admin.Used.Deployment(depl.ID)

	hostURL, err := url.Parse(depl.RuntimeHost)
	if err != nil {
		return nil, err
	}
	hostURL.Scheme = "postgres"
	hostURL.Host = hostURL.Hostname() + ":" + strconv.FormatInt(int64(15432), 10)
	hostURL.User = url.UserPassword("postgres", fmt.Sprintf("Bearer %s", jwt))
	hostURL.Path = depl.RuntimeInstanceID
	conn, err := connectionPool(ctx, hostURL.String())
	if err != nil {
		s.logger.Info("error in get connection pool", zap.Error(err))
		return nil, err
	}

	rows, err := conn.Query(ctx, query) // query to underlying host
	if err != nil {
		s.logger.Info("error in query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	// handle schema
	fds := rows.FieldDescriptions()
	cols := make([]wire.Column, 0, len(fds))
	for _, fd := range fds {
		cols = append(cols, wire.Column{
			Table: int32(fd.TableOID),
			Name:  fd.Name,
			Oid:   oid.Oid(fd.DataTypeOID),
			Width: fd.DataTypeSize,
			Attr:  int16(fd.TableAttributeNumber),
		})
	}

	// handle data
	// NOTE :: This creates a copy of data and stores this till client starts reading data. This is required so that we
	// can close runtime connection and not wait for client to complete reading whole data which can leak connection.
	// We can improve this logic in future.
	var data [][]any
	for rows.Next() {
		d, err := rows.Values()
		if err != nil {
			s.logger.Info("error in fetching next row", zap.Error(err))
			return nil, err
		}
		data = append(data, d)
	}
	if rows.Err() != nil {
		s.logger.Info("error in fetching rows", zap.Error(err))
		return nil, err
	}

	handle := func(ctx context.Context, writer wire.DataWriter, parameters []wire.Parameter) error {
		for i := 0; i < len(data); i++ {
			if err := writer.Row(data[i]); err != nil {
				return err
			}
		}
		return writer.Complete("OK")
	}
	return wire.Prepared(wire.NewStatement(handle, wire.WithColumns(cols))), nil
}

var (
	runtimePool map[string]*pgxpool.Pool = make(map[string]*pgxpool.Pool)
	mu          sync.Mutex
)

func connectionPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	mu.Lock()
	defer mu.Unlock()

	pool, ok := runtimePool[dsn]
	if ok {
		return pool, nil
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	// Runtime JWts are valid for 30 minutes only
	config.MaxConnLifetime = time.Minute * 29
	// since runtimes get restarted more often than actual DB servers. Consider if this should be reduced to even less time
	// also consider if we should add some health check on connection acquisition
	config.HealthCheckPeriod = time.Minute

	return pgxpool.NewWithConfig(ctx, config)
}

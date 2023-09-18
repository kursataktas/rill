package sources

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/pkg/duckdbsql"
	"github.com/rilldata/rill/runtime/pkg/fileutil"
	"github.com/rilldata/rill/runtime/services/catalog/migrator"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

const _defaultIngestTimeout = 60 * time.Minute

func init() {
	migrator.Register(drivers.ObjectTypeSource, &sourceMigrator{})
}

type sourceMigrator struct{}

func (m *sourceMigrator) Create(
	ctx context.Context,
	olap drivers.OLAPStore,
	repo drivers.RepoStore,
	opts migrator.Options,
	catalogObj *drivers.CatalogEntry,
	logger *zap.Logger,
	ac activity.Client,
) error {
	name := catalogObj.GetSource().Name
	newTableName := fmt.Sprintf("__%s_%s", name, fmt.Sprint(time.Now().UnixMilli()))
	err := ingestSource(ctx, olap, repo, opts, catalogObj, "rill_sources", newTableName, logger, ac)
	if err != nil {
		return err
	}

	// drop existing table to ensure create or replace view succeeds
	err = dropIfExists(ctx, olap, name, false)
	if err != nil {
		return err
	}

	fullTableName := fmt.Sprintf("rill_sources.%s", safeName(newTableName))
	// create view on the ingested table
	err = olap.Exec(ctx, &drivers.Statement{
		Query:    fmt.Sprintf("CREATE OR REPLACE VIEW %s AS SELECT * FROM %s", safeName(name), fullTableName),
		Priority: 1,
	})
	if err != nil {
		// cleanup of temp table if view creation failed
		_ = olap.Exec(ctx, &drivers.Statement{
			Query:    fmt.Sprintf("DROP TABLE IF EXISTS %s", fullTableName),
			Priority: 100,
		})
		// return the original error. error for dropping is less important for the user
		return err
	}
	return nil
}

func (m *sourceMigrator) Update(ctx context.Context,
	olap drivers.OLAPStore,
	repo drivers.RepoStore,
	opts migrator.Options,
	oldCatalogObj, newCatalogObj *drivers.CatalogEntry,
	logger *zap.Logger,
	ac activity.Client,
) error {
	apiSource := newCatalogObj.GetSource()
	newTableName := fmt.Sprintf("__%s_%s", apiSource.Name, fmt.Sprint(time.Now().UnixMilli()))
	fullTableName := fmt.Sprintf("rill_sources.%s", safeName(newTableName))
	err := ingestSource(ctx, olap, repo, opts, newCatalogObj, "rill_sources", newTableName, logger, ac)
	if err != nil {
		// cleanup of temp table. can exist and still error out in incremental ingestion
		_ = olap.Exec(ctx, &drivers.Statement{
			Query:    fmt.Sprintf("DROP TABLE IF EXISTS %s", fullTableName),
			Priority: 100,
		})
		// return the original error. error for dropping is less important for the user
		return err
	}

	// drop existing table to ensure create or replace view succeeds
	err = dropIfExists(ctx, olap, apiSource.Name, false)
	if err != nil {
		return err
	}
	err = olap.WithConnection(ctx, 100, true, true, func(ctx, ensuredCtx context.Context, conn *sql.Conn) error {
		// create view on the ingested table
		_, err = conn.ExecContext(ctx, fmt.Sprintf("CREATE OR REPLACE VIEW %s AS SELECT * FROM %s", safeName(apiSource.Name), fullTableName))
		if err != nil {
			return err
		}

		// query all previous tables and drop tables, ignore any error, it is okay for these queries to fail
		rows, err := conn.QueryContext(ensuredCtx, fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema='rill_sources' AND regexp_full_match(table_name, '__%s_[0-9]*')", apiSource.Name))
		if err == nil {
			var tableName string
			for rows.Next() {
				if err := rows.Scan(&tableName); err != nil {
					break
				}
				if tableName == newTableName {
					continue
				}
				_, err = conn.ExecContext(ensuredCtx, fmt.Sprintf("DROP TABLE IF EXISTS rill_sources.%s", safeName(tableName)))
				if err != nil {
					logger.Info("drop table failed", zap.String("table", tableName), zap.Error(err))
				}
			}
		}
		return nil
	})
	if err != nil {
		// cleanup of temp table if view creation failed
		_ = olap.Exec(ctx, &drivers.Statement{
			Query:    fmt.Sprintf("DROP TABLE IF EXISTS %s", fullTableName),
			Priority: 100,
		})
		// return the original error. error for dropping is less important for the user
		return err
	}
	return nil
}

func (m *sourceMigrator) Rename(ctx context.Context, olap drivers.OLAPStore, from string, catalogObj *drivers.CatalogEntry) error {
	if strings.EqualFold(from, catalogObj.Name) {
		tempName := fmt.Sprintf("__rill_temp_%s", from)
		err := olap.Exec(ctx, &drivers.Statement{
			Query:    fmt.Sprintf("ALTER TABLE %s RENAME TO %s", from, tempName),
			Priority: 100,
		})
		if err != nil {
			return err
		}
		from = tempName
	}

	return olap.Exec(ctx, &drivers.Statement{
		Query:    fmt.Sprintf("ALTER TABLE %s RENAME TO %s", from, catalogObj.Name),
		Priority: 100,
	})
}

func (m *sourceMigrator) Delete(ctx context.Context, olap drivers.OLAPStore, catalogObj *drivers.CatalogEntry) error {
	return dropIfExists(ctx, olap, catalogObj.Name, true)
}

func (m *sourceMigrator) GetDependencies(ctx context.Context, olap drivers.OLAPStore, catalog *drivers.CatalogEntry) ([]string, []*drivers.CatalogEntry) {
	return []string{}, nil
}

func (m *sourceMigrator) Validate(ctx context.Context, olap drivers.OLAPStore, catalog *drivers.CatalogEntry) []*runtimev1.ReconcileError {
	// TODO - Details needs to be added here
	return nil
}

func (m *sourceMigrator) IsEqual(ctx context.Context, cat1, cat2 *drivers.CatalogEntry) bool {
	// TODO: This is hopefully not needed in the new reconcile where parse and changing of connector happens before equals is called
	isSQLSource := cat1.GetSource().Connector == "duckdb" && (cat2.GetSource().Connector == "s3" || cat2.GetSource().Connector == "gcs")

	if !isSQLSource && cat1.GetSource().Connector != cat2.GetSource().Connector {
		return false
	}

	map2 := cat2.GetSource().Properties.AsMap()
	if isSQLSource {
		delete(map2, "path")
		return equal(cat1.GetSource().Properties.AsMap(), map2)
	}

	return equal(cat1.GetSource().Properties.AsMap(), map2)
}

func (m *sourceMigrator) ExistsInOlap(ctx context.Context, olap drivers.OLAPStore, catalog *drivers.CatalogEntry) (bool, error) {
	_, err := olap.InformationSchema().Lookup(ctx, catalog.Name)
	if errors.Is(err, drivers.ErrNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func convertLower(in map[string]string) map[string]string {
	m := make(map[string]string, len(in))
	for key, value := range in {
		m[strings.ToLower(key)] = value
	}
	return m
}

func ingestSource(ctx context.Context, olap drivers.OLAPStore, repo drivers.RepoStore, opts migrator.Options,
	catalogObj *drivers.CatalogEntry, schema, name string, logger *zap.Logger, ac activity.Client,
) (outErr error) {
	apiSource := catalogObj.GetSource()
	if name == "" {
		name = apiSource.Name
	}

	var err error
	// TODO: this should go in the parser in the new reconcile
	if apiSource.Connector == "duckdb" {
		err = mergeFromParsedQuery(apiSource, convertLower(opts.InstanceEnv), repo.Root())
		if err != nil {
			return err
		}
	}

	logger = logger.With(zap.String("source", name))
	var srcConnector drivers.Handle

	if apiSource.Connector == "duckdb" {
		srcConnector = olap.(drivers.Handle)
	} else {
		var err error
		variables := convertLower(opts.InstanceEnv)
		srcConnector, err = drivers.Open(apiSource.Connector, connectorVariables(apiSource, variables, repo.Root()), false, activity.NewNoopClient(), logger)
		if err != nil {
			return fmt.Errorf("failed to open driver %w", err)
		}
		defer srcConnector.Close()
	}

	olapConnection := olap.(drivers.Handle)
	t, ok := olapConnection.AsTransporter(srcConnector, olapConnection)
	if !ok {
		t, ok = srcConnector.AsTransporter(srcConnector, olapConnection)
		if !ok {
			return fmt.Errorf("data transfer not possible from %q to %q", srcConnector.Driver(), olapConnection.Driver())
		}
	}

	src, err := source(apiSource.Connector, apiSource)
	if err != nil {
		return err
	}

	sink := sink(olapConnection.Driver(), schema, name)

	timeout := _defaultIngestTimeout
	if apiSource.GetTimeoutSeconds() > 0 {
		timeout = time.Duration(apiSource.GetTimeoutSeconds()) * time.Second
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ingestionLimit := opts.IngestStorageLimitInBytes
	limitExceeded := false
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctxWithTimeout.Done():
				return
			case <-ticker.C:
				olap, _ := olapConnection.AsOLAP("") // todo :: check this
				if size, ok := olap.EstimateSize(); ok && size > ingestionLimit {
					limitExceeded = true
					cancel()
				}
			}
		}
	}()

	env := convertLower(opts.InstanceEnv)
	allowHostAccess := strings.EqualFold(env["allow_host_access"], "true")

	p := &progress{}
	transferOpts := &drivers.TransferOptions{
		AcquireConnector: func(name string) (drivers.Handle, func(), error) {
			return nil, nil, fmt.Errorf("this reconciler can't resolve connectors")
		},
		Progress:        p,
		LimitInBytes:    ingestionLimit,
		RepoRoot:        repo.Root(),
		AllowHostAccess: allowHostAccess,
	}

	transferStart := time.Now()
	defer func() {
		transferLatency := time.Since(transferStart).Milliseconds()
		commonDims := []attribute.KeyValue{
			attribute.String("source", srcConnector.Driver()),
			attribute.String("destination", olapConnection.Driver()),
			attribute.Bool("cancelled", errors.Is(outErr, context.Canceled)),
			attribute.Bool("failed", outErr != nil),
			attribute.Bool("limit_exceeded", limitExceeded),
			attribute.Int64("limit_bytes", ingestionLimit),
		}
		ac.Emit(ctx, "ingestion_ms", float64(transferLatency), commonDims...)
		if p.unit == drivers.ProgressUnitByte {
			ac.Emit(ctx, "ingestion_bytes", float64(p.catalogObj.BytesIngested), commonDims...)
		}
	}()

	err = t.Transfer(ctxWithTimeout, src, sink, transferOpts)
	if limitExceeded {
		return drivers.ErrIngestionLimitExceeded
	}
	return err
}

func mergeFromParsedQuery(apiSource *runtimev1.Source, env map[string]string, repoRoot string) error {
	props := apiSource.Properties.AsMap()
	query, ok := props["sql"]
	if !ok {
		return nil
	}
	queryStr, ok := query.(string)
	if !ok {
		return errors.New("query should be a string")
	}

	// raw sql query
	ast, err := duckdbsql.Parse(queryStr)
	if err != nil {
		return err
	}
	refs := ast.GetTableRefs()
	if len(refs) != 1 {
		return errors.New("sql source should have exactly one table reference")
	}
	ref := refs[0]

	if len(ref.Paths) == 0 {
		return errors.New("only read_* functions with a single path is supported")
	}
	if len(ref.Paths) > 1 {
		return errors.New("invalid source, only a single path for source is supported")
	}

	p, c, ok := parseEmbeddedSourceConnector(ref.Paths[0])
	if !ok {
		return errors.New("unknown source")
	}
	switch c {
	case "local_file":
		queryStr, err = rewriteLocalRelativePath(ast, repoRoot, strings.EqualFold(env["allow_host_access"], "true"))
		if err != nil {
			return err
		}
	case "s3", "gcs":
		apiSource.Connector = c
		props["path"] = p
	default:
		return nil
	}

	props["sql"] = queryStr

	pbProps, err := structpb.NewStruct(props)
	if err != nil {
		return err
	}
	apiSource.Properties = pbProps
	return nil
}

func rewriteLocalRelativePath(ast *duckdbsql.AST, repoRoot string, allowRootAccess bool) (string, error) {
	var resolveErr error
	err := ast.RewriteTableRefs(func(table *duckdbsql.TableRef) (*duckdbsql.TableRef, bool) {
		newPaths := make([]string, 0)
		for _, p := range table.Paths {
			lp, err := fileutil.ResolveLocalPath(p, repoRoot, allowRootAccess)
			if err != nil {
				resolveErr = err
				return nil, false
			}
			newPaths = append(newPaths, lp)
		}

		return &duckdbsql.TableRef{
			Function:   table.Function,
			Paths:      newPaths,
			Properties: table.Properties,
		}, true
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	if err != nil {
		return "", err
	}

	return ast.Format()
}

type progress struct {
	catalogObj drivers.CatalogEntry
	unit       drivers.ProgressUnit
}

func (p *progress) Target(val int64, unit drivers.ProgressUnit) {
	p.unit = unit
}

func (p *progress) Observe(val int64, unit drivers.ProgressUnit) {
	if unit == drivers.ProgressUnitByte {
		p.catalogObj.BytesIngested += val
	}
}

func source(connector string, src *runtimev1.Source) (map[string]any, error) {
	props := src.Properties.AsMap()
	return props, nil
}

func sink(connector, schemaName, tableName string) map[string]any {
	return map[string]any{"table": tableName, "schema": schemaName}
}

func connectorVariables(src *runtimev1.Source, env map[string]string, repoRoot string) map[string]any {
	connector := src.Connector
	vars := map[string]any{
		"allow_host_access": strings.EqualFold(env["allow_host_access"], "true"),
	}
	switch connector {
	case "s3":
		vars["aws_access_key_id"] = env["aws_access_key_id"]
		vars["aws_secret_access_key"] = env["aws_secret_access_key"]
		vars["aws_session_token"] = env["aws_session_token"]
	case "gcs":
		vars["google_application_credentials"] = env["google_application_credentials"]
	case "motherduck":
		vars["token"] = env["token"]
		vars["dsn"] = ""
	case "local_file":
		vars["dsn"] = repoRoot
	case "bigquery":
		vars["google_application_credentials"] = env["google_application_credentials"]
	}
	return vars
}

func equal(s, o map[string]any) bool {
	return reflect.DeepEqual(s, o)
}

func safeName(name string) string {
	if name == "" {
		return name
	}
	return fmt.Sprintf("\"%s\"", strings.ReplaceAll(name, "\"", "\"\""))
}

func dropIfExists(ctx context.Context, olap drivers.OLAPStore, name string, dropView bool) error {
	tbl, err := olap.InformationSchema().Lookup(ctx, name)
	if err != nil {
		if errors.Is(err, drivers.ErrNotFound) {
			return nil
		}
		return err
	}

	return olap.WithConnection(ctx, 100, false, false, func(ctx, ensuredCtx context.Context, conn *sql.Conn) error {
		if !tbl.View {
			// table
			_, err = conn.ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", safeName(name)))
			return err
		} else if dropView {
			_, err = conn.ExecContext(ctx, fmt.Sprintf("DROP VIEW IF EXISTS %s", safeName(name)))
			if err != nil {
				return err
			}
			rows, err := conn.QueryContext(ensuredCtx, fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema='rill_sources' AND regexp_full_match(table_name, '__%s_[0-9]*')", name))
			if err != nil {
				return err
			}

			var tableName string
			var lastErr error
			for rows.Next() {
				if err := rows.Scan(&tableName); err != nil {
					break
				}
				if _, err = conn.ExecContext(ensuredCtx, fmt.Sprintf("DROP TABLE IF EXISTS rill_sources.%s", tableName)); err != nil {
					lastErr = err
				}
			}
			return lastErr
		}
		return nil
	})
}

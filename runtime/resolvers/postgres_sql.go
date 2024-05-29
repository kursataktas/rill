package resolvers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	metricssqlparser "github.com/rilldata/rill/runtime/pkg/metricssql"
	"go.uber.org/zap"
)

var errForbidden = errors.New("access to metrics view is forbidden")

func init() {
	runtime.RegisterResolverInitializer("builtin_pg_catalog_sql", newBuiltinPostgresSQL)
	runtime.RegisterBuiltinAPI("pg-catalog-sql", "builtin_pg_catalog_sql", nil)
}

type args struct {
	SQL      string `mapstructure:"sql"`
	Priority int    `mapstructure:"priority"`
	DataDir  string `mapstructure:"temp_dir"`
}

// newBuiltinPostgresSQL is the resolver for queries that are required for supporting postgres protocol
// It is supposed to be used internally only
func newBuiltinPostgresSQL(ctx context.Context, opts *runtime.ResolverOptions) (runtime.Resolver, error) {
	// Decode the args
	args := &args{}
	if err := mapstructure.Decode(opts.Args, args); err != nil {
		return nil, err
	}

	// hacks for working with superset
	// replaces part of queries that are not supported in duckDB
	args.SQL = strings.ReplaceAll(args.SQL, "ix.indrelid = c.conrelid and\n                                ix.indexrelid = c.conindid and\n                                c.contype in ('p', 'u', 'x')", "ix.indrelid = c.conrelid")
	args.SQL = strings.ReplaceAll(args.SQL, "t.oid = a.attrelid and a.attnum = ANY(ix.indkey)", "t.oid = a.attrelid")
	args.SQL = strings.ReplaceAll(args.SQL, "pg_get_constraintdef(cons.oid)", "pg_get_constraintdef(cons.oid, false)")
	// duckdb reports type hugeint postgres supports bigint
	args.SQL = strings.ReplaceAll(args.SQL, "pg_catalog.format_type(a.atttypid, a.atttypmod)", "CASE WHEN pg_catalog.format_type(a.atttypid, a.atttypmod) == 'hugeint' THEN 'bigint' ELSE pg_catalog.format_type(a.atttypid, a.atttypmod) END")
	if args.SQL == "SELECT nspname FROM pg_namespace WHERE nspname NOT LIKE 'pg_%' ORDER BY nspname" {
		args.SQL = "SELECT nspname FROM pg_namespace WHERE nspname NOT IN ('pg_catalog', 'information_schema', 'main') ORDER BY nspname"
	}

	// hacks for working with metabase
	args.SQL = strings.ReplaceAll(args.SQL, "t.schemaname <> 'information_schema'", "t.schemaname <> 'information_schema' AND t.schemaname <> 'pg_catalog' AND t.schemaname <> 'main'")
	args.SQL = strings.ReplaceAll(args.SQL, "(information_schema._pg_expandarray(i.indkey)).n", "generate_subscripts(i.indkey, 1)")
	args.SQL = extraCharRe.ReplaceAllString(args.SQL, "\n")
	// check if its a non catalog query like `SHOW variable` or `select 1`
	resolver, ok := resolveNonCatalog(args.SQL)
	if ok {
		return resolver, nil
	}

	ctrl, err := opts.Runtime.Controller(ctx, opts.InstanceID)
	if err != nil {
		return nil, err
	}

	resources, err := ctrl.List(ctx, runtime.ResourceKindMetricsView, "", false)
	if err != nil {
		return nil, err
	}

	// We create a db in temporary location, attach it to main db, update catalog with metric_view resources
	// detach the temp db and use it to run queries
	dbDir := filepath.Join(args.DataDir, strconv.FormatInt(time.Now().UnixMilli(), 10))
	if err := os.Mkdir(dbDir, fs.ModePerm); err != nil {
		return nil, err
	}
	dbPath := filepath.Join(dbDir, "catalog.db")
	// nolint:gosec // We don't need cryptographically secure random numbers
	dbName := fmt.Sprintf("pg_catalog_db_%v", rand.Int())
	args.SQL = rewriteSQL(args.SQL)
	// loop over all resources and create corresponding table in duckdb so that these can be queried with information_schema
	for _, resource := range resources {
		metricSQL, err := fromQueryForMetricsView(ctx, ctrl, opts, resource)
		if err != nil {
			if errors.Is(err, errForbidden) {
				continue
			}
			return nil, err
		}

		compiler := metricssqlparser.New(ctrl, opts.InstanceID, opts.UserAttributes, args.Priority)
		createTableSQL, connector, _, err := compiler.Compile(ctx, metricSQL)
		if err != nil {
			return nil, err
		}

		if err := populateMetricView(ctx, opts, connector, createTableSQL, dbPath, dbName, resource.Meta.Name.Name); err != nil {
			return nil, err
		}
	}

	handle, err := drivers.Open("duckdb", opts.InstanceID, map[string]any{"path": dbPath}, activity.NewNoopClient(), zap.NewNop())
	if err != nil {
		return nil, err
	}

	olap, ok := handle.AsOLAP(opts.InstanceID)
	if !ok {
		return nil, fmt.Errorf("developer error : handle is not an OLAP driver")
	}
	if err := olap.Exec(ctx, &drivers.Statement{Query: "USE catalog.public"}); err != nil {
		return nil, err
	}

	if err := populateCatalogTables(ctx, olap); err != nil {
		return nil, err
	}
	return &catalogSQLResolver{
		olap:  olap,
		sql:   args.SQL,
		dbDir: dbDir,
	}, nil
}

func populateMetricView(ctx context.Context, opts *runtime.ResolverOptions, connector, query, path, dbName, mvName string) error {
	olap, release, err := opts.Runtime.OLAP(ctx, opts.InstanceID, connector)
	if err != nil {
		return err
	}
	defer release()

	if olap.Dialect() != drivers.DialectDuckDB {
		// other OLAP engines can be supported with a schema translation layer from engine -> duckDB
		// and creating table with DDL statements
		return fmt.Errorf("only duckdb is supported")
	}

	err = olap.WithConnection(ctx, 1, false, false, func(ctx, ensuredCtx context.Context, _ *sql.Conn) error {
		if err := olap.Exec(ctx, &drivers.Statement{Query: fmt.Sprintf("ATTACH '%s' AS %s", path, dbName)}); err != nil {
			return err
		}

		// postgres's default schema is public
		if err := olap.Exec(ctx, &drivers.Statement{Query: fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s.public", dbName)}); err != nil {
			return err
		}

		defer func() {
			_ = olap.Exec(ensuredCtx, &drivers.Statement{Query: fmt.Sprintf("DETACH %s", dbName)})
		}()

		qry := fmt.Sprintf("CREATE TABLE %s.public.%s AS SELECT * FROM (%s) LIMIT 0", dbName, olap.Dialect().EscapeIdentifier(mvName), query)
		return olap.Exec(ctx, &drivers.Statement{Query: qry})
	})
	return err
}

func populateCatalogTables(ctx context.Context, olap drivers.OLAPStore) error {
	// create missing catalog tables
	return olap.Exec(ctx, &drivers.Statement{
		Query: "CREATE TABLE catalog.pg_catalog.pg_matviews(schemaname VARCHAR, matviewname VARCHAR, matviewowner VARCHAR, tablespace VARCHAR, hasindexes BOOLEAN, ispopulated BOOLEAN, definition VARCHAR)",
	})
}

type catalogSQLResolver struct {
	olap  drivers.OLAPStore
	sql   string
	dbDir string
}

func (r *catalogSQLResolver) Close() error {
	return os.RemoveAll(r.dbDir)
}

func (r *catalogSQLResolver) Key() string {
	return r.sql
}

func (r *catalogSQLResolver) Refs() []*runtimev1.ResourceName {
	return nil
}

func (r *catalogSQLResolver) Validate(ctx context.Context) error {
	_, err := r.olap.Execute(ctx, &drivers.Statement{
		Query:  r.sql,
		DryRun: true,
	})
	return err
}

func (r *catalogSQLResolver) ResolveInteractive(ctx context.Context, opts *runtime.ResolverInteractiveOptions) (*runtime.ResolverResult, error) {
	res, err := r.olap.Execute(ctx, &drivers.Statement{
		Query: r.sql,
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if opts != nil && opts.Format == runtime.GOOBJECTS {
		return r.scanAsGoObjects(res)
	}

	var out []map[string]any
	for res.Rows.Next() {
		row := make(map[string]any)
		err = res.Rows.MapScan(row)
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}

	data, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}

	return &runtime.ResolverResult{
		Data:   data,
		Schema: res.Schema,
		Cache:  false, // never cache information schema queries
	}, nil
}

func (r *catalogSQLResolver) scanAsGoObjects(res *drivers.Result) (*runtime.ResolverResult, error) {
	var out [][]any
	for res.Rows.Next() {
		row, err := res.Rows.SliceScan()
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}

	return &runtime.ResolverResult{
		Rows:   out,
		Schema: res.Schema,
		Cache:  false,
	}, nil
}

func (r *catalogSQLResolver) ResolveExport(ctx context.Context, w io.Writer, opts *runtime.ResolverExportOptions) error {
	return fmt.Errorf("not implemented")
}

func fromQueryForMetricsView(ctx context.Context, ctrl *runtime.Controller, opts *runtime.ResolverOptions, mv *runtimev1.Resource) (string, error) {
	spec := mv.GetMetricsView().State.ValidSpec
	if spec == nil {
		return "", fmt.Errorf("metrics view %q is not ready for querying, reconcile status: %q", mv.Meta.GetName(), mv.Meta.ReconcileStatus)
	}

	olap, release, err := ctrl.Runtime.OLAP(ctx, opts.InstanceID, spec.Connector)
	if err != nil {
		return "", err
	}
	defer release()
	dialect := olap.Dialect()

	security, err := ctrl.Runtime.ResolveMetricsViewSecurity(opts.UserAttributes, opts.InstanceID, spec, mv.Meta.StateUpdatedOn.AsTime())
	if err != nil {
		return "", err
	}

	var cols []string
	for _, measure := range spec.Measures {
		cols = append(cols, measure.Name)
	}
	for _, dim := range spec.Dimensions {
		cols = append(cols, dim.Name)
	}

	if security == nil {
		if spec.TimeDimension != "" {
			cols = append(cols, spec.TimeDimension)
		}
		return fmt.Sprintf("SELECT %s FROM %s", strings.Join(cols, ","), dialect.EscapeIdentifier(mv.Meta.Name.Name)), nil
	}

	if !security.Access || security.ExcludeAll {
		return "", errForbidden
	}

	var final []string
	if len(security.Include) != 0 {
		for _, measure := range cols {
			if slices.Contains(security.Include, measure) { // only include the included cols if include is set
				final = append(final, measure)
			}
		}
	}
	if len(final) > 0 {
		cols = final
	}

	for _, col := range cols {
		if !slices.Contains(security.Exclude, col) {
			final = append(final, col)
		}
	}

	if spec.TimeDimension != "" {
		final = append(final, spec.TimeDimension)
	}

	var sqlStr strings.Builder
	sqlStr.WriteString("SELECT ")
	sqlStr.WriteString(strings.Join(final, ","))
	sqlStr.WriteString(" FROM ")
	sqlStr.WriteString(dialect.EscapeIdentifier(mv.Meta.Name.Name))
	if security.RowFilter != "" {
		sqlStr.WriteString(" WHERE ")
		sqlStr.WriteString(security.RowFilter)
	}
	return sqlStr.String(), nil
}

var (
	functions         = "has_any_column_privilege|has_column_privilege|has_database_privilege|has_foreign_data_wrapper_privilege|has_function_privilege|has_language_privilege|has_parameter_privilege|has_schema_privilege|has_sequence_privilege|has_server_privilege|has_table_privilege|has_tablespace_privilege|has_type_privilege|pg_has_role"
	re                = regexp.MustCompile(fmt.Sprintf(`pg_catalog.(%s)\(([^,]+), ([^,]+), ([^)]+)\)`, functions))
	dbRe              = regexp.MustCompile(`pg_catalog\.(\w+)`)
	regclassRe        = regexp.MustCompile(`'pg_class'::regclass`)
	versionRe         = regexp.MustCompile(`pg_catalog\.version\(\)`)
	pgBackendPid      = regexp.MustCompile(`(?:pg_catalog\.)?pg_backend_pid\([^)]*\)`)
	indexRe           = regexp.MustCompile(`(?:pg_catalog\.)?pg_get_indexdef\([^)]*\)`)
	identifyOptionsRe = regexp.MustCompile(`(?is)\(SELECT\s+json_build_object\([^)]*\)\s*FROM[^)]*\)\s+as\s+identity_options`)
	serialSequenceRe  = regexp.MustCompile(`pg_catalog\.pg_get_serial_sequence\([^\)]*\)`)
	extraCharRe       = regexp.MustCompile(`[\n\t\r]`)
	showVarRe         = regexp.MustCompile(`(?i)SHOW\s+(.+)`)
)

func rewriteSQL(input string) string {
	// DuckDB does not support user optional argument in `functions`. We need to remove that.
	result := re.ReplaceAllString(input, `(select pg_catalog.$1($3, $4))`)
	// pg_get_serial_sequence not supported
	result = serialSequenceRe.ReplaceAllString(result, "NULL")
	// setting fixed pg_backend_pid
	result = pgBackendPid.ReplaceAllString(result, `(SELECT 1234) AS pg_backend_pid`)
	// pg_get_indexdef not supported
	result = indexRe.ReplaceAllString(result, "NULL")
	// postgres version
	result = versionRe.ReplaceAllString(result, `(SELECT 'PostgreSQL 16.3 (Debian 16.3-1.pgdg120+1) on aarch64-unknown-linux-gnu, compiled by gcc (Debian 12.2.0-14) 12.2.0, 64-bit') AS version`)
	// duckdb executes catalog queries in system schema by default. We want to execute in user database's public schema.
	result = dbRe.ReplaceAllString(result, `catalog.pg_catalog.$1`)
	// duckdb does not have `regclass` typecast
	result = regclassRe.ReplaceAllString(result, `(SELECT oid FROM pg_class WHERE relname = 'pg_class')`)
	// json_build_object is not supported. It is used in indexes for metabase so we directly set it as NULL.
	result = identifyOptionsRe.ReplaceAllString(result, " NULL AS identity_options")
	return result
}

func resolveNonCatalog(sqlStr string) (runtime.Resolver, bool) {
	sqlStr = strings.TrimSuffix(sqlStr, ";")
	matches := showVarRe.FindStringSubmatch(sqlStr)
	if len(matches) <= 1 {
		return nil, false
	}
	return &postgresSQLResolver{variable: matches[1]}, true
}

type postgresSQLResolver struct {
	variable string
}

// Close implements runtime.Resolver.
func (p *postgresSQLResolver) Close() error {
	return nil
}

// Key implements runtime.Resolver.
func (p *postgresSQLResolver) Key() string {
	return ""
}

// Refs implements runtime.Resolver.
func (p *postgresSQLResolver) Refs() []*runtimev1.ResourceName {
	return nil
}

// ResolveExport implements runtime.Resolver.
func (p *postgresSQLResolver) ResolveExport(ctx context.Context, w io.Writer, opts *runtime.ResolverExportOptions) error {
	panic("unimplemented")
}

// ResolveInteractive implements runtime.Resolver.
func (p *postgresSQLResolver) ResolveInteractive(ctx context.Context, opts *runtime.ResolverInteractiveOptions) (*runtime.ResolverResult, error) {
	fields := make([]*runtimev1.StructType_Field, 1)
	fields[0] = &runtimev1.StructType_Field{
		Name: name(p.variable),
		Type: &runtimev1.Type{Code: runtimev1.Type_CODE_STRING, Nullable: false},
	}

	row := make([][]any, 1)
	row[0] = make([]any, 1)
	row[0][0] = value(p.variable)
	return &runtime.ResolverResult{
		Rows:   row,
		Schema: &runtimev1.StructType{Fields: fields},
		Cache:  false,
	}, nil
}

// Validate implements runtime.Resolver.
func (p *postgresSQLResolver) Validate(ctx context.Context) error {
	return nil
}

func name(variable string) string {
	switch strings.ToLower(variable) {
	case "transaction isolation level":
		return "transaction_isolation"
	default:
		return variable
	}
}

func value(variable string) string {
	switch strings.ToLower(variable) {
	case "standard_conforming_string", "standard_conforming_strings":
		return "on"
	case "transaction isolation level":
		return "read committed"
	case "timezone":
		return "Etc/UTC"
	default:
		return "tbd"
	}
}

var _ runtime.Resolver = &postgresSQLResolver{}

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/jmoiron/sqlx"
	"github.com/rilldata/rill/admin/database"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	// Load postgres driver
	_ "github.com/jackc/pgx/v4/stdlib"
)

func init() {
	database.Register("postgres", driver{})
}

type driver struct{}

func (d driver) Open(dsn string) (database.DB, error) {
	db, err := otelsql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
	if err != nil {
		return nil, err
	}

	dbx := sqlx.NewDb(db, "pgx")
	return &connection{db: dbx}, nil
}

type connection struct {
	db *sqlx.DB
}

func (c *connection) Close() error {
	return c.db.Close()
}

func (c *connection) FindOrganizations(ctx context.Context) ([]*database.Organization, error) {
	var res []*database.Organization
	err := c.getDB(ctx).SelectContext(ctx, &res, "SELECT * FROM organizations ORDER BY lower(name)")
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindOrganizationByName(ctx context.Context, name string) (*database.Organization, error) {
	res := &database.Organization{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT * FROM organizations WHERE lower(name)=lower($1)", name).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindOrganizationByID(ctx context.Context, orgID string) (*database.Organization, error) {
	res := &database.Organization{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT * FROM organizations WHERE id = $1", orgID).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) InsertOrganization(ctx context.Context, opts *database.InsertOrganizationOptions) (*database.Organization, error) {
	if err := database.Validate(opts); err != nil {
		return nil, err
	}

	res := &database.Organization{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "INSERT INTO organizations(name, description) VALUES ($1, $2) RETURNING *", opts.Name, opts.Description).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) UpdateOrganization(ctx context.Context, id string, opts *database.UpdateOrganizationOptions) (*database.Organization, error) {
	if err := database.Validate(opts); err != nil {
		return nil, err
	}

	res := &database.Organization{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "UPDATE organizations SET name=$1, description=$2, updated_on=now() WHERE id=$3 RETURNING *", opts.Name, opts.Description, id).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) DeleteOrganization(ctx context.Context, name string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM organizations WHERE lower(name)=lower($1)", name)
	return parseErr(err)
}

func (c *connection) FindProjects(ctx context.Context, orgName string) ([]*database.Project, error) {
	var res []*database.Project
	err := c.getDB(ctx).SelectContext(ctx, &res, "SELECT p.* FROM projects p JOIN organizations o ON p.org_id = o.id WHERE lower(o.name)=lower($1) ORDER BY lower(p.name)", orgName)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindProjectByID(ctx context.Context, id string) (*database.Project, error) {
	res := &database.Project{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT * FROM projects WHERE id=$1", id).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindProjectByName(ctx context.Context, orgName, name string) (*database.Project, error) {
	res := &database.Project{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT p.* FROM projects p JOIN organizations o ON p.org_id = o.id WHERE lower(p.name)=lower($1) AND lower(o.name)=lower($2)", name, orgName).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindProjectByGithubURL(ctx context.Context, githubURL string) (*database.Project, error) {
	res := &database.Project{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT p.* FROM projects p WHERE lower(p.github_url)=lower($1)", githubURL).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) InsertProject(ctx context.Context, opts *database.InsertProjectOptions) (*database.Project, error) {
	if err := database.Validate(opts); err != nil {
		return nil, err
	}

	res := &database.Project{}
	err := c.getDB(ctx).QueryRowxContext(ctx, `
		INSERT INTO projects (org_id, name, description, public, region, production_olap_driver, production_olap_dsn, production_slots, production_branch, production_variables, github_url, github_installation_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING *`,
		opts.OrganizationID, opts.Name, opts.Description, opts.Public, opts.Region, opts.ProductionOLAPDriver, opts.ProductionOLAPDSN, opts.ProductionSlots, opts.ProductionBranch, database.Variables(opts.ProductionVariables), opts.GithubURL, opts.GithubInstallationID,
	).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) UpdateProject(ctx context.Context, id string, opts *database.UpdateProjectOptions) (*database.Project, error) {
	if err := database.Validate(opts); err != nil {
		return nil, err
	}

	res := &database.Project{}
	err := c.getDB(ctx).QueryRowxContext(ctx, `
		UPDATE projects SET name=$1, description=$2, public=$3, production_branch=$4, production_variables=$5, github_url=$6, github_installation_id=$7, production_deployment_id=$8, updated_on=now()
		WHERE id=$9 RETURNING *`,
		opts.Name, opts.Description, opts.Public, opts.ProductionBranch, database.Variables(opts.ProductionVariables), opts.GithubURL, opts.GithubInstallationID, opts.ProductionDeploymentID, id,
	).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) DeleteProject(ctx context.Context, id string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM projects WHERE id=$1", id)
	return parseErr(err)
}

func (c *connection) FindUsers(ctx context.Context) ([]*database.User, error) {
	var res []*database.User
	err := c.getDB(ctx).SelectContext(ctx, &res, "SELECT u.* FROM users u")
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindUser(ctx context.Context, id string) (*database.User, error) {
	res := &database.User{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT u.* FROM users u WHERE u.id=$1", id).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindUserByEmail(ctx context.Context, email string) (*database.User, error) {
	res := &database.User{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT u.* FROM users u WHERE lower(u.email)=lower($1)", email).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) InsertUser(ctx context.Context, opts *database.InsertUserOptions) (*database.User, error) {
	if err := database.Validate(opts); err != nil {
		return nil, err
	}

	res := &database.User{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "INSERT INTO users (email, display_name, photo_url) VALUES ($1, $2, $3) RETURNING *", opts.Email, opts.DisplayName, opts.PhotoURL).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) UpdateUser(ctx context.Context, id string, opts *database.UpdateUserOptions) (*database.User, error) {
	if err := database.Validate(opts); err != nil {
		return nil, err
	}

	res := &database.User{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "UPDATE users SET display_name=$2, photo_url=$3, github_username=$4, updated_on=now() WHERE id=$1 RETURNING *",
		id,
		opts.DisplayName,
		opts.PhotoURL,
		opts.GithubUsername).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) DeleteUser(ctx context.Context, id string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	return parseErr(err)
}

func (c *connection) FindUserAuthTokens(ctx context.Context, userID string) ([]*database.UserAuthToken, error) {
	var res []*database.UserAuthToken
	err := c.getDB(ctx).SelectContext(ctx, &res, "SELECT t.* FROM user_auth_tokens t WHERE t.user_id=$1", userID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindUserAuthToken(ctx context.Context, id string) (*database.UserAuthToken, error) {
	res := &database.UserAuthToken{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT t.* FROM user_auth_tokens t WHERE t.id=$1", id).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) InsertUserAuthToken(ctx context.Context, opts *database.InsertUserAuthTokenOptions) (*database.UserAuthToken, error) {
	if err := database.Validate(opts); err != nil {
		return nil, err
	}

	res := &database.UserAuthToken{}
	err := c.getDB(ctx).QueryRowxContext(ctx, `
		INSERT INTO user_auth_tokens (id, secret_hash, user_id, display_name, auth_client_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		opts.ID, opts.SecretHash, opts.UserID, opts.DisplayName, opts.AuthClientID,
	).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) DeleteUserAuthToken(ctx context.Context, id string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM user_auth_tokens WHERE id=$1", id)
	return parseErr(err)
}

func (c *connection) InsertAuthCode(ctx context.Context, deviceCode, userCode, clientID string, expiresOn time.Time) (*database.AuthCode, error) {
	res := &database.AuthCode{}
	err := c.getDB(ctx).QueryRowxContext(ctx,
		`INSERT INTO device_code_auth (device_code, user_code, expires_on, approval_state, client_id)
		VALUES ($1, $2, $3, $4, $5)  RETURNING *`, deviceCode, userCode, expiresOn, database.AuthCodeStatePending, clientID).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindAuthCodeByDeviceCode(ctx context.Context, deviceCode string) (*database.AuthCode, error) {
	authCode := &database.AuthCode{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT * FROM device_code_auth WHERE device_code = $1", deviceCode).StructScan(authCode)
	if err != nil {
		return nil, parseErr(err)
	}
	return authCode, nil
}

func (c *connection) FindAuthCodeByUserCode(ctx context.Context, userCode string) (*database.AuthCode, error) {
	authCode := &database.AuthCode{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT * FROM device_code_auth WHERE user_code = $1", userCode).StructScan(authCode)
	if err != nil {
		return nil, parseErr(err)
	}
	return authCode, nil
}

func (c *connection) UpdateAuthCode(ctx context.Context, userCode, userID string, approvalState database.AuthCodeState) error {
	res, err := c.getDB(ctx).ExecContext(ctx, "UPDATE device_code_auth SET approval_state=$1, user_id=$2, updated_on=now() WHERE user_code=$3",
		approvalState, userID, userCode)
	if err != nil {
		return parseErr(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return database.ErrNotFound
	}
	if rows != 1 {
		return fmt.Errorf("problem in updating auth code, expected 1 row to be affected, got %d", rows)
	}
	return nil
}

func (c *connection) DeleteAuthCode(ctx context.Context, deviceCode string) error {
	res, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM device_code_auth WHERE device_code=$1", deviceCode)
	if err != nil {
		return parseErr(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return database.ErrNotFound
	}
	if rows != 1 {
		return fmt.Errorf("problem in deleting auth code, expected 1 row to be affected, got %d", rows)
	}
	return nil
}

func (c *connection) FindDeployments(ctx context.Context, projectID string) ([]*database.Deployment, error) {
	var res []*database.Deployment
	err := c.getDB(ctx).SelectContext(ctx, &res, "SELECT * FROM deployments d WHERE d.project_id=$1", projectID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindDeployment(ctx context.Context, id string) (*database.Deployment, error) {
	res := &database.Deployment{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT d.* FROM deployments d WHERE d.id=$1", id).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) InsertDeployment(ctx context.Context, opts *database.InsertDeploymentOptions) (*database.Deployment, error) {
	if err := database.Validate(opts); err != nil {
		return nil, err
	}

	res := &database.Deployment{}
	err := c.getDB(ctx).QueryRowxContext(ctx, `
		INSERT INTO deployments (project_id, slots, branch, runtime_host, runtime_instance_id, runtime_audience, status, logs)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`,
		opts.ProjectID, opts.Slots, opts.Branch, opts.RuntimeHost, opts.RuntimeInstanceID, opts.RuntimeAudience, opts.Status, opts.Logs,
	).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) UpdateDeploymentStatus(ctx context.Context, id string, status database.DeploymentStatus, logs string) (*database.Deployment, error) {
	res := &database.Deployment{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "UPDATE deployments SET status=$1, logs=$2, updated_on=now() WHERE id=$3 RETURNING *", status, logs, id).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) DeleteDeployment(ctx context.Context, id string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM deployments WHERE id=$1", id)
	return parseErr(err)
}

func (c *connection) QueryRuntimeSlotsUsed(ctx context.Context) ([]*database.RuntimeSlotsUsed, error) {
	var res []*database.RuntimeSlotsUsed
	err := c.getDB(ctx).SelectContext(ctx, &res, "SELECT d.runtime_host, SUM(d.slots) AS slots_used FROM deployments d GROUP BY d.runtime_host")
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func parseErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return database.ErrNotFound
	}
	if strings.Contains(strings.ToLower(err.Error()), "violates unique constraint") {
		return database.ErrNotUnique
	}
	return err
}

func (c *connection) FindOrganizationMemberUsers(ctx context.Context, orgID string) ([]*database.Member, error) {
	var res []*database.Member
	err := c.getDB(ctx).SelectContext(ctx, &res, `SELECT u.id, u.email, u.display_name, u.created_on, u.updated_on, r.name FROM users u 
    	JOIN users_orgs_roles uor ON u.id = uor.user_id
		JOIN org_roles r ON r.id = uor.org_role_id WHERE uor.org_id=$1`, orgID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindOrganizationMemberUsersByRole(ctx context.Context, orgID, roleID string) ([]*database.User, error) {
	var res []*database.User
	err := c.getDB(ctx).SelectContext(
		ctx, &res, "SELECT u.* FROM users u JOIN users_orgs_roles uor on u.id = uor.user_id WHERE uor.org_id=$1 AND uor.org_role_id=$2", orgID, roleID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) InsertOrganizationMemberUser(ctx context.Context, orgID, userID, roleID string) error {
	res, err := c.getDB(ctx).ExecContext(ctx, "INSERT INTO users_orgs_roles (user_id, org_id, org_role_id) VALUES ($1, $2, $3)", userID, orgID, roleID)
	if err != nil {
		return parseErr(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no rows affected when adding user to organization")
	}
	return nil
}

func (c *connection) DeleteOrganizationMemberUser(ctx context.Context, orgID, userID string) error {
	res, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM users_orgs_roles WHERE user_id = $1 AND org_id = $2", userID, orgID)
	if err != nil {
		return parseErr(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no rows affected when removing user from organization")
	}
	return nil
}

func (c *connection) UpdateOrganizationMemberUserRole(ctx context.Context, orgID, userID, roleID string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, `UPDATE users_orgs_roles SET org_role_id = $1 WHERE user_id = $2 AND org_id = $3`,
		roleID, userID, orgID)
	return parseErr(err)
}

func (c *connection) FindProjectMemberUsers(ctx context.Context, projectID string) ([]*database.Member, error) {
	var res []*database.Member
	err := c.getDB(ctx).SelectContext(ctx, &res, `SELECT u.id, u.email, u.display_name, u.created_on, u.updated_on, r.name FROM users u 
    	JOIN users_projects_roles upr ON u.id = upr.user_id
		JOIN project_roles r ON r.id = upr.project_role_id WHERE upr.project_id=$1`, projectID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) InsertProjectMemberUser(ctx context.Context, projectID, userID, roleID string) error {
	res, err := c.getDB(ctx).ExecContext(ctx, "INSERT INTO users_projects_roles (user_id, project_id, project_role_id) VALUES ($1, $2, $3)", userID, projectID, roleID)
	if err != nil {
		return parseErr(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no rows affected when adding user to project")
	}
	return nil
}

func (c *connection) DeleteProjectMemberUser(ctx context.Context, projectID, userID string) error {
	res, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM users_projects_roles WHERE user_id = $1 AND project_id = $2", userID, projectID)
	if err != nil {
		return parseErr(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no rows affected when removing user from project")
	}
	return nil
}

func (c *connection) UpdateProjectMemberUserRole(ctx context.Context, projectID, userID, roleID string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, `UPDATE users_projects_roles SET project_role_id = $1 WHERE user_id = $2 AND project_id = $3`,
		roleID, userID, projectID)
	return parseErr(err)
}

func (c *connection) FindOrganizationRole(ctx context.Context, name string) (*database.OrganizationRole, error) {
	role := &database.OrganizationRole{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT * FROM org_roles WHERE lower(name)=lower($1)", name).StructScan(role)
	if err != nil {
		return nil, parseErr(err)
	}
	return role, nil
}

func (c *connection) ResolveOrganizationMemberUserRoles(ctx context.Context, userID, orgID string) ([]*database.OrganizationRole, error) {
	var res []*database.OrganizationRole
	err := c.getDB(ctx).SelectContext(ctx, &res, `
		SELECT r.* FROM users_orgs_roles uor
		JOIN org_roles r ON uor.org_role_id = r.id
		WHERE uor.user_id = $1 AND uor.org_id = $2
		UNION
		SELECT * FROM org_roles WHERE id IN (
			SELECT org_role_id FROM usergroups_orgs_roles uor JOIN users_usergroups uug 
			ON uor.usergroup_id = uug.usergroup_id WHERE uug.user_id = $1 AND uor.org_id = $2
		)`, userID, orgID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) ResolveProjectMemberUserRoles(ctx context.Context, userID, projectID string) ([]*database.ProjectRole, error) {
	var res []*database.ProjectRole
	err := c.getDB(ctx).SelectContext(ctx, &res, `
		SELECT r.* FROM users_projects_roles upr
		JOIN project_roles r ON upr.project_role_id = r.id
		WHERE upr.user_id = $1 AND upr.project_id = $2
		UNION
		SELECT * FROM project_roles WHERE id IN (
			SELECT project_role_id FROM usergroups_projects_roles upr JOIN users_usergroups uug 
			ON upr.usergroup_id = uug.usergroup_id WHERE uug.user_id = $1 AND upr.project_id = $2
		)`, userID, projectID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindProjectRole(ctx context.Context, name string) (*database.ProjectRole, error) {
	role := &database.ProjectRole{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT * FROM project_roles WHERE lower(name)=lower($1)", name).StructScan(role)
	if err != nil {
		return nil, parseErr(err)
	}
	return role, nil
}

func (c *connection) InsertUsergroup(ctx context.Context, opts *database.InsertUsergroupOptions) (*database.Usergroup, error) {
	if err := database.Validate(opts); err != nil {
		return nil, err
	}

	res := &database.Usergroup{}
	err := c.getDB(ctx).QueryRowxContext(ctx, `
		INSERT INTO usergroups (org_id, name) VALUES ($1, $2) RETURNING *
	`, opts.OrgID, opts.Name).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) UpdateOrganizationAllUsergroup(ctx context.Context, orgID, groupID string) (*database.Organization, error) {
	res := &database.Organization{}
	err := c.getDB(ctx).QueryRowxContext(ctx, `
		UPDATE organizations SET all_usergroup_id = $1 WHERE id = $2 RETURNING *
	`, groupID, orgID).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) InsertUserInUsergroup(ctx context.Context, userID, groupID string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "INSERT INTO users_usergroups (user_id, usergroup_id) VALUES ($1, $2)", userID, groupID)
	if err != nil {
		return parseErr(err)
	}
	return nil
}

func (c *connection) DeleteUserFromUsergroup(ctx context.Context, userID, groupID string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM users_usergroups WHERE user_id = $1 AND usergroup_id = $2", userID, groupID)
	if err != nil {
		return parseErr(err)
	}
	return nil
}

func (c *connection) InsertProjectMemberUsergroup(ctx context.Context, groupID, projectID, roleID string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "INSERT INTO usergroups_projects_roles (usergroup_id, project_id, project_role_id) VALUES ($1, $2, $3)", groupID, projectID, roleID)
	if err != nil {
		return parseErr(err)
	}
	return nil
}

func (c *connection) FindUsersUsergroups(ctx context.Context, userID, orgID string) ([]*database.Usergroup, error) {
	var res []*database.Usergroup
	err := c.getDB(ctx).SelectContext(ctx, &res, `
		SELECT ug.* FROM usergroups ug JOIN users_usergroups uug ON ug.id = uug.usergroup_id
		WHERE uug.user_id = $1 AND ug.org_id = $2
	`, userID, orgID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindOrganizationsForUser(ctx context.Context, userID string) ([]*database.Organization, error) {
	var res []*database.Organization
	err := c.getDB(ctx).SelectContext(ctx, &res, `
		SELECT o.* FROM organizations o JOIN users_orgs_roles uor ON o.id = uor.org_id
		WHERE uor.user_id = $1
		UNION
		SELECT o.* FROM organizations o JOIN usergroups_orgs_roles ugor ON o.id = ugor.org_id
		JOIN users_usergroups uug ON ugor.usergroup_id = uug.usergroup_id
		WHERE uug.user_id = $1
		UNION
		SELECT o.* FROM organizations o JOIN projects p ON o.id = p.org_id
		JOIN users_projects_roles upr ON p.id = upr.project_id
		WHERE upr.user_id = $1
	`, userID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindProjectsForUser(ctx context.Context, userID string) ([]*database.Project, error) {
	var res []*database.Project
	err := c.getDB(ctx).SelectContext(ctx, &res, `
		SELECT p.* FROM projects p JOIN users_projects_roles upr ON p.id = upr.project_id
		WHERE upr.user_id = $1
		UNION
		SELECT p.* FROM projects p JOIN usergroups_projects_roles upgr ON p.id = upgr.project_id
		JOIN users_usergroups uug ON upgr.usergroup_id = uug.usergroup_id
		WHERE uug.user_id = $1
	`, userID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindProjectsForOrganization(ctx context.Context, orgID string) ([]*database.Project, error) {
	var res []*database.Project
	err := c.getDB(ctx).SelectContext(ctx, &res, "SELECT p.* FROM projects p WHERE p.org_id=$1 ORDER BY lower(p.name)", orgID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) CheckOrganizationProjectsHasMemberUser(ctx context.Context, orgID, userID string) (bool, error) {
	var res bool
	err := c.getDB(ctx).QueryRowxContext(ctx,
		"SELECT EXISTS (SELECT 1 FROM projects p JOIN users_projects_roles upr ON p.id = upr.project_id WHERE p.org_id = $1 AND upr.user_id = $2 limit 1)", orgID, userID).StructScan(&res)
	if err != nil {
		return false, parseErr(err)
	}
	return res, nil
}

func (c *connection) CheckOrganizationHasPublicProjects(ctx context.Context, orgID string) (bool, error) {
	var res bool
	err := c.getDB(ctx).QueryRowxContext(ctx,
		"SELECT EXISTS (SELECT 1 FROM projects p WHERE p.org_id = $1 AND p.public = true limit 1)", orgID).StructScan(&res)
	if err != nil {
		return false, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindProjectsForProjectMemberUser(ctx context.Context, orgID, userID string) ([]*database.Project, error) {
	var res []*database.Project
	err := c.getDB(ctx).SelectContext(ctx, &res, "SELECT p.* FROM projects p JOIN users_projects_roles upr ON p.id = upr.project_id WHERE p.org_id = $1 AND upr.user_id = $2", orgID, userID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindPublicProjectsInOrganization(ctx context.Context, orgID string) ([]*database.Project, error) {
	var res []*database.Project
	err := c.getDB(ctx).SelectContext(ctx, &res, "SELECT p.* FROM projects p WHERE p.org_id = $1 AND p.public = true", orgID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) InsertOrganizationMemberUserInvitation(ctx context.Context, email, invitedByID, orgID, roleID string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "INSERT INTO user_org_invites (email, invited_by_user_id, org_id, org_role_id) VALUES ($1, $2, $3, $4)", email, invitedByID, orgID, roleID)
	if err != nil {
		return parseErr(err)
	}
	return nil
}

func (c *connection) InsertProjectMemberUserInvitation(ctx context.Context, email, invitedByID, projectID, roleID string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "INSERT INTO user_project_invites (email, invited_by_user_id, project_id, project_role_id) VALUES ($1, $2, $3, $4)", email, invitedByID, projectID, roleID)
	if err != nil {
		return parseErr(err)
	}
	return nil
}

func (c *connection) FindOrganizationMemberInvitations(ctx context.Context, orgID string) ([]*database.UserInvite, error) {
	var res []*database.UserInvite
	err := c.getDB(ctx).SelectContext(ctx, &res, `
			SELECT uoi.email, ur.name as role, u.email as invited_by 
			FROM user_org_invites uoi JOIN org_roles ur ON uoi.org_role_id = ur.id JOIN users u ON uoi.invited_by_user_id = u.id WHERE uoi.org_id = $1`, orgID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindOrganizationMemberUserInvitations(ctx context.Context, userEmail string) ([]*database.OrganizationMemberUserInvitation, error) {
	var res []*database.OrganizationMemberUserInvitation
	err := c.getDB(ctx).SelectContext(ctx, &res, "SELECT * FROM user_org_invites WHERE lower(email) = lower($1)", userEmail)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindOrganizationMemberUserInvitation(ctx context.Context, orgID, userEmail string) (*database.OrganizationMemberUserInvitation, error) {
	res := &database.OrganizationMemberUserInvitation{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT * FROM user_org_invites WHERE lower(email) = lower($1) AND org_id = $2", userEmail, orgID).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindProjectMemberInvitations(ctx context.Context, projectID string) ([]*database.UserInvite, error) {
	var res []*database.UserInvite
	err := c.getDB(ctx).SelectContext(ctx, &res, `
			SELECT upi.email, ur.name as role, u.email as invited_by 
			FROM user_project_invites upi JOIN project_roles ur ON upi.project_role_id = ur.id JOIN users u ON upi.invited_by_user_id = u.id WHERE upi.project_id = $1`, projectID)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindProjectMemberUserInvitations(ctx context.Context, userEmail string) ([]*database.ProjectMemberUserInvitation, error) {
	var res []*database.ProjectMemberUserInvitation
	err := c.getDB(ctx).SelectContext(ctx, &res, "SELECT * FROM user_project_invites WHERE lower(email) = lower($1)", userEmail)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) FindProjectMemberUserInvitation(ctx context.Context, projectID, userEmail string) (*database.ProjectMemberUserInvitation, error) {
	res := &database.ProjectMemberUserInvitation{}
	err := c.getDB(ctx).QueryRowxContext(ctx, "SELECT * FROM user_project_invites WHERE lower(email) = lower($1) AND project_id = $2", userEmail, projectID).StructScan(res)
	if err != nil {
		return nil, parseErr(err)
	}
	return res, nil
}

func (c *connection) DeleteOrganizationMemberUserInvitation(ctx context.Context, id string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM user_org_invites WHERE id = $1", id)
	if err != nil {
		return parseErr(err)
	}
	return nil
}

func (c *connection) DeleteProjectMemberUserInvitation(ctx context.Context, id string) error {
	_, err := c.getDB(ctx).ExecContext(ctx, "DELETE FROM user_project_invites WHERE id = $1", id)
	if err != nil {
		return parseErr(err)
	}
	return nil
}

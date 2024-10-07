/**
 * Generated by orval v6.12.0 🍺
 * Do not edit manually.
 * rill/admin/v1/ai.proto
 * OpenAPI spec version: version not set
 */
export type AdminServiceSearchUsersParams = {
  emailPattern?: string;
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceListBookmarksParams = {
  projectId?: string;
  resourceKind?: string;
  resourceName?: string;
};

export type AdminServiceGetUserParams = { email?: string };

export type AdminServiceSudoGetResourceParams = {
  userId?: string;
  orgId?: string;
  projectId?: string;
  deploymentId?: string;
  instanceId?: string;
};

export type AdminServiceSearchProjectNamesParams = {
  namePattern?: string;
  annotations?: string;
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceGetReportMetaBodyAnnotations = {
  [key: string]: string;
};

export type AdminServiceGetReportMetaBody = {
  branch?: string;
  report?: string;
  annotations?: AdminServiceGetReportMetaBodyAnnotations;
  executionTime?: string;
};

export type AdminServicePullVirtualRepoParams = {
  branch?: string;
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceGetRepoMetaParams = { branch?: string };

export type AdminServiceGetAlertMetaBodyAnnotations = { [key: string]: string };

export type AdminServiceGetAlertMetaBody = {
  branch?: string;
  alert?: string;
  annotations?: AdminServiceGetAlertMetaBodyAnnotations;
  queryForUserId?: string;
  queryForUserEmail?: string;
};

export type AdminServiceUpdateServiceBody = {
  newName?: string;
};

export type AdminServiceCreateServiceParams = { name?: string };

export type AdminServiceUpdateProjectVariablesBodyVariables = {
  [key: string]: string;
};

export type AdminServiceUpdateProjectVariablesBody = {
  variables?: AdminServiceUpdateProjectVariablesBodyVariables;
};

export type AdminServiceUpdateProjectBody = {
  description?: string;
  public?: boolean;
  prodBranch?: string;
  githubUrl?: string;
  subpath?: string;
  archiveAssetId?: string;
  prodSlots?: string;
  provisioner?: string;
  newName?: string;
  prodTtlSeconds?: string;
  prodVersion?: string;
};

export type AdminServiceGetProjectParams = {
  accessTokenTtlSeconds?: number;
  issueSuperuserToken?: boolean;
};

export type AdminServiceCreateProjectBodyVariables = { [key: string]: string };

export type AdminServiceCreateProjectBody = {
  name?: string;
  description?: string;
  public?: boolean;
  provisioner?: string;
  prodOlapDriver?: string;
  prodOlapDsn?: string;
  prodSlots?: string;
  subpath?: string;
  prodBranch?: string;
  /** github_url is set for projects whose project files are stored in github. This is set to a github repo url.
Either github_url or archive_asset_id should be set. */
  githubUrl?: string;
  /** archive_asset_id is set for projects whose project files are not stored in github but are managed by rill. */
  archiveAssetId?: string;
  variables?: AdminServiceCreateProjectBodyVariables;
  prodVersion?: string;
};

export type AdminServiceListProjectsForOrganizationParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceCreateAssetBody = {
  type?: string;
  name?: string;
  extension?: string;
};

export type AdminServiceListUsergroupMemberUsersParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceEditUsergroupBody = {
  description?: string;
};

export type AdminServiceGetUsergroupParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceListOrganizationMemberUsergroupsParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceSearchProjectUsersParams = {
  emailQuery?: string;
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceIssueMagicAuthTokenBody = {
  /** TTL for the token in minutes. Set to 0 for no expiry. Defaults to no expiry. */
  ttlMinutes?: string;
  /** Type of resource to grant access to. Currently only supports "rill.runtime.v1.Explore". */
  resourceType?: string;
  /** Name of the resource to grant access to. */
  resourceName?: string;
  filter?: V1Expression;
  /** Optional list of fields to limit access to. If empty, no field access rule will be added.
This will be translated to a rill.runtime.v1.SecurityRuleFieldAccess, which currently applies to dimension and measure names for explores and metrics views. */
  fields?: string[];
  /** Optional state to store with the token. Can be fetched with GetCurrentMagicAuthToken. */
  state?: string;
  /** Optional public url title to store with the token. */
  title?: string;
};

export type AdminServiceListMagicAuthTokensParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceAddProjectMemberUserBody = {
  email?: string;
  role?: string;
};

export type AdminServiceListProjectMemberUsersParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceListProjectInvitesParams = {
  pageSize?: number;
  pageToken?: string;
};

/**
 * DEPRECATED: Additional parameters to set outright in the generated URL query.
 */
export type AdminServiceGetIFrameBodyQuery = { [key: string]: string };

/**
 * If set, will use the provided attributes outright.
 */
export type AdminServiceGetIFrameBodyAttributes = { [key: string]: any };

/**
 * GetIFrameRequest is the request payload for AdminService.GetIFrame.
 */
export type AdminServiceGetIFrameBody = {
  /** Branch to embed. If not set, the production branch is used. */
  branch?: string;
  /** TTL for the iframe's access token. If not set, defaults to 24 hours. */
  ttlSeconds?: number;
  /** If set, will use the attributes of the user with this ID. */
  userId?: string;
  /** If set, will generate attributes corresponding to a user with this email. */
  userEmail?: string;
  /** If set, will use the provided attributes outright. */
  attributes?: AdminServiceGetIFrameBodyAttributes;
  /** Type of resource to embed. If not set, defaults to "rill.runtime.v1.Explore". */
  type?: string;
  /** Deprecated: Alias for `type`. */
  kind?: string;
  /** Name of the resource to embed. This should identify a resource that is valid for embedding, such as a dashboard or component. */
  resource?: string;
  /** Theme to use for the embedded resource. */
  theme?: string;
  /** Navigation denotes whether navigation between different resources should be enabled in the embed. */
  navigation?: boolean;
  /** Blob containing UI state for rendering the initial embed. Not currently supported. */
  state?: string;
  /** DEPRECATED: Additional parameters to set outright in the generated URL query. */
  query?: AdminServiceGetIFrameBodyQuery;
};

export type AdminServiceGetDeploymentCredentialsBodyAttributes = {
  [key: string]: any;
};

export type AdminServiceGetDeploymentCredentialsBody = {
  branch?: string;
  ttlSeconds?: number;
  userId?: string;
  userEmail?: string;
  attributes?: AdminServiceGetDeploymentCredentialsBodyAttributes;
};

export type AdminServiceConnectProjectToGithubBody = {
  repo?: string;
  branch?: string;
  subpath?: string;
  force?: boolean;
};

export type AdminServiceListProjectMemberUsergroupsParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceRemoveOrganizationMemberUserParams = {
  keepProjectRoles?: boolean;
};

export type AdminServiceAddOrganizationMemberUserBody = {
  email?: string;
  role?: string;
  superuserForceAccess?: boolean;
};

export type AdminServiceListOrganizationMemberUsersParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceListOrganizationInvitesParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceGetPaymentsPortalURLParams = { returnUrl?: string };

export type AdminServiceUpdateOrganizationBody = {
  description?: string;
  newName?: string;
  displayName?: string;
  billingEmail?: string;
};

export type AdminServiceListOrganizationsParams = {
  pageSize?: number;
  pageToken?: string;
};

export type AdminServiceGetGithubRepoStatusParams = { githubUrl?: string };

export type AdminServiceTriggerRefreshSourcesBody = {
  sources?: string[];
};

export type AdminServiceCreateUsergroupBodyBody = {
  name?: string;
};

export type AdminServiceCreateReportBodyBody = {
  options?: V1ReportOptions;
};

export type AdminServiceCreateAlertBodyBody = {
  options?: V1AlertOptions;
};

export type AdminServiceCreateProjectWhitelistedDomainBodyBody = {
  domain?: string;
  role?: string;
};

export type AdminServiceSetOrganizationMemberUserRoleBodyBody = {
  role?: string;
};

export type AdminServiceTriggerReconcileBodyBody = { [key: string]: any };

export type AdminServiceUpdateBillingSubscriptionBodyBody = {
  planName?: string;
  superuserForceAccess?: boolean;
};

export interface V1WhitelistedDomain {
  domain?: string;
  role?: string;
}

export interface V1VirtualFile {
  path?: string;
  data?: string;
  deleted?: boolean;
  updatedOn?: string;
}

export interface V1Usergroup {
  groupId?: string;
  groupName?: string;
  groupDescription?: string;
  createdOn?: string;
  updatedOn?: string;
}

export interface V1UserQuotas {
  singleuserOrgs?: number;
  trialOrgs?: number;
}

export interface V1UserPreferences {
  timeZone?: string;
}

export interface V1UserInvite {
  email?: string;
  role?: string;
  invitedBy?: string;
}

export interface V1User {
  id?: string;
  email?: string;
  displayName?: string;
  photoUrl?: string;
  quotas?: V1UserQuotas;
  createdOn?: string;
  updatedOn?: string;
}

export interface V1UploadProjectAssetsResponse {
  [key: string]: any;
}

export interface V1UpdateUserPreferencesResponse {
  preferences?: V1UserPreferences;
}

export interface V1UpdateUserPreferencesRequest {
  preferences?: V1UserPreferences;
}

export interface V1UpdateServiceResponse {
  service?: V1Service;
}

export type V1UpdateProjectVariablesResponseVariables = {
  [key: string]: string;
};

export interface V1UpdateProjectVariablesResponse {
  variables?: V1UpdateProjectVariablesResponseVariables;
}

export interface V1UpdateProjectResponse {
  project?: V1Project;
}

export interface V1UpdateOrganizationResponse {
  organization?: V1Organization;
}

export interface V1UpdateBookmarkResponse {
  [key: string]: any;
}

export interface V1UpdateBookmarkRequest {
  bookmarkId?: string;
  displayName?: string;
  description?: string;
  data?: string;
  default?: boolean;
  shared?: boolean;
}

export interface V1UpdateBillingSubscriptionResponse {
  organization?: V1Organization;
  subscription?: V1Subscription;
}

export interface V1UnsubscribeReportResponse {
  [key: string]: any;
}

export interface V1UnsubscribeAlertResponse {
  [key: string]: any;
}

export interface V1TriggerReportResponse {
  [key: string]: any;
}

export interface V1TriggerRefreshSourcesResponse {
  [key: string]: any;
}

export interface V1TriggerRedeployResponse {
  [key: string]: any;
}

export interface V1TriggerRedeployRequest {
  organization?: string;
  project?: string;
  deploymentId?: string;
}

export interface V1TriggerReconcileResponse {
  [key: string]: any;
}

export interface V1SudoUpdateUserQuotasResponse {
  user?: V1User;
}

export interface V1SudoUpdateUserQuotasRequest {
  email?: string;
  singleuserOrgs?: number;
  trialOrgs?: number;
}

export interface V1SudoUpdateOrganizationQuotasResponse {
  organization?: V1Organization;
}

export interface V1SudoUpdateOrganizationQuotasRequest {
  organization?: string;
  projects?: number;
  deployments?: number;
  slotsTotal?: number;
  slotsPerDeployment?: number;
  outstandingInvites?: number;
  storageLimitBytesPerDeployment?: string;
}

export interface V1SudoUpdateOrganizationCustomDomainResponse {
  organization?: V1Organization;
}

export interface V1SudoUpdateOrganizationCustomDomainRequest {
  name?: string;
  customDomain?: string;
}

export interface V1SudoUpdateOrganizationBillingCustomerResponse {
  organization?: V1Organization;
  subscription?: V1Subscription;
}

export interface V1SudoUpdateOrganizationBillingCustomerRequest {
  organization?: string;
  billingCustomerId?: string;
}

export interface V1SudoUpdateAnnotationsResponse {
  project?: V1Project;
}

export type V1SudoUpdateAnnotationsRequestAnnotations = {
  [key: string]: string;
};

export interface V1SudoUpdateAnnotationsRequest {
  organization?: string;
  project?: string;
  annotations?: V1SudoUpdateAnnotationsRequestAnnotations;
}

export interface V1SudoIssueRuntimeManagerTokenResponse {
  token?: string;
}

export interface V1SudoIssueRuntimeManagerTokenRequest {
  host?: string;
}

export interface V1SudoGetResourceResponse {
  user?: V1User;
  org?: V1Organization;
  project?: V1Project;
  deployment?: V1Deployment;
  instance?: V1Deployment;
}

export interface V1SudoDeleteOrganizationBillingIssueResponse {
  [key: string]: any;
}

export interface V1Subscription {
  id?: string;
  plan?: V1BillingPlan;
  startDate?: string;
  endDate?: string;
  currentBillingCycleStartDate?: string;
  currentBillingCycleEndDate?: string;
  trialEndDate?: string;
}

export interface V1Subquery {
  dimension?: string;
  measures?: string[];
  where?: V1Expression;
  having?: V1Expression;
}

export interface V1SetSuperuserResponse {
  [key: string]: any;
}

export interface V1SetSuperuserRequest {
  email?: string;
  superuser?: boolean;
}

export interface V1SetProjectMemberUsergroupRoleResponse {
  [key: string]: any;
}

export interface V1SetProjectMemberUserRoleResponse {
  [key: string]: any;
}

export interface V1SetOrganizationMemberUsergroupRoleResponse {
  [key: string]: any;
}

export interface V1SetOrganizationMemberUserRoleResponse {
  [key: string]: any;
}

export interface V1ServiceToken {
  id?: string;
  createdOn?: string;
  expiresOn?: string;
}

export interface V1Service {
  id?: string;
  name?: string;
  orgId?: string;
  orgName?: string;
  createdOn?: string;
  updatedOn?: string;
}

export interface V1SearchUsersResponse {
  users?: V1User[];
  nextPageToken?: string;
}

export interface V1SearchProjectUsersResponse {
  users?: V1User[];
  nextPageToken?: string;
}

export interface V1SearchProjectNamesResponse {
  names?: string[];
  nextPageToken?: string;
}

export interface V1RevokeServiceAuthTokenResponse {
  [key: string]: any;
}

export interface V1RevokeMagicAuthTokenResponse {
  [key: string]: any;
}

export interface V1RevokeCurrentAuthTokenResponse {
  tokenId?: string;
}

export interface V1RequestProjectAccessResponse {
  [key: string]: any;
}

export interface V1ReportOptions {
  title?: string;
  refreshCron?: string;
  refreshTimeZone?: string;
  intervalDuration?: string;
  queryName?: string;
  queryArgsJson?: string;
  exportLimit?: string;
  exportFormat?: V1ExportFormat;
  emailRecipients?: string[];
  slackUsers?: string[];
  slackChannels?: string[];
  slackWebhooks?: string[];
  /** Annotation for the subpath of <UI host>/org/project to open for the report. */
  webOpenPath?: string;
  /** Annotation for the base64-encoded UI state to open for the report. */
  webOpenState?: string;
}

export interface V1RenewBillingSubscriptionResponse {
  organization?: V1Organization;
  subscription?: V1Subscription;
}

export interface V1RenameUsergroupResponse {
  [key: string]: any;
}

export interface V1RemoveWhitelistedDomainResponse {
  [key: string]: any;
}

export interface V1RemoveUsergroupMemberUserResponse {
  [key: string]: any;
}

export interface V1RemoveProjectWhitelistedDomainResponse {
  [key: string]: any;
}

export interface V1RemoveProjectMemberUsergroupResponse {
  [key: string]: any;
}

export interface V1RemoveProjectMemberUserResponse {
  [key: string]: any;
}

export interface V1RemoveOrganizationMemberUsergroupResponse {
  [key: string]: any;
}

export interface V1RemoveOrganizationMemberUserResponse {
  [key: string]: any;
}

export interface V1RemoveBookmarkResponse {
  [key: string]: any;
}

export interface V1RedeployProjectResponse {
  [key: string]: any;
}

export interface V1RecordEventsResponse {
  [key: string]: any;
}

export type V1RecordEventsRequestEventsItem = { [key: string]: any };

export interface V1RecordEventsRequest {
  events?: V1RecordEventsRequestEventsItem[];
}

export interface V1Quotas {
  projects?: string;
  deployments?: string;
  slotsTotal?: string;
  slotsPerDeployment?: string;
  outstandingInvites?: string;
  storageLimitBytesPerDeployment?: string;
}

export interface V1PullVirtualRepoResponse {
  files?: V1VirtualFile[];
  nextPageToken?: string;
}

export interface V1ProjectPermissions {
  readProject?: boolean;
  manageProject?: boolean;
  readProd?: boolean;
  readProdStatus?: boolean;
  manageProd?: boolean;
  readDev?: boolean;
  readDevStatus?: boolean;
  manageDev?: boolean;
  readProjectMembers?: boolean;
  manageProjectMembers?: boolean;
  createMagicAuthTokens?: boolean;
  manageMagicAuthTokens?: boolean;
  createReports?: boolean;
  manageReports?: boolean;
  createAlerts?: boolean;
  manageAlerts?: boolean;
  createBookmarks?: boolean;
  manageBookmarks?: boolean;
}

export type V1ProjectAnnotations = { [key: string]: string };

export interface V1Project {
  id?: string;
  name?: string;
  orgId?: string;
  orgName?: string;
  description?: string;
  public?: boolean;
  createdByUserId?: string;
  provisioner?: string;
  githubUrl?: string;
  subpath?: string;
  prodBranch?: string;
  archiveAssetId?: string;
  prodOlapDriver?: string;
  prodOlapDsn?: string;
  prodSlots?: string;
  prodDeploymentId?: string;
  /** Note: Does NOT incorporate the parent org's custom domain. */
  frontendUrl?: string;
  prodTtlSeconds?: string;
  annotations?: V1ProjectAnnotations;
  prodVersion?: string;
  createdOn?: string;
  updatedOn?: string;
}

export interface V1PingResponse {
  version?: string;
  time?: string;
}

export interface V1OrganizationQuotas {
  projects?: number;
  deployments?: number;
  slotsTotal?: number;
  slotsPerDeployment?: number;
  outstandingInvites?: number;
  storageLimitBytesPerDeployment?: string;
}

export interface V1OrganizationPermissions {
  readOrg?: boolean;
  manageOrg?: boolean;
  readProjects?: boolean;
  createProjects?: boolean;
  manageProjects?: boolean;
  readOrgMembers?: boolean;
  manageOrgMembers?: boolean;
}

export interface V1Organization {
  id?: string;
  name?: string;
  displayName?: string;
  description?: string;
  customDomain?: string;
  quotas?: V1OrganizationQuotas;
  billingCustomerId?: string;
  paymentCustomerId?: string;
  billingEmail?: string;
  createdOn?: string;
  updatedOn?: string;
}

export type V1Operation = (typeof V1Operation)[keyof typeof V1Operation];

// eslint-disable-next-line @typescript-eslint/no-redeclare
export const V1Operation = {
  OPERATION_UNSPECIFIED: "OPERATION_UNSPECIFIED",
  OPERATION_EQ: "OPERATION_EQ",
  OPERATION_NEQ: "OPERATION_NEQ",
  OPERATION_LT: "OPERATION_LT",
  OPERATION_LTE: "OPERATION_LTE",
  OPERATION_GT: "OPERATION_GT",
  OPERATION_GTE: "OPERATION_GTE",
  OPERATION_OR: "OPERATION_OR",
  OPERATION_AND: "OPERATION_AND",
  OPERATION_IN: "OPERATION_IN",
  OPERATION_NIN: "OPERATION_NIN",
  OPERATION_LIKE: "OPERATION_LIKE",
  OPERATION_NLIKE: "OPERATION_NLIKE",
} as const;

export interface V1MemberUsergroup {
  groupId?: string;
  groupName?: string;
  roleName?: string;
  createdOn?: string;
  updatedOn?: string;
}

export interface V1MemberUser {
  userId?: string;
  userEmail?: string;
  userName?: string;
  roleName?: string;
  createdOn?: string;
  updatedOn?: string;
}

export type V1MagicAuthTokenAttributes = { [key: string]: any };

export interface V1MagicAuthToken {
  id?: string;
  projectId?: string;
  url?: string;
  token?: string;
  createdOn?: string;
  expiresOn?: string;
  usedOn?: string;
  createdByUserId?: string;
  createdByUserEmail?: string;
  attributes?: V1MagicAuthTokenAttributes;
  resourceType?: string;
  resourceName?: string;
  filter?: V1Expression;
  fields?: string[];
  state?: string;
  title?: string;
}

export interface V1ListWhitelistedDomainsResponse {
  domains?: V1WhitelistedDomain[];
}

export interface V1ListUsergroupMemberUsersResponse {
  members?: V1MemberUser[];
  nextPageToken?: string;
}

export interface V1ListSuperusersResponse {
  users?: V1User[];
}

export interface V1ListServicesResponse {
  services?: V1Service[];
}

export interface V1ListServiceAuthTokensResponse {
  tokens?: V1ServiceToken[];
}

export interface V1ListPublicBillingPlansResponse {
  plans?: V1BillingPlan[];
}

export interface V1ListProjectsForOrganizationResponse {
  projects?: V1Project[];
  nextPageToken?: string;
}

export interface V1ListProjectWhitelistedDomainsResponse {
  domains?: V1WhitelistedDomain[];
}

export interface V1ListProjectMemberUsersResponse {
  members?: V1MemberUser[];
  nextPageToken?: string;
}

export interface V1ListProjectMemberUsergroupsResponse {
  members?: V1MemberUsergroup[];
  nextPageToken?: string;
}

export interface V1ListProjectInvitesResponse {
  invites?: V1UserInvite[];
  nextPageToken?: string;
}

export interface V1ListOrganizationsResponse {
  organizations?: V1Organization[];
  nextPageToken?: string;
}

export interface V1ListOrganizationMemberUsersResponse {
  members?: V1MemberUser[];
  nextPageToken?: string;
}

export interface V1ListOrganizationMemberUsergroupsResponse {
  members?: V1MemberUsergroup[];
  nextPageToken?: string;
}

export interface V1ListOrganizationInvitesResponse {
  invites?: V1UserInvite[];
  nextPageToken?: string;
}

export interface V1ListOrganizationBillingIssuesResponse {
  issues?: V1BillingIssue[];
}

export interface V1ListMagicAuthTokensResponse {
  tokens?: V1MagicAuthToken[];
  nextPageToken?: string;
}

export interface V1ListGithubUserReposResponse {
  repos?: ListGithubUserReposResponseRepo[];
}

export interface V1ListBookmarksResponse {
  bookmarks?: V1Bookmark[];
}

export interface V1LeaveOrganizationResponse {
  [key: string]: any;
}

export interface V1IssueServiceAuthTokenResponse {
  token?: string;
}

export interface V1IssueRepresentativeAuthTokenResponse {
  token?: string;
}

export interface V1IssueRepresentativeAuthTokenRequest {
  email?: string;
  ttlMinutes?: string;
}

export interface V1IssueMagicAuthTokenResponse {
  token?: string;
  url?: string;
}

export interface V1HibernateProjectResponse {
  [key: string]: any;
}

export type V1GithubPermission =
  (typeof V1GithubPermission)[keyof typeof V1GithubPermission];

// eslint-disable-next-line @typescript-eslint/no-redeclare
export const V1GithubPermission = {
  GITHUB_PERMISSION_UNSPECIFIED: "GITHUB_PERMISSION_UNSPECIFIED",
  GITHUB_PERMISSION_READ: "GITHUB_PERMISSION_READ",
  GITHUB_PERMISSION_WRITE: "GITHUB_PERMISSION_WRITE",
} as const;

export interface V1GetUsergroupResponse {
  usergroup?: V1Usergroup;
  nextPageToken?: string;
}

export interface V1GetUserResponse {
  user?: V1User;
}

export interface V1GetReportMetaResponse {
  openUrl?: string;
  exportUrl?: string;
  editUrl?: string;
}

export interface V1GetRepoMetaResponse {
  gitUrl?: string;
  gitUrlExpiresOn?: string;
  gitSubpath?: string;
  archiveDownloadUrl?: string;
}

export type V1GetProjectVariablesResponseVariables = { [key: string]: string };

export interface V1GetProjectVariablesResponse {
  variables?: V1GetProjectVariablesResponseVariables;
}

export interface V1GetProjectResponse {
  project?: V1Project;
  prodDeployment?: V1Deployment;
  jwt?: string;
  projectPermissions?: V1ProjectPermissions;
}

export interface V1GetProjectByIDResponse {
  project?: V1Project;
}

export interface V1GetProjectAccessRequestResponse {
  email?: string;
}

export interface V1GetPaymentsPortalURLResponse {
  url?: string;
}

export interface V1GetOrganizationResponse {
  organization?: V1Organization;
  permissions?: V1OrganizationPermissions;
}

export interface V1GetOrganizationNameForDomainResponse {
  name?: string;
}

export interface V1GetIFrameResponse {
  iframeSrc?: string;
  runtimeHost?: string;
  instanceId?: string;
  accessToken?: string;
  ttlSeconds?: number;
}

export type V1GetGithubUserStatusResponseOrganizationInstallationPermissions = {
  [key: string]: V1GithubPermission;
};

export interface V1GetGithubUserStatusResponse {
  hasAccess?: boolean;
  grantAccessUrl?: string;
  accessToken?: string;
  account?: string;
  userInstallationPermission?: V1GithubPermission;
  organizationInstallationPermissions?: V1GetGithubUserStatusResponseOrganizationInstallationPermissions;
  /** DEPRECATED: Use organization_installation_permissions instead. */
  organizations?: string[];
}

export interface V1GetGithubRepoStatusResponse {
  hasAccess?: boolean;
  grantAccessUrl?: string;
  defaultBranch?: string;
}

export interface V1GetDeploymentCredentialsResponse {
  runtimeHost?: string;
  instanceId?: string;
  accessToken?: string;
  ttlSeconds?: number;
}

export interface V1GetCurrentUserResponse {
  user?: V1User;
  preferences?: V1UserPreferences;
}

export interface V1GetCurrentMagicAuthTokenResponse {
  token?: V1MagicAuthToken;
}

export interface V1GetCloneCredentialsResponse {
  gitRepoUrl?: string;
  gitUsername?: string;
  gitPassword?: string;
  gitSubpath?: string;
  gitProdBranch?: string;
  archiveDownloadUrl?: string;
}

export interface V1GetBookmarkResponse {
  bookmark?: V1Bookmark;
}

export interface V1GetBillingSubscriptionResponse {
  organization?: V1Organization;
  subscription?: V1Subscription;
  billingPortalUrl?: string;
}

export interface V1GetAlertYAMLResponse {
  yaml?: string;
}

export type V1GetAlertMetaResponseQueryForAttributes = { [key: string]: any };

export interface V1GetAlertMetaResponse {
  openUrl?: string;
  editUrl?: string;
  queryForAttributes?: V1GetAlertMetaResponseQueryForAttributes;
}

export interface V1GenerateReportYAMLResponse {
  yaml?: string;
}

export interface V1GenerateAlertYAMLResponse {
  yaml?: string;
}

export interface V1Expression {
  ident?: string;
  val?: unknown;
  cond?: V1Condition;
  subquery?: V1Subquery;
}

export type V1ExportFormat =
  (typeof V1ExportFormat)[keyof typeof V1ExportFormat];

// eslint-disable-next-line @typescript-eslint/no-redeclare
export const V1ExportFormat = {
  EXPORT_FORMAT_UNSPECIFIED: "EXPORT_FORMAT_UNSPECIFIED",
  EXPORT_FORMAT_CSV: "EXPORT_FORMAT_CSV",
  EXPORT_FORMAT_XLSX: "EXPORT_FORMAT_XLSX",
  EXPORT_FORMAT_PARQUET: "EXPORT_FORMAT_PARQUET",
} as const;

export interface V1EditUsergroupResponse {
  [key: string]: any;
}

export interface V1EditReportResponse {
  [key: string]: any;
}

export interface V1EditAlertResponse {
  [key: string]: any;
}

export type V1DeploymentStatus =
  (typeof V1DeploymentStatus)[keyof typeof V1DeploymentStatus];

// eslint-disable-next-line @typescript-eslint/no-redeclare
export const V1DeploymentStatus = {
  DEPLOYMENT_STATUS_UNSPECIFIED: "DEPLOYMENT_STATUS_UNSPECIFIED",
  DEPLOYMENT_STATUS_PENDING: "DEPLOYMENT_STATUS_PENDING",
  DEPLOYMENT_STATUS_OK: "DEPLOYMENT_STATUS_OK",
  DEPLOYMENT_STATUS_ERROR: "DEPLOYMENT_STATUS_ERROR",
} as const;

export interface V1Deployment {
  id?: string;
  projectId?: string;
  slots?: string;
  branch?: string;
  runtimeHost?: string;
  runtimeInstanceId?: string;
  status?: V1DeploymentStatus;
  statusMessage?: string;
  createdOn?: string;
  updatedOn?: string;
}

export interface V1DenyProjectAccessResponse {
  [key: string]: any;
}

export interface V1DeleteUsergroupResponse {
  [key: string]: any;
}

export interface V1DeleteServiceResponse {
  service?: V1Service;
}

export interface V1DeleteReportResponse {
  [key: string]: any;
}

export interface V1DeleteProjectResponse {
  id?: string;
}

export interface V1DeleteOrganizationResponse {
  [key: string]: any;
}

export interface V1DeleteAlertResponse {
  [key: string]: any;
}

export interface V1CreateWhitelistedDomainResponse {
  [key: string]: any;
}

export interface V1CreateUsergroupResponse {
  [key: string]: any;
}

export interface V1CreateServiceResponse {
  service?: V1Service;
}

export interface V1CreateReportResponse {
  name?: string;
}

export interface V1CreateProjectWhitelistedDomainResponse {
  [key: string]: any;
}

export interface V1CreateProjectResponse {
  project?: V1Project;
}

export interface V1CreateOrganizationResponse {
  organization?: V1Organization;
}

export interface V1CreateOrganizationRequest {
  name?: string;
  description?: string;
}

export interface V1CreateBookmarkResponse {
  bookmark?: V1Bookmark;
}

export interface V1CreateBookmarkRequest {
  displayName?: string;
  description?: string;
  data?: string;
  resourceKind?: string;
  resourceName?: string;
  projectId?: string;
  default?: boolean;
  shared?: boolean;
}

export type V1CreateAssetResponseSigningHeaders = { [key: string]: string };

export interface V1CreateAssetResponse {
  assetId?: string;
  signedUrl?: string;
  signingHeaders?: V1CreateAssetResponseSigningHeaders;
}

export interface V1CreateAlertResponse {
  name?: string;
}

export interface V1ConnectProjectToGithubResponse {
  [key: string]: any;
}

export interface V1Condition {
  op?: V1Operation;
  exprs?: V1Expression[];
}

export interface V1CompletionMessage {
  role?: string;
  data?: string;
}

export interface V1CompleteResponse {
  message?: V1CompletionMessage;
}

export interface V1CompleteRequest {
  messages?: V1CompletionMessage[];
}

export interface V1CancelBillingSubscriptionResponse {
  [key: string]: any;
}

export interface V1Bookmark {
  id?: string;
  displayName?: string;
  description?: string;
  data?: string;
  resourceKind?: string;
  resourceName?: string;
  projectId?: string;
  userId?: string;
  default?: boolean;
  shared?: boolean;
  createdOn?: string;
  updatedOn?: string;
}

export interface V1BillingPlan {
  id?: string;
  name?: string;
  displayName?: string;
  description?: string;
  trialPeriodDays?: number;
  default?: boolean;
  quotas?: V1Quotas;
}

export type V1BillingIssueType =
  (typeof V1BillingIssueType)[keyof typeof V1BillingIssueType];

// eslint-disable-next-line @typescript-eslint/no-redeclare
export const V1BillingIssueType = {
  BILLING_ISSUE_TYPE_UNSPECIFIED: "BILLING_ISSUE_TYPE_UNSPECIFIED",
  BILLING_ISSUE_TYPE_ON_TRIAL: "BILLING_ISSUE_TYPE_ON_TRIAL",
  BILLING_ISSUE_TYPE_TRIAL_ENDED: "BILLING_ISSUE_TYPE_TRIAL_ENDED",
  BILLING_ISSUE_TYPE_NO_PAYMENT_METHOD: "BILLING_ISSUE_TYPE_NO_PAYMENT_METHOD",
  BILLING_ISSUE_TYPE_NO_BILLABLE_ADDRESS:
    "BILLING_ISSUE_TYPE_NO_BILLABLE_ADDRESS",
  BILLING_ISSUE_TYPE_PAYMENT_FAILED: "BILLING_ISSUE_TYPE_PAYMENT_FAILED",
  BILLING_ISSUE_TYPE_SUBSCRIPTION_CANCELLED:
    "BILLING_ISSUE_TYPE_SUBSCRIPTION_CANCELLED",
  BILLING_ISSUE_TYPE_NEVER_SUBSCRIBED: "BILLING_ISSUE_TYPE_NEVER_SUBSCRIBED",
} as const;

export interface V1BillingIssueMetadataTrialEnded {
  gracePeriodEndDate?: string;
}

export interface V1BillingIssueMetadataSubscriptionCancelled {
  endDate?: string;
}

export interface V1BillingIssueMetadataPaymentFailedMeta {
  invoiceId?: string;
  invoiceNumber?: string;
  invoiceUrl?: string;
  amountDue?: string;
  currency?: string;
  dueDate?: string;
  failedOn?: string;
  gracePeriodEndDate?: string;
}

export interface V1BillingIssueMetadataPaymentFailed {
  invoices?: V1BillingIssueMetadataPaymentFailedMeta[];
}

export interface V1BillingIssueMetadataOnTrial {
  endDate?: string;
}

export interface V1BillingIssueMetadataNoPaymentMethod {
  [key: string]: any;
}

export interface V1BillingIssueMetadataNoBillableAddress {
  [key: string]: any;
}

export interface V1BillingIssueMetadataNeverSubscribed {
  [key: string]: any;
}

export interface V1BillingIssueMetadata {
  onTrial?: V1BillingIssueMetadataOnTrial;
  trialEnded?: V1BillingIssueMetadataTrialEnded;
  noPaymentMethod?: V1BillingIssueMetadataNoPaymentMethod;
  noBillableAddress?: V1BillingIssueMetadataNoBillableAddress;
  paymentFailed?: V1BillingIssueMetadataPaymentFailed;
  subscriptionCancelled?: V1BillingIssueMetadataSubscriptionCancelled;
  neverSubscribed?: V1BillingIssueMetadataNeverSubscribed;
}

export type V1BillingIssueLevel =
  (typeof V1BillingIssueLevel)[keyof typeof V1BillingIssueLevel];

// eslint-disable-next-line @typescript-eslint/no-redeclare
export const V1BillingIssueLevel = {
  BILLING_ISSUE_LEVEL_UNSPECIFIED: "BILLING_ISSUE_LEVEL_UNSPECIFIED",
  BILLING_ISSUE_LEVEL_WARNING: "BILLING_ISSUE_LEVEL_WARNING",
  BILLING_ISSUE_LEVEL_ERROR: "BILLING_ISSUE_LEVEL_ERROR",
} as const;

export interface V1BillingIssue {
  organization?: string;
  type?: V1BillingIssueType;
  level?: V1BillingIssueLevel;
  metadata?: V1BillingIssueMetadata;
  eventTime?: string;
  createdOn?: string;
}

export interface V1ApproveProjectAccessResponse {
  [key: string]: any;
}

export type V1AlertOptionsResolverProperties = { [key: string]: any };

export interface V1AlertOptions {
  title?: string;
  intervalDuration?: string;
  resolver?: string;
  resolverProperties?: V1AlertOptionsResolverProperties;
  /** DEPRECATED: Use resolver and resolver_properties instead. */
  queryName?: string;
  /** DEPRECATED: Use resolver and resolver_properties instead. */
  queryArgsJson?: string;
  metricsViewName?: string;
  renotify?: boolean;
  renotifyAfterSeconds?: number;
  emailRecipients?: string[];
  slackUsers?: string[];
  slackChannels?: string[];
  slackWebhooks?: string[];
  /** Annotation for the subpath of <UI host>/org/project to open for the report. */
  webOpenPath?: string;
  /** Annotation for the base64-encoded UI state to open for the report. */
  webOpenState?: string;
}

export interface V1AddUsergroupMemberUserResponse {
  [key: string]: any;
}

export interface V1AddProjectMemberUsergroupResponse {
  [key: string]: any;
}

export interface V1AddProjectMemberUserResponse {
  pendingSignup?: boolean;
}

export interface V1AddOrganizationMemberUsergroupResponse {
  [key: string]: any;
}

export interface V1AddOrganizationMemberUserResponse {
  pendingSignup?: boolean;
}

export interface RpcStatus {
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

/**
 * `NullValue` is a singleton enumeration to represent the null value for the
`Value` type union.

The JSON representation for `NullValue` is JSON `null`.

 - NULL_VALUE: Null value.
 */
export type ProtobufNullValue =
  (typeof ProtobufNullValue)[keyof typeof ProtobufNullValue];

// eslint-disable-next-line @typescript-eslint/no-redeclare
export const ProtobufNullValue = {
  NULL_VALUE: "NULL_VALUE",
} as const;

export interface ProtobufAny {
  "@type"?: string;
  [key: string]: unknown;
}

export interface ListGithubUserReposResponseRepo {
  name?: string;
  owner?: string;
  description?: string;
  url?: string;
  defaultBranch?: string;
}

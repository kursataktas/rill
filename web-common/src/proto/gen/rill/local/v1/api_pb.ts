// @generated by protoc-gen-es v1.10.0 with parameter "target=ts"
// @generated from file rill/local/v1/api.proto (package rill.local.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, Timestamp } from "@bufbuild/protobuf";
import { GithubPermission, User } from "../../admin/v1/api_pb.js";

/**
 * @generated from message rill.local.v1.PingRequest
 */
export class PingRequest extends Message<PingRequest> {
  constructor(data?: PartialMessage<PingRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.PingRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PingRequest {
    return new PingRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PingRequest {
    return new PingRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PingRequest {
    return new PingRequest().fromJsonString(jsonString, options);
  }

  static equals(a: PingRequest | PlainMessage<PingRequest> | undefined, b: PingRequest | PlainMessage<PingRequest> | undefined): boolean {
    return proto3.util.equals(PingRequest, a, b);
  }
}

/**
 * @generated from message rill.local.v1.PingResponse
 */
export class PingResponse extends Message<PingResponse> {
  /**
   * @generated from field: google.protobuf.Timestamp time = 1;
   */
  time?: Timestamp;

  constructor(data?: PartialMessage<PingResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.PingResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "time", kind: "message", T: Timestamp },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PingResponse {
    return new PingResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PingResponse {
    return new PingResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PingResponse {
    return new PingResponse().fromJsonString(jsonString, options);
  }

  static equals(a: PingResponse | PlainMessage<PingResponse> | undefined, b: PingResponse | PlainMessage<PingResponse> | undefined): boolean {
    return proto3.util.equals(PingResponse, a, b);
  }
}

/**
 * @generated from message rill.local.v1.GetMetadataRequest
 */
export class GetMetadataRequest extends Message<GetMetadataRequest> {
  constructor(data?: PartialMessage<GetMetadataRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.GetMetadataRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetMetadataRequest {
    return new GetMetadataRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetMetadataRequest {
    return new GetMetadataRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetMetadataRequest {
    return new GetMetadataRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetMetadataRequest | PlainMessage<GetMetadataRequest> | undefined, b: GetMetadataRequest | PlainMessage<GetMetadataRequest> | undefined): boolean {
    return proto3.util.equals(GetMetadataRequest, a, b);
  }
}

/**
 * @generated from message rill.local.v1.GetMetadataResponse
 */
export class GetMetadataResponse extends Message<GetMetadataResponse> {
  /**
   * @generated from field: string instance_id = 1;
   */
  instanceId = "";

  /**
   * @generated from field: string project_path = 2;
   */
  projectPath = "";

  /**
   * @generated from field: string install_id = 3;
   */
  installId = "";

  /**
   * @generated from field: string user_id = 4;
   */
  userId = "";

  /**
   * @generated from field: string version = 5;
   */
  version = "";

  /**
   * @generated from field: string build_commit = 6;
   */
  buildCommit = "";

  /**
   * @generated from field: string build_time = 7;
   */
  buildTime = "";

  /**
   * @generated from field: bool is_dev = 8;
   */
  isDev = false;

  /**
   * @generated from field: bool analytics_enabled = 9;
   */
  analyticsEnabled = false;

  /**
   * @generated from field: bool readonly = 10;
   */
  readonly = false;

  /**
   * @generated from field: int32 grpc_port = 11;
   */
  grpcPort = 0;

  constructor(data?: PartialMessage<GetMetadataResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.GetMetadataResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "instance_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "project_path", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "install_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "user_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "version", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "build_commit", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 7, name: "build_time", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 8, name: "is_dev", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 9, name: "analytics_enabled", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 10, name: "readonly", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 11, name: "grpc_port", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetMetadataResponse {
    return new GetMetadataResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetMetadataResponse {
    return new GetMetadataResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetMetadataResponse {
    return new GetMetadataResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetMetadataResponse | PlainMessage<GetMetadataResponse> | undefined, b: GetMetadataResponse | PlainMessage<GetMetadataResponse> | undefined): boolean {
    return proto3.util.equals(GetMetadataResponse, a, b);
  }
}

/**
 * @generated from message rill.local.v1.GetVersionRequest
 */
export class GetVersionRequest extends Message<GetVersionRequest> {
  constructor(data?: PartialMessage<GetVersionRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.GetVersionRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetVersionRequest {
    return new GetVersionRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetVersionRequest {
    return new GetVersionRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetVersionRequest {
    return new GetVersionRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetVersionRequest | PlainMessage<GetVersionRequest> | undefined, b: GetVersionRequest | PlainMessage<GetVersionRequest> | undefined): boolean {
    return proto3.util.equals(GetVersionRequest, a, b);
  }
}

/**
 * @generated from message rill.local.v1.GetVersionResponse
 */
export class GetVersionResponse extends Message<GetVersionResponse> {
  /**
   * @generated from field: string current = 1;
   */
  current = "";

  /**
   * @generated from field: string latest = 2;
   */
  latest = "";

  constructor(data?: PartialMessage<GetVersionResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.GetVersionResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "current", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "latest", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetVersionResponse {
    return new GetVersionResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetVersionResponse {
    return new GetVersionResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetVersionResponse {
    return new GetVersionResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetVersionResponse | PlainMessage<GetVersionResponse> | undefined, b: GetVersionResponse | PlainMessage<GetVersionResponse> | undefined): boolean {
    return proto3.util.equals(GetVersionResponse, a, b);
  }
}

/**
 * @generated from message rill.local.v1.DeployValidationRequest
 */
export class DeployValidationRequest extends Message<DeployValidationRequest> {
  constructor(data?: PartialMessage<DeployValidationRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.DeployValidationRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeployValidationRequest {
    return new DeployValidationRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeployValidationRequest {
    return new DeployValidationRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeployValidationRequest {
    return new DeployValidationRequest().fromJsonString(jsonString, options);
  }

  static equals(a: DeployValidationRequest | PlainMessage<DeployValidationRequest> | undefined, b: DeployValidationRequest | PlainMessage<DeployValidationRequest> | undefined): boolean {
    return proto3.util.equals(DeployValidationRequest, a, b);
  }
}

/**
 * @generated from message rill.local.v1.DeployValidationResponse
 */
export class DeployValidationResponse extends Message<DeployValidationResponse> {
  /**
   * if true below fields are relevant after login
   *
   * @generated from field: bool is_authenticated = 1;
   */
  isAuthenticated = false;

  /**
   * redirect to this if is_authenticated is false
   *
   * @generated from field: string login_url = 2;
   */
  loginUrl = "";

  /**
   * if true below fields are relevant after github install
   *
   * @generated from field: bool is_github_connected = 3;
   */
  isGithubConnected = false;

  /**
   * redirect to this if is_github_connected or is_github_repo_access_granted is false
   *
   * @generated from field: string github_grant_access_url = 4;
   */
  githubGrantAccessUrl = "";

  /**
   * @generated from field: string github_user_name = 5;
   */
  githubUserName = "";

  /**
   * if unspecified then github app not installed on user account
   *
   * @generated from field: rill.admin.v1.GithubPermission github_user_permission = 6;
   */
  githubUserPermission = GithubPermission.UNSPECIFIED;

  /**
   * @generated from field: map<string, rill.admin.v1.GithubPermission> github_organization_permissions = 7;
   */
  githubOrganizationPermissions: { [key: string]: GithubPermission } = {};

  /**
   * @generated from field: bool is_github_repo = 8;
   */
  isGithubRepo = false;

  /**
   * only applicable when is_github_repo is true
   *
   * @generated from field: bool is_github_remote_found = 9;
   */
  isGithubRemoteFound = false;

  /**
   * relevant only when is_github_repo is true and remote found, if false redirect to github_grant_access_url
   *
   * @generated from field: bool is_github_repo_access_granted = 10;
   */
  isGithubRepoAccessGranted = false;

  /**
   * only applicable when is_github_repo is true and remote found
   *
   * @generated from field: string github_url = 11;
   */
  githubUrl = "";

  /**
   * only applicable when is_github_repo is true and remote found
   *
   * @generated from field: optional bool has_uncommitted_changes = 12;
   */
  hasUncommittedChanges?: boolean;

  /**
   * only applicable when user does not have any orgs
   *
   * @generated from field: bool rill_org_exists_as_github_user_name = 13;
   */
  rillOrgExistsAsGithubUserName = false;

  /**
   * @generated from field: repeated string rill_user_orgs = 14;
   */
  rillUserOrgs: string[] = [];

  /**
   * @generated from field: string local_project_name = 15;
   */
  localProjectName = "";

  /**
   * @generated from field: string deployed_project_id = 16;
   */
  deployedProjectId = "";

  constructor(data?: PartialMessage<DeployValidationResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.DeployValidationResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "is_authenticated", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 2, name: "login_url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "is_github_connected", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 4, name: "github_grant_access_url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "github_user_name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "github_user_permission", kind: "enum", T: proto3.getEnumType(GithubPermission) },
    { no: 7, name: "github_organization_permissions", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "enum", T: proto3.getEnumType(GithubPermission)} },
    { no: 8, name: "is_github_repo", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 9, name: "is_github_remote_found", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 10, name: "is_github_repo_access_granted", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 11, name: "github_url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 12, name: "has_uncommitted_changes", kind: "scalar", T: 8 /* ScalarType.BOOL */, opt: true },
    { no: 13, name: "rill_org_exists_as_github_user_name", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 14, name: "rill_user_orgs", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 15, name: "local_project_name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 16, name: "deployed_project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeployValidationResponse {
    return new DeployValidationResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeployValidationResponse {
    return new DeployValidationResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeployValidationResponse {
    return new DeployValidationResponse().fromJsonString(jsonString, options);
  }

  static equals(a: DeployValidationResponse | PlainMessage<DeployValidationResponse> | undefined, b: DeployValidationResponse | PlainMessage<DeployValidationResponse> | undefined): boolean {
    return proto3.util.equals(DeployValidationResponse, a, b);
  }
}

/**
 * @generated from message rill.local.v1.PushToGithubRequest
 */
export class PushToGithubRequest extends Message<PushToGithubRequest> {
  /**
   * @generated from field: string account = 1;
   */
  account = "";

  /**
   * @generated from field: string repo = 2;
   */
  repo = "";

  constructor(data?: PartialMessage<PushToGithubRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.PushToGithubRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "account", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "repo", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PushToGithubRequest {
    return new PushToGithubRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PushToGithubRequest {
    return new PushToGithubRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PushToGithubRequest {
    return new PushToGithubRequest().fromJsonString(jsonString, options);
  }

  static equals(a: PushToGithubRequest | PlainMessage<PushToGithubRequest> | undefined, b: PushToGithubRequest | PlainMessage<PushToGithubRequest> | undefined): boolean {
    return proto3.util.equals(PushToGithubRequest, a, b);
  }
}

/**
 * @generated from message rill.local.v1.PushToGithubResponse
 */
export class PushToGithubResponse extends Message<PushToGithubResponse> {
  /**
   * @generated from field: string github_url = 1;
   */
  githubUrl = "";

  /**
   * @generated from field: string account = 2;
   */
  account = "";

  /**
   * @generated from field: string repo = 3;
   */
  repo = "";

  constructor(data?: PartialMessage<PushToGithubResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.PushToGithubResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "github_url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "account", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "repo", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PushToGithubResponse {
    return new PushToGithubResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PushToGithubResponse {
    return new PushToGithubResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PushToGithubResponse {
    return new PushToGithubResponse().fromJsonString(jsonString, options);
  }

  static equals(a: PushToGithubResponse | PlainMessage<PushToGithubResponse> | undefined, b: PushToGithubResponse | PlainMessage<PushToGithubResponse> | undefined): boolean {
    return proto3.util.equals(PushToGithubResponse, a, b);
  }
}

/**
 * @generated from message rill.local.v1.DeployProjectRequest
 */
export class DeployProjectRequest extends Message<DeployProjectRequest> {
  /**
   * @generated from field: string org = 1;
   */
  org = "";

  /**
   * @generated from field: string project_name = 2;
   */
  projectName = "";

  /**
   * @generated from field: bool upload = 3;
   */
  upload = false;

  constructor(data?: PartialMessage<DeployProjectRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.DeployProjectRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "org", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "project_name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "upload", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeployProjectRequest {
    return new DeployProjectRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeployProjectRequest {
    return new DeployProjectRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeployProjectRequest {
    return new DeployProjectRequest().fromJsonString(jsonString, options);
  }

  static equals(a: DeployProjectRequest | PlainMessage<DeployProjectRequest> | undefined, b: DeployProjectRequest | PlainMessage<DeployProjectRequest> | undefined): boolean {
    return proto3.util.equals(DeployProjectRequest, a, b);
  }
}

/**
 * @generated from message rill.local.v1.DeployProjectResponse
 */
export class DeployProjectResponse extends Message<DeployProjectResponse> {
  /**
   * @generated from field: string deploy_id = 1;
   */
  deployId = "";

  /**
   * @generated from field: string org = 2;
   */
  org = "";

  /**
   * @generated from field: string project = 3;
   */
  project = "";

  /**
   * @generated from field: string frontend_url = 4;
   */
  frontendUrl = "";

  constructor(data?: PartialMessage<DeployProjectResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.DeployProjectResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "deploy_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "org", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "project", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "frontend_url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeployProjectResponse {
    return new DeployProjectResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeployProjectResponse {
    return new DeployProjectResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeployProjectResponse {
    return new DeployProjectResponse().fromJsonString(jsonString, options);
  }

  static equals(a: DeployProjectResponse | PlainMessage<DeployProjectResponse> | undefined, b: DeployProjectResponse | PlainMessage<DeployProjectResponse> | undefined): boolean {
    return proto3.util.equals(DeployProjectResponse, a, b);
  }
}

/**
 * @generated from message rill.local.v1.RedeployProjectRequest
 */
export class RedeployProjectRequest extends Message<RedeployProjectRequest> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId = "";

  /**
   * @generated from field: bool reupload = 2;
   */
  reupload = false;

  constructor(data?: PartialMessage<RedeployProjectRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.RedeployProjectRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "project_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "reupload", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RedeployProjectRequest {
    return new RedeployProjectRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RedeployProjectRequest {
    return new RedeployProjectRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RedeployProjectRequest {
    return new RedeployProjectRequest().fromJsonString(jsonString, options);
  }

  static equals(a: RedeployProjectRequest | PlainMessage<RedeployProjectRequest> | undefined, b: RedeployProjectRequest | PlainMessage<RedeployProjectRequest> | undefined): boolean {
    return proto3.util.equals(RedeployProjectRequest, a, b);
  }
}

/**
 * @generated from message rill.local.v1.RedeployProjectResponse
 */
export class RedeployProjectResponse extends Message<RedeployProjectResponse> {
  constructor(data?: PartialMessage<RedeployProjectResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.RedeployProjectResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RedeployProjectResponse {
    return new RedeployProjectResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RedeployProjectResponse {
    return new RedeployProjectResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RedeployProjectResponse {
    return new RedeployProjectResponse().fromJsonString(jsonString, options);
  }

  static equals(a: RedeployProjectResponse | PlainMessage<RedeployProjectResponse> | undefined, b: RedeployProjectResponse | PlainMessage<RedeployProjectResponse> | undefined): boolean {
    return proto3.util.equals(RedeployProjectResponse, a, b);
  }
}

/**
 * @generated from message rill.local.v1.GetCurrentUserRequest
 */
export class GetCurrentUserRequest extends Message<GetCurrentUserRequest> {
  constructor(data?: PartialMessage<GetCurrentUserRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.GetCurrentUserRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetCurrentUserRequest {
    return new GetCurrentUserRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetCurrentUserRequest {
    return new GetCurrentUserRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetCurrentUserRequest {
    return new GetCurrentUserRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetCurrentUserRequest | PlainMessage<GetCurrentUserRequest> | undefined, b: GetCurrentUserRequest | PlainMessage<GetCurrentUserRequest> | undefined): boolean {
    return proto3.util.equals(GetCurrentUserRequest, a, b);
  }
}

/**
 * @generated from message rill.local.v1.GetCurrentUserResponse
 */
export class GetCurrentUserResponse extends Message<GetCurrentUserResponse> {
  /**
   * @generated from field: rill.admin.v1.User user = 1;
   */
  user?: User;

  constructor(data?: PartialMessage<GetCurrentUserResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.local.v1.GetCurrentUserResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "user", kind: "message", T: User },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetCurrentUserResponse {
    return new GetCurrentUserResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetCurrentUserResponse {
    return new GetCurrentUserResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetCurrentUserResponse {
    return new GetCurrentUserResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetCurrentUserResponse | PlainMessage<GetCurrentUserResponse> | undefined, b: GetCurrentUserResponse | PlainMessage<GetCurrentUserResponse> | undefined): boolean {
    return proto3.util.equals(GetCurrentUserResponse, a, b);
  }
}


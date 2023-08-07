// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: rill/runtime/v1/connectors.proto

package runtimev1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ConnectorService_S3ListBuckets_FullMethodName         = "/rill.runtime.v1.ConnectorService/S3ListBuckets"
	ConnectorService_S3ListObjects_FullMethodName         = "/rill.runtime.v1.ConnectorService/S3ListObjects"
	ConnectorService_S3GetBucketMetadata_FullMethodName   = "/rill.runtime.v1.ConnectorService/S3GetBucketMetadata"
	ConnectorService_S3GetCredentialsInfo_FullMethodName  = "/rill.runtime.v1.ConnectorService/S3GetCredentialsInfo"
	ConnectorService_GCSListBuckets_FullMethodName        = "/rill.runtime.v1.ConnectorService/GCSListBuckets"
	ConnectorService_GCSListObjects_FullMethodName        = "/rill.runtime.v1.ConnectorService/GCSListObjects"
	ConnectorService_GCSGetCredentialsInfo_FullMethodName = "/rill.runtime.v1.ConnectorService/GCSGetCredentialsInfo"
	ConnectorService_OLAPListTables_FullMethodName        = "/rill.runtime.v1.ConnectorService/OLAPListTables"
	ConnectorService_BigQueryListDatasets_FullMethodName  = "/rill.runtime.v1.ConnectorService/BigQueryListDatasets"
	ConnectorService_BigQueryListTables_FullMethodName    = "/rill.runtime.v1.ConnectorService/BigQueryListTables"
)

// ConnectorServiceClient is the client API for ConnectorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConnectorServiceClient interface {
	// S3ListBuckets lists buckets accessible with the configured credentials.
	S3ListBuckets(ctx context.Context, in *S3ListBucketsRequest, opts ...grpc.CallOption) (*S3ListBucketsResponse, error)
	// S3ListBuckets lists objects for the given bucket.
	S3ListObjects(ctx context.Context, in *S3ListObjectsRequest, opts ...grpc.CallOption) (*S3ListObjectsResponse, error)
	// S3GetBucketMetadata returns metadata for the given bucket.
	S3GetBucketMetadata(ctx context.Context, in *S3GetBucketMetadataRequest, opts ...grpc.CallOption) (*S3GetBucketMetadataResponse, error)
	// S3GetCredentialsInfo returns metadata for the given bucket.
	S3GetCredentialsInfo(ctx context.Context, in *S3GetCredentialsInfoRequest, opts ...grpc.CallOption) (*S3GetCredentialsInfoResponse, error)
	// GCSListBuckets lists buckets accessible with the configured credentials.
	GCSListBuckets(ctx context.Context, in *GCSListBucketsRequest, opts ...grpc.CallOption) (*GCSListBucketsResponse, error)
	// GCSListObjects lists objects for the given bucket.
	GCSListObjects(ctx context.Context, in *GCSListObjectsRequest, opts ...grpc.CallOption) (*GCSListObjectsResponse, error)
	// GCSGetCredentialsInfo returns metadata for the given bucket.
	GCSGetCredentialsInfo(ctx context.Context, in *GCSGetCredentialsInfoRequest, opts ...grpc.CallOption) (*GCSGetCredentialsInfoResponse, error)
	// OLAPListTables list all tables across all databases on motherduck
	OLAPListTables(ctx context.Context, in *OLAPListTablesRequest, opts ...grpc.CallOption) (*OLAPListTablesResponse, error)
	// BigQueryListDatasets list all datasets in a bigquery project
	BigQueryListDatasets(ctx context.Context, in *BigQueryListDatasetsRequest, opts ...grpc.CallOption) (*BigQueryListDatasetsResponse, error)
	// BigQueryListTables list all tables in a bigquery project:dataset
	BigQueryListTables(ctx context.Context, in *BigQueryListTablesRequest, opts ...grpc.CallOption) (*BigQueryListTablesResponse, error)
}

type connectorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConnectorServiceClient(cc grpc.ClientConnInterface) ConnectorServiceClient {
	return &connectorServiceClient{cc}
}

func (c *connectorServiceClient) S3ListBuckets(ctx context.Context, in *S3ListBucketsRequest, opts ...grpc.CallOption) (*S3ListBucketsResponse, error) {
	out := new(S3ListBucketsResponse)
	err := c.cc.Invoke(ctx, ConnectorService_S3ListBuckets_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectorServiceClient) S3ListObjects(ctx context.Context, in *S3ListObjectsRequest, opts ...grpc.CallOption) (*S3ListObjectsResponse, error) {
	out := new(S3ListObjectsResponse)
	err := c.cc.Invoke(ctx, ConnectorService_S3ListObjects_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectorServiceClient) S3GetBucketMetadata(ctx context.Context, in *S3GetBucketMetadataRequest, opts ...grpc.CallOption) (*S3GetBucketMetadataResponse, error) {
	out := new(S3GetBucketMetadataResponse)
	err := c.cc.Invoke(ctx, ConnectorService_S3GetBucketMetadata_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectorServiceClient) S3GetCredentialsInfo(ctx context.Context, in *S3GetCredentialsInfoRequest, opts ...grpc.CallOption) (*S3GetCredentialsInfoResponse, error) {
	out := new(S3GetCredentialsInfoResponse)
	err := c.cc.Invoke(ctx, ConnectorService_S3GetCredentialsInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectorServiceClient) GCSListBuckets(ctx context.Context, in *GCSListBucketsRequest, opts ...grpc.CallOption) (*GCSListBucketsResponse, error) {
	out := new(GCSListBucketsResponse)
	err := c.cc.Invoke(ctx, ConnectorService_GCSListBuckets_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectorServiceClient) GCSListObjects(ctx context.Context, in *GCSListObjectsRequest, opts ...grpc.CallOption) (*GCSListObjectsResponse, error) {
	out := new(GCSListObjectsResponse)
	err := c.cc.Invoke(ctx, ConnectorService_GCSListObjects_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectorServiceClient) GCSGetCredentialsInfo(ctx context.Context, in *GCSGetCredentialsInfoRequest, opts ...grpc.CallOption) (*GCSGetCredentialsInfoResponse, error) {
	out := new(GCSGetCredentialsInfoResponse)
	err := c.cc.Invoke(ctx, ConnectorService_GCSGetCredentialsInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectorServiceClient) OLAPListTables(ctx context.Context, in *OLAPListTablesRequest, opts ...grpc.CallOption) (*OLAPListTablesResponse, error) {
	out := new(OLAPListTablesResponse)
	err := c.cc.Invoke(ctx, ConnectorService_OLAPListTables_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectorServiceClient) BigQueryListDatasets(ctx context.Context, in *BigQueryListDatasetsRequest, opts ...grpc.CallOption) (*BigQueryListDatasetsResponse, error) {
	out := new(BigQueryListDatasetsResponse)
	err := c.cc.Invoke(ctx, ConnectorService_BigQueryListDatasets_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectorServiceClient) BigQueryListTables(ctx context.Context, in *BigQueryListTablesRequest, opts ...grpc.CallOption) (*BigQueryListTablesResponse, error) {
	out := new(BigQueryListTablesResponse)
	err := c.cc.Invoke(ctx, ConnectorService_BigQueryListTables_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConnectorServiceServer is the server API for ConnectorService service.
// All implementations must embed UnimplementedConnectorServiceServer
// for forward compatibility
type ConnectorServiceServer interface {
	// S3ListBuckets lists buckets accessible with the configured credentials.
	S3ListBuckets(context.Context, *S3ListBucketsRequest) (*S3ListBucketsResponse, error)
	// S3ListBuckets lists objects for the given bucket.
	S3ListObjects(context.Context, *S3ListObjectsRequest) (*S3ListObjectsResponse, error)
	// S3GetBucketMetadata returns metadata for the given bucket.
	S3GetBucketMetadata(context.Context, *S3GetBucketMetadataRequest) (*S3GetBucketMetadataResponse, error)
	// S3GetCredentialsInfo returns metadata for the given bucket.
	S3GetCredentialsInfo(context.Context, *S3GetCredentialsInfoRequest) (*S3GetCredentialsInfoResponse, error)
	// GCSListBuckets lists buckets accessible with the configured credentials.
	GCSListBuckets(context.Context, *GCSListBucketsRequest) (*GCSListBucketsResponse, error)
	// GCSListObjects lists objects for the given bucket.
	GCSListObjects(context.Context, *GCSListObjectsRequest) (*GCSListObjectsResponse, error)
	// GCSGetCredentialsInfo returns metadata for the given bucket.
	GCSGetCredentialsInfo(context.Context, *GCSGetCredentialsInfoRequest) (*GCSGetCredentialsInfoResponse, error)
	// OLAPListTables list all tables across all databases on motherduck
	OLAPListTables(context.Context, *OLAPListTablesRequest) (*OLAPListTablesResponse, error)
	// BigQueryListDatasets list all datasets in a bigquery project
	BigQueryListDatasets(context.Context, *BigQueryListDatasetsRequest) (*BigQueryListDatasetsResponse, error)
	// BigQueryListTables list all tables in a bigquery project:dataset
	BigQueryListTables(context.Context, *BigQueryListTablesRequest) (*BigQueryListTablesResponse, error)
	mustEmbedUnimplementedConnectorServiceServer()
}

// UnimplementedConnectorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConnectorServiceServer struct {
}

func (UnimplementedConnectorServiceServer) S3ListBuckets(context.Context, *S3ListBucketsRequest) (*S3ListBucketsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method S3ListBuckets not implemented")
}
func (UnimplementedConnectorServiceServer) S3ListObjects(context.Context, *S3ListObjectsRequest) (*S3ListObjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method S3ListObjects not implemented")
}
func (UnimplementedConnectorServiceServer) S3GetBucketMetadata(context.Context, *S3GetBucketMetadataRequest) (*S3GetBucketMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method S3GetBucketMetadata not implemented")
}
func (UnimplementedConnectorServiceServer) S3GetCredentialsInfo(context.Context, *S3GetCredentialsInfoRequest) (*S3GetCredentialsInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method S3GetCredentialsInfo not implemented")
}
func (UnimplementedConnectorServiceServer) GCSListBuckets(context.Context, *GCSListBucketsRequest) (*GCSListBucketsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GCSListBuckets not implemented")
}
func (UnimplementedConnectorServiceServer) GCSListObjects(context.Context, *GCSListObjectsRequest) (*GCSListObjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GCSListObjects not implemented")
}
func (UnimplementedConnectorServiceServer) GCSGetCredentialsInfo(context.Context, *GCSGetCredentialsInfoRequest) (*GCSGetCredentialsInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GCSGetCredentialsInfo not implemented")
}
func (UnimplementedConnectorServiceServer) OLAPListTables(context.Context, *OLAPListTablesRequest) (*OLAPListTablesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OLAPListTables not implemented")
}
func (UnimplementedConnectorServiceServer) BigQueryListDatasets(context.Context, *BigQueryListDatasetsRequest) (*BigQueryListDatasetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BigQueryListDatasets not implemented")
}
func (UnimplementedConnectorServiceServer) BigQueryListTables(context.Context, *BigQueryListTablesRequest) (*BigQueryListTablesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BigQueryListTables not implemented")
}
func (UnimplementedConnectorServiceServer) mustEmbedUnimplementedConnectorServiceServer() {}

// UnsafeConnectorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConnectorServiceServer will
// result in compilation errors.
type UnsafeConnectorServiceServer interface {
	mustEmbedUnimplementedConnectorServiceServer()
}

func RegisterConnectorServiceServer(s grpc.ServiceRegistrar, srv ConnectorServiceServer) {
	s.RegisterService(&ConnectorService_ServiceDesc, srv)
}

func _ConnectorService_S3ListBuckets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(S3ListBucketsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectorServiceServer).S3ListBuckets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConnectorService_S3ListBuckets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectorServiceServer).S3ListBuckets(ctx, req.(*S3ListBucketsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectorService_S3ListObjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(S3ListObjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectorServiceServer).S3ListObjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConnectorService_S3ListObjects_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectorServiceServer).S3ListObjects(ctx, req.(*S3ListObjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectorService_S3GetBucketMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(S3GetBucketMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectorServiceServer).S3GetBucketMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConnectorService_S3GetBucketMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectorServiceServer).S3GetBucketMetadata(ctx, req.(*S3GetBucketMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectorService_S3GetCredentialsInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(S3GetCredentialsInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectorServiceServer).S3GetCredentialsInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConnectorService_S3GetCredentialsInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectorServiceServer).S3GetCredentialsInfo(ctx, req.(*S3GetCredentialsInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectorService_GCSListBuckets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GCSListBucketsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectorServiceServer).GCSListBuckets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConnectorService_GCSListBuckets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectorServiceServer).GCSListBuckets(ctx, req.(*GCSListBucketsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectorService_GCSListObjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GCSListObjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectorServiceServer).GCSListObjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConnectorService_GCSListObjects_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectorServiceServer).GCSListObjects(ctx, req.(*GCSListObjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectorService_GCSGetCredentialsInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GCSGetCredentialsInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectorServiceServer).GCSGetCredentialsInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConnectorService_GCSGetCredentialsInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectorServiceServer).GCSGetCredentialsInfo(ctx, req.(*GCSGetCredentialsInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectorService_OLAPListTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OLAPListTablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectorServiceServer).OLAPListTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConnectorService_OLAPListTables_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectorServiceServer).OLAPListTables(ctx, req.(*OLAPListTablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectorService_BigQueryListDatasets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BigQueryListDatasetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectorServiceServer).BigQueryListDatasets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConnectorService_BigQueryListDatasets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectorServiceServer).BigQueryListDatasets(ctx, req.(*BigQueryListDatasetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectorService_BigQueryListTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BigQueryListTablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectorServiceServer).BigQueryListTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConnectorService_BigQueryListTables_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectorServiceServer).BigQueryListTables(ctx, req.(*BigQueryListTablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConnectorService_ServiceDesc is the grpc.ServiceDesc for ConnectorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConnectorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rill.runtime.v1.ConnectorService",
	HandlerType: (*ConnectorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "S3ListBuckets",
			Handler:    _ConnectorService_S3ListBuckets_Handler,
		},
		{
			MethodName: "S3ListObjects",
			Handler:    _ConnectorService_S3ListObjects_Handler,
		},
		{
			MethodName: "S3GetBucketMetadata",
			Handler:    _ConnectorService_S3GetBucketMetadata_Handler,
		},
		{
			MethodName: "S3GetCredentialsInfo",
			Handler:    _ConnectorService_S3GetCredentialsInfo_Handler,
		},
		{
			MethodName: "GCSListBuckets",
			Handler:    _ConnectorService_GCSListBuckets_Handler,
		},
		{
			MethodName: "GCSListObjects",
			Handler:    _ConnectorService_GCSListObjects_Handler,
		},
		{
			MethodName: "GCSGetCredentialsInfo",
			Handler:    _ConnectorService_GCSGetCredentialsInfo_Handler,
		},
		{
			MethodName: "OLAPListTables",
			Handler:    _ConnectorService_OLAPListTables_Handler,
		},
		{
			MethodName: "BigQueryListDatasets",
			Handler:    _ConnectorService_BigQueryListDatasets_Handler,
		},
		{
			MethodName: "BigQueryListTables",
			Handler:    _ConnectorService_BigQueryListTables_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rill/runtime/v1/connectors.proto",
}

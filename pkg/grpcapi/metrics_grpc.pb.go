// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: proto/metrics.proto

package grpcapi

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
	Metrics_Get_FullMethodName         = "/go.metric.Metrics/Get"
	Metrics_Update_FullMethodName      = "/go.metric.Metrics/Update"
	Metrics_BatchUpdate_FullMethodName = "/go.metric.Metrics/BatchUpdate"
)

// MetricsClient is the client API for Metrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsClient interface {
	Get(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error)
	Update(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*UpdateMetricResponse, error)
	BatchUpdate(ctx context.Context, in *BatchUpdateMetricRequest, opts ...grpc.CallOption) (*BatchUpdateMetricResponse, error)
}

type metricsClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsClient(cc grpc.ClientConnInterface) MetricsClient {
	return &metricsClient{cc}
}

func (c *metricsClient) Get(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error) {
	out := new(GetMetricResponse)
	err := c.cc.Invoke(ctx, Metrics_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) Update(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*UpdateMetricResponse, error) {
	out := new(UpdateMetricResponse)
	err := c.cc.Invoke(ctx, Metrics_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) BatchUpdate(ctx context.Context, in *BatchUpdateMetricRequest, opts ...grpc.CallOption) (*BatchUpdateMetricResponse, error) {
	out := new(BatchUpdateMetricResponse)
	err := c.cc.Invoke(ctx, Metrics_BatchUpdate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServer is the server API for Metrics service.
// All implementations must embed UnimplementedMetricsServer
// for forward compatibility
type MetricsServer interface {
	Get(context.Context, *GetMetricRequest) (*GetMetricResponse, error)
	Update(context.Context, *UpdateMetricRequest) (*UpdateMetricResponse, error)
	BatchUpdate(context.Context, *BatchUpdateMetricRequest) (*BatchUpdateMetricResponse, error)
	mustEmbedUnimplementedMetricsServer()
}

// UnimplementedMetricsServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsServer struct {
}

func (UnimplementedMetricsServer) Get(context.Context, *GetMetricRequest) (*GetMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedMetricsServer) Update(context.Context, *UpdateMetricRequest) (*UpdateMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedMetricsServer) BatchUpdate(context.Context, *BatchUpdateMetricRequest) (*BatchUpdateMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchUpdate not implemented")
}
func (UnimplementedMetricsServer) mustEmbedUnimplementedMetricsServer() {}

// UnsafeMetricsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServer will
// result in compilation errors.
type UnsafeMetricsServer interface {
	mustEmbedUnimplementedMetricsServer()
}

func RegisterMetricsServer(s grpc.ServiceRegistrar, srv MetricsServer) {
	s.RegisterService(&Metrics_ServiceDesc, srv)
}

func _Metrics_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).Get(ctx, req.(*GetMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).Update(ctx, req.(*UpdateMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_BatchUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchUpdateMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).BatchUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metrics_BatchUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).BatchUpdate(ctx, req.(*BatchUpdateMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Metrics_ServiceDesc is the grpc.ServiceDesc for Metrics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metrics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "go.metric.Metrics",
	HandlerType: (*MetricsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Metrics_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Metrics_Update_Handler,
		},
		{
			MethodName: "BatchUpdate",
			Handler:    _Metrics_BatchUpdate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/metrics.proto",
}

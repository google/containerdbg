// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: node_api.proto

package proto

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

// NodeDaemonServiceClient is the client API for NodeDaemonService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeDaemonServiceClient interface {
	Monitor(ctx context.Context, in *MonitorPodRequest, opts ...grpc.CallOption) (*MonitorPodResponse, error)
	ReportDnsQuery(ctx context.Context, in *ReportDnsQueryResultRequest, opts ...grpc.CallOption) (*ReportDnsQueryResultResponse, error)
	ReportDnsSearchValues(ctx context.Context, in *ReportDnsSearchValuesRequest, opts ...grpc.CallOption) (*ReportDnsSearchValuesResponse, error)
}

type nodeDaemonServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeDaemonServiceClient(cc grpc.ClientConnInterface) NodeDaemonServiceClient {
	return &nodeDaemonServiceClient{cc}
}

func (c *nodeDaemonServiceClient) Monitor(ctx context.Context, in *MonitorPodRequest, opts ...grpc.CallOption) (*MonitorPodResponse, error) {
	out := new(MonitorPodResponse)
	err := c.cc.Invoke(ctx, "/NodeDaemonService/Monitor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeDaemonServiceClient) ReportDnsQuery(ctx context.Context, in *ReportDnsQueryResultRequest, opts ...grpc.CallOption) (*ReportDnsQueryResultResponse, error) {
	out := new(ReportDnsQueryResultResponse)
	err := c.cc.Invoke(ctx, "/NodeDaemonService/ReportDnsQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeDaemonServiceClient) ReportDnsSearchValues(ctx context.Context, in *ReportDnsSearchValuesRequest, opts ...grpc.CallOption) (*ReportDnsSearchValuesResponse, error) {
	out := new(ReportDnsSearchValuesResponse)
	err := c.cc.Invoke(ctx, "/NodeDaemonService/ReportDnsSearchValues", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeDaemonServiceServer is the server API for NodeDaemonService service.
// All implementations must embed UnimplementedNodeDaemonServiceServer
// for forward compatibility
type NodeDaemonServiceServer interface {
	Monitor(context.Context, *MonitorPodRequest) (*MonitorPodResponse, error)
	ReportDnsQuery(context.Context, *ReportDnsQueryResultRequest) (*ReportDnsQueryResultResponse, error)
	ReportDnsSearchValues(context.Context, *ReportDnsSearchValuesRequest) (*ReportDnsSearchValuesResponse, error)
	mustEmbedUnimplementedNodeDaemonServiceServer()
}

// UnimplementedNodeDaemonServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNodeDaemonServiceServer struct {
}

func (UnimplementedNodeDaemonServiceServer) Monitor(context.Context, *MonitorPodRequest) (*MonitorPodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Monitor not implemented")
}
func (UnimplementedNodeDaemonServiceServer) ReportDnsQuery(context.Context, *ReportDnsQueryResultRequest) (*ReportDnsQueryResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportDnsQuery not implemented")
}
func (UnimplementedNodeDaemonServiceServer) ReportDnsSearchValues(context.Context, *ReportDnsSearchValuesRequest) (*ReportDnsSearchValuesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportDnsSearchValues not implemented")
}
func (UnimplementedNodeDaemonServiceServer) mustEmbedUnimplementedNodeDaemonServiceServer() {}

// UnsafeNodeDaemonServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeDaemonServiceServer will
// result in compilation errors.
type UnsafeNodeDaemonServiceServer interface {
	mustEmbedUnimplementedNodeDaemonServiceServer()
}

func RegisterNodeDaemonServiceServer(s grpc.ServiceRegistrar, srv NodeDaemonServiceServer) {
	s.RegisterService(&NodeDaemonService_ServiceDesc, srv)
}

func _NodeDaemonService_Monitor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MonitorPodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeDaemonServiceServer).Monitor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NodeDaemonService/Monitor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeDaemonServiceServer).Monitor(ctx, req.(*MonitorPodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeDaemonService_ReportDnsQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportDnsQueryResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeDaemonServiceServer).ReportDnsQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NodeDaemonService/ReportDnsQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeDaemonServiceServer).ReportDnsQuery(ctx, req.(*ReportDnsQueryResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeDaemonService_ReportDnsSearchValues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportDnsSearchValuesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeDaemonServiceServer).ReportDnsSearchValues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NodeDaemonService/ReportDnsSearchValues",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeDaemonServiceServer).ReportDnsSearchValues(ctx, req.(*ReportDnsSearchValuesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NodeDaemonService_ServiceDesc is the grpc.ServiceDesc for NodeDaemonService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NodeDaemonService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "NodeDaemonService",
	HandlerType: (*NodeDaemonServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Monitor",
			Handler:    _NodeDaemonService_Monitor_Handler,
		},
		{
			MethodName: "ReportDnsQuery",
			Handler:    _NodeDaemonService_ReportDnsQuery_Handler,
		},
		{
			MethodName: "ReportDnsSearchValues",
			Handler:    _NodeDaemonService_ReportDnsSearchValues_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node_api.proto",
}

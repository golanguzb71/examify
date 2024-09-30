// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: proto/integration.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	IntegrationService_GetResultWritingTask_FullMethodName  = "/integration.IntegrationService/GetResultWritingTask"
	IntegrationService_GetResultSpeakingPart_FullMethodName = "/integration.IntegrationService/GetResultSpeakingPart"
)

// IntegrationServiceClient is the client API for IntegrationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IntegrationServiceClient interface {
	GetResultWritingTask(ctx context.Context, in *WritingTaskAbsRequest, opts ...grpc.CallOption) (*WritingTaskAbsResponse, error)
	GetResultSpeakingPart(ctx context.Context, in *SpeakingPartAbsRequest, opts ...grpc.CallOption) (*SpeakingPartAbsResponse, error)
}

type integrationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIntegrationServiceClient(cc grpc.ClientConnInterface) IntegrationServiceClient {
	return &integrationServiceClient{cc}
}

func (c *integrationServiceClient) GetResultWritingTask(ctx context.Context, in *WritingTaskAbsRequest, opts ...grpc.CallOption) (*WritingTaskAbsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WritingTaskAbsResponse)
	err := c.cc.Invoke(ctx, IntegrationService_GetResultWritingTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *integrationServiceClient) GetResultSpeakingPart(ctx context.Context, in *SpeakingPartAbsRequest, opts ...grpc.CallOption) (*SpeakingPartAbsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SpeakingPartAbsResponse)
	err := c.cc.Invoke(ctx, IntegrationService_GetResultSpeakingPart_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IntegrationServiceServer is the server API for IntegrationService service.
// All implementations must embed UnimplementedIntegrationServiceServer
// for forward compatibility.
type IntegrationServiceServer interface {
	GetResultWritingTask(context.Context, *WritingTaskAbsRequest) (*WritingTaskAbsResponse, error)
	GetResultSpeakingPart(context.Context, *SpeakingPartAbsRequest) (*SpeakingPartAbsResponse, error)
	mustEmbedUnimplementedIntegrationServiceServer()
}

// UnimplementedIntegrationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedIntegrationServiceServer struct{}

func (UnimplementedIntegrationServiceServer) GetResultWritingTask(context.Context, *WritingTaskAbsRequest) (*WritingTaskAbsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResultWritingTask not implemented")
}
func (UnimplementedIntegrationServiceServer) GetResultSpeakingPart(context.Context, *SpeakingPartAbsRequest) (*SpeakingPartAbsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResultSpeakingPart not implemented")
}
func (UnimplementedIntegrationServiceServer) mustEmbedUnimplementedIntegrationServiceServer() {}
func (UnimplementedIntegrationServiceServer) testEmbeddedByValue()                            {}

// UnsafeIntegrationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IntegrationServiceServer will
// result in compilation errors.
type UnsafeIntegrationServiceServer interface {
	mustEmbedUnimplementedIntegrationServiceServer()
}

func RegisterIntegrationServiceServer(s grpc.ServiceRegistrar, srv IntegrationServiceServer) {
	// If the following call pancis, it indicates UnimplementedIntegrationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&IntegrationService_ServiceDesc, srv)
}

func _IntegrationService_GetResultWritingTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WritingTaskAbsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IntegrationServiceServer).GetResultWritingTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IntegrationService_GetResultWritingTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IntegrationServiceServer).GetResultWritingTask(ctx, req.(*WritingTaskAbsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IntegrationService_GetResultSpeakingPart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SpeakingPartAbsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IntegrationServiceServer).GetResultSpeakingPart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IntegrationService_GetResultSpeakingPart_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IntegrationServiceServer).GetResultSpeakingPart(ctx, req.(*SpeakingPartAbsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IntegrationService_ServiceDesc is the grpc.ServiceDesc for IntegrationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IntegrationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "integration.IntegrationService",
	HandlerType: (*IntegrationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetResultWritingTask",
			Handler:    _IntegrationService_GetResultWritingTask_Handler,
		},
		{
			MethodName: "GetResultSpeakingPart",
			Handler:    _IntegrationService_GetResultSpeakingPart_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/integration.proto",
}
// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: auth.proto

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
	BonusService_UseBonusAttempt_FullMethodName = "/auth.BonusService/UseBonusAttempt"
)

// BonusServiceClient is the client API for BonusService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BonusServiceClient interface {
	UseBonusAttempt(ctx context.Context, in *UseBonusAttemptRequest, opts ...grpc.CallOption) (*UseBonusAttemptResponse, error)
}

type bonusServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBonusServiceClient(cc grpc.ClientConnInterface) BonusServiceClient {
	return &bonusServiceClient{cc}
}

func (c *bonusServiceClient) UseBonusAttempt(ctx context.Context, in *UseBonusAttemptRequest, opts ...grpc.CallOption) (*UseBonusAttemptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UseBonusAttemptResponse)
	err := c.cc.Invoke(ctx, BonusService_UseBonusAttempt_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BonusServiceServer is the server API for BonusService service.
// All implementations must embed UnimplementedBonusServiceServer
// for forward compatibility.
type BonusServiceServer interface {
	UseBonusAttempt(context.Context, *UseBonusAttemptRequest) (*UseBonusAttemptResponse, error)
	mustEmbedUnimplementedBonusServiceServer()
}

// UnimplementedBonusServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBonusServiceServer struct{}

func (UnimplementedBonusServiceServer) UseBonusAttempt(context.Context, *UseBonusAttemptRequest) (*UseBonusAttemptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseBonusAttempt not implemented")
}
func (UnimplementedBonusServiceServer) mustEmbedUnimplementedBonusServiceServer() {}
func (UnimplementedBonusServiceServer) testEmbeddedByValue()                      {}

// UnsafeBonusServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BonusServiceServer will
// result in compilation errors.
type UnsafeBonusServiceServer interface {
	mustEmbedUnimplementedBonusServiceServer()
}

func RegisterBonusServiceServer(s grpc.ServiceRegistrar, srv BonusServiceServer) {
	// If the following call pancis, it indicates UnimplementedBonusServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BonusService_ServiceDesc, srv)
}

func _BonusService_UseBonusAttempt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UseBonusAttemptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BonusServiceServer).UseBonusAttempt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BonusService_UseBonusAttempt_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BonusServiceServer).UseBonusAttempt(ctx, req.(*UseBonusAttemptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BonusService_ServiceDesc is the grpc.ServiceDesc for BonusService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BonusService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.BonusService",
	HandlerType: (*BonusServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UseBonusAttempt",
			Handler:    _BonusService_UseBonusAttempt_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}

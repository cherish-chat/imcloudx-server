// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: gateway.proto

package pb

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

// GatewayServiceClient is the client API for GatewayService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GatewayServiceClient interface {
	GetSessionDescription(ctx context.Context, in *GetSessionDescriptionReq, opts ...grpc.CallOption) (*GetSessionDescriptionResp, error)
	Node(ctx context.Context, in *NodeReq, opts ...grpc.CallOption) (*NodeResp, error)
}

type gatewayServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGatewayServiceClient(cc grpc.ClientConnInterface) GatewayServiceClient {
	return &gatewayServiceClient{cc}
}

func (c *gatewayServiceClient) GetSessionDescription(ctx context.Context, in *GetSessionDescriptionReq, opts ...grpc.CallOption) (*GetSessionDescriptionResp, error) {
	out := new(GetSessionDescriptionResp)
	err := c.cc.Invoke(ctx, "/pb.gatewayService/GetSessionDescription", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayServiceClient) Node(ctx context.Context, in *NodeReq, opts ...grpc.CallOption) (*NodeResp, error) {
	out := new(NodeResp)
	err := c.cc.Invoke(ctx, "/pb.gatewayService/Node", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GatewayServiceServer is the server API for GatewayService service.
// All implementations must embed UnimplementedGatewayServiceServer
// for forward compatibility
type GatewayServiceServer interface {
	GetSessionDescription(context.Context, *GetSessionDescriptionReq) (*GetSessionDescriptionResp, error)
	Node(context.Context, *NodeReq) (*NodeResp, error)
	mustEmbedUnimplementedGatewayServiceServer()
}

// UnimplementedGatewayServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGatewayServiceServer struct {
}

func (UnimplementedGatewayServiceServer) GetSessionDescription(context.Context, *GetSessionDescriptionReq) (*GetSessionDescriptionResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSessionDescription not implemented")
}
func (UnimplementedGatewayServiceServer) Node(context.Context, *NodeReq) (*NodeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Node not implemented")
}
func (UnimplementedGatewayServiceServer) mustEmbedUnimplementedGatewayServiceServer() {}

// UnsafeGatewayServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GatewayServiceServer will
// result in compilation errors.
type UnsafeGatewayServiceServer interface {
	mustEmbedUnimplementedGatewayServiceServer()
}

func RegisterGatewayServiceServer(s grpc.ServiceRegistrar, srv GatewayServiceServer) {
	s.RegisterService(&GatewayService_ServiceDesc, srv)
}

func _GatewayService_GetSessionDescription_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSessionDescriptionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServiceServer).GetSessionDescription(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.gatewayService/GetSessionDescription",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServiceServer).GetSessionDescription(ctx, req.(*GetSessionDescriptionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GatewayService_Node_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServiceServer).Node(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.gatewayService/Node",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServiceServer).Node(ctx, req.(*NodeReq))
	}
	return interceptor(ctx, in, info, handler)
}

// GatewayService_ServiceDesc is the grpc.ServiceDesc for GatewayService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GatewayService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.gatewayService",
	HandlerType: (*GatewayServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSessionDescription",
			Handler:    _GatewayService_GetSessionDescription_Handler,
		},
		{
			MethodName: "Node",
			Handler:    _GatewayService_Node_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gateway.proto",
}

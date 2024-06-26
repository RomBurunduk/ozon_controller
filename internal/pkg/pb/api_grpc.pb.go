// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: api.proto

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

const (
	PVZService_CreatePVZ_FullMethodName     = "/pvz.PVZService/CreatePVZ"
	PVZService_GetPVZ_FullMethodName        = "/pvz.PVZService/GetPVZ"
	PVZService_UpdatePVZ_FullMethodName     = "/pvz.PVZService/UpdatePVZ"
	PVZService_DeletePVZ_FullMethodName     = "/pvz.PVZService/DeletePVZ"
	PVZService_ListAllPVZ_FullMethodName    = "/pvz.PVZService/ListAllPVZ"
	PVZService_DeleteListPVZ_FullMethodName = "/pvz.PVZService/DeleteListPVZ"
)

// PVZServiceClient is the client API for PVZService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PVZServiceClient interface {
	CreatePVZ(ctx context.Context, in *CreatePVZRequest, opts ...grpc.CallOption) (*CreatePVZResponse, error)
	GetPVZ(ctx context.Context, in *GetPVZRequest, opts ...grpc.CallOption) (*GetPVZResponse, error)
	UpdatePVZ(ctx context.Context, in *UpdatePVZRequest, opts ...grpc.CallOption) (*UpdatePVZResponse, error)
	DeletePVZ(ctx context.Context, in *DeletePVZRequest, opts ...grpc.CallOption) (*DeletePVZResponse, error)
	ListAllPVZ(ctx context.Context, in *ListAllPVZRequest, opts ...grpc.CallOption) (*ListAllPVZResponse, error)
	DeleteListPVZ(ctx context.Context, in *DeleteListPVZRequest, opts ...grpc.CallOption) (*DeleteListPVZResponse, error)
}

type pVZServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPVZServiceClient(cc grpc.ClientConnInterface) PVZServiceClient {
	return &pVZServiceClient{cc}
}

func (c *pVZServiceClient) CreatePVZ(ctx context.Context, in *CreatePVZRequest, opts ...grpc.CallOption) (*CreatePVZResponse, error) {
	out := new(CreatePVZResponse)
	err := c.cc.Invoke(ctx, PVZService_CreatePVZ_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pVZServiceClient) GetPVZ(ctx context.Context, in *GetPVZRequest, opts ...grpc.CallOption) (*GetPVZResponse, error) {
	out := new(GetPVZResponse)
	err := c.cc.Invoke(ctx, PVZService_GetPVZ_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pVZServiceClient) UpdatePVZ(ctx context.Context, in *UpdatePVZRequest, opts ...grpc.CallOption) (*UpdatePVZResponse, error) {
	out := new(UpdatePVZResponse)
	err := c.cc.Invoke(ctx, PVZService_UpdatePVZ_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pVZServiceClient) DeletePVZ(ctx context.Context, in *DeletePVZRequest, opts ...grpc.CallOption) (*DeletePVZResponse, error) {
	out := new(DeletePVZResponse)
	err := c.cc.Invoke(ctx, PVZService_DeletePVZ_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pVZServiceClient) ListAllPVZ(ctx context.Context, in *ListAllPVZRequest, opts ...grpc.CallOption) (*ListAllPVZResponse, error) {
	out := new(ListAllPVZResponse)
	err := c.cc.Invoke(ctx, PVZService_ListAllPVZ_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pVZServiceClient) DeleteListPVZ(ctx context.Context, in *DeleteListPVZRequest, opts ...grpc.CallOption) (*DeleteListPVZResponse, error) {
	out := new(DeleteListPVZResponse)
	err := c.cc.Invoke(ctx, PVZService_DeleteListPVZ_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PVZServiceServer is the server API for PVZService service.
// All implementations must embed UnimplementedPVZServiceServer
// for forward compatibility
type PVZServiceServer interface {
	CreatePVZ(context.Context, *CreatePVZRequest) (*CreatePVZResponse, error)
	GetPVZ(context.Context, *GetPVZRequest) (*GetPVZResponse, error)
	UpdatePVZ(context.Context, *UpdatePVZRequest) (*UpdatePVZResponse, error)
	DeletePVZ(context.Context, *DeletePVZRequest) (*DeletePVZResponse, error)
	ListAllPVZ(context.Context, *ListAllPVZRequest) (*ListAllPVZResponse, error)
	DeleteListPVZ(context.Context, *DeleteListPVZRequest) (*DeleteListPVZResponse, error)
	mustEmbedUnimplementedPVZServiceServer()
}

// UnimplementedPVZServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPVZServiceServer struct {
}

func (UnimplementedPVZServiceServer) CreatePVZ(context.Context, *CreatePVZRequest) (*CreatePVZResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePVZ not implemented")
}
func (UnimplementedPVZServiceServer) GetPVZ(context.Context, *GetPVZRequest) (*GetPVZResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPVZ not implemented")
}
func (UnimplementedPVZServiceServer) UpdatePVZ(context.Context, *UpdatePVZRequest) (*UpdatePVZResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePVZ not implemented")
}
func (UnimplementedPVZServiceServer) DeletePVZ(context.Context, *DeletePVZRequest) (*DeletePVZResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePVZ not implemented")
}
func (UnimplementedPVZServiceServer) ListAllPVZ(context.Context, *ListAllPVZRequest) (*ListAllPVZResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAllPVZ not implemented")
}
func (UnimplementedPVZServiceServer) DeleteListPVZ(context.Context, *DeleteListPVZRequest) (*DeleteListPVZResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteListPVZ not implemented")
}
func (UnimplementedPVZServiceServer) mustEmbedUnimplementedPVZServiceServer() {}

// UnsafePVZServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PVZServiceServer will
// result in compilation errors.
type UnsafePVZServiceServer interface {
	mustEmbedUnimplementedPVZServiceServer()
}

func RegisterPVZServiceServer(s grpc.ServiceRegistrar, srv PVZServiceServer) {
	s.RegisterService(&PVZService_ServiceDesc, srv)
}

func _PVZService_CreatePVZ_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePVZRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).CreatePVZ(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_CreatePVZ_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).CreatePVZ(ctx, req.(*CreatePVZRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PVZService_GetPVZ_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPVZRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).GetPVZ(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_GetPVZ_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).GetPVZ(ctx, req.(*GetPVZRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PVZService_UpdatePVZ_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePVZRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).UpdatePVZ(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_UpdatePVZ_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).UpdatePVZ(ctx, req.(*UpdatePVZRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PVZService_DeletePVZ_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePVZRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).DeletePVZ(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_DeletePVZ_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).DeletePVZ(ctx, req.(*DeletePVZRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PVZService_ListAllPVZ_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAllPVZRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).ListAllPVZ(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_ListAllPVZ_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).ListAllPVZ(ctx, req.(*ListAllPVZRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PVZService_DeleteListPVZ_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteListPVZRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).DeleteListPVZ(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_DeleteListPVZ_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).DeleteListPVZ(ctx, req.(*DeleteListPVZRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PVZService_ServiceDesc is the grpc.ServiceDesc for PVZService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PVZService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pvz.PVZService",
	HandlerType: (*PVZServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePVZ",
			Handler:    _PVZService_CreatePVZ_Handler,
		},
		{
			MethodName: "GetPVZ",
			Handler:    _PVZService_GetPVZ_Handler,
		},
		{
			MethodName: "UpdatePVZ",
			Handler:    _PVZService_UpdatePVZ_Handler,
		},
		{
			MethodName: "DeletePVZ",
			Handler:    _PVZService_DeletePVZ_Handler,
		},
		{
			MethodName: "ListAllPVZ",
			Handler:    _PVZService_ListAllPVZ_Handler,
		},
		{
			MethodName: "DeleteListPVZ",
			Handler:    _PVZService_DeleteListPVZ_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

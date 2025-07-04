// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.31.0
// source: api/proto/v1/service.proto

package v1

import (
	context "context"
	rpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
	user "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc/user"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	GophKeeper_Login_FullMethodName      = "/api.proto.v1.GophKeeper/Login"
	GophKeeper_Signup_FullMethodName     = "/api.proto.v1.GophKeeper/Signup"
	GophKeeper_Ping_FullMethodName       = "/api.proto.v1.GophKeeper/Ping"
	GophKeeper_DataSave_FullMethodName   = "/api.proto.v1.GophKeeper/DataSave"
	GophKeeper_DataDelete_FullMethodName = "/api.proto.v1.GophKeeper/DataDelete"
	GophKeeper_DataList_FullMethodName   = "/api.proto.v1.GophKeeper/DataList"
	GophKeeper_DataView_FullMethodName   = "/api.proto.v1.GophKeeper/DataView"
)

// GophKeeperClient is the client API for GophKeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophKeeperClient interface {
	Login(ctx context.Context, in *user.LoginRequest, opts ...grpc.CallOption) (*user.LoginResponse, error)
	Signup(ctx context.Context, in *user.SignupRequest, opts ...grpc.CallOption) (*user.SignupResponse, error)
	Ping(ctx context.Context, in *rpc.PingRequest, opts ...grpc.CallOption) (*rpc.PingResponse, error)
	DataSave(ctx context.Context, in *rpc.DataSaveRequest, opts ...grpc.CallOption) (*rpc.DataSaveResponse, error)
	DataDelete(ctx context.Context, in *rpc.DataDeleteRequest, opts ...grpc.CallOption) (*rpc.DataDeleteResponse, error)
	DataList(ctx context.Context, in *rpc.DataListRequest, opts ...grpc.CallOption) (*rpc.DataListResponse, error)
	DataView(ctx context.Context, in *rpc.DataViewRequest, opts ...grpc.CallOption) (*rpc.DataViewResponse, error)
}

type gophKeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewGophKeeperClient(cc grpc.ClientConnInterface) GophKeeperClient {
	return &gophKeeperClient{cc}
}

func (c *gophKeeperClient) Login(ctx context.Context, in *user.LoginRequest, opts ...grpc.CallOption) (*user.LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(user.LoginResponse)
	err := c.cc.Invoke(ctx, GophKeeper_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) Signup(ctx context.Context, in *user.SignupRequest, opts ...grpc.CallOption) (*user.SignupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(user.SignupResponse)
	err := c.cc.Invoke(ctx, GophKeeper_Signup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) Ping(ctx context.Context, in *rpc.PingRequest, opts ...grpc.CallOption) (*rpc.PingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(rpc.PingResponse)
	err := c.cc.Invoke(ctx, GophKeeper_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) DataSave(ctx context.Context, in *rpc.DataSaveRequest, opts ...grpc.CallOption) (*rpc.DataSaveResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(rpc.DataSaveResponse)
	err := c.cc.Invoke(ctx, GophKeeper_DataSave_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) DataDelete(ctx context.Context, in *rpc.DataDeleteRequest, opts ...grpc.CallOption) (*rpc.DataDeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(rpc.DataDeleteResponse)
	err := c.cc.Invoke(ctx, GophKeeper_DataDelete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) DataList(ctx context.Context, in *rpc.DataListRequest, opts ...grpc.CallOption) (*rpc.DataListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(rpc.DataListResponse)
	err := c.cc.Invoke(ctx, GophKeeper_DataList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) DataView(ctx context.Context, in *rpc.DataViewRequest, opts ...grpc.CallOption) (*rpc.DataViewResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(rpc.DataViewResponse)
	err := c.cc.Invoke(ctx, GophKeeper_DataView_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophKeeperServer is the server API for GophKeeper service.
// All implementations must embed UnimplementedGophKeeperServer
// for forward compatibility.
type GophKeeperServer interface {
	Login(context.Context, *user.LoginRequest) (*user.LoginResponse, error)
	Signup(context.Context, *user.SignupRequest) (*user.SignupResponse, error)
	Ping(context.Context, *rpc.PingRequest) (*rpc.PingResponse, error)
	DataSave(context.Context, *rpc.DataSaveRequest) (*rpc.DataSaveResponse, error)
	DataDelete(context.Context, *rpc.DataDeleteRequest) (*rpc.DataDeleteResponse, error)
	DataList(context.Context, *rpc.DataListRequest) (*rpc.DataListResponse, error)
	DataView(context.Context, *rpc.DataViewRequest) (*rpc.DataViewResponse, error)
	mustEmbedUnimplementedGophKeeperServer()
}

// UnimplementedGophKeeperServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGophKeeperServer struct{}

func (UnimplementedGophKeeperServer) Login(context.Context, *user.LoginRequest) (*user.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedGophKeeperServer) Signup(context.Context, *user.SignupRequest) (*user.SignupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signup not implemented")
}
func (UnimplementedGophKeeperServer) Ping(context.Context, *rpc.PingRequest) (*rpc.PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedGophKeeperServer) DataSave(context.Context, *rpc.DataSaveRequest) (*rpc.DataSaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DataSave not implemented")
}
func (UnimplementedGophKeeperServer) DataDelete(context.Context, *rpc.DataDeleteRequest) (*rpc.DataDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DataDelete not implemented")
}
func (UnimplementedGophKeeperServer) DataList(context.Context, *rpc.DataListRequest) (*rpc.DataListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DataList not implemented")
}
func (UnimplementedGophKeeperServer) DataView(context.Context, *rpc.DataViewRequest) (*rpc.DataViewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DataView not implemented")
}
func (UnimplementedGophKeeperServer) mustEmbedUnimplementedGophKeeperServer() {}
func (UnimplementedGophKeeperServer) testEmbeddedByValue()                    {}

// UnsafeGophKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophKeeperServer will
// result in compilation errors.
type UnsafeGophKeeperServer interface {
	mustEmbedUnimplementedGophKeeperServer()
}

func RegisterGophKeeperServer(s grpc.ServiceRegistrar, srv GophKeeperServer) {
	// If the following call pancis, it indicates UnimplementedGophKeeperServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GophKeeper_ServiceDesc, srv)
}

func _GophKeeper_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(user.LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Login(ctx, req.(*user.LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_Signup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(user.SignupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Signup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_Signup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Signup(ctx, req.(*user.SignupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(rpc.PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Ping(ctx, req.(*rpc.PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_DataSave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(rpc.DataSaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).DataSave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_DataSave_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).DataSave(ctx, req.(*rpc.DataSaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_DataDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(rpc.DataDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).DataDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_DataDelete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).DataDelete(ctx, req.(*rpc.DataDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_DataList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(rpc.DataListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).DataList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_DataList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).DataList(ctx, req.(*rpc.DataListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_DataView_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(rpc.DataViewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).DataView(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_DataView_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).DataView(ctx, req.(*rpc.DataViewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GophKeeper_ServiceDesc is the grpc.ServiceDesc for GophKeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GophKeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.proto.v1.GophKeeper",
	HandlerType: (*GophKeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _GophKeeper_Login_Handler,
		},
		{
			MethodName: "Signup",
			Handler:    _GophKeeper_Signup_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _GophKeeper_Ping_Handler,
		},
		{
			MethodName: "DataSave",
			Handler:    _GophKeeper_DataSave_Handler,
		},
		{
			MethodName: "DataDelete",
			Handler:    _GophKeeper_DataDelete_Handler,
		},
		{
			MethodName: "DataList",
			Handler:    _GophKeeper_DataList_Handler,
		},
		{
			MethodName: "DataView",
			Handler:    _GophKeeper_DataView_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/v1/service.proto",
}

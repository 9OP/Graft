// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: p2pRpc.proto

package p2pRpc

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

// RpcClient is the client API for Rpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RpcClient interface {
	AppendEntries(ctx context.Context, in *AppendEntriesInput, opts ...grpc.CallOption) (*AppendEntriesOutput, error)
	RequestVote(ctx context.Context, in *RequestVoteInput, opts ...grpc.CallOption) (*RequestVoteOutput, error)
	PreVote(ctx context.Context, in *RequestVoteInput, opts ...grpc.CallOption) (*RequestVoteOutput, error)
	InstallSnapshot(ctx context.Context, in *InstallSnapshotInput, opts ...grpc.CallOption) (*InstallSnapshotOutput, error)
	// Move to different service named Cluster
	Execute(ctx context.Context, in *ExecuteInput, opts ...grpc.CallOption) (*ExecuteOutput, error)
	// Rename Configuration
	ClusterConfiguration(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*ClusterConfigurationOutput, error)
	Shutdown(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*Nil, error)
	Ping(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*Nil, error)
}

type rpcClient struct {
	cc grpc.ClientConnInterface
}

func NewRpcClient(cc grpc.ClientConnInterface) RpcClient {
	return &rpcClient{cc}
}

func (c *rpcClient) AppendEntries(ctx context.Context, in *AppendEntriesInput, opts ...grpc.CallOption) (*AppendEntriesOutput, error) {
	out := new(AppendEntriesOutput)
	err := c.cc.Invoke(ctx, "/p2pRpc.Rpc/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcClient) RequestVote(ctx context.Context, in *RequestVoteInput, opts ...grpc.CallOption) (*RequestVoteOutput, error) {
	out := new(RequestVoteOutput)
	err := c.cc.Invoke(ctx, "/p2pRpc.Rpc/RequestVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcClient) PreVote(ctx context.Context, in *RequestVoteInput, opts ...grpc.CallOption) (*RequestVoteOutput, error) {
	out := new(RequestVoteOutput)
	err := c.cc.Invoke(ctx, "/p2pRpc.Rpc/PreVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcClient) InstallSnapshot(ctx context.Context, in *InstallSnapshotInput, opts ...grpc.CallOption) (*InstallSnapshotOutput, error) {
	out := new(InstallSnapshotOutput)
	err := c.cc.Invoke(ctx, "/p2pRpc.Rpc/InstallSnapshot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcClient) Execute(ctx context.Context, in *ExecuteInput, opts ...grpc.CallOption) (*ExecuteOutput, error) {
	out := new(ExecuteOutput)
	err := c.cc.Invoke(ctx, "/p2pRpc.Rpc/Execute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcClient) ClusterConfiguration(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*ClusterConfigurationOutput, error) {
	out := new(ClusterConfigurationOutput)
	err := c.cc.Invoke(ctx, "/p2pRpc.Rpc/ClusterConfiguration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcClient) Shutdown(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*Nil, error) {
	out := new(Nil)
	err := c.cc.Invoke(ctx, "/p2pRpc.Rpc/Shutdown", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcClient) Ping(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*Nil, error) {
	out := new(Nil)
	err := c.cc.Invoke(ctx, "/p2pRpc.Rpc/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RpcServer is the server API for Rpc service.
// All implementations must embed UnimplementedRpcServer
// for forward compatibility
type RpcServer interface {
	AppendEntries(context.Context, *AppendEntriesInput) (*AppendEntriesOutput, error)
	RequestVote(context.Context, *RequestVoteInput) (*RequestVoteOutput, error)
	PreVote(context.Context, *RequestVoteInput) (*RequestVoteOutput, error)
	InstallSnapshot(context.Context, *InstallSnapshotInput) (*InstallSnapshotOutput, error)
	// Move to different service named Cluster
	Execute(context.Context, *ExecuteInput) (*ExecuteOutput, error)
	// Rename Configuration
	ClusterConfiguration(context.Context, *Nil) (*ClusterConfigurationOutput, error)
	Shutdown(context.Context, *Nil) (*Nil, error)
	Ping(context.Context, *Nil) (*Nil, error)
	mustEmbedUnimplementedRpcServer()
}

// UnimplementedRpcServer must be embedded to have forward compatible implementations.
type UnimplementedRpcServer struct{}

func (UnimplementedRpcServer) AppendEntries(context.Context, *AppendEntriesInput) (*AppendEntriesOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}

func (UnimplementedRpcServer) RequestVote(context.Context, *RequestVoteInput) (*RequestVoteOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}

func (UnimplementedRpcServer) PreVote(context.Context, *RequestVoteInput) (*RequestVoteOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreVote not implemented")
}

func (UnimplementedRpcServer) InstallSnapshot(context.Context, *InstallSnapshotInput) (*InstallSnapshotOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InstallSnapshot not implemented")
}

func (UnimplementedRpcServer) Execute(context.Context, *ExecuteInput) (*ExecuteOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}

func (UnimplementedRpcServer) ClusterConfiguration(context.Context, *Nil) (*ClusterConfigurationOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClusterConfiguration not implemented")
}

func (UnimplementedRpcServer) Shutdown(context.Context, *Nil) (*Nil, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shutdown not implemented")
}

func (UnimplementedRpcServer) Ping(context.Context, *Nil) (*Nil, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedRpcServer) mustEmbedUnimplementedRpcServer() {}

// UnsafeRpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RpcServer will
// result in compilation errors.
type UnsafeRpcServer interface {
	mustEmbedUnimplementedRpcServer()
}

func RegisterRpcServer(s grpc.ServiceRegistrar, srv RpcServer) {
	s.RegisterService(&Rpc_ServiceDesc, srv)
}

func _Rpc_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2pRpc.Rpc/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServer).AppendEntries(ctx, req.(*AppendEntriesInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rpc_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVoteInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2pRpc.Rpc/RequestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServer).RequestVote(ctx, req.(*RequestVoteInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rpc_PreVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVoteInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServer).PreVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2pRpc.Rpc/PreVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServer).PreVote(ctx, req.(*RequestVoteInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rpc_InstallSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstallSnapshotInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServer).InstallSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2pRpc.Rpc/InstallSnapshot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServer).InstallSnapshot(ctx, req.(*InstallSnapshotInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rpc_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2pRpc.Rpc/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServer).Execute(ctx, req.(*ExecuteInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rpc_ClusterConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nil)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServer).ClusterConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2pRpc.Rpc/ClusterConfiguration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServer).ClusterConfiguration(ctx, req.(*Nil))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rpc_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nil)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2pRpc.Rpc/Shutdown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServer).Shutdown(ctx, req.(*Nil))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rpc_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nil)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2pRpc.Rpc/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServer).Ping(ctx, req.(*Nil))
	}
	return interceptor(ctx, in, info, handler)
}

// Rpc_ServiceDesc is the grpc.ServiceDesc for Rpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Rpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "p2pRpc.Rpc",
	HandlerType: (*RpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppendEntries",
			Handler:    _Rpc_AppendEntries_Handler,
		},
		{
			MethodName: "RequestVote",
			Handler:    _Rpc_RequestVote_Handler,
		},
		{
			MethodName: "PreVote",
			Handler:    _Rpc_PreVote_Handler,
		},
		{
			MethodName: "InstallSnapshot",
			Handler:    _Rpc_InstallSnapshot_Handler,
		},
		{
			MethodName: "Execute",
			Handler:    _Rpc_Execute_Handler,
		},
		{
			MethodName: "ClusterConfiguration",
			Handler:    _Rpc_ClusterConfiguration_Handler,
		},
		{
			MethodName: "Shutdown",
			Handler:    _Rpc_Shutdown_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Rpc_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "p2pRpc.proto",
}

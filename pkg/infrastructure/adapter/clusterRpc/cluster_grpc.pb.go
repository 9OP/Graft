// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: cluster.proto

package clusterRpc

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

// ClusterClient is the client API for Cluster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClusterClient interface {
	// P2P
	AppendEntries(ctx context.Context, in *AppendEntriesInput, opts ...grpc.CallOption) (*AppendEntriesOutput, error)
	RequestVote(ctx context.Context, in *RequestVoteInput, opts ...grpc.CallOption) (*RequestVoteOutput, error)
	PreVote(ctx context.Context, in *RequestVoteInput, opts ...grpc.CallOption) (*RequestVoteOutput, error)
	InstallSnapshot(ctx context.Context, in *InstallSnapshotInput, opts ...grpc.CallOption) (*InstallSnapshotOutput, error)
	Configuration(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*ConfigurationOutput, error)
	LeadershipTransfer(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*Nil, error)
	Execute(ctx context.Context, in *ExecuteInput, opts ...grpc.CallOption) (*ExecuteOutput, error)
	Shutdown(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*Nil, error)
	Ping(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*Nil, error)
}

type clusterClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterClient(cc grpc.ClientConnInterface) ClusterClient {
	return &clusterClient{cc}
}

func (c *clusterClient) AppendEntries(ctx context.Context, in *AppendEntriesInput, opts ...grpc.CallOption) (*AppendEntriesOutput, error) {
	out := new(AppendEntriesOutput)
	err := c.cc.Invoke(ctx, "/clusterRpc.Cluster/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) RequestVote(ctx context.Context, in *RequestVoteInput, opts ...grpc.CallOption) (*RequestVoteOutput, error) {
	out := new(RequestVoteOutput)
	err := c.cc.Invoke(ctx, "/clusterRpc.Cluster/RequestVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) PreVote(ctx context.Context, in *RequestVoteInput, opts ...grpc.CallOption) (*RequestVoteOutput, error) {
	out := new(RequestVoteOutput)
	err := c.cc.Invoke(ctx, "/clusterRpc.Cluster/PreVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) InstallSnapshot(ctx context.Context, in *InstallSnapshotInput, opts ...grpc.CallOption) (*InstallSnapshotOutput, error) {
	out := new(InstallSnapshotOutput)
	err := c.cc.Invoke(ctx, "/clusterRpc.Cluster/InstallSnapshot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) Configuration(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*ConfigurationOutput, error) {
	out := new(ConfigurationOutput)
	err := c.cc.Invoke(ctx, "/clusterRpc.Cluster/Configuration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) LeadershipTransfer(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*Nil, error) {
	out := new(Nil)
	err := c.cc.Invoke(ctx, "/clusterRpc.Cluster/LeadershipTransfer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) Execute(ctx context.Context, in *ExecuteInput, opts ...grpc.CallOption) (*ExecuteOutput, error) {
	out := new(ExecuteOutput)
	err := c.cc.Invoke(ctx, "/clusterRpc.Cluster/Execute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) Shutdown(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*Nil, error) {
	out := new(Nil)
	err := c.cc.Invoke(ctx, "/clusterRpc.Cluster/Shutdown", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) Ping(ctx context.Context, in *Nil, opts ...grpc.CallOption) (*Nil, error) {
	out := new(Nil)
	err := c.cc.Invoke(ctx, "/clusterRpc.Cluster/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterServer is the server API for Cluster service.
// All implementations must embed UnimplementedClusterServer
// for forward compatibility
type ClusterServer interface {
	// P2P
	AppendEntries(context.Context, *AppendEntriesInput) (*AppendEntriesOutput, error)
	RequestVote(context.Context, *RequestVoteInput) (*RequestVoteOutput, error)
	PreVote(context.Context, *RequestVoteInput) (*RequestVoteOutput, error)
	InstallSnapshot(context.Context, *InstallSnapshotInput) (*InstallSnapshotOutput, error)
	Configuration(context.Context, *Nil) (*ConfigurationOutput, error)
	LeadershipTransfer(context.Context, *Nil) (*Nil, error)
	Execute(context.Context, *ExecuteInput) (*ExecuteOutput, error)
	Shutdown(context.Context, *Nil) (*Nil, error)
	Ping(context.Context, *Nil) (*Nil, error)
	mustEmbedUnimplementedClusterServer()
}

// UnimplementedClusterServer must be embedded to have forward compatible implementations.
type UnimplementedClusterServer struct{}

func (UnimplementedClusterServer) AppendEntries(context.Context, *AppendEntriesInput) (*AppendEntriesOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}

func (UnimplementedClusterServer) RequestVote(context.Context, *RequestVoteInput) (*RequestVoteOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}

func (UnimplementedClusterServer) PreVote(context.Context, *RequestVoteInput) (*RequestVoteOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreVote not implemented")
}

func (UnimplementedClusterServer) InstallSnapshot(context.Context, *InstallSnapshotInput) (*InstallSnapshotOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InstallSnapshot not implemented")
}

func (UnimplementedClusterServer) Configuration(context.Context, *Nil) (*ConfigurationOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Configuration not implemented")
}

func (UnimplementedClusterServer) LeadershipTransfer(context.Context, *Nil) (*Nil, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeadershipTransfer not implemented")
}

func (UnimplementedClusterServer) Execute(context.Context, *ExecuteInput) (*ExecuteOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}

func (UnimplementedClusterServer) Shutdown(context.Context, *Nil) (*Nil, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shutdown not implemented")
}

func (UnimplementedClusterServer) Ping(context.Context, *Nil) (*Nil, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedClusterServer) mustEmbedUnimplementedClusterServer() {}

// UnsafeClusterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClusterServer will
// result in compilation errors.
type UnsafeClusterServer interface {
	mustEmbedUnimplementedClusterServer()
}

func RegisterClusterServer(s grpc.ServiceRegistrar, srv ClusterServer) {
	s.RegisterService(&Cluster_ServiceDesc, srv)
}

func _Cluster_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clusterRpc.Cluster/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).AppendEntries(ctx, req.(*AppendEntriesInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVoteInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clusterRpc.Cluster/RequestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).RequestVote(ctx, req.(*RequestVoteInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_PreVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVoteInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).PreVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clusterRpc.Cluster/PreVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).PreVote(ctx, req.(*RequestVoteInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_InstallSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstallSnapshotInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).InstallSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clusterRpc.Cluster/InstallSnapshot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).InstallSnapshot(ctx, req.(*InstallSnapshotInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_Configuration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nil)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).Configuration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clusterRpc.Cluster/Configuration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).Configuration(ctx, req.(*Nil))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_LeadershipTransfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nil)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).LeadershipTransfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clusterRpc.Cluster/LeadershipTransfer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).LeadershipTransfer(ctx, req.(*Nil))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clusterRpc.Cluster/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).Execute(ctx, req.(*ExecuteInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nil)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clusterRpc.Cluster/Shutdown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).Shutdown(ctx, req.(*Nil))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nil)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clusterRpc.Cluster/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).Ping(ctx, req.(*Nil))
	}
	return interceptor(ctx, in, info, handler)
}

// Cluster_ServiceDesc is the grpc.ServiceDesc for Cluster service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cluster_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "clusterRpc.Cluster",
	HandlerType: (*ClusterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppendEntries",
			Handler:    _Cluster_AppendEntries_Handler,
		},
		{
			MethodName: "RequestVote",
			Handler:    _Cluster_RequestVote_Handler,
		},
		{
			MethodName: "PreVote",
			Handler:    _Cluster_PreVote_Handler,
		},
		{
			MethodName: "InstallSnapshot",
			Handler:    _Cluster_InstallSnapshot_Handler,
		},
		{
			MethodName: "Configuration",
			Handler:    _Cluster_Configuration_Handler,
		},
		{
			MethodName: "LeadershipTransfer",
			Handler:    _Cluster_LeadershipTransfer_Handler,
		},
		{
			MethodName: "Execute",
			Handler:    _Cluster_Execute_Handler,
		},
		{
			MethodName: "Shutdown",
			Handler:    _Cluster_Shutdown_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Cluster_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cluster.proto",
}

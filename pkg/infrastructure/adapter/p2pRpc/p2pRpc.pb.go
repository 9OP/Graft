// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.8
// source: p2pRpc.proto

package p2pRpc

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LogEntry_LogType int32

const (
	LogEntry_LogCommand       LogEntry_LogType = 0
	LogEntry_LogQuery         LogEntry_LogType = 1
	LogEntry_LogNoop          LogEntry_LogType = 2
	LogEntry_LogConfiguration LogEntry_LogType = 3
)

// Enum value maps for LogEntry_LogType.
var (
	LogEntry_LogType_name = map[int32]string{
		0: "LogCommand",
		1: "LogQuery",
		2: "LogNoop",
		3: "LogConfiguration",
	}
	LogEntry_LogType_value = map[string]int32{
		"LogCommand":       0,
		"LogQuery":         1,
		"LogNoop":          2,
		"LogConfiguration": 3,
	}
)

func (x LogEntry_LogType) Enum() *LogEntry_LogType {
	p := new(LogEntry_LogType)
	*p = x
	return p
}

func (x LogEntry_LogType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LogEntry_LogType) Descriptor() protoreflect.EnumDescriptor {
	return file_p2pRpc_proto_enumTypes[0].Descriptor()
}

func (LogEntry_LogType) Type() protoreflect.EnumType {
	return &file_p2pRpc_proto_enumTypes[0]
}

func (x LogEntry_LogType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LogEntry_LogType.Descriptor instead.
func (LogEntry_LogType) EnumDescriptor() ([]byte, []int) {
	return file_p2pRpc_proto_rawDescGZIP(), []int{0, 0}
}

type LogEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index uint64           `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Term  uint32           `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	Data  []byte           `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Type  LogEntry_LogType `protobuf:"varint,4,opt,name=type,proto3,enum=p2pRpc.LogEntry_LogType" json:"type,omitempty"`
}

func (x *LogEntry) Reset() {
	*x = LogEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2pRpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntry) ProtoMessage() {}

func (x *LogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_p2pRpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntry.ProtoReflect.Descriptor instead.
func (*LogEntry) Descriptor() ([]byte, []int) {
	return file_p2pRpc_proto_rawDescGZIP(), []int{0}
}

func (x *LogEntry) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *LogEntry) GetTerm() uint32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *LogEntry) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *LogEntry) GetType() LogEntry_LogType {
	if x != nil {
		return x.Type
	}
	return LogEntry_LogCommand
}

type AppendEntriesInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         uint32      `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId     string      `protobuf:"bytes,2,opt,name=leaderId,proto3" json:"leaderId,omitempty"`
	PrevLogIndex uint32      `protobuf:"varint,3,opt,name=prevLogIndex,proto3" json:"prevLogIndex,omitempty"`
	PrevLogTerm  uint32      `protobuf:"varint,4,opt,name=prevLogTerm,proto3" json:"prevLogTerm,omitempty"`
	Entries      []*LogEntry `protobuf:"bytes,5,rep,name=entries,proto3" json:"entries,omitempty"`
	LeaderCommit uint32      `protobuf:"varint,6,opt,name=leaderCommit,proto3" json:"leaderCommit,omitempty"`
}

func (x *AppendEntriesInput) Reset() {
	*x = AppendEntriesInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2pRpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntriesInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesInput) ProtoMessage() {}

func (x *AppendEntriesInput) ProtoReflect() protoreflect.Message {
	mi := &file_p2pRpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesInput.ProtoReflect.Descriptor instead.
func (*AppendEntriesInput) Descriptor() ([]byte, []int) {
	return file_p2pRpc_proto_rawDescGZIP(), []int{1}
}

func (x *AppendEntriesInput) GetTerm() uint32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntriesInput) GetLeaderId() string {
	if x != nil {
		return x.LeaderId
	}
	return ""
}

func (x *AppendEntriesInput) GetPrevLogIndex() uint32 {
	if x != nil {
		return x.PrevLogIndex
	}
	return 0
}

func (x *AppendEntriesInput) GetPrevLogTerm() uint32 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

func (x *AppendEntriesInput) GetEntries() []*LogEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

func (x *AppendEntriesInput) GetLeaderCommit() uint32 {
	if x != nil {
		return x.LeaderCommit
	}
	return 0
}

type AppendEntriesOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term    uint32 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Success bool   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *AppendEntriesOutput) Reset() {
	*x = AppendEntriesOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2pRpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntriesOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesOutput) ProtoMessage() {}

func (x *AppendEntriesOutput) ProtoReflect() protoreflect.Message {
	mi := &file_p2pRpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesOutput.ProtoReflect.Descriptor instead.
func (*AppendEntriesOutput) Descriptor() ([]byte, []int) {
	return file_p2pRpc_proto_rawDescGZIP(), []int{2}
}

func (x *AppendEntriesOutput) GetTerm() uint32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntriesOutput) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type RequestVoteInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         uint32 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	CandidateId  string `protobuf:"bytes,2,opt,name=candidateId,proto3" json:"candidateId,omitempty"`
	LastLogIndex uint32 `protobuf:"varint,3,opt,name=lastLogIndex,proto3" json:"lastLogIndex,omitempty"`
	LastLogTerm  uint32 `protobuf:"varint,4,opt,name=lastLogTerm,proto3" json:"lastLogTerm,omitempty"`
}

func (x *RequestVoteInput) Reset() {
	*x = RequestVoteInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2pRpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestVoteInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVoteInput) ProtoMessage() {}

func (x *RequestVoteInput) ProtoReflect() protoreflect.Message {
	mi := &file_p2pRpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVoteInput.ProtoReflect.Descriptor instead.
func (*RequestVoteInput) Descriptor() ([]byte, []int) {
	return file_p2pRpc_proto_rawDescGZIP(), []int{3}
}

func (x *RequestVoteInput) GetTerm() uint32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RequestVoteInput) GetCandidateId() string {
	if x != nil {
		return x.CandidateId
	}
	return ""
}

func (x *RequestVoteInput) GetLastLogIndex() uint32 {
	if x != nil {
		return x.LastLogIndex
	}
	return 0
}

func (x *RequestVoteInput) GetLastLogTerm() uint32 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

type RequestVoteOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term        uint32 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	VoteGranted bool   `protobuf:"varint,2,opt,name=voteGranted,proto3" json:"voteGranted,omitempty"`
}

func (x *RequestVoteOutput) Reset() {
	*x = RequestVoteOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2pRpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestVoteOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVoteOutput) ProtoMessage() {}

func (x *RequestVoteOutput) ProtoReflect() protoreflect.Message {
	mi := &file_p2pRpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVoteOutput.ProtoReflect.Descriptor instead.
func (*RequestVoteOutput) Descriptor() ([]byte, []int) {
	return file_p2pRpc_proto_rawDescGZIP(), []int{4}
}

func (x *RequestVoteOutput) GetTerm() uint32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RequestVoteOutput) GetVoteGranted() bool {
	if x != nil {
		return x.VoteGranted
	}
	return false
}

type InstallSnapshotInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term              uint32 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId          string `protobuf:"bytes,2,opt,name=leaderId,proto3" json:"leaderId,omitempty"`
	LastIncludedIndex uint32 `protobuf:"varint,3,opt,name=lastIncludedIndex,proto3" json:"lastIncludedIndex,omitempty"`
	LastIncludedTerm  uint32 `protobuf:"varint,4,opt,name=lastIncludedTerm,proto3" json:"lastIncludedTerm,omitempty"`
	Offset            uint32 `protobuf:"varint,5,opt,name=offset,proto3" json:"offset,omitempty"`
	Data              []byte `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *InstallSnapshotInput) Reset() {
	*x = InstallSnapshotInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2pRpc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstallSnapshotInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstallSnapshotInput) ProtoMessage() {}

func (x *InstallSnapshotInput) ProtoReflect() protoreflect.Message {
	mi := &file_p2pRpc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstallSnapshotInput.ProtoReflect.Descriptor instead.
func (*InstallSnapshotInput) Descriptor() ([]byte, []int) {
	return file_p2pRpc_proto_rawDescGZIP(), []int{5}
}

func (x *InstallSnapshotInput) GetTerm() uint32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *InstallSnapshotInput) GetLeaderId() string {
	if x != nil {
		return x.LeaderId
	}
	return ""
}

func (x *InstallSnapshotInput) GetLastIncludedIndex() uint32 {
	if x != nil {
		return x.LastIncludedIndex
	}
	return 0
}

func (x *InstallSnapshotInput) GetLastIncludedTerm() uint32 {
	if x != nil {
		return x.LastIncludedTerm
	}
	return 0
}

func (x *InstallSnapshotInput) GetOffset() uint32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *InstallSnapshotInput) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type InstallSnapshotOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term uint32 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
}

func (x *InstallSnapshotOutput) Reset() {
	*x = InstallSnapshotOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2pRpc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstallSnapshotOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstallSnapshotOutput) ProtoMessage() {}

func (x *InstallSnapshotOutput) ProtoReflect() protoreflect.Message {
	mi := &file_p2pRpc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstallSnapshotOutput.ProtoReflect.Descriptor instead.
func (*InstallSnapshotOutput) Descriptor() ([]byte, []int) {
	return file_p2pRpc_proto_rawDescGZIP(), []int{6}
}

func (x *InstallSnapshotOutput) GetTerm() uint32 {
	if x != nil {
		return x.Term
	}
	return 0
}

var File_p2pRpc_proto protoreflect.FileDescriptor

var file_p2pRpc_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x32, 0x70, 0x52, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x70, 0x32, 0x70, 0x52, 0x70, 0x63, 0x22, 0xc2, 0x01, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x2c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x18, 0x2e, 0x70, 0x32, 0x70, 0x52, 0x70, 0x63, 0x2e, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22,
	0x4a, 0x0a, 0x07, 0x4c, 0x6f, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x4c, 0x6f,
	0x67, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x4c, 0x6f,
	0x67, 0x51, 0x75, 0x65, 0x72, 0x79, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x4c, 0x6f, 0x67, 0x4e,
	0x6f, 0x6f, 0x70, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x4c, 0x6f, 0x67, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x03, 0x22, 0xda, 0x01, 0x0a, 0x12,
	0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f,
	0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f,
	0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x70, 0x72, 0x65,
	0x76, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x2a, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72,
	0x69, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x32, 0x70, 0x52,
	0x70, 0x63, 0x2e, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74,
	0x72, 0x69, 0x65, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f,
	0x6d, 0x6d, 0x69, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x6c, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x22, 0x43, 0x0a, 0x13, 0x41, 0x70, 0x70, 0x65,
	0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74,
	0x65, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x8e, 0x01,
	0x0a, 0x10, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x61, 0x6e,
	0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74,
	0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c,
	0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b,
	0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x49,
	0x0a, 0x11, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x4f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x76, 0x6f, 0x74, 0x65, 0x47,
	0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x76, 0x6f,
	0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x22, 0xcc, 0x01, 0x0a, 0x14, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6c, 0x6c, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x2c, 0x0a, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x49, 0x6e, 0x63, 0x6c, 0x75, 0x64,
	0x65, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x6c,
	0x61, 0x73, 0x74, 0x49, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x12, 0x2a, 0x0a, 0x10, 0x6c, 0x61, 0x73, 0x74, 0x49, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x64,
	0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x10, 0x6c, 0x61, 0x73, 0x74,
	0x49, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x64, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x16, 0x0a, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x2b, 0x0a, 0x15, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6c, 0x6c, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x4f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x04, 0x74, 0x65, 0x72, 0x6d, 0x32, 0xab, 0x02, 0x0a, 0x03, 0x52, 0x70, 0x63, 0x12, 0x4a, 0x0a,
	0x0d, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x1a,
	0x2e, 0x70, 0x32, 0x70, 0x52, 0x70, 0x63, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e,
	0x74, 0x72, 0x69, 0x65, 0x73, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x32, 0x70,
	0x52, 0x70, 0x63, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65,
	0x73, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x0b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x70, 0x32, 0x70, 0x52, 0x70,
	0x63, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x32, 0x70, 0x52, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x12,
	0x40, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x70, 0x32, 0x70,
	0x52, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x49,
	0x6e, 0x70, 0x75, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x32, 0x70, 0x52, 0x70, 0x63, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22,
	0x00, 0x12, 0x50, 0x0a, 0x0f, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x53, 0x6e, 0x61, 0x70,
	0x73, 0x68, 0x6f, 0x74, 0x12, 0x1c, 0x2e, 0x70, 0x32, 0x70, 0x52, 0x70, 0x63, 0x2e, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6c, 0x6c, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x32, 0x70, 0x52, 0x70, 0x63, 0x2e, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6c, 0x6c, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x4f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x70, 0x32, 0x70, 0x52, 0x70, 0x63, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_p2pRpc_proto_rawDescOnce sync.Once
	file_p2pRpc_proto_rawDescData = file_p2pRpc_proto_rawDesc
)

func file_p2pRpc_proto_rawDescGZIP() []byte {
	file_p2pRpc_proto_rawDescOnce.Do(func() {
		file_p2pRpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_p2pRpc_proto_rawDescData)
	})
	return file_p2pRpc_proto_rawDescData
}

var (
	file_p2pRpc_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
	file_p2pRpc_proto_msgTypes  = make([]protoimpl.MessageInfo, 7)
	file_p2pRpc_proto_goTypes   = []interface{}{
		(LogEntry_LogType)(0),         // 0: p2pRpc.LogEntry.LogType
		(*LogEntry)(nil),              // 1: p2pRpc.LogEntry
		(*AppendEntriesInput)(nil),    // 2: p2pRpc.AppendEntriesInput
		(*AppendEntriesOutput)(nil),   // 3: p2pRpc.AppendEntriesOutput
		(*RequestVoteInput)(nil),      // 4: p2pRpc.RequestVoteInput
		(*RequestVoteOutput)(nil),     // 5: p2pRpc.RequestVoteOutput
		(*InstallSnapshotInput)(nil),  // 6: p2pRpc.InstallSnapshotInput
		(*InstallSnapshotOutput)(nil), // 7: p2pRpc.InstallSnapshotOutput
	}
)
var file_p2pRpc_proto_depIdxs = []int32{
	0, // 0: p2pRpc.LogEntry.type:type_name -> p2pRpc.LogEntry.LogType
	1, // 1: p2pRpc.AppendEntriesInput.entries:type_name -> p2pRpc.LogEntry
	2, // 2: p2pRpc.Rpc.AppendEntries:input_type -> p2pRpc.AppendEntriesInput
	4, // 3: p2pRpc.Rpc.RequestVote:input_type -> p2pRpc.RequestVoteInput
	4, // 4: p2pRpc.Rpc.PreVote:input_type -> p2pRpc.RequestVoteInput
	6, // 5: p2pRpc.Rpc.InstallSnapshot:input_type -> p2pRpc.InstallSnapshotInput
	3, // 6: p2pRpc.Rpc.AppendEntries:output_type -> p2pRpc.AppendEntriesOutput
	5, // 7: p2pRpc.Rpc.RequestVote:output_type -> p2pRpc.RequestVoteOutput
	5, // 8: p2pRpc.Rpc.PreVote:output_type -> p2pRpc.RequestVoteOutput
	7, // 9: p2pRpc.Rpc.InstallSnapshot:output_type -> p2pRpc.InstallSnapshotOutput
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_p2pRpc_proto_init() }
func file_p2pRpc_proto_init() {
	if File_p2pRpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_p2pRpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_p2pRpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntriesInput); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_p2pRpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntriesOutput); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_p2pRpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestVoteInput); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_p2pRpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestVoteOutput); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_p2pRpc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstallSnapshotInput); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_p2pRpc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstallSnapshotOutput); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_p2pRpc_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_p2pRpc_proto_goTypes,
		DependencyIndexes: file_p2pRpc_proto_depIdxs,
		EnumInfos:         file_p2pRpc_proto_enumTypes,
		MessageInfos:      file_p2pRpc_proto_msgTypes,
	}.Build()
	File_p2pRpc_proto = out.File
	file_p2pRpc_proto_rawDesc = nil
	file_p2pRpc_proto_goTypes = nil
	file_p2pRpc_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: internal/grpc/logPassword/logPassword.proto

package grpc_logPassword

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LogPasswordEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	LoginHash    string `protobuf:"bytes,2,opt,name=login_hash,json=loginHash,proto3" json:"login_hash,omitempty"`
	PasswordHash string `protobuf:"bytes,3,opt,name=password_hash,json=passwordHash,proto3" json:"password_hash,omitempty"`
	ResourceName string `protobuf:"bytes,4,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	EntryHash    string `protobuf:"bytes,5,opt,name=entry_hash,json=entryHash,proto3" json:"entry_hash,omitempty"`
}

func (x *LogPasswordEntry) Reset() {
	*x = LogPasswordEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogPasswordEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogPasswordEntry) ProtoMessage() {}

func (x *LogPasswordEntry) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogPasswordEntry.ProtoReflect.Descriptor instead.
func (*LogPasswordEntry) Descriptor() ([]byte, []int) {
	return file_internal_grpc_logPassword_logPassword_proto_rawDescGZIP(), []int{0}
}

func (x *LogPasswordEntry) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *LogPasswordEntry) GetLoginHash() string {
	if x != nil {
		return x.LoginHash
	}
	return ""
}

func (x *LogPasswordEntry) GetPasswordHash() string {
	if x != nil {
		return x.PasswordHash
	}
	return ""
}

func (x *LogPasswordEntry) GetResourceName() string {
	if x != nil {
		return x.ResourceName
	}
	return ""
}

func (x *LogPasswordEntry) GetEntryHash() string {
	if x != nil {
		return x.EntryHash
	}
	return ""
}

type ExistingLogPasswordHash struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	EntryHash string `protobuf:"bytes,2,opt,name=entry_hash,json=entryHash,proto3" json:"entry_hash,omitempty"`
}

func (x *ExistingLogPasswordHash) Reset() {
	*x = ExistingLogPasswordHash{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExistingLogPasswordHash) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExistingLogPasswordHash) ProtoMessage() {}

func (x *ExistingLogPasswordHash) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExistingLogPasswordHash.ProtoReflect.Descriptor instead.
func (*ExistingLogPasswordHash) Descriptor() ([]byte, []int) {
	return file_internal_grpc_logPassword_logPassword_proto_rawDescGZIP(), []int{1}
}

func (x *ExistingLogPasswordHash) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ExistingLogPasswordHash) GetEntryHash() string {
	if x != nil {
		return x.EntryHash
	}
	return ""
}

type ListLogPasswordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExistingHashes []*ExistingLogPasswordHash `protobuf:"bytes,1,rep,name=existing_hashes,json=existingHashes,proto3" json:"existing_hashes,omitempty"`
}

func (x *ListLogPasswordRequest) Reset() {
	*x = ListLogPasswordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListLogPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLogPasswordRequest) ProtoMessage() {}

func (x *ListLogPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLogPasswordRequest.ProtoReflect.Descriptor instead.
func (*ListLogPasswordRequest) Descriptor() ([]byte, []int) {
	return file_internal_grpc_logPassword_logPassword_proto_rawDescGZIP(), []int{2}
}

func (x *ListLogPasswordRequest) GetExistingHashes() []*ExistingLogPasswordHash {
	if x != nil {
		return x.ExistingHashes
	}
	return nil
}

type ListLogPasswordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error        string              `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	LogPasswords []*LogPasswordEntry `protobuf:"bytes,2,rep,name=log_passwords,json=logPasswords,proto3" json:"log_passwords,omitempty"`
}

func (x *ListLogPasswordResponse) Reset() {
	*x = ListLogPasswordResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListLogPasswordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLogPasswordResponse) ProtoMessage() {}

func (x *ListLogPasswordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLogPasswordResponse.ProtoReflect.Descriptor instead.
func (*ListLogPasswordResponse) Descriptor() ([]byte, []int) {
	return file_internal_grpc_logPassword_logPassword_proto_rawDescGZIP(), []int{3}
}

func (x *ListLogPasswordResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *ListLogPasswordResponse) GetLogPasswords() []*LogPasswordEntry {
	if x != nil {
		return x.LogPasswords
	}
	return nil
}

type CreateLogPasswordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LoginHash    string `protobuf:"bytes,1,opt,name=login_hash,json=loginHash,proto3" json:"login_hash,omitempty"`
	PasswordHash string `protobuf:"bytes,2,opt,name=password_hash,json=passwordHash,proto3" json:"password_hash,omitempty"`
	ResourceName string `protobuf:"bytes,3,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
}

func (x *CreateLogPasswordRequest) Reset() {
	*x = CreateLogPasswordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLogPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLogPasswordRequest) ProtoMessage() {}

func (x *CreateLogPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLogPasswordRequest.ProtoReflect.Descriptor instead.
func (*CreateLogPasswordRequest) Descriptor() ([]byte, []int) {
	return file_internal_grpc_logPassword_logPassword_proto_rawDescGZIP(), []int{4}
}

func (x *CreateLogPasswordRequest) GetLoginHash() string {
	if x != nil {
		return x.LoginHash
	}
	return ""
}

func (x *CreateLogPasswordRequest) GetPasswordHash() string {
	if x != nil {
		return x.PasswordHash
	}
	return ""
}

func (x *CreateLogPasswordRequest) GetResourceName() string {
	if x != nil {
		return x.ResourceName
	}
	return ""
}

type CreateLogPasswordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Id    string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CreateLogPasswordResponse) Reset() {
	*x = CreateLogPasswordResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLogPasswordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLogPasswordResponse) ProtoMessage() {}

func (x *CreateLogPasswordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLogPasswordResponse.ProtoReflect.Descriptor instead.
func (*CreateLogPasswordResponse) Descriptor() ([]byte, []int) {
	return file_internal_grpc_logPassword_logPassword_proto_rawDescGZIP(), []int{5}
}

func (x *CreateLogPasswordResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *CreateLogPasswordResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type UpdateLogPasswordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	LoginHash    string `protobuf:"bytes,2,opt,name=login_hash,json=loginHash,proto3" json:"login_hash,omitempty"`
	PasswordHash string `protobuf:"bytes,3,opt,name=password_hash,json=passwordHash,proto3" json:"password_hash,omitempty"`
	ResourceName string `protobuf:"bytes,4,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	PreviousHash string `protobuf:"bytes,5,opt,name=previous_hash,json=previousHash,proto3" json:"previous_hash,omitempty"`
	IsDeleted    bool   `protobuf:"varint,6,opt,name=is_deleted,json=isDeleted,proto3" json:"is_deleted,omitempty"`
	ForceUpdate  bool   `protobuf:"varint,7,opt,name=force_update,json=forceUpdate,proto3" json:"force_update,omitempty"`
}

func (x *UpdateLogPasswordRequest) Reset() {
	*x = UpdateLogPasswordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateLogPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateLogPasswordRequest) ProtoMessage() {}

func (x *UpdateLogPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateLogPasswordRequest.ProtoReflect.Descriptor instead.
func (*UpdateLogPasswordRequest) Descriptor() ([]byte, []int) {
	return file_internal_grpc_logPassword_logPassword_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateLogPasswordRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateLogPasswordRequest) GetLoginHash() string {
	if x != nil {
		return x.LoginHash
	}
	return ""
}

func (x *UpdateLogPasswordRequest) GetPasswordHash() string {
	if x != nil {
		return x.PasswordHash
	}
	return ""
}

func (x *UpdateLogPasswordRequest) GetResourceName() string {
	if x != nil {
		return x.ResourceName
	}
	return ""
}

func (x *UpdateLogPasswordRequest) GetPreviousHash() string {
	if x != nil {
		return x.PreviousHash
	}
	return ""
}

func (x *UpdateLogPasswordRequest) GetIsDeleted() bool {
	if x != nil {
		return x.IsDeleted
	}
	return false
}

func (x *UpdateLogPasswordRequest) GetForceUpdate() bool {
	if x != nil {
		return x.ForceUpdate
	}
	return false
}

type UpdateLogPasswordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UpdateLogPasswordResponse) Reset() {
	*x = UpdateLogPasswordResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateLogPasswordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateLogPasswordResponse) ProtoMessage() {}

func (x *UpdateLogPasswordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpc_logPassword_logPassword_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateLogPasswordResponse.ProtoReflect.Descriptor instead.
func (*UpdateLogPasswordResponse) Descriptor() ([]byte, []int) {
	return file_internal_grpc_logPassword_logPassword_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateLogPasswordResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_internal_grpc_logPassword_logPassword_proto protoreflect.FileDescriptor

var file_internal_grpc_logPassword_logPassword_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f,
	0x6c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2f, 0x6c, 0x6f, 0x67, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaa, 0x01,
	0x0a, 0x10, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x48, 0x61, 0x73, 0x68, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x65,
	0x6e, 0x74, 0x72, 0x79, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x48, 0x61, 0x73, 0x68, 0x22, 0x48, 0x0a, 0x17, 0x45, 0x78,
	0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x48, 0x61, 0x73, 0x68, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x6e, 0x74, 0x72, 0x79,
	0x48, 0x61, 0x73, 0x68, 0x22, 0x5b, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x41,
	0x0a, 0x0f, 0x65, 0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x45, 0x78, 0x69, 0x73, 0x74, 0x69,
	0x6e, 0x67, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x48, 0x61, 0x73,
	0x68, 0x52, 0x0e, 0x65, 0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x48, 0x61, 0x73, 0x68, 0x65,
	0x73, 0x22, 0x67, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x12, 0x36, 0x0a, 0x0d, 0x6c, 0x6f, 0x67, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x4c, 0x6f, 0x67, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0c, 0x6c, 0x6f,
	0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x22, 0x83, 0x01, 0x0a, 0x18, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x67, 0x69, 0x6e,
	0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f, 0x67,
	0x69, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x48, 0x61, 0x73, 0x68, 0x12, 0x23, 0x0a, 0x0d, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0x41, 0x0a, 0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0xfa, 0x01, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x6f,
	0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x23, 0x0a, 0x0d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72, 0x65,
	0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d,
	0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x21, 0x0a,
	0x0c, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0b, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x22, 0x31, 0x0a, 0x19, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x32, 0xeb, 0x01, 0x0a, 0x0b, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x12, 0x44, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x17, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x6f, 0x67,
	0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x18, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x11, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x19,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c,
	0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x19, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x4c, 0x6f, 0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x6f,
	0x67, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x12, 0x5a, 0x10, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x6c, 0x6f, 0x67, 0x50, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_grpc_logPassword_logPassword_proto_rawDescOnce sync.Once
	file_internal_grpc_logPassword_logPassword_proto_rawDescData = file_internal_grpc_logPassword_logPassword_proto_rawDesc
)

func file_internal_grpc_logPassword_logPassword_proto_rawDescGZIP() []byte {
	file_internal_grpc_logPassword_logPassword_proto_rawDescOnce.Do(func() {
		file_internal_grpc_logPassword_logPassword_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_grpc_logPassword_logPassword_proto_rawDescData)
	})
	return file_internal_grpc_logPassword_logPassword_proto_rawDescData
}

var file_internal_grpc_logPassword_logPassword_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_internal_grpc_logPassword_logPassword_proto_goTypes = []interface{}{
	(*LogPasswordEntry)(nil),          // 0: LogPasswordEntry
	(*ExistingLogPasswordHash)(nil),   // 1: ExistingLogPasswordHash
	(*ListLogPasswordRequest)(nil),    // 2: ListLogPasswordRequest
	(*ListLogPasswordResponse)(nil),   // 3: ListLogPasswordResponse
	(*CreateLogPasswordRequest)(nil),  // 4: CreateLogPasswordRequest
	(*CreateLogPasswordResponse)(nil), // 5: CreateLogPasswordResponse
	(*UpdateLogPasswordRequest)(nil),  // 6: UpdateLogPasswordRequest
	(*UpdateLogPasswordResponse)(nil), // 7: UpdateLogPasswordResponse
}
var file_internal_grpc_logPassword_logPassword_proto_depIdxs = []int32{
	1, // 0: ListLogPasswordRequest.existing_hashes:type_name -> ExistingLogPasswordHash
	0, // 1: ListLogPasswordResponse.log_passwords:type_name -> LogPasswordEntry
	2, // 2: LogPassword.ListLogPassword:input_type -> ListLogPasswordRequest
	4, // 3: LogPassword.CreateLogPassword:input_type -> CreateLogPasswordRequest
	6, // 4: LogPassword.UpdateLogPassword:input_type -> UpdateLogPasswordRequest
	3, // 5: LogPassword.ListLogPassword:output_type -> ListLogPasswordResponse
	5, // 6: LogPassword.CreateLogPassword:output_type -> CreateLogPasswordResponse
	7, // 7: LogPassword.UpdateLogPassword:output_type -> UpdateLogPasswordResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_grpc_logPassword_logPassword_proto_init() }
func file_internal_grpc_logPassword_logPassword_proto_init() {
	if File_internal_grpc_logPassword_logPassword_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_grpc_logPassword_logPassword_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogPasswordEntry); i {
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
		file_internal_grpc_logPassword_logPassword_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExistingLogPasswordHash); i {
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
		file_internal_grpc_logPassword_logPassword_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListLogPasswordRequest); i {
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
		file_internal_grpc_logPassword_logPassword_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListLogPasswordResponse); i {
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
		file_internal_grpc_logPassword_logPassword_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLogPasswordRequest); i {
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
		file_internal_grpc_logPassword_logPassword_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLogPasswordResponse); i {
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
		file_internal_grpc_logPassword_logPassword_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateLogPasswordRequest); i {
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
		file_internal_grpc_logPassword_logPassword_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateLogPasswordResponse); i {
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
			RawDescriptor: file_internal_grpc_logPassword_logPassword_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_grpc_logPassword_logPassword_proto_goTypes,
		DependencyIndexes: file_internal_grpc_logPassword_logPassword_proto_depIdxs,
		MessageInfos:      file_internal_grpc_logPassword_logPassword_proto_msgTypes,
	}.Build()
	File_internal_grpc_logPassword_logPassword_proto = out.File
	file_internal_grpc_logPassword_logPassword_proto_rawDesc = nil
	file_internal_grpc_logPassword_logPassword_proto_goTypes = nil
	file_internal_grpc_logPassword_logPassword_proto_depIdxs = nil
}

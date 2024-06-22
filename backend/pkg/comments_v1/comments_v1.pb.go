// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.25.3
// source: protos/comments_v1/comments_v1.proto

package comments_v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Запрос комментариев к фильму
type GetCommentsByIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MovieId int32 `protobuf:"varint,1,opt,name=movieId,proto3" json:"movieId,omitempty"`
	Limit   int32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Page    int32 `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`
}

func (x *GetCommentsByIdRequest) Reset() {
	*x = GetCommentsByIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCommentsByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentsByIdRequest) ProtoMessage() {}

func (x *GetCommentsByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentsByIdRequest.ProtoReflect.Descriptor instead.
func (*GetCommentsByIdRequest) Descriptor() ([]byte, []int) {
	return file_protos_comments_v1_comments_v1_proto_rawDescGZIP(), []int{0}
}

func (x *GetCommentsByIdRequest) GetMovieId() int32 {
	if x != nil {
		return x.MovieId
	}
	return 0
}

func (x *GetCommentsByIdRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *GetCommentsByIdRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

// Запрос на добавление комментария
type AddCommentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ParentId  int32                  `protobuf:"varint,1,opt,name=parentId,proto3" json:"parentId,omitempty"`
	MovieId   int32                  `protobuf:"varint,2,opt,name=movieId,proto3" json:"movieId,omitempty"`
	UserId    int32                  `protobuf:"varint,3,opt,name=userId,proto3" json:"userId,omitempty"`
	Text      string                 `protobuf:"bytes,4,opt,name=text,proto3" json:"text,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
}

func (x *AddCommentRequest) Reset() {
	*x = AddCommentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddCommentRequest) ProtoMessage() {}

func (x *AddCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddCommentRequest.ProtoReflect.Descriptor instead.
func (*AddCommentRequest) Descriptor() ([]byte, []int) {
	return file_protos_comments_v1_comments_v1_proto_rawDescGZIP(), []int{1}
}

func (x *AddCommentRequest) GetParentId() int32 {
	if x != nil {
		return x.ParentId
	}
	return 0
}

func (x *AddCommentRequest) GetMovieId() int32 {
	if x != nil {
		return x.MovieId
	}
	return 0
}

func (x *AddCommentRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *AddCommentRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *AddCommentRequest) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

// Запрос на обновление комментария
type UpdateCommentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Text      string                 `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
}

func (x *UpdateCommentRequest) Reset() {
	*x = UpdateCommentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCommentRequest) ProtoMessage() {}

func (x *UpdateCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCommentRequest.ProtoReflect.Descriptor instead.
func (*UpdateCommentRequest) Descriptor() ([]byte, []int) {
	return file_protos_comments_v1_comments_v1_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateCommentRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateCommentRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *UpdateCommentRequest) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// Запрос на удаление комментария
type DelCommentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ParentId int32 `protobuf:"varint,2,opt,name=parentId,proto3" json:"parentId,omitempty"`
}

func (x *DelCommentRequest) Reset() {
	*x = DelCommentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelCommentRequest) ProtoMessage() {}

func (x *DelCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelCommentRequest.ProtoReflect.Descriptor instead.
func (*DelCommentRequest) Descriptor() ([]byte, []int) {
	return file_protos_comments_v1_comments_v1_proto_rawDescGZIP(), []int{3}
}

func (x *DelCommentRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DelCommentRequest) GetParentId() int32 {
	if x != nil {
		return x.ParentId
	}
	return 0
}

// Вывод комметариев к фильму
type GetCommentsByIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Comments []*GetCommentsByIdItem `protobuf:"bytes,1,rep,name=comments,proto3" json:"comments,omitempty"`
}

func (x *GetCommentsByIdResponse) Reset() {
	*x = GetCommentsByIdResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCommentsByIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentsByIdResponse) ProtoMessage() {}

func (x *GetCommentsByIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentsByIdResponse.ProtoReflect.Descriptor instead.
func (*GetCommentsByIdResponse) Descriptor() ([]byte, []int) {
	return file_protos_comments_v1_comments_v1_proto_rawDescGZIP(), []int{4}
}

func (x *GetCommentsByIdResponse) GetComments() []*GetCommentsByIdItem {
	if x != nil {
		return x.Comments
	}
	return nil
}

type GetCommentsByIdItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ParentId  int32                  `protobuf:"varint,2,opt,name=parentId,proto3" json:"parentId,omitempty"`
	UserId    int32                  `protobuf:"varint,3,opt,name=userId,proto3" json:"userId,omitempty"`
	Text      string                 `protobuf:"bytes,4,opt,name=text,proto3" json:"text,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
}

func (x *GetCommentsByIdItem) Reset() {
	*x = GetCommentsByIdItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCommentsByIdItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentsByIdItem) ProtoMessage() {}

func (x *GetCommentsByIdItem) ProtoReflect() protoreflect.Message {
	mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentsByIdItem.ProtoReflect.Descriptor instead.
func (*GetCommentsByIdItem) Descriptor() ([]byte, []int) {
	return file_protos_comments_v1_comments_v1_proto_rawDescGZIP(), []int{5}
}

func (x *GetCommentsByIdItem) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GetCommentsByIdItem) GetParentId() int32 {
	if x != nil {
		return x.ParentId
	}
	return 0
}

func (x *GetCommentsByIdItem) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *GetCommentsByIdItem) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *GetCommentsByIdItem) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *GetCommentsByIdItem) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// Добавление комментария
type AddCommentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Err string `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *AddCommentResponse) Reset() {
	*x = AddCommentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddCommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddCommentResponse) ProtoMessage() {}

func (x *AddCommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddCommentResponse.ProtoReflect.Descriptor instead.
func (*AddCommentResponse) Descriptor() ([]byte, []int) {
	return file_protos_comments_v1_comments_v1_proto_rawDescGZIP(), []int{6}
}

func (x *AddCommentResponse) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AddCommentResponse) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

// Обновление комментария
type UpdateCommentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *UpdateCommentResponse) Reset() {
	*x = UpdateCommentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateCommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCommentResponse) ProtoMessage() {}

func (x *UpdateCommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCommentResponse.ProtoReflect.Descriptor instead.
func (*UpdateCommentResponse) Descriptor() ([]byte, []int) {
	return file_protos_comments_v1_comments_v1_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateCommentResponse) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

// Удаление комментария
type DelCommentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *DelCommentResponse) Reset() {
	*x = DelCommentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelCommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelCommentResponse) ProtoMessage() {}

func (x *DelCommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_comments_v1_comments_v1_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelCommentResponse.ProtoReflect.Descriptor instead.
func (*DelCommentResponse) Descriptor() ([]byte, []int) {
	return file_protos_comments_v1_comments_v1_proto_rawDescGZIP(), []int{8}
}

func (x *DelCommentResponse) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

var File_protos_comments_v1_comments_v1_proto protoreflect.FileDescriptor

var file_protos_comments_v1_comments_v1_proto_rawDesc = []byte{
	0x0a, 0x24, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x5f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x76, 0x31,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x5f, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5c, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x07, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61,
	0x67, 0x65, 0x22, 0xaf, 0x01, 0x0a, 0x11, 0x41, 0x64, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x49, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x49, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x22, 0x74, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74,
	0x12, 0x38, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x3f, 0x0a, 0x11, 0x44, 0x65,
	0x6c, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x08, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x57, 0x0a, 0x17, 0x47,
	0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x5f, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x42, 0x79, 0x49, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x22, 0xe1, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x42, 0x79, 0x49, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x65, 0x78, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x38,
	0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x36, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x43,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10,
	0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72,
	0x22, 0x29, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x22, 0x26, 0x0a, 0x12, 0x44,
	0x65, 0x6c, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x65, 0x72, 0x72, 0x32, 0xe0, 0x02, 0x0a, 0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x56, 0x31, 0x12, 0x5c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x42, 0x79, 0x49, 0x64, 0x12, 0x23, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x5f, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x42,
	0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x4d, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1e,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64,
	0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64,
	0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x56, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x21, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x76, 0x31, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x76,
	0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x43, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x5f, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x5f, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_comments_v1_comments_v1_proto_rawDescOnce sync.Once
	file_protos_comments_v1_comments_v1_proto_rawDescData = file_protos_comments_v1_comments_v1_proto_rawDesc
)

func file_protos_comments_v1_comments_v1_proto_rawDescGZIP() []byte {
	file_protos_comments_v1_comments_v1_proto_rawDescOnce.Do(func() {
		file_protos_comments_v1_comments_v1_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_comments_v1_comments_v1_proto_rawDescData)
	})
	return file_protos_comments_v1_comments_v1_proto_rawDescData
}

var file_protos_comments_v1_comments_v1_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_protos_comments_v1_comments_v1_proto_goTypes = []interface{}{
	(*GetCommentsByIdRequest)(nil),  // 0: comments_v1.GetCommentsByIdRequest
	(*AddCommentRequest)(nil),       // 1: comments_v1.AddCommentRequest
	(*UpdateCommentRequest)(nil),    // 2: comments_v1.UpdateCommentRequest
	(*DelCommentRequest)(nil),       // 3: comments_v1.DelCommentRequest
	(*GetCommentsByIdResponse)(nil), // 4: comments_v1.GetCommentsByIdResponse
	(*GetCommentsByIdItem)(nil),     // 5: comments_v1.GetCommentsByIdItem
	(*AddCommentResponse)(nil),      // 6: comments_v1.AddCommentResponse
	(*UpdateCommentResponse)(nil),   // 7: comments_v1.UpdateCommentResponse
	(*DelCommentResponse)(nil),      // 8: comments_v1.DelCommentResponse
	(*timestamppb.Timestamp)(nil),   // 9: google.protobuf.Timestamp
}
var file_protos_comments_v1_comments_v1_proto_depIdxs = []int32{
	9, // 0: comments_v1.AddCommentRequest.createdAt:type_name -> google.protobuf.Timestamp
	9, // 1: comments_v1.UpdateCommentRequest.updatedAt:type_name -> google.protobuf.Timestamp
	5, // 2: comments_v1.GetCommentsByIdResponse.comments:type_name -> comments_v1.GetCommentsByIdItem
	9, // 3: comments_v1.GetCommentsByIdItem.createdAt:type_name -> google.protobuf.Timestamp
	9, // 4: comments_v1.GetCommentsByIdItem.updatedAt:type_name -> google.protobuf.Timestamp
	0, // 5: comments_v1.CommentsV1.GetCommentsById:input_type -> comments_v1.GetCommentsByIdRequest
	1, // 6: comments_v1.CommentsV1.AddComment:input_type -> comments_v1.AddCommentRequest
	2, // 7: comments_v1.CommentsV1.UpdateComment:input_type -> comments_v1.UpdateCommentRequest
	3, // 8: comments_v1.CommentsV1.DelComment:input_type -> comments_v1.DelCommentRequest
	4, // 9: comments_v1.CommentsV1.GetCommentsById:output_type -> comments_v1.GetCommentsByIdResponse
	6, // 10: comments_v1.CommentsV1.AddComment:output_type -> comments_v1.AddCommentResponse
	7, // 11: comments_v1.CommentsV1.UpdateComment:output_type -> comments_v1.UpdateCommentResponse
	8, // 12: comments_v1.CommentsV1.DelComment:output_type -> comments_v1.DelCommentResponse
	9, // [9:13] is the sub-list for method output_type
	5, // [5:9] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_protos_comments_v1_comments_v1_proto_init() }
func file_protos_comments_v1_comments_v1_proto_init() {
	if File_protos_comments_v1_comments_v1_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_comments_v1_comments_v1_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCommentsByIdRequest); i {
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
		file_protos_comments_v1_comments_v1_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddCommentRequest); i {
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
		file_protos_comments_v1_comments_v1_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateCommentRequest); i {
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
		file_protos_comments_v1_comments_v1_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelCommentRequest); i {
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
		file_protos_comments_v1_comments_v1_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCommentsByIdResponse); i {
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
		file_protos_comments_v1_comments_v1_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCommentsByIdItem); i {
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
		file_protos_comments_v1_comments_v1_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddCommentResponse); i {
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
		file_protos_comments_v1_comments_v1_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateCommentResponse); i {
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
		file_protos_comments_v1_comments_v1_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelCommentResponse); i {
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
			RawDescriptor: file_protos_comments_v1_comments_v1_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_comments_v1_comments_v1_proto_goTypes,
		DependencyIndexes: file_protos_comments_v1_comments_v1_proto_depIdxs,
		MessageInfos:      file_protos_comments_v1_comments_v1_proto_msgTypes,
	}.Build()
	File_protos_comments_v1_comments_v1_proto = out.File
	file_protos_comments_v1_comments_v1_proto_rawDesc = nil
	file_protos_comments_v1_comments_v1_proto_goTypes = nil
	file_protos_comments_v1_comments_v1_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.21.12
// source: auth.proto

package pb

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

type ValidateCodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code"`
}

func (x *ValidateCodeRequest) Reset() {
	*x = ValidateCodeRequest{}
	mi := &file_auth_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateCodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateCodeRequest) ProtoMessage() {}

func (x *ValidateCodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateCodeRequest.ProtoReflect.Descriptor instead.
func (*ValidateCodeRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{0}
}

func (x *ValidateCodeRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type ValidateCodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User  *User  `protobuf:"bytes,1,opt,name=user,proto3" json:"user"`
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token"`
}

func (x *ValidateCodeResponse) Reset() {
	*x = ValidateCodeResponse{}
	mi := &file_auth_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateCodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateCodeResponse) ProtoMessage() {}

func (x *ValidateCodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateCodeResponse.ProtoReflect.Descriptor instead.
func (*ValidateCodeResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{1}
}

func (x *ValidateCodeResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *ValidateCodeResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type ValidateTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token         string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token"`
	RequiredRoles []string `protobuf:"bytes,2,rep,name=requiredRoles,proto3" json:"requiredRoles"`
}

func (x *ValidateTokenRequest) Reset() {
	*x = ValidateTokenRequest{}
	mi := &file_auth_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateTokenRequest) ProtoMessage() {}

func (x *ValidateTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateTokenRequest.ProtoReflect.Descriptor instead.
func (*ValidateTokenRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{2}
}

func (x *ValidateTokenRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ValidateTokenRequest) GetRequiredRoles() []string {
	if x != nil {
		return x.RequiredRoles
	}
	return nil
}

type GetBonusInformationByChatIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string                     `protobuf:"bytes,1,opt,name=name,proto3" json:"name"`
	RefLink      string                     `protobuf:"bytes,2,opt,name=refLink,proto3" json:"refLink"`
	WelcomeCount int32                      `protobuf:"varint,3,opt,name=welcomeCount,proto3" json:"welcomeCount"`
	BonusCount   int32                      `protobuf:"varint,4,opt,name=bonusCount,proto3" json:"bonusCount"`
	RegisteredAt string                     `protobuf:"bytes,5,opt,name=registeredAt,proto3" json:"registeredAt"`
	More         []*GetMoreBonusInformation `protobuf:"bytes,6,rep,name=more,proto3" json:"more"`
}

func (x *GetBonusInformationByChatIdResponse) Reset() {
	*x = GetBonusInformationByChatIdResponse{}
	mi := &file_auth_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBonusInformationByChatIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBonusInformationByChatIdResponse) ProtoMessage() {}

func (x *GetBonusInformationByChatIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBonusInformationByChatIdResponse.ProtoReflect.Descriptor instead.
func (*GetBonusInformationByChatIdResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{3}
}

func (x *GetBonusInformationByChatIdResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetBonusInformationByChatIdResponse) GetRefLink() string {
	if x != nil {
		return x.RefLink
	}
	return ""
}

func (x *GetBonusInformationByChatIdResponse) GetWelcomeCount() int32 {
	if x != nil {
		return x.WelcomeCount
	}
	return 0
}

func (x *GetBonusInformationByChatIdResponse) GetBonusCount() int32 {
	if x != nil {
		return x.BonusCount
	}
	return 0
}

func (x *GetBonusInformationByChatIdResponse) GetRegisteredAt() string {
	if x != nil {
		return x.RegisteredAt
	}
	return ""
}

func (x *GetBonusInformationByChatIdResponse) GetMore() []*GetMoreBonusInformation {
	if x != nil {
		return x.More
	}
	return nil
}

type GetMoreBonusInformation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GuestName         string `protobuf:"bytes,1,opt,name=guestName,proto3" json:"guestName"`
	GuestChatId       string `protobuf:"bytes,2,opt,name=guestChatId,proto3" json:"guestChatId"`
	GuestRegisteredAt string `protobuf:"bytes,3,opt,name=guestRegisteredAt,proto3" json:"guestRegisteredAt"`
}

func (x *GetMoreBonusInformation) Reset() {
	*x = GetMoreBonusInformation{}
	mi := &file_auth_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMoreBonusInformation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMoreBonusInformation) ProtoMessage() {}

func (x *GetMoreBonusInformation) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMoreBonusInformation.ProtoReflect.Descriptor instead.
func (*GetMoreBonusInformation) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{4}
}

func (x *GetMoreBonusInformation) GetGuestName() string {
	if x != nil {
		return x.GuestName
	}
	return ""
}

func (x *GetMoreBonusInformation) GetGuestChatId() string {
	if x != nil {
		return x.GuestChatId
	}
	return ""
}

func (x *GetMoreBonusInformation) GetGuestRegisteredAt() string {
	if x != nil {
		return x.GuestRegisteredAt
	}
	return ""
}

type BonusServiceAbsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChatId string `protobuf:"bytes,1,opt,name=chatId,proto3" json:"chatId"`
}

func (x *BonusServiceAbsRequest) Reset() {
	*x = BonusServiceAbsRequest{}
	mi := &file_auth_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BonusServiceAbsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BonusServiceAbsRequest) ProtoMessage() {}

func (x *BonusServiceAbsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BonusServiceAbsRequest.ProtoReflect.Descriptor instead.
func (*BonusServiceAbsRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{5}
}

func (x *BonusServiceAbsRequest) GetChatId() string {
	if x != nil {
		return x.ChatId
	}
	return ""
}

type CalculateBonusByChatIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int32 `protobuf:"varint,1,opt,name=count,proto3" json:"count"`
}

func (x *CalculateBonusByChatIdResponse) Reset() {
	*x = CalculateBonusByChatIdResponse{}
	mi := &file_auth_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CalculateBonusByChatIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CalculateBonusByChatIdResponse) ProtoMessage() {}

func (x *CalculateBonusByChatIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CalculateBonusByChatIdResponse.ProtoReflect.Descriptor instead.
func (*CalculateBonusByChatIdResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{6}
}

func (x *CalculateBonusByChatIdResponse) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type UseBonusAttemptRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChatId string `protobuf:"bytes,1,opt,name=chatId,proto3" json:"chatId"`
}

func (x *UseBonusAttemptRequest) Reset() {
	*x = UseBonusAttemptRequest{}
	mi := &file_auth_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UseBonusAttemptRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UseBonusAttemptRequest) ProtoMessage() {}

func (x *UseBonusAttemptRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UseBonusAttemptRequest.ProtoReflect.Descriptor instead.
func (*UseBonusAttemptRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{7}
}

func (x *UseBonusAttemptRequest) GetChatId() string {
	if x != nil {
		return x.ChatId
	}
	return ""
}

type UseBonusAttemptResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response bool `protobuf:"varint,1,opt,name=response,proto3" json:"response"`
}

func (x *UseBonusAttemptResponse) Reset() {
	*x = UseBonusAttemptResponse{}
	mi := &file_auth_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UseBonusAttemptResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UseBonusAttemptResponse) ProtoMessage() {}

func (x *UseBonusAttemptResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UseBonusAttemptResponse.ProtoReflect.Descriptor instead.
func (*UseBonusAttemptResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{8}
}

func (x *UseBonusAttemptResponse) GetResponse() bool {
	if x != nil {
		return x.Response
	}
	return false
}

var File_auth_proto protoreflect.FileDescriptor

var file_auth_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x61, 0x75,
	0x74, 0x68, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x29, 0x0a, 0x13, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x4e, 0x0a, 0x14, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x52, 0x0a, 0x14, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65, 0x71,
	0x75, 0x69, 0x72, 0x65, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x22,
	0xee, 0x01, 0x0a, 0x23, 0x47, 0x65, 0x74, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x49, 0x6e, 0x66, 0x6f,
	0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x79, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72,
	0x65, 0x66, 0x4c, 0x69, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65,
	0x66, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x22, 0x0a, 0x0c, 0x77, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x77, 0x65, 0x6c,
	0x63, 0x6f, 0x6d, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x62, 0x6f, 0x6e,
	0x75, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x62,
	0x6f, 0x6e, 0x75, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x41, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x41, 0x74, 0x12, 0x31, 0x0a,
	0x04, 0x6d, 0x6f, 0x72, 0x65, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x72, 0x65, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x49,
	0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x6d, 0x6f, 0x72, 0x65,
	0x22, 0x87, 0x01, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x72, 0x65, 0x42, 0x6f, 0x6e, 0x75,
	0x73, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09,
	0x67, 0x75, 0x65, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x67, 0x75, 0x65, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x67, 0x75,
	0x65, 0x73, 0x74, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x67, 0x75, 0x65, 0x73, 0x74, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x11,
	0x67, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x41,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x67, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x41, 0x74, 0x22, 0x30, 0x0a, 0x16, 0x42, 0x6f,
	0x6e, 0x75, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x41, 0x62, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x22, 0x36, 0x0a, 0x1e,
	0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x42, 0x79,
	0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x22, 0x30, 0x0a, 0x16, 0x55, 0x73, 0x65, 0x42, 0x6f, 0x6e, 0x75, 0x73,
	0x41, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x22, 0x35, 0x0a, 0x17, 0x55, 0x73, 0x65, 0x42, 0x6f, 0x6e,
	0x75, 0x73, 0x41, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x8f, 0x01,
	0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a,
	0x0c, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x19, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x0d, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x56, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0c, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x32,
	0xa4, 0x02, 0x0a, 0x0c, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x4e, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x41, 0x74, 0x74, 0x65,
	0x6d, 0x70, 0x74, 0x12, 0x1c, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x55, 0x73, 0x65, 0x42, 0x6f,
	0x6e, 0x75, 0x73, 0x41, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x55, 0x73, 0x65, 0x42, 0x6f, 0x6e, 0x75,
	0x73, 0x41, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x5c, 0x0a, 0x16, 0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x6e,
	0x75, 0x73, 0x42, 0x79, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x12, 0x1c, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x41, 0x62,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e,
	0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x42, 0x79,
	0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x66,
	0x0a, 0x1b, 0x47, 0x65, 0x74, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x79, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x12, 0x1c, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x2e, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x41, 0x62, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x72,
	0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x79, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_proto_rawDescOnce sync.Once
	file_auth_proto_rawDescData = file_auth_proto_rawDesc
)

func file_auth_proto_rawDescGZIP() []byte {
	file_auth_proto_rawDescOnce.Do(func() {
		file_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_proto_rawDescData)
	})
	return file_auth_proto_rawDescData
}

var file_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_auth_proto_goTypes = []any{
	(*ValidateCodeRequest)(nil),                 // 0: auth.ValidateCodeRequest
	(*ValidateCodeResponse)(nil),                // 1: auth.ValidateCodeResponse
	(*ValidateTokenRequest)(nil),                // 2: auth.ValidateTokenRequest
	(*GetBonusInformationByChatIdResponse)(nil), // 3: auth.GetBonusInformationByChatIdResponse
	(*GetMoreBonusInformation)(nil),             // 4: auth.GetMoreBonusInformation
	(*BonusServiceAbsRequest)(nil),              // 5: auth.BonusServiceAbsRequest
	(*CalculateBonusByChatIdResponse)(nil),      // 6: auth.CalculateBonusByChatIdResponse
	(*UseBonusAttemptRequest)(nil),              // 7: auth.UseBonusAttemptRequest
	(*UseBonusAttemptResponse)(nil),             // 8: auth.UseBonusAttemptResponse
	(*User)(nil),                                // 9: common.User
}
var file_auth_proto_depIdxs = []int32{
	9, // 0: auth.ValidateCodeResponse.user:type_name -> common.User
	4, // 1: auth.GetBonusInformationByChatIdResponse.more:type_name -> auth.GetMoreBonusInformation
	0, // 2: auth.AuthService.ValidateCode:input_type -> auth.ValidateCodeRequest
	2, // 3: auth.AuthService.ValidateToken:input_type -> auth.ValidateTokenRequest
	7, // 4: auth.BonusService.UseBonusAttempt:input_type -> auth.UseBonusAttemptRequest
	5, // 5: auth.BonusService.CalculateBonusByChatId:input_type -> auth.BonusServiceAbsRequest
	5, // 6: auth.BonusService.GetBonusInformationByChatId:input_type -> auth.BonusServiceAbsRequest
	1, // 7: auth.AuthService.ValidateCode:output_type -> auth.ValidateCodeResponse
	9, // 8: auth.AuthService.ValidateToken:output_type -> common.User
	8, // 9: auth.BonusService.UseBonusAttempt:output_type -> auth.UseBonusAttemptResponse
	6, // 10: auth.BonusService.CalculateBonusByChatId:output_type -> auth.CalculateBonusByChatIdResponse
	3, // 11: auth.BonusService.GetBonusInformationByChatId:output_type -> auth.GetBonusInformationByChatIdResponse
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_auth_proto_init() }
func file_auth_proto_init() {
	if File_auth_proto != nil {
		return
	}
	file_common_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_auth_proto_goTypes,
		DependencyIndexes: file_auth_proto_depIdxs,
		MessageInfos:      file_auth_proto_msgTypes,
	}.Build()
	File_auth_proto = out.File
	file_auth_proto_rawDesc = nil
	file_auth_proto_goTypes = nil
	file_auth_proto_depIdxs = nil
}

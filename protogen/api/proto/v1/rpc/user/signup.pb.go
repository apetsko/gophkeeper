// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.0
// source: api/proto/v1/rpc/user/signup.proto

package user

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SignupRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SignupRequest) Reset() {
	*x = SignupRequest{}
	mi := &file_api_proto_v1_rpc_user_signup_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignupRequest) ProtoMessage() {}

func (x *SignupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1_rpc_user_signup_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignupRequest.ProtoReflect.Descriptor instead.
func (*SignupRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_v1_rpc_user_signup_proto_rawDescGZIP(), []int{0}
}

func (x *SignupRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *SignupRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type SignupResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Token         string                 `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SignupResponse) Reset() {
	*x = SignupResponse{}
	mi := &file_api_proto_v1_rpc_user_signup_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignupResponse) ProtoMessage() {}

func (x *SignupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_v1_rpc_user_signup_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignupResponse.ProtoReflect.Descriptor instead.
func (*SignupResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_v1_rpc_user_signup_proto_rawDescGZIP(), []int{1}
}

func (x *SignupResponse) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SignupResponse) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *SignupResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

var File_api_proto_v1_rpc_user_signup_proto protoreflect.FileDescriptor

const file_api_proto_v1_rpc_user_signup_proto_rawDesc = "" +
	"\n" +
	"\"api/proto/v1/rpc/user/signup.proto\x12\x15api.proto.v1.rpc.user\"G\n" +
	"\rSignupRequest\x12\x1a\n" +
	"\busername\x18\x01 \x01(\tR\busername\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"R\n" +
	"\x0eSignupResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\x12\x1a\n" +
	"\busername\x18\x02 \x01(\tR\busername\x12\x14\n" +
	"\x05token\x18\x03 \x01(\tR\x05tokenB>Z<github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc/userb\x06proto3"

var (
	file_api_proto_v1_rpc_user_signup_proto_rawDescOnce sync.Once
	file_api_proto_v1_rpc_user_signup_proto_rawDescData []byte
)

func file_api_proto_v1_rpc_user_signup_proto_rawDescGZIP() []byte {
	file_api_proto_v1_rpc_user_signup_proto_rawDescOnce.Do(func() {
		file_api_proto_v1_rpc_user_signup_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_proto_v1_rpc_user_signup_proto_rawDesc), len(file_api_proto_v1_rpc_user_signup_proto_rawDesc)))
	})
	return file_api_proto_v1_rpc_user_signup_proto_rawDescData
}

var file_api_proto_v1_rpc_user_signup_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_proto_v1_rpc_user_signup_proto_goTypes = []any{
	(*SignupRequest)(nil),  // 0: api.proto.v1.rpc.user.SignupRequest
	(*SignupResponse)(nil), // 1: api.proto.v1.rpc.user.SignupResponse
}
var file_api_proto_v1_rpc_user_signup_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_v1_rpc_user_signup_proto_init() }
func file_api_proto_v1_rpc_user_signup_proto_init() {
	if File_api_proto_v1_rpc_user_signup_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_proto_v1_rpc_user_signup_proto_rawDesc), len(file_api_proto_v1_rpc_user_signup_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_v1_rpc_user_signup_proto_goTypes,
		DependencyIndexes: file_api_proto_v1_rpc_user_signup_proto_depIdxs,
		MessageInfos:      file_api_proto_v1_rpc_user_signup_proto_msgTypes,
	}.Build()
	File_api_proto_v1_rpc_user_signup_proto = out.File
	file_api_proto_v1_rpc_user_signup_proto_goTypes = nil
	file_api_proto_v1_rpc_user_signup_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc.proto

package shard

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type User struct {
	UserId               int64    `protobuf:"varint,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Birthday             string   `protobuf:"bytes,3,opt,name=Birthday,proto3" json:"Birthday,omitempty"`
	Gender               string   `protobuf:"bytes,4,opt,name=Gender,proto3" json:"Gender,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_rpc_ed8f5045742dc650, []int{0}
}
func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (dst *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(dst, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *User) GetBirthday() string {
	if m != nil {
		return m.Birthday
	}
	return ""
}

func (m *User) GetGender() string {
	if m != nil {
		return m.Gender
	}
	return ""
}

type UserReq struct {
	ReqUser              int64    `protobuf:"varint,1,opt,name=ReqUser,proto3" json:"ReqUser,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserReq) Reset()         { *m = UserReq{} }
func (m *UserReq) String() string { return proto.CompactTextString(m) }
func (*UserReq) ProtoMessage()    {}
func (*UserReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_rpc_ed8f5045742dc650, []int{1}
}
func (m *UserReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserReq.Unmarshal(m, b)
}
func (m *UserReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserReq.Marshal(b, m, deterministic)
}
func (dst *UserReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserReq.Merge(dst, src)
}
func (m *UserReq) XXX_Size() int {
	return xxx_messageInfo_UserReq.Size(m)
}
func (m *UserReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UserReq.DiscardUnknown(m)
}

var xxx_messageInfo_UserReq proto.InternalMessageInfo

func (m *UserReq) GetReqUser() int64 {
	if m != nil {
		return m.ReqUser
	}
	return 0
}

type TokenReq struct {
	Token                string   `protobuf:"bytes,1,opt,name=Token,proto3" json:"Token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenReq) Reset()         { *m = TokenReq{} }
func (m *TokenReq) String() string { return proto.CompactTextString(m) }
func (*TokenReq) ProtoMessage()    {}
func (*TokenReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_rpc_ed8f5045742dc650, []int{2}
}
func (m *TokenReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenReq.Unmarshal(m, b)
}
func (m *TokenReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenReq.Marshal(b, m, deterministic)
}
func (dst *TokenReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenReq.Merge(dst, src)
}
func (m *TokenReq) XXX_Size() int {
	return xxx_messageInfo_TokenReq.Size(m)
}
func (m *TokenReq) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenReq.DiscardUnknown(m)
}

var xxx_messageInfo_TokenReq proto.InternalMessageInfo

func (m *TokenReq) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type UserRsp struct {
	OtherUsers           []*User  `protobuf:"bytes,1,rep,name=OtherUsers,proto3" json:"OtherUsers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserRsp) Reset()         { *m = UserRsp{} }
func (m *UserRsp) String() string { return proto.CompactTextString(m) }
func (*UserRsp) ProtoMessage()    {}
func (*UserRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_rpc_ed8f5045742dc650, []int{3}
}
func (m *UserRsp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserRsp.Unmarshal(m, b)
}
func (m *UserRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserRsp.Marshal(b, m, deterministic)
}
func (dst *UserRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserRsp.Merge(dst, src)
}
func (m *UserRsp) XXX_Size() int {
	return xxx_messageInfo_UserRsp.Size(m)
}
func (m *UserRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_UserRsp.DiscardUnknown(m)
}

var xxx_messageInfo_UserRsp proto.InternalMessageInfo

func (m *UserRsp) GetOtherUsers() []*User {
	if m != nil {
		return m.OtherUsers
	}
	return nil
}

type UserPairRequest struct {
	Uid1                 int64    `protobuf:"varint,1,opt,name=Uid1,proto3" json:"Uid1,omitempty"`
	Uid2                 int64    `protobuf:"varint,2,opt,name=Uid2,proto3" json:"Uid2,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserPairRequest) Reset()         { *m = UserPairRequest{} }
func (m *UserPairRequest) String() string { return proto.CompactTextString(m) }
func (*UserPairRequest) ProtoMessage()    {}
func (*UserPairRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_rpc_ed8f5045742dc650, []int{4}
}
func (m *UserPairRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserPairRequest.Unmarshal(m, b)
}
func (m *UserPairRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserPairRequest.Marshal(b, m, deterministic)
}
func (dst *UserPairRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserPairRequest.Merge(dst, src)
}
func (m *UserPairRequest) XXX_Size() int {
	return xxx_messageInfo_UserPairRequest.Size(m)
}
func (m *UserPairRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UserPairRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UserPairRequest proto.InternalMessageInfo

func (m *UserPairRequest) GetUid1() int64 {
	if m != nil {
		return m.Uid1
	}
	return 0
}

func (m *UserPairRequest) GetUid2() int64 {
	if m != nil {
		return m.Uid2
	}
	return 0
}

type IsFiendResp struct {
	IsFriend             bool     `protobuf:"varint,1,opt,name=IsFriend,proto3" json:"IsFriend,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IsFiendResp) Reset()         { *m = IsFiendResp{} }
func (m *IsFiendResp) String() string { return proto.CompactTextString(m) }
func (*IsFiendResp) ProtoMessage()    {}
func (*IsFiendResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_rpc_ed8f5045742dc650, []int{5}
}
func (m *IsFiendResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IsFiendResp.Unmarshal(m, b)
}
func (m *IsFiendResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IsFiendResp.Marshal(b, m, deterministic)
}
func (dst *IsFiendResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IsFiendResp.Merge(dst, src)
}
func (m *IsFiendResp) XXX_Size() int {
	return xxx_messageInfo_IsFiendResp.Size(m)
}
func (m *IsFiendResp) XXX_DiscardUnknown() {
	xxx_messageInfo_IsFiendResp.DiscardUnknown(m)
}

var xxx_messageInfo_IsFiendResp proto.InternalMessageInfo

func (m *IsFiendResp) GetIsFriend() bool {
	if m != nil {
		return m.IsFriend
	}
	return false
}

func init() {
	proto.RegisterType((*User)(nil), "shard.User")
	proto.RegisterType((*UserReq)(nil), "shard.UserReq")
	proto.RegisterType((*TokenReq)(nil), "shard.TokenReq")
	proto.RegisterType((*UserRsp)(nil), "shard.UserRsp")
	proto.RegisterType((*UserPairRequest)(nil), "shard.UserPairRequest")
	proto.RegisterType((*IsFiendResp)(nil), "shard.IsFiendResp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UserSvcClient is the client API for UserSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserSvcClient interface {
	UserInfo(ctx context.Context, in *UserReq, opts ...grpc.CallOption) (*User, error)
	UserByToken(ctx context.Context, in *TokenReq, opts ...grpc.CallOption) (*User, error)
}

type userSvcClient struct {
	cc *grpc.ClientConn
}

func NewUserSvcClient(cc *grpc.ClientConn) UserSvcClient {
	return &userSvcClient{cc}
}

func (c *userSvcClient) UserInfo(ctx context.Context, in *UserReq, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/shard.UserSvc/UserInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userSvcClient) UserByToken(ctx context.Context, in *TokenReq, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/shard.UserSvc/UserByToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserSvcServer is the server API for UserSvc service.
type UserSvcServer interface {
	UserInfo(context.Context, *UserReq) (*User, error)
	UserByToken(context.Context, *TokenReq) (*User, error)
}

func RegisterUserSvcServer(s *grpc.Server, srv UserSvcServer) {
	s.RegisterService(&_UserSvc_serviceDesc, srv)
}

func _UserSvc_UserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserSvcServer).UserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shard.UserSvc/UserInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserSvcServer).UserInfo(ctx, req.(*UserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserSvc_UserByToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserSvcServer).UserByToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shard.UserSvc/UserByToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserSvcServer).UserByToken(ctx, req.(*TokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserSvc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "shard.UserSvc",
	HandlerType: (*UserSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserInfo",
			Handler:    _UserSvc_UserInfo_Handler,
		},
		{
			MethodName: "UserByToken",
			Handler:    _UserSvc_UserByToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}

// RelSvcClient is the client API for RelSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RelSvcClient interface {
	IsFriend(ctx context.Context, in *UserPairRequest, opts ...grpc.CallOption) (*IsFiendResp, error)
}

type relSvcClient struct {
	cc *grpc.ClientConn
}

func NewRelSvcClient(cc *grpc.ClientConn) RelSvcClient {
	return &relSvcClient{cc}
}

func (c *relSvcClient) IsFriend(ctx context.Context, in *UserPairRequest, opts ...grpc.CallOption) (*IsFiendResp, error) {
	out := new(IsFiendResp)
	err := c.cc.Invoke(ctx, "/shard.RelSvc/IsFriend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RelSvcServer is the server API for RelSvc service.
type RelSvcServer interface {
	IsFriend(context.Context, *UserPairRequest) (*IsFiendResp, error)
}

func RegisterRelSvcServer(s *grpc.Server, srv RelSvcServer) {
	s.RegisterService(&_RelSvc_serviceDesc, srv)
}

func _RelSvc_IsFriend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserPairRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelSvcServer).IsFriend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shard.RelSvc/IsFriend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelSvcServer).IsFriend(ctx, req.(*UserPairRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RelSvc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "shard.RelSvc",
	HandlerType: (*RelSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsFriend",
			Handler:    _RelSvc_IsFriend_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}

func init() { proto.RegisterFile("rpc.proto", fileDescriptor_rpc_ed8f5045742dc650) }

var fileDescriptor_rpc_ed8f5045742dc650 = []byte{
	// 312 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0x1b, 0xfb, 0x2f, 0x9d, 0x80, 0x85, 0x41, 0x4a, 0xe8, 0x29, 0xac, 0x97, 0x16, 0xb1,
	0x60, 0x04, 0xd1, 0x6b, 0x0f, 0x95, 0x5e, 0x54, 0x56, 0xfb, 0x01, 0x62, 0x33, 0xa5, 0x41, 0x4d,
	0xd2, 0xdd, 0x28, 0xf4, 0xdb, 0xcb, 0x4c, 0xb7, 0x75, 0xf1, 0x94, 0xf9, 0xcd, 0xbc, 0xcc, 0xbe,
	0xb7, 0x2c, 0x0c, 0x4c, 0xbd, 0x9e, 0xd5, 0xa6, 0x6a, 0x2a, 0xec, 0xda, 0x6d, 0x66, 0x72, 0xb5,
	0x81, 0xce, 0xca, 0x92, 0xc1, 0x11, 0xf4, 0xf8, 0xbb, 0xcc, 0xe3, 0x20, 0x09, 0x26, 0x6d, 0xed,
	0x08, 0x11, 0x3a, 0x4f, 0xd9, 0x17, 0xc5, 0x67, 0x49, 0x30, 0x19, 0x68, 0xa9, 0x71, 0x0c, 0xe1,
	0xbc, 0x30, 0xcd, 0x36, 0xcf, 0xf6, 0x71, 0x5b, 0xfa, 0x27, 0xe6, 0x3d, 0x8f, 0x54, 0xe6, 0x64,
	0xe2, 0x8e, 0x4c, 0x1c, 0xa9, 0x4b, 0xe8, 0xf3, 0x46, 0x4d, 0x3b, 0x8c, 0xa1, 0xaf, 0x69, 0xc7,
	0xe4, 0xce, 0x3a, 0xa2, 0x4a, 0x20, 0x7c, 0xab, 0x3e, 0xa8, 0x64, 0xd5, 0x05, 0x74, 0xa5, 0x16,
	0xcd, 0x40, 0x1f, 0x40, 0xdd, 0xb9, 0x35, 0xb6, 0xc6, 0x2b, 0x80, 0xe7, 0x66, 0x4b, 0x86, 0xd9,
	0xc6, 0x41, 0xd2, 0x9e, 0x44, 0x69, 0x34, 0x93, 0x54, 0x33, 0xd1, 0x78, 0x63, 0xf5, 0x00, 0x43,
	0x2e, 0x5e, 0xb2, 0x82, 0x2d, 0x7c, 0x93, 0x6d, 0x38, 0xd9, 0xaa, 0xc8, 0x6f, 0x9c, 0x07, 0xa9,
	0x5d, 0x2f, 0x95, 0xb4, 0x87, 0x5e, 0xaa, 0xa6, 0x10, 0x2d, 0xed, 0xa2, 0xa0, 0x32, 0xd7, 0x64,
	0x6b, 0x0e, 0xbf, 0xb4, 0x0b, 0xc3, 0x2c, 0xbf, 0x86, 0xfa, 0xc4, 0xe9, 0xfa, 0xe0, 0xee, 0xf5,
	0x67, 0x8d, 0x53, 0x08, 0xe5, 0x06, 0xcb, 0x4d, 0x85, 0xe7, 0xbe, 0x2b, 0xda, 0x8d, 0x7d, 0x97,
	0xaa, 0x85, 0xd7, 0x10, 0x71, 0x35, 0xdf, 0x4b, 0x44, 0x1c, 0xba, 0xe9, 0xf1, 0x26, 0xfe, 0xc9,
	0xd3, 0x39, 0xf4, 0x34, 0x7d, 0xf2, 0x19, 0xf7, 0x7f, 0x56, 0x70, 0xe4, 0x89, 0xbc, 0x94, 0x63,
	0x74, 0x7d, 0x2f, 0x82, 0x6a, 0xbd, 0xf7, 0xe4, 0x0d, 0xdc, 0xfe, 0x06, 0x00, 0x00, 0xff, 0xff,
	0x7d, 0x80, 0x14, 0x42, 0x10, 0x02, 0x00, 0x00,
}

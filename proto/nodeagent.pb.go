// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nodeagent.proto

package nodeagent

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GetDirectorySizeRequest struct {
	Path                 string   `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetDirectorySizeRequest) Reset()         { *m = GetDirectorySizeRequest{} }
func (m *GetDirectorySizeRequest) String() string { return proto.CompactTextString(m) }
func (*GetDirectorySizeRequest) ProtoMessage()    {}
func (*GetDirectorySizeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b5331ae4b8115762, []int{0}
}

func (m *GetDirectorySizeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDirectorySizeRequest.Unmarshal(m, b)
}
func (m *GetDirectorySizeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDirectorySizeRequest.Marshal(b, m, deterministic)
}
func (m *GetDirectorySizeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDirectorySizeRequest.Merge(m, src)
}
func (m *GetDirectorySizeRequest) XXX_Size() int {
	return xxx_messageInfo_GetDirectorySizeRequest.Size(m)
}
func (m *GetDirectorySizeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDirectorySizeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetDirectorySizeRequest proto.InternalMessageInfo

func (m *GetDirectorySizeRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

type GetDirectorySizeReply struct {
	Size                 int64    `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetDirectorySizeReply) Reset()         { *m = GetDirectorySizeReply{} }
func (m *GetDirectorySizeReply) String() string { return proto.CompactTextString(m) }
func (*GetDirectorySizeReply) ProtoMessage()    {}
func (*GetDirectorySizeReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_b5331ae4b8115762, []int{1}
}

func (m *GetDirectorySizeReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDirectorySizeReply.Unmarshal(m, b)
}
func (m *GetDirectorySizeReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDirectorySizeReply.Marshal(b, m, deterministic)
}
func (m *GetDirectorySizeReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDirectorySizeReply.Merge(m, src)
}
func (m *GetDirectorySizeReply) XXX_Size() int {
	return xxx_messageInfo_GetDirectorySizeReply.Size(m)
}
func (m *GetDirectorySizeReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDirectorySizeReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetDirectorySizeReply proto.InternalMessageInfo

func (m *GetDirectorySizeReply) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func init() {
	proto.RegisterType((*GetDirectorySizeRequest)(nil), "nodeagent.GetDirectorySizeRequest")
	proto.RegisterType((*GetDirectorySizeReply)(nil), "nodeagent.GetDirectorySizeReply")
}

func init() { proto.RegisterFile("nodeagent.proto", fileDescriptor_b5331ae4b8115762) }

var fileDescriptor_b5331ae4b8115762 = []byte{
	// 149 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0xcb, 0x4f, 0x49,
	0x4d, 0x4c, 0x4f, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x84, 0x0b, 0x28,
	0xe9, 0x72, 0x89, 0xbb, 0xa7, 0x96, 0xb8, 0x64, 0x16, 0xa5, 0x26, 0x97, 0xe4, 0x17, 0x55, 0x06,
	0x67, 0x56, 0xa5, 0x06, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x09, 0x71, 0xb1, 0x14, 0x24,
	0x96, 0x64, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x4a, 0xda, 0x5c, 0xa2, 0x98,
	0xca, 0x0b, 0x72, 0x2a, 0x41, 0x8a, 0x8b, 0x33, 0xab, 0x52, 0xc1, 0x8a, 0x99, 0x83, 0xc0, 0x6c,
	0xa3, 0x74, 0x2e, 0x4e, 0xbf, 0xfc, 0x94, 0x54, 0x47, 0x90, 0x45, 0x42, 0x51, 0x5c, 0x02, 0xe8,
	0x3a, 0x85, 0x94, 0xf4, 0x10, 0x2e, 0xc3, 0xe1, 0x0a, 0x29, 0x05, 0xbc, 0x6a, 0x0a, 0x72, 0x2a,
	0x95, 0x18, 0x92, 0xd8, 0xc0, 0xde, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x21, 0x07, 0xef,
	0xf8, 0xe9, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// NodeAgentClient is the client API for NodeAgent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NodeAgentClient interface {
	GetDirectorySize(ctx context.Context, in *GetDirectorySizeRequest, opts ...grpc.CallOption) (*GetDirectorySizeReply, error)
}

type nodeAgentClient struct {
	cc *grpc.ClientConn
}

func NewNodeAgentClient(cc *grpc.ClientConn) NodeAgentClient {
	return &nodeAgentClient{cc}
}

func (c *nodeAgentClient) GetDirectorySize(ctx context.Context, in *GetDirectorySizeRequest, opts ...grpc.CallOption) (*GetDirectorySizeReply, error) {
	out := new(GetDirectorySizeReply)
	err := c.cc.Invoke(ctx, "/nodeagent.NodeAgent/GetDirectorySize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeAgentServer is the server API for NodeAgent service.
type NodeAgentServer interface {
	GetDirectorySize(context.Context, *GetDirectorySizeRequest) (*GetDirectorySizeReply, error)
}

// UnimplementedNodeAgentServer can be embedded to have forward compatible implementations.
type UnimplementedNodeAgentServer struct {
}

func (*UnimplementedNodeAgentServer) GetDirectorySize(ctx context.Context, req *GetDirectorySizeRequest) (*GetDirectorySizeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDirectorySize not implemented")
}

func RegisterNodeAgentServer(s *grpc.Server, srv NodeAgentServer) {
	s.RegisterService(&_NodeAgent_serviceDesc, srv)
}

func _NodeAgent_GetDirectorySize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDirectorySizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeAgentServer).GetDirectorySize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeagent.NodeAgent/GetDirectorySize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeAgentServer).GetDirectorySize(ctx, req.(*GetDirectorySizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NodeAgent_serviceDesc = grpc.ServiceDesc{
	ServiceName: "nodeagent.NodeAgent",
	HandlerType: (*NodeAgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDirectorySize",
			Handler:    _NodeAgent_GetDirectorySize_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nodeagent.proto",
}

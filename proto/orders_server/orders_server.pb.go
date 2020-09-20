// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/orders_server/orders_server.proto

package orders_server

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

type Order struct {
	Email                string       `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	RestEmail            string       `protobuf:"bytes,2,opt,name=restEmail,proto3" json:"restEmail,omitempty"`
	RestPhone            string       `protobuf:"bytes,3,opt,name=restPhone,proto3" json:"restPhone,omitempty"`
	OrderItem            []*OrderItem `protobuf:"bytes,4,rep,name=orderItem,proto3" json:"orderItem,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Order) Reset()         { *m = Order{} }
func (m *Order) String() string { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()    {}
func (*Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_734b15a465776565, []int{0}
}

func (m *Order) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Order.Unmarshal(m, b)
}
func (m *Order) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Order.Marshal(b, m, deterministic)
}
func (m *Order) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Order.Merge(m, src)
}
func (m *Order) XXX_Size() int {
	return xxx_messageInfo_Order.Size(m)
}
func (m *Order) XXX_DiscardUnknown() {
	xxx_messageInfo_Order.DiscardUnknown(m)
}

var xxx_messageInfo_Order proto.InternalMessageInfo

func (m *Order) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Order) GetRestEmail() string {
	if m != nil {
		return m.RestEmail
	}
	return ""
}

func (m *Order) GetRestPhone() string {
	if m != nil {
		return m.RestPhone
	}
	return ""
}

func (m *Order) GetOrderItem() []*OrderItem {
	if m != nil {
		return m.OrderItem
	}
	return nil
}

type OrderItem struct {
	MenuId               int64    `protobuf:"varint,1,opt,name=menuId,proto3" json:"menuId,omitempty"`
	Quantity             int64    `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderItem) Reset()         { *m = OrderItem{} }
func (m *OrderItem) String() string { return proto.CompactTextString(m) }
func (*OrderItem) ProtoMessage()    {}
func (*OrderItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_734b15a465776565, []int{1}
}

func (m *OrderItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderItem.Unmarshal(m, b)
}
func (m *OrderItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderItem.Marshal(b, m, deterministic)
}
func (m *OrderItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderItem.Merge(m, src)
}
func (m *OrderItem) XXX_Size() int {
	return xxx_messageInfo_OrderItem.Size(m)
}
func (m *OrderItem) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderItem.DiscardUnknown(m)
}

var xxx_messageInfo_OrderItem proto.InternalMessageInfo

func (m *OrderItem) GetMenuId() int64 {
	if m != nil {
		return m.MenuId
	}
	return 0
}

func (m *OrderItem) GetQuantity() int64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

type CreateOrderRequest struct {
	Order                *Order   `protobuf:"bytes,1,opt,name=order,proto3" json:"order,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateOrderRequest) Reset()         { *m = CreateOrderRequest{} }
func (m *CreateOrderRequest) String() string { return proto.CompactTextString(m) }
func (*CreateOrderRequest) ProtoMessage()    {}
func (*CreateOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_734b15a465776565, []int{2}
}

func (m *CreateOrderRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateOrderRequest.Unmarshal(m, b)
}
func (m *CreateOrderRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateOrderRequest.Marshal(b, m, deterministic)
}
func (m *CreateOrderRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateOrderRequest.Merge(m, src)
}
func (m *CreateOrderRequest) XXX_Size() int {
	return xxx_messageInfo_CreateOrderRequest.Size(m)
}
func (m *CreateOrderRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateOrderRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateOrderRequest proto.InternalMessageInfo

func (m *CreateOrderRequest) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

type CreateOrderResponse struct {
	Response             string   `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateOrderResponse) Reset()         { *m = CreateOrderResponse{} }
func (m *CreateOrderResponse) String() string { return proto.CompactTextString(m) }
func (*CreateOrderResponse) ProtoMessage()    {}
func (*CreateOrderResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_734b15a465776565, []int{3}
}

func (m *CreateOrderResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateOrderResponse.Unmarshal(m, b)
}
func (m *CreateOrderResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateOrderResponse.Marshal(b, m, deterministic)
}
func (m *CreateOrderResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateOrderResponse.Merge(m, src)
}
func (m *CreateOrderResponse) XXX_Size() int {
	return xxx_messageInfo_CreateOrderResponse.Size(m)
}
func (m *CreateOrderResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateOrderResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateOrderResponse proto.InternalMessageInfo

func (m *CreateOrderResponse) GetResponse() string {
	if m != nil {
		return m.Response
	}
	return ""
}

func init() {
	proto.RegisterType((*Order)(nil), "orders_server.Order")
	proto.RegisterType((*OrderItem)(nil), "orders_server.OrderItem")
	proto.RegisterType((*CreateOrderRequest)(nil), "orders_server.CreateOrderRequest")
	proto.RegisterType((*CreateOrderResponse)(nil), "orders_server.CreateOrderResponse")
}

func init() {
	proto.RegisterFile("proto/orders_server/orders_server.proto", fileDescriptor_734b15a465776565)
}

var fileDescriptor_734b15a465776565 = []byte{
	// 267 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x51, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x35, 0xc6, 0x04, 0x33, 0xc1, 0xcb, 0xb4, 0x4a, 0x28, 0x1e, 0xea, 0x5e, 0x2c, 0x1e, 0x2a,
	0x46, 0xf0, 0xaa, 0x28, 0x1e, 0x7a, 0xaa, 0xec, 0xc1, 0x83, 0x17, 0x8d, 0x74, 0xc0, 0x82, 0xc9,
	0xb6, 0xb3, 0x1b, 0xc1, 0x9f, 0xe1, 0x3f, 0x96, 0x4e, 0xd2, 0xc4, 0x68, 0xf0, 0xb6, 0xef, 0x63,
	0xdf, 0x7b, 0xec, 0xc2, 0xe9, 0x8a, 0x8d, 0x33, 0xe7, 0x86, 0x17, 0xc4, 0xf6, 0xd9, 0x12, 0x7f,
	0x10, 0x77, 0xd1, 0x54, 0x1c, 0x78, 0xd0, 0x21, 0xd5, 0x97, 0x07, 0xc1, 0x7c, 0xc3, 0xe0, 0x10,
	0x02, 0xca, 0xb3, 0xe5, 0x7b, 0xe2, 0x8d, 0xbd, 0x49, 0xa4, 0x2b, 0x80, 0xc7, 0x10, 0x31, 0x59,
	0x77, 0x2f, 0xca, 0xae, 0x28, 0x2d, 0xb1, 0x55, 0x1f, 0xde, 0x4c, 0x41, 0x89, 0xdf, 0xaa, 0x42,
	0xe0, 0x15, 0x44, 0x52, 0x36, 0x73, 0x94, 0x27, 0x7b, 0x63, 0x7f, 0x12, 0xa7, 0xc9, 0xb4, 0xbb,
	0x69, 0xbe, 0xd5, 0x75, 0x6b, 0x55, 0xd7, 0x10, 0x35, 0x3c, 0x1e, 0x41, 0x98, 0x53, 0x51, 0xce,
	0x16, 0xb2, 0xcb, 0xd7, 0x35, 0xc2, 0x11, 0xec, 0xaf, 0xcb, 0xac, 0x70, 0x4b, 0xf7, 0x29, 0xbb,
	0x7c, 0xdd, 0x60, 0x75, 0x03, 0x78, 0xc7, 0x94, 0x39, 0x92, 0x18, 0x4d, 0xeb, 0x92, 0xac, 0xc3,
	0x33, 0x08, 0xa4, 0x43, 0x82, 0xe2, 0x74, 0xd8, 0x37, 0x45, 0x57, 0x16, 0x75, 0x01, 0x83, 0x4e,
	0x82, 0x5d, 0x99, 0xc2, 0xd2, 0xa6, 0x94, 0xeb, 0x73, 0xfd, 0x4c, 0x0d, 0x4e, 0x5f, 0x20, 0x14,
	0xb3, 0xc5, 0x47, 0x88, 0x7f, 0x5c, 0xc6, 0x93, 0x5f, 0x45, 0x7f, 0xa7, 0x8d, 0xd4, 0x7f, 0x96,
	0x2a, 0x5f, 0xed, 0xdc, 0x1e, 0x3e, 0x0d, 0x7a, 0x7e, 0xf9, 0x35, 0x14, 0xf2, 0xf2, 0x3b, 0x00,
	0x00, 0xff, 0xff, 0x4a, 0x11, 0xa8, 0x74, 0x03, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OrdersClient is the client API for Orders service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrdersClient interface {
	// // Unary API to create orders
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
}

type ordersClient struct {
	cc *grpc.ClientConn
}

func NewOrdersClient(cc *grpc.ClientConn) OrdersClient {
	return &ordersClient{cc}
}

func (c *ordersClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, "/orders_server.Orders/CreateOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrdersServer is the server API for Orders service.
type OrdersServer interface {
	// // Unary API to create orders
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
}

// UnimplementedOrdersServer can be embedded to have forward compatible implementations.
type UnimplementedOrdersServer struct {
}

func (*UnimplementedOrdersServer) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}

func RegisterOrdersServer(s *grpc.Server, srv OrdersServer) {
	s.RegisterService(&_Orders_serviceDesc, srv)
}

func _Orders_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdersServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/orders_server.Orders/CreateOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdersServer).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Orders_serviceDesc = grpc.ServiceDesc{
	ServiceName: "orders_server.Orders",
	HandlerType: (*OrdersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrder",
			Handler:    _Orders_CreateOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/orders_server/orders_server.proto",
}

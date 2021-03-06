// Code generated by protoc-gen-gogo.
// source: dashprovider.proto
// DO NOT EDIT!

package otsimopb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ProviderGetRequest struct {
	Request    *DashboardGetRequest `protobuf:"bytes,1,opt,name=request" json:"request,omitempty"`
	UserGroups []string             `protobuf:"bytes,2,rep,name=user_groups,json=userGroups" json:"user_groups,omitempty"`
}

func (m *ProviderGetRequest) Reset()                    { *m = ProviderGetRequest{} }
func (m *ProviderGetRequest) String() string            { return proto.CompactTextString(m) }
func (*ProviderGetRequest) ProtoMessage()               {}
func (*ProviderGetRequest) Descriptor() ([]byte, []int) { return fileDescriptorDashprovider, []int{0} }

type ProviderItem struct {
	Cacheable bool  `protobuf:"varint,1,opt,name=cacheable,proto3" json:"cacheable,omitempty"`
	Ttl       int64 `protobuf:"varint,2,opt,name=ttl,proto3" json:"ttl,omitempty"`
	Item      *Card `protobuf:"bytes,4,opt,name=item" json:"item,omitempty"`
}

func (m *ProviderItem) Reset()                    { *m = ProviderItem{} }
func (m *ProviderItem) String() string            { return proto.CompactTextString(m) }
func (*ProviderItem) ProtoMessage()               {}
func (*ProviderItem) Descriptor() ([]byte, []int) { return fileDescriptorDashprovider, []int{1} }

type ProviderItems struct {
	// ProfileId
	ProfileId string `protobuf:"bytes,1,opt,name=profile_id,json=profileId,proto3" json:"profile_id,omitempty"`
	// ChildId
	ChildId string `protobuf:"bytes,2,opt,name=child_id,json=childId,proto3" json:"child_id,omitempty"`
	// CreatedAt
	CreatedAt int64 `protobuf:"varint,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Cacheable
	Cacheable bool `protobuf:"varint,4,opt,name=cacheable,proto3" json:"cacheable,omitempty"`
	// TTL is titme to live duration
	Ttl int64 `protobuf:"varint,5,opt,name=ttl,proto3" json:"ttl,omitempty"`
	// Items
	Items []*ProviderItem `protobuf:"bytes,8,rep,name=items" json:"items,omitempty"`
}

func (m *ProviderItems) Reset()                    { *m = ProviderItems{} }
func (m *ProviderItems) String() string            { return proto.CompactTextString(m) }
func (*ProviderItems) ProtoMessage()               {}
func (*ProviderItems) Descriptor() ([]byte, []int) { return fileDescriptorDashprovider, []int{2} }

type ProviderInfoRequest struct {
}

func (m *ProviderInfoRequest) Reset()                    { *m = ProviderInfoRequest{} }
func (m *ProviderInfoRequest) String() string            { return proto.CompactTextString(m) }
func (*ProviderInfoRequest) ProtoMessage()               {}
func (*ProviderInfoRequest) Descriptor() ([]byte, []int) { return fileDescriptorDashprovider, []int{3} }

type ProviderInfo struct {
	ServesFor []string `protobuf:"bytes,1,rep,name=serves_for,json=servesFor" json:"serves_for,omitempty"`
}

func (m *ProviderInfo) Reset()                    { *m = ProviderInfo{} }
func (m *ProviderInfo) String() string            { return proto.CompactTextString(m) }
func (*ProviderInfo) ProtoMessage()               {}
func (*ProviderInfo) Descriptor() ([]byte, []int) { return fileDescriptorDashprovider, []int{4} }

func init() {
	proto.RegisterType((*ProviderGetRequest)(nil), "otsimo.ProviderGetRequest")
	proto.RegisterType((*ProviderItem)(nil), "otsimo.ProviderItem")
	proto.RegisterType((*ProviderItems)(nil), "otsimo.ProviderItems")
	proto.RegisterType((*ProviderInfoRequest)(nil), "otsimo.ProviderInfoRequest")
	proto.RegisterType((*ProviderInfo)(nil), "otsimo.ProviderInfo")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for DashboardProvider service

type DashboardProviderClient interface {
	Info(ctx context.Context, in *ProviderInfoRequest, opts ...grpc.CallOption) (*ProviderInfo, error)
	Get(ctx context.Context, in *ProviderGetRequest, opts ...grpc.CallOption) (*ProviderItems, error)
}

type dashboardProviderClient struct {
	cc *grpc.ClientConn
}

func NewDashboardProviderClient(cc *grpc.ClientConn) DashboardProviderClient {
	return &dashboardProviderClient{cc}
}

func (c *dashboardProviderClient) Info(ctx context.Context, in *ProviderInfoRequest, opts ...grpc.CallOption) (*ProviderInfo, error) {
	out := new(ProviderInfo)
	err := grpc.Invoke(ctx, "/otsimo.DashboardProvider/Info", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardProviderClient) Get(ctx context.Context, in *ProviderGetRequest, opts ...grpc.CallOption) (*ProviderItems, error) {
	out := new(ProviderItems)
	err := grpc.Invoke(ctx, "/otsimo.DashboardProvider/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DashboardProvider service

type DashboardProviderServer interface {
	Info(context.Context, *ProviderInfoRequest) (*ProviderInfo, error)
	Get(context.Context, *ProviderGetRequest) (*ProviderItems, error)
}

func RegisterDashboardProviderServer(s *grpc.Server, srv DashboardProviderServer) {
	s.RegisterService(&_DashboardProvider_serviceDesc, srv)
}

func _DashboardProvider_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProviderInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardProviderServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otsimo.DashboardProvider/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardProviderServer).Info(ctx, req.(*ProviderInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DashboardProvider_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProviderGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardProviderServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/otsimo.DashboardProvider/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardProviderServer).Get(ctx, req.(*ProviderGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DashboardProvider_serviceDesc = grpc.ServiceDesc{
	ServiceName: "otsimo.DashboardProvider",
	HandlerType: (*DashboardProviderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Info",
			Handler:    _DashboardProvider_Info_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _DashboardProvider_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptorDashprovider,
}

func (m *ProviderGetRequest) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *ProviderGetRequest) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Request != nil {
		data[i] = 0xa
		i++
		i = encodeVarintDashprovider(data, i, uint64(m.Request.Size()))
		n1, err := m.Request.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.UserGroups) > 0 {
		for _, s := range m.UserGroups {
			data[i] = 0x12
			i++
			l = len(s)
			for l >= 1<<7 {
				data[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			data[i] = uint8(l)
			i++
			i += copy(data[i:], s)
		}
	}
	return i, nil
}

func (m *ProviderItem) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *ProviderItem) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Cacheable {
		data[i] = 0x8
		i++
		if m.Cacheable {
			data[i] = 1
		} else {
			data[i] = 0
		}
		i++
	}
	if m.Ttl != 0 {
		data[i] = 0x10
		i++
		i = encodeVarintDashprovider(data, i, uint64(m.Ttl))
	}
	if m.Item != nil {
		data[i] = 0x22
		i++
		i = encodeVarintDashprovider(data, i, uint64(m.Item.Size()))
		n2, err := m.Item.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func (m *ProviderItems) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *ProviderItems) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ProfileId) > 0 {
		data[i] = 0xa
		i++
		i = encodeVarintDashprovider(data, i, uint64(len(m.ProfileId)))
		i += copy(data[i:], m.ProfileId)
	}
	if len(m.ChildId) > 0 {
		data[i] = 0x12
		i++
		i = encodeVarintDashprovider(data, i, uint64(len(m.ChildId)))
		i += copy(data[i:], m.ChildId)
	}
	if m.CreatedAt != 0 {
		data[i] = 0x18
		i++
		i = encodeVarintDashprovider(data, i, uint64(m.CreatedAt))
	}
	if m.Cacheable {
		data[i] = 0x20
		i++
		if m.Cacheable {
			data[i] = 1
		} else {
			data[i] = 0
		}
		i++
	}
	if m.Ttl != 0 {
		data[i] = 0x28
		i++
		i = encodeVarintDashprovider(data, i, uint64(m.Ttl))
	}
	if len(m.Items) > 0 {
		for _, msg := range m.Items {
			data[i] = 0x42
			i++
			i = encodeVarintDashprovider(data, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *ProviderInfoRequest) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *ProviderInfoRequest) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *ProviderInfo) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *ProviderInfo) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ServesFor) > 0 {
		for _, s := range m.ServesFor {
			data[i] = 0xa
			i++
			l = len(s)
			for l >= 1<<7 {
				data[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			data[i] = uint8(l)
			i++
			i += copy(data[i:], s)
		}
	}
	return i, nil
}

func encodeFixed64Dashprovider(data []byte, offset int, v uint64) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	data[offset+4] = uint8(v >> 32)
	data[offset+5] = uint8(v >> 40)
	data[offset+6] = uint8(v >> 48)
	data[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Dashprovider(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintDashprovider(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}
func (m *ProviderGetRequest) Size() (n int) {
	var l int
	_ = l
	if m.Request != nil {
		l = m.Request.Size()
		n += 1 + l + sovDashprovider(uint64(l))
	}
	if len(m.UserGroups) > 0 {
		for _, s := range m.UserGroups {
			l = len(s)
			n += 1 + l + sovDashprovider(uint64(l))
		}
	}
	return n
}

func (m *ProviderItem) Size() (n int) {
	var l int
	_ = l
	if m.Cacheable {
		n += 2
	}
	if m.Ttl != 0 {
		n += 1 + sovDashprovider(uint64(m.Ttl))
	}
	if m.Item != nil {
		l = m.Item.Size()
		n += 1 + l + sovDashprovider(uint64(l))
	}
	return n
}

func (m *ProviderItems) Size() (n int) {
	var l int
	_ = l
	l = len(m.ProfileId)
	if l > 0 {
		n += 1 + l + sovDashprovider(uint64(l))
	}
	l = len(m.ChildId)
	if l > 0 {
		n += 1 + l + sovDashprovider(uint64(l))
	}
	if m.CreatedAt != 0 {
		n += 1 + sovDashprovider(uint64(m.CreatedAt))
	}
	if m.Cacheable {
		n += 2
	}
	if m.Ttl != 0 {
		n += 1 + sovDashprovider(uint64(m.Ttl))
	}
	if len(m.Items) > 0 {
		for _, e := range m.Items {
			l = e.Size()
			n += 1 + l + sovDashprovider(uint64(l))
		}
	}
	return n
}

func (m *ProviderInfoRequest) Size() (n int) {
	var l int
	_ = l
	return n
}

func (m *ProviderInfo) Size() (n int) {
	var l int
	_ = l
	if len(m.ServesFor) > 0 {
		for _, s := range m.ServesFor {
			l = len(s)
			n += 1 + l + sovDashprovider(uint64(l))
		}
	}
	return n
}

func sovDashprovider(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozDashprovider(x uint64) (n int) {
	return sovDashprovider(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProviderGetRequest) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDashprovider
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProviderGetRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProviderGetRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Request", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDashprovider
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Request == nil {
				m.Request = &DashboardGetRequest{}
			}
			if err := m.Request.Unmarshal(data[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserGroups", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDashprovider
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UserGroups = append(m.UserGroups, string(data[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDashprovider(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDashprovider
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ProviderItem) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDashprovider
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProviderItem: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProviderItem: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cacheable", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Cacheable = bool(v != 0)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ttl", wireType)
			}
			m.Ttl = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.Ttl |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Item", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDashprovider
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Item == nil {
				m.Item = &Card{}
			}
			if err := m.Item.Unmarshal(data[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDashprovider(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDashprovider
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ProviderItems) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDashprovider
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProviderItems: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProviderItems: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProfileId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDashprovider
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProfileId = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChildId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDashprovider
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChildId = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			m.CreatedAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.CreatedAt |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cacheable", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Cacheable = bool(v != 0)
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ttl", wireType)
			}
			m.Ttl = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.Ttl |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Items", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthDashprovider
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Items = append(m.Items, &ProviderItem{})
			if err := m.Items[len(m.Items)-1].Unmarshal(data[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDashprovider(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDashprovider
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ProviderInfoRequest) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDashprovider
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProviderInfoRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProviderInfoRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipDashprovider(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDashprovider
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ProviderInfo) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDashprovider
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProviderInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProviderInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ServesFor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDashprovider
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ServesFor = append(m.ServesFor, string(data[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDashprovider(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDashprovider
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipDashprovider(data []byte) (n int, err error) {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDashprovider
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if data[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDashprovider
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthDashprovider
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowDashprovider
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := data[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipDashprovider(data[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthDashprovider = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDashprovider   = fmt.Errorf("proto: integer overflow")
)

var fileDescriptorDashprovider = []byte{
	// 443 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x92, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xb3, 0x71, 0xda, 0xc6, 0x93, 0xa2, 0xc2, 0xd2, 0x4a, 0xc6, 0x05, 0x63, 0xf9, 0x14,
	0x21, 0x35, 0x95, 0x82, 0x2a, 0xc1, 0x91, 0x3f, 0x22, 0xca, 0xad, 0xda, 0x23, 0x97, 0x60, 0x7b,
	0xd7, 0x7f, 0x24, 0xbb, 0x63, 0x76, 0xd7, 0x7d, 0x06, 0x8e, 0x3c, 0x14, 0x87, 0x1e, 0xfb, 0x08,
	0x34, 0xbc, 0x08, 0xf2, 0xda, 0xae, 0x03, 0x94, 0xdb, 0xcc, 0xf7, 0xcd, 0xe4, 0xfb, 0x65, 0xbc,
	0x40, 0x79, 0xa8, 0xb2, 0x4a, 0xe2, 0x75, 0xce, 0x85, 0x5c, 0x54, 0x12, 0x35, 0xd2, 0x7d, 0xd4,
	0x2a, 0x2f, 0xd1, 0x3d, 0x6a, 0xbc, 0x08, 0x43, 0xc9, 0x5b, 0xc3, 0x3d, 0x4b, 0x73, 0x9d, 0xd5,
	0xd1, 0x22, 0xc6, 0xf2, 0x3c, 0xc5, 0x14, 0xcf, 0x8d, 0x1c, 0xd5, 0x89, 0xe9, 0x4c, 0x63, 0xaa,
	0x76, 0x3c, 0x28, 0x80, 0x5e, 0x76, 0xbf, 0xbc, 0x12, 0x9a, 0x89, 0xaf, 0xb5, 0x50, 0x9a, 0x5e,
	0xc0, 0x81, 0x6c, 0x4b, 0x87, 0xf8, 0x64, 0x3e, 0x5b, 0x9e, 0x2e, 0xda, 0xbc, 0xc5, 0xc7, 0x3e,
	0x6e, 0x98, 0x66, 0xfd, 0x2c, 0x7d, 0x09, 0xb3, 0x5a, 0x09, 0xb9, 0x49, 0x25, 0xd6, 0x95, 0x72,
	0xc6, 0xbe, 0x35, 0xb7, 0x19, 0x34, 0xd2, 0xca, 0x28, 0xc1, 0x17, 0x38, 0xec, 0xd3, 0xd6, 0x5a,
	0x94, 0xf4, 0x39, 0xd8, 0x71, 0x18, 0x67, 0x22, 0x8c, 0x0a, 0x61, 0x92, 0xa6, 0x6c, 0x10, 0xe8,
	0x63, 0xb0, 0xb4, 0x2e, 0x9c, 0xb1, 0x4f, 0xe6, 0x16, 0x6b, 0x4a, 0xea, 0xc3, 0x24, 0xd7, 0xa2,
	0x74, 0x26, 0x06, 0xea, 0xb0, 0x87, 0xfa, 0x10, 0x4a, 0xce, 0x8c, 0x13, 0xfc, 0x20, 0xf0, 0x68,
	0x37, 0x42, 0xd1, 0x17, 0x00, 0x95, 0xc4, 0x24, 0x2f, 0xc4, 0x26, 0xe7, 0x26, 0xc4, 0x66, 0x76,
	0xa7, 0xac, 0x39, 0x7d, 0x06, 0xd3, 0x38, 0xcb, 0x0b, 0xde, 0x98, 0x63, 0x63, 0x1e, 0x98, 0x7e,
	0xcd, 0x9b, 0xcd, 0x58, 0x8a, 0x50, 0x0b, 0xbe, 0x09, 0xb5, 0x63, 0x19, 0x0c, 0xbb, 0x53, 0xde,
	0xe9, 0x3f, 0xe1, 0x27, 0xff, 0x81, 0xdf, 0x1b, 0xe0, 0x5f, 0xc1, 0x5e, 0x83, 0xa8, 0x9c, 0xa9,
	0x6f, 0xcd, 0x67, 0xcb, 0xe3, 0x9e, 0x7e, 0x17, 0x97, 0xb5, 0x23, 0xc1, 0x09, 0x3c, 0xbd, 0x97,
	0xaf, 0x12, 0xec, 0x2e, 0x1d, 0x9c, 0xed, 0xdc, 0xef, 0x2a, 0xc1, 0x86, 0x50, 0x09, 0x79, 0x2d,
	0xd4, 0x26, 0x41, 0xe9, 0x10, 0x73, 0x6f, 0xbb, 0x55, 0x3e, 0xa1, 0x5c, 0x7e, 0x23, 0xf0, 0xe4,
	0xfe, 0x83, 0xf5, 0x8b, 0xf4, 0x2d, 0x4c, 0xcc, 0xf2, 0xe9, 0x3f, 0x00, 0x43, 0x92, 0x7b, 0xfc,
	0x90, 0x49, 0xdf, 0x80, 0xb5, 0x12, 0x9a, 0xba, 0x7f, 0x9b, 0xc3, 0x63, 0x70, 0x4f, 0x1e, 0xfa,
	0x5b, 0xea, 0xfd, 0xc5, 0xcd, 0x9d, 0x37, 0xba, 0xbd, 0xf3, 0x46, 0x37, 0x5b, 0x8f, 0xdc, 0x6e,
	0x3d, 0xf2, 0x73, 0xeb, 0x91, 0xef, 0xbf, 0xbc, 0x11, 0x1c, 0xc5, 0x58, 0xf6, 0x3b, 0xa9, 0xac,
	0xe2, 0x4b, 0xf2, 0x79, 0xda, 0xb6, 0x55, 0x14, 0xed, 0x9b, 0x57, 0xfa, 0xfa, 0x77, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x08, 0xda, 0xf6, 0x56, 0x03, 0x03, 0x00, 0x00,
}

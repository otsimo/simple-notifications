// Code generated by protoc-gen-gogo.
// source: registry.proto
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

type AllGameReleases struct {
	GameId   string                         `protobuf:"bytes,1,opt,name=game_id,json=gameId,proto3" json:"game_id,omitempty"`
	Releases []*AllGameReleases_MiniRelease `protobuf:"bytes,2,rep,name=releases" json:"releases,omitempty"`
}

func (m *AllGameReleases) Reset()                    { *m = AllGameReleases{} }
func (m *AllGameReleases) String() string            { return proto.CompactTextString(m) }
func (*AllGameReleases) ProtoMessage()               {}
func (*AllGameReleases) Descriptor() ([]byte, []int) { return fileDescriptorRegistry, []int{0} }

type AllGameReleases_MiniRelease struct {
	Version      string       `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	ReleasedAt   int64        `protobuf:"varint,2,opt,name=released_at,json=releasedAt,proto3" json:"released_at,omitempty"`
	ReleaseState ReleaseState `protobuf:"varint,3,opt,name=release_state,json=releaseState,proto3,enum=apipb.ReleaseState" json:"release_state,omitempty"`
}

func (m *AllGameReleases_MiniRelease) Reset()         { *m = AllGameReleases_MiniRelease{} }
func (m *AllGameReleases_MiniRelease) String() string { return proto.CompactTextString(m) }
func (*AllGameReleases_MiniRelease) ProtoMessage()    {}
func (*AllGameReleases_MiniRelease) Descriptor() ([]byte, []int) {
	return fileDescriptorRegistry, []int{0, 0}
}

func init() {
	proto.RegisterType((*AllGameReleases)(nil), "apipb.AllGameReleases")
	proto.RegisterType((*AllGameReleases_MiniRelease)(nil), "apipb.AllGameReleases.MiniRelease")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for RegistryService service

type RegistryServiceClient interface {
	// Get returns game
	Get(ctx context.Context, in *GetGameRequest, opts ...grpc.CallOption) (*Game, error)
	// GetRelease returns GameRelease by given game id and version
	GetRelease(ctx context.Context, in *GetGameReleaseRequest, opts ...grpc.CallOption) (*GameRelease, error)
	// Publish tries to create a new GameRelease by given manifest
	Publish(ctx context.Context, in *GameManifest, opts ...grpc.CallOption) (*PublishResponse, error)
	// ChangeReleaseState changes state of a release, If user is admin than s/he can change
	// from WAITING to REJECTED or VALIDATED, developers can change to any except VALIDATED
	ChangeReleaseState(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*Response, error)
	// GetLatestVersions returns latest versions of given game ids
	GetLatestVersions(ctx context.Context, in *GetLatestVersionsRequest, opts ...grpc.CallOption) (*GameVersionsResponse, error)
	// Search does search
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
	// ListGames returns all games
	ListGames(ctx context.Context, in *ListGamesRequest, opts ...grpc.CallOption) (RegistryService_ListGamesClient, error)
	// AllReleases returns all releases of the given game
	AllReleases(ctx context.Context, in *GetGameRequest, opts ...grpc.CallOption) (*AllGameReleases, error)
}

type registryServiceClient struct {
	cc *grpc.ClientConn
}

func NewRegistryServiceClient(cc *grpc.ClientConn) RegistryServiceClient {
	return &registryServiceClient{cc}
}

func (c *registryServiceClient) Get(ctx context.Context, in *GetGameRequest, opts ...grpc.CallOption) (*Game, error) {
	out := new(Game)
	err := grpc.Invoke(ctx, "/apipb.RegistryService/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryServiceClient) GetRelease(ctx context.Context, in *GetGameReleaseRequest, opts ...grpc.CallOption) (*GameRelease, error) {
	out := new(GameRelease)
	err := grpc.Invoke(ctx, "/apipb.RegistryService/GetRelease", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryServiceClient) Publish(ctx context.Context, in *GameManifest, opts ...grpc.CallOption) (*PublishResponse, error) {
	out := new(PublishResponse)
	err := grpc.Invoke(ctx, "/apipb.RegistryService/Publish", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryServiceClient) ChangeReleaseState(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/apipb.RegistryService/ChangeReleaseState", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryServiceClient) GetLatestVersions(ctx context.Context, in *GetLatestVersionsRequest, opts ...grpc.CallOption) (*GameVersionsResponse, error) {
	out := new(GameVersionsResponse)
	err := grpc.Invoke(ctx, "/apipb.RegistryService/GetLatestVersions", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryServiceClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := grpc.Invoke(ctx, "/apipb.RegistryService/Search", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryServiceClient) ListGames(ctx context.Context, in *ListGamesRequest, opts ...grpc.CallOption) (RegistryService_ListGamesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_RegistryService_serviceDesc.Streams[0], c.cc, "/apipb.RegistryService/ListGames", opts...)
	if err != nil {
		return nil, err
	}
	x := &registryServiceListGamesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RegistryService_ListGamesClient interface {
	Recv() (*ListItem, error)
	grpc.ClientStream
}

type registryServiceListGamesClient struct {
	grpc.ClientStream
}

func (x *registryServiceListGamesClient) Recv() (*ListItem, error) {
	m := new(ListItem)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *registryServiceClient) AllReleases(ctx context.Context, in *GetGameRequest, opts ...grpc.CallOption) (*AllGameReleases, error) {
	out := new(AllGameReleases)
	err := grpc.Invoke(ctx, "/apipb.RegistryService/AllReleases", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RegistryService service

type RegistryServiceServer interface {
	// Get returns game
	Get(context.Context, *GetGameRequest) (*Game, error)
	// GetRelease returns GameRelease by given game id and version
	GetRelease(context.Context, *GetGameReleaseRequest) (*GameRelease, error)
	// Publish tries to create a new GameRelease by given manifest
	Publish(context.Context, *GameManifest) (*PublishResponse, error)
	// ChangeReleaseState changes state of a release, If user is admin than s/he can change
	// from WAITING to REJECTED or VALIDATED, developers can change to any except VALIDATED
	ChangeReleaseState(context.Context, *ValidateRequest) (*Response, error)
	// GetLatestVersions returns latest versions of given game ids
	GetLatestVersions(context.Context, *GetLatestVersionsRequest) (*GameVersionsResponse, error)
	// Search does search
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	// ListGames returns all games
	ListGames(*ListGamesRequest, RegistryService_ListGamesServer) error
	// AllReleases returns all releases of the given game
	AllReleases(context.Context, *GetGameRequest) (*AllGameReleases, error)
}

func RegisterRegistryServiceServer(s *grpc.Server, srv RegistryServiceServer) {
	s.RegisterService(&_RegistryService_serviceDesc, srv)
}

func _RegistryService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apipb.RegistryService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServiceServer).Get(ctx, req.(*GetGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryService_GetRelease_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGameReleaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServiceServer).GetRelease(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apipb.RegistryService/GetRelease",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServiceServer).GetRelease(ctx, req.(*GetGameReleaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GameManifest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apipb.RegistryService/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServiceServer).Publish(ctx, req.(*GameManifest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryService_ChangeReleaseState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServiceServer).ChangeReleaseState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apipb.RegistryService/ChangeReleaseState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServiceServer).ChangeReleaseState(ctx, req.(*ValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryService_GetLatestVersions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLatestVersionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServiceServer).GetLatestVersions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apipb.RegistryService/GetLatestVersions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServiceServer).GetLatestVersions(ctx, req.(*GetLatestVersionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apipb.RegistryService/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServiceServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryService_ListGames_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListGamesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RegistryServiceServer).ListGames(m, &registryServiceListGamesServer{stream})
}

type RegistryService_ListGamesServer interface {
	Send(*ListItem) error
	grpc.ServerStream
}

type registryServiceListGamesServer struct {
	grpc.ServerStream
}

func (x *registryServiceListGamesServer) Send(m *ListItem) error {
	return x.ServerStream.SendMsg(m)
}

func _RegistryService_AllReleases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServiceServer).AllReleases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apipb.RegistryService/AllReleases",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServiceServer).AllReleases(ctx, req.(*GetGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RegistryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apipb.RegistryService",
	HandlerType: (*RegistryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _RegistryService_Get_Handler,
		},
		{
			MethodName: "GetRelease",
			Handler:    _RegistryService_GetRelease_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _RegistryService_Publish_Handler,
		},
		{
			MethodName: "ChangeReleaseState",
			Handler:    _RegistryService_ChangeReleaseState_Handler,
		},
		{
			MethodName: "GetLatestVersions",
			Handler:    _RegistryService_GetLatestVersions_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _RegistryService_Search_Handler,
		},
		{
			MethodName: "AllReleases",
			Handler:    _RegistryService_AllReleases_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListGames",
			Handler:       _RegistryService_ListGames_Handler,
			ServerStreams: true,
		},
	},
	Metadata: fileDescriptorRegistry,
}

func (m *AllGameReleases) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *AllGameReleases) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.GameId) > 0 {
		data[i] = 0xa
		i++
		i = encodeVarintRegistry(data, i, uint64(len(m.GameId)))
		i += copy(data[i:], m.GameId)
	}
	if len(m.Releases) > 0 {
		for _, msg := range m.Releases {
			data[i] = 0x12
			i++
			i = encodeVarintRegistry(data, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *AllGameReleases_MiniRelease) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *AllGameReleases_MiniRelease) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Version) > 0 {
		data[i] = 0xa
		i++
		i = encodeVarintRegistry(data, i, uint64(len(m.Version)))
		i += copy(data[i:], m.Version)
	}
	if m.ReleasedAt != 0 {
		data[i] = 0x10
		i++
		i = encodeVarintRegistry(data, i, uint64(m.ReleasedAt))
	}
	if m.ReleaseState != 0 {
		data[i] = 0x18
		i++
		i = encodeVarintRegistry(data, i, uint64(m.ReleaseState))
	}
	return i, nil
}

func encodeFixed64Registry(data []byte, offset int, v uint64) int {
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
func encodeFixed32Registry(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintRegistry(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}
func (m *AllGameReleases) Size() (n int) {
	var l int
	_ = l
	l = len(m.GameId)
	if l > 0 {
		n += 1 + l + sovRegistry(uint64(l))
	}
	if len(m.Releases) > 0 {
		for _, e := range m.Releases {
			l = e.Size()
			n += 1 + l + sovRegistry(uint64(l))
		}
	}
	return n
}

func (m *AllGameReleases_MiniRelease) Size() (n int) {
	var l int
	_ = l
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovRegistry(uint64(l))
	}
	if m.ReleasedAt != 0 {
		n += 1 + sovRegistry(uint64(m.ReleasedAt))
	}
	if m.ReleaseState != 0 {
		n += 1 + sovRegistry(uint64(m.ReleaseState))
	}
	return n
}

func sovRegistry(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozRegistry(x uint64) (n int) {
	return sovRegistry(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AllGameReleases) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRegistry
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
			return fmt.Errorf("proto: AllGameReleases: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AllGameReleases: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GameId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRegistry
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
				return ErrInvalidLengthRegistry
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GameId = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Releases", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRegistry
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
				return ErrInvalidLengthRegistry
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Releases = append(m.Releases, &AllGameReleases_MiniRelease{})
			if err := m.Releases[len(m.Releases)-1].Unmarshal(data[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRegistry(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRegistry
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
func (m *AllGameReleases_MiniRelease) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRegistry
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
			return fmt.Errorf("proto: MiniRelease: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MiniRelease: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRegistry
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
				return ErrInvalidLengthRegistry
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Version = string(data[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReleasedAt", wireType)
			}
			m.ReleasedAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRegistry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.ReleasedAt |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReleaseState", wireType)
			}
			m.ReleaseState = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRegistry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				m.ReleaseState |= (ReleaseState(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRegistry(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRegistry
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
func skipRegistry(data []byte) (n int, err error) {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRegistry
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
					return 0, ErrIntOverflowRegistry
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
					return 0, ErrIntOverflowRegistry
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
				return 0, ErrInvalidLengthRegistry
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowRegistry
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
				next, err := skipRegistry(data[start:])
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
	ErrInvalidLengthRegistry = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRegistry   = fmt.Errorf("proto: integer overflow")
)

var fileDescriptorRegistry = []byte{
	// 497 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x53, 0x4d, 0x6e, 0xd3, 0x40,
	0x14, 0xee, 0xc4, 0x22, 0x69, 0x5f, 0x4a, 0x22, 0xa6, 0xb4, 0xb5, 0x02, 0x72, 0xa3, 0xac, 0xc2,
	0x02, 0x17, 0x85, 0x45, 0xbb, 0x40, 0x45, 0x81, 0x85, 0x55, 0xa9, 0x15, 0x95, 0x8d, 0xba, 0x60,
	0x13, 0x8d, 0xed, 0x57, 0x67, 0x24, 0xdb, 0x63, 0x3c, 0x93, 0x4a, 0x6c, 0x39, 0x01, 0x67, 0xe0,
	0x34, 0x5d, 0x72, 0x04, 0x08, 0x07, 0xe0, 0x02, 0x2c, 0x90, 0xed, 0xb1, 0x6b, 0x22, 0xd8, 0xe5,
	0xfb, 0x7b, 0x33, 0xdf, 0xcb, 0x18, 0x06, 0x39, 0x46, 0x5c, 0xaa, 0xfc, 0x93, 0x9d, 0xe5, 0x42,
	0x09, 0xfa, 0x80, 0x65, 0x3c, 0xf3, 0x47, 0x83, 0x04, 0xa5, 0x64, 0x11, 0xca, 0x8a, 0x1e, 0xed,
	0x26, 0x22, 0xc4, 0xb8, 0x46, 0xcf, 0x23, 0xae, 0x96, 0x2b, 0xdf, 0x0e, 0x44, 0x72, 0x1c, 0x89,
	0x48, 0x1c, 0x97, 0xb4, 0xbf, 0xba, 0x29, 0x51, 0x09, 0xca, 0x5f, 0x95, 0x7d, 0xf2, 0x8b, 0xc0,
	0x70, 0x1e, 0xc7, 0x0e, 0x4b, 0xd0, 0xc5, 0x18, 0x99, 0x44, 0x49, 0x0f, 0xa1, 0x17, 0xb1, 0x04,
	0x17, 0x3c, 0x34, 0xc9, 0x98, 0x4c, 0x77, 0xdc, 0x6e, 0x01, 0xcf, 0x43, 0x7a, 0x06, 0xdb, 0xb9,
	0x36, 0x99, 0x9d, 0xb1, 0x31, 0xed, 0xcf, 0x26, 0x76, 0x79, 0x27, 0x7b, 0x63, 0x84, 0x7d, 0xc9,
	0x53, 0xae, 0x81, 0xdb, 0x64, 0x46, 0x9f, 0x09, 0xf4, 0x5b, 0x0a, 0x35, 0xa1, 0x77, 0x8b, 0xb9,
	0xe4, 0x22, 0xd5, 0x07, 0xd5, 0x90, 0x1e, 0x41, 0x5f, 0xa7, 0xc2, 0x05, 0x53, 0x66, 0x67, 0x4c,
	0xa6, 0x86, 0x0b, 0x35, 0x35, 0x57, 0xf4, 0x14, 0x1e, 0x6a, 0xb4, 0x90, 0x8a, 0x29, 0x34, 0x8d,
	0x31, 0x99, 0x0e, 0x66, 0x7b, 0xfa, 0x3e, 0xfa, 0x04, 0xaf, 0x90, 0xdc, 0xdd, 0xbc, 0x85, 0x66,
	0xbf, 0x0d, 0x18, 0xba, 0x7a, 0xb1, 0x1e, 0xe6, 0xb7, 0x3c, 0x40, 0xfa, 0x0c, 0x0c, 0x07, 0x15,
	0xdd, 0xd7, 0x69, 0x07, 0x55, 0xd5, 0xe6, 0xe3, 0x0a, 0xa5, 0x1a, 0xf5, 0x6b, 0x9a, 0x25, 0x48,
	0xcf, 0x00, 0x1c, 0x54, 0x75, 0x83, 0xa7, 0x9b, 0x89, 0xaa, 0xb2, 0x0e, 0xd2, 0x56, 0xb0, 0x4e,
	0x9c, 0x42, 0xef, 0x6a, 0xe5, 0xc7, 0x5c, 0x2e, 0xe9, 0x5e, 0x4b, 0xbe, 0x64, 0x29, 0xbf, 0x29,
	0x32, 0x07, 0x9a, 0xd4, 0x26, 0x17, 0x65, 0x26, 0x52, 0x89, 0x93, 0x2d, 0xfa, 0x1a, 0xe8, 0xdb,
	0x25, 0x4b, 0x23, 0x6c, 0x97, 0xa3, 0xb5, 0xff, 0x9a, 0xc5, 0x3c, 0x2c, 0xda, 0xea, 0xb3, 0x87,
	0xcd, 0x26, 0x9a, 0x01, 0x1e, 0x3c, 0x72, 0x50, 0x5d, 0x30, 0x85, 0x52, 0x5d, 0x57, 0x8b, 0x96,
	0xf4, 0xe8, 0xbe, 0xc1, 0xdf, 0x4a, 0x3d, 0xe8, 0x49, 0xeb, 0x96, 0xf7, 0x5a, 0x33, 0xf4, 0x04,
	0xba, 0x1e, 0xb2, 0x3c, 0x58, 0xd2, 0xc7, 0xda, 0x58, 0xc1, 0x3a, 0xbe, 0xbf, 0xc1, 0xb6, 0x82,
	0x3b, 0x17, 0x5c, 0x96, 0x6b, 0x93, 0xf4, 0x50, 0xbb, 0x1a, 0x66, 0xb3, 0x46, 0x21, 0x9c, 0x2b,
	0x4c, 0x5e, 0x10, 0xfa, 0x0a, 0xfa, 0xf3, 0x38, 0x6e, 0x5e, 0xeb, 0x7f, 0xfe, 0xb4, 0x83, 0x7f,
	0xbf, 0xcc, 0x37, 0x27, 0x77, 0x3f, 0xac, 0xad, 0xbb, 0xb5, 0x45, 0xbe, 0xad, 0x2d, 0xf2, 0x7d,
	0x6d, 0x91, 0x2f, 0x3f, 0xad, 0x2d, 0x18, 0x06, 0x22, 0xb1, 0x85, 0x92, 0x3c, 0x11, 0x76, 0x94,
	0x67, 0xc1, 0x15, 0xf9, 0xb0, 0x5d, 0xc1, 0xcc, 0xff, 0xda, 0x31, 0xde, 0xbd, 0xf7, 0xfc, 0x6e,
	0xf9, 0xc1, 0xbc, 0xfc, 0x13, 0x00, 0x00, 0xff, 0xff, 0x76, 0x20, 0xd5, 0x94, 0x96, 0x03, 0x00,
	0x00,
}
// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ruler.proto

package ruler

import (
	context "context"
	fmt "fmt"
	io "io"
	math "math"
	reflect "reflect"
	strings "strings"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"

	rules "github.com/cortexproject/cortex/pkg/ruler/rules"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type RulesRequest struct {
}

func (m *RulesRequest) Reset()      { *m = RulesRequest{} }
func (*RulesRequest) ProtoMessage() {}
func (*RulesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9ecbec0a4cfddea6, []int{0}
}
func (m *RulesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RulesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RulesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RulesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RulesRequest.Merge(m, src)
}
func (m *RulesRequest) XXX_Size() int {
	return m.Size()
}
func (m *RulesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RulesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RulesRequest proto.InternalMessageInfo

type RulesResponse struct {
	Groups []*rules.RuleGroupDesc `protobuf:"bytes,1,rep,name=groups,proto3" json:"groups,omitempty"`
}

func (m *RulesResponse) Reset()      { *m = RulesResponse{} }
func (*RulesResponse) ProtoMessage() {}
func (*RulesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9ecbec0a4cfddea6, []int{1}
}
func (m *RulesResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RulesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RulesResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RulesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RulesResponse.Merge(m, src)
}
func (m *RulesResponse) XXX_Size() int {
	return m.Size()
}
func (m *RulesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RulesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RulesResponse proto.InternalMessageInfo

func (m *RulesResponse) GetGroups() []*rules.RuleGroupDesc {
	if m != nil {
		return m.Groups
	}
	return nil
}

func init() {
	proto.RegisterType((*RulesRequest)(nil), "ruler.RulesRequest")
	proto.RegisterType((*RulesResponse)(nil), "ruler.RulesResponse")
}

func init() { proto.RegisterFile("ruler.proto", fileDescriptor_9ecbec0a4cfddea6) }

var fileDescriptor_9ecbec0a4cfddea6 = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8f, 0x3d, 0x4e, 0xc4, 0x30,
	0x10, 0x85, 0x6d, 0xa1, 0xdd, 0xc2, 0x0b, 0x14, 0x61, 0x0b, 0x94, 0x62, 0x84, 0x52, 0x51, 0x40,
	0x22, 0x2d, 0xdb, 0xa1, 0x6d, 0x10, 0x12, 0x7d, 0x8e, 0x90, 0xc8, 0x98, 0xff, 0x31, 0xfe, 0x91,
	0x28, 0x39, 0x02, 0xc7, 0xe0, 0x28, 0x94, 0x29, 0xb7, 0x24, 0x4e, 0x43, 0xb9, 0x47, 0x40, 0x1e,
	0xa7, 0xc8, 0x36, 0xa3, 0xf9, 0xfc, 0xde, 0xf3, 0xe8, 0x89, 0x85, 0xf1, 0x2f, 0xd2, 0x94, 0xda,
	0xa0, 0xc3, 0x6c, 0x46, 0x90, 0x5f, 0xaa, 0x47, 0xf7, 0xe0, 0x9b, 0xb2, 0xc5, 0xd7, 0x4a, 0xa1,
	0xc2, 0x8a, 0xd4, 0xc6, 0xdf, 0x13, 0x11, 0xd0, 0x96, 0x52, 0xf9, 0xf5, 0xc4, 0xde, 0xa2, 0x71,
	0xf2, 0x43, 0x1b, 0x7c, 0x92, 0xad, 0x1b, 0xa9, 0xd2, 0xcf, 0xaa, 0xa2, 0x9f, 0x69, 0xda, 0x34,
	0x53, 0xb8, 0x38, 0x16, 0x87, 0x75, 0xc4, 0x5a, 0xbe, 0x7b, 0x69, 0x5d, 0xb1, 0x11, 0x47, 0x23,
	0x5b, 0x8d, 0x6f, 0x56, 0x66, 0x17, 0x62, 0xae, 0x0c, 0x7a, 0x6d, 0x4f, 0xf9, 0xd9, 0xc1, 0xf9,
	0x62, 0xb5, 0x2c, 0x53, 0x3c, 0xba, 0xee, 0xa2, 0x70, 0x2b, 0x6d, 0x5b, 0x8f, 0x9e, 0xd5, 0x46,
	0xcc, 0xa2, 0x60, 0xb2, 0x75, 0x5a, 0x6c, 0x76, 0x52, 0xa6, 0x86, 0xd3, 0x2b, 0xf9, 0x72, 0xff,
	0x31, 0x9d, 0x2a, 0xd8, 0xcd, 0xba, 0xeb, 0x81, 0x6d, 0x7b, 0x60, 0xbb, 0x1e, 0xf8, 0x67, 0x00,
	0xfe, 0x1d, 0x80, 0xff, 0x04, 0xe0, 0x5d, 0x00, 0xfe, 0x1b, 0x80, 0xff, 0x05, 0x60, 0xbb, 0x00,
	0xfc, 0x6b, 0x00, 0xd6, 0x0d, 0xc0, 0xb6, 0x03, 0xb0, 0x66, 0x4e, 0x55, 0xae, 0xfe, 0x03, 0x00,
	0x00, 0xff, 0xff, 0xec, 0x23, 0x45, 0x0d, 0x4c, 0x01, 0x00, 0x00,
}

func (this *RulesRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RulesRequest)
	if !ok {
		that2, ok := that.(RulesRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	return true
}
func (this *RulesResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RulesResponse)
	if !ok {
		that2, ok := that.(RulesResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Groups) != len(that1.Groups) {
		return false
	}
	for i := range this.Groups {
		if !this.Groups[i].Equal(that1.Groups[i]) {
			return false
		}
	}
	return true
}
func (this *RulesRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 4)
	s = append(s, "&ruler.RulesRequest{")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *RulesResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&ruler.RulesResponse{")
	if this.Groups != nil {
		s = append(s, "Groups: "+fmt.Sprintf("%#v", this.Groups)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringRuler(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RulerClient is the client API for Ruler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RulerClient interface {
	Rules(ctx context.Context, in *RulesRequest, opts ...grpc.CallOption) (*RulesResponse, error)
}

type rulerClient struct {
	cc *grpc.ClientConn
}

func NewRulerClient(cc *grpc.ClientConn) RulerClient {
	return &rulerClient{cc}
}

func (c *rulerClient) Rules(ctx context.Context, in *RulesRequest, opts ...grpc.CallOption) (*RulesResponse, error) {
	out := new(RulesResponse)
	err := c.cc.Invoke(ctx, "/ruler.Ruler/Rules", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RulerServer is the server API for Ruler service.
type RulerServer interface {
	Rules(context.Context, *RulesRequest) (*RulesResponse, error)
}

func RegisterRulerServer(s *grpc.Server, srv RulerServer) {
	s.RegisterService(&_Ruler_serviceDesc, srv)
}

func _Ruler_Rules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RulerServer).Rules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ruler.Ruler/Rules",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RulerServer).Rules(ctx, req.(*RulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Ruler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ruler.Ruler",
	HandlerType: (*RulerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Rules",
			Handler:    _Ruler_Rules_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ruler.proto",
}

func (m *RulesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RulesRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *RulesResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RulesResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Groups) > 0 {
		for _, msg := range m.Groups {
			dAtA[i] = 0xa
			i++
			i = encodeVarintRuler(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeVarintRuler(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *RulesRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *RulesResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Groups) > 0 {
		for _, e := range m.Groups {
			l = e.Size()
			n += 1 + l + sovRuler(uint64(l))
		}
	}
	return n
}

func sovRuler(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozRuler(x uint64) (n int) {
	return sovRuler(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *RulesRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&RulesRequest{`,
		`}`,
	}, "")
	return s
}
func (this *RulesResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&RulesResponse{`,
		`Groups:` + strings.Replace(fmt.Sprintf("%v", this.Groups), "RuleGroupDesc", "rules.RuleGroupDesc", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringRuler(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *RulesRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRuler
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RulesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RulesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipRuler(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRuler
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRuler
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
func (m *RulesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRuler
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RulesResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RulesResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Groups", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRuler
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRuler
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRuler
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Groups = append(m.Groups, &rules.RuleGroupDesc{})
			if err := m.Groups[len(m.Groups)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRuler(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRuler
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRuler
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
func skipRuler(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRuler
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
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
					return 0, ErrIntOverflowRuler
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
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
					return 0, ErrIntOverflowRuler
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthRuler
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthRuler
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowRuler
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
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
				next, err := skipRuler(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthRuler
				}
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
	ErrInvalidLengthRuler = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRuler   = fmt.Errorf("proto: integer overflow")
)

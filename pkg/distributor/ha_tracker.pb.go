// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ha_tracker.proto

package distributor

import (
	fmt "fmt"
	io "io"
	math "math"
	reflect "reflect"
	strings "strings"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type ReplicaDesc struct {
	Replica    string `protobuf:"bytes,1,opt,name=replica,proto3" json:"replica,omitempty"`
	ReceivedAt int64  `protobuf:"varint,2,opt,name=receivedAt,proto3" json:"receivedAt,omitempty"`
}

func (m *ReplicaDesc) Reset()      { *m = ReplicaDesc{} }
func (*ReplicaDesc) ProtoMessage() {}
func (*ReplicaDesc) Descriptor() ([]byte, []int) {
	return fileDescriptor_86f0e7bcf71d860b, []int{0}
}
func (m *ReplicaDesc) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReplicaDesc) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReplicaDesc.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReplicaDesc) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReplicaDesc.Merge(m, src)
}
func (m *ReplicaDesc) XXX_Size() int {
	return m.Size()
}
func (m *ReplicaDesc) XXX_DiscardUnknown() {
	xxx_messageInfo_ReplicaDesc.DiscardUnknown(m)
}

var xxx_messageInfo_ReplicaDesc proto.InternalMessageInfo

func (m *ReplicaDesc) GetReplica() string {
	if m != nil {
		return m.Replica
	}
	return ""
}

func (m *ReplicaDesc) GetReceivedAt() int64 {
	if m != nil {
		return m.ReceivedAt
	}
	return 0
}

func init() {
	proto.RegisterType((*ReplicaDesc)(nil), "distributor.ReplicaDesc")
}

func init() { proto.RegisterFile("ha_tracker.proto", fileDescriptor_86f0e7bcf71d860b) }

var fileDescriptor_86f0e7bcf71d860b = []byte{
	// 201 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0x48, 0x8c, 0x2f,
	0x29, 0x4a, 0x4c, 0xce, 0x4e, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4e, 0xc9,
	0x2c, 0x2e, 0x29, 0xca, 0x4c, 0x2a, 0x2d, 0xc9, 0x2f, 0x92, 0xd2, 0x4d, 0xcf, 0x2c, 0xc9, 0x28,
	0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xcf, 0x4f, 0xcf, 0xd7, 0x07, 0xab, 0x49, 0x2a, 0x4d,
	0x03, 0xf3, 0xc0, 0x1c, 0x30, 0x0b, 0xa2, 0x57, 0xc9, 0x9d, 0x8b, 0x3b, 0x28, 0xb5, 0x20, 0x27,
	0x33, 0x39, 0xd1, 0x25, 0xb5, 0x38, 0x59, 0x48, 0x82, 0x8b, 0xbd, 0x08, 0xc2, 0x95, 0x60, 0x54,
	0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x71, 0x85, 0xe4, 0xb8, 0xb8, 0x8a, 0x52, 0x93, 0x53, 0x33, 0xcb,
	0x52, 0x53, 0x1c, 0x4b, 0x24, 0x98, 0x14, 0x18, 0x35, 0x98, 0x83, 0x90, 0x44, 0x9c, 0x4c, 0x2e,
	0x3c, 0x94, 0x63, 0xb8, 0xf1, 0x50, 0x8e, 0xe1, 0xc3, 0x43, 0x39, 0xc6, 0x86, 0x47, 0x72, 0x8c,
	0x2b, 0x1e, 0xc9, 0x31, 0x9e, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72,
	0x8c, 0x2f, 0x1e, 0xc9, 0x31, 0x7c, 0x78, 0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3, 0x85, 0xc7,
	0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x24, 0xb1, 0x81, 0x5d, 0x61, 0x0c, 0x08, 0x00, 0x00, 0xff,
	0xff, 0xa4, 0xe2, 0xd2, 0xff, 0xd5, 0x00, 0x00, 0x00,
}

func (this *ReplicaDesc) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ReplicaDesc)
	if !ok {
		that2, ok := that.(ReplicaDesc)
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
	if this.Replica != that1.Replica {
		return false
	}
	if this.ReceivedAt != that1.ReceivedAt {
		return false
	}
	return true
}
func (this *ReplicaDesc) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&distributor.ReplicaDesc{")
	s = append(s, "Replica: "+fmt.Sprintf("%#v", this.Replica)+",\n")
	s = append(s, "ReceivedAt: "+fmt.Sprintf("%#v", this.ReceivedAt)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringHaTracker(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *ReplicaDesc) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReplicaDesc) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Replica) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintHaTracker(dAtA, i, uint64(len(m.Replica)))
		i += copy(dAtA[i:], m.Replica)
	}
	if m.ReceivedAt != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintHaTracker(dAtA, i, uint64(m.ReceivedAt))
	}
	return i, nil
}

func encodeVarintHaTracker(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ReplicaDesc) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Replica)
	if l > 0 {
		n += 1 + l + sovHaTracker(uint64(l))
	}
	if m.ReceivedAt != 0 {
		n += 1 + sovHaTracker(uint64(m.ReceivedAt))
	}
	return n
}

func sovHaTracker(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozHaTracker(x uint64) (n int) {
	return sovHaTracker(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *ReplicaDesc) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ReplicaDesc{`,
		`Replica:` + fmt.Sprintf("%v", this.Replica) + `,`,
		`ReceivedAt:` + fmt.Sprintf("%v", this.ReceivedAt) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringHaTracker(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *ReplicaDesc) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHaTracker
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
			return fmt.Errorf("proto: ReplicaDesc: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReplicaDesc: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Replica", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHaTracker
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthHaTracker
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthHaTracker
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Replica = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceivedAt", wireType)
			}
			m.ReceivedAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHaTracker
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReceivedAt |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipHaTracker(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHaTracker
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthHaTracker
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
func skipHaTracker(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowHaTracker
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
					return 0, ErrIntOverflowHaTracker
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
					return 0, ErrIntOverflowHaTracker
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
				return 0, ErrInvalidLengthHaTracker
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthHaTracker
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowHaTracker
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
				next, err := skipHaTracker(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthHaTracker
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
	ErrInvalidLengthHaTracker = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowHaTracker   = fmt.Errorf("proto: integer overflow")
)

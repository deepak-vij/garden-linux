// Code generated by protoc-gen-gogo.
// source: net_out.proto
// DO NOT EDIT!

package garden

import proto "github.com/gogo/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type NetOutRequest_Protocol int32

const (
	NetOutRequest_ALL NetOutRequest_Protocol = 0
	NetOutRequest_TCP NetOutRequest_Protocol = 1
)

var NetOutRequest_Protocol_name = map[int32]string{
	0: "ALL",
	1: "TCP",
}
var NetOutRequest_Protocol_value = map[string]int32{
	"ALL": 0,
	"TCP": 1,
}

func (x NetOutRequest_Protocol) Enum() *NetOutRequest_Protocol {
	p := new(NetOutRequest_Protocol)
	*p = x
	return p
}
func (x NetOutRequest_Protocol) String() string {
	return proto.EnumName(NetOutRequest_Protocol_name, int32(x))
}
func (x *NetOutRequest_Protocol) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(NetOutRequest_Protocol_value, data, "NetOutRequest_Protocol")
	if err != nil {
		return err
	}
	*x = NetOutRequest_Protocol(value)
	return nil
}

type NetOutRequest struct {
	Handle           *string                 `protobuf:"bytes,1,req,name=handle" json:"handle,omitempty"`
	Network          *string                 `protobuf:"bytes,2,opt,name=network" json:"network,omitempty"`
	Port             *uint32                 `protobuf:"varint,3,opt,name=port" json:"port,omitempty"`
	Protocol         *NetOutRequest_Protocol `protobuf:"varint,5,opt,name=protocol,enum=garden.NetOutRequest_Protocol" json:"protocol,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *NetOutRequest) Reset()         { *m = NetOutRequest{} }
func (m *NetOutRequest) String() string { return proto.CompactTextString(m) }
func (*NetOutRequest) ProtoMessage()    {}

func (m *NetOutRequest) GetHandle() string {
	if m != nil && m.Handle != nil {
		return *m.Handle
	}
	return ""
}

func (m *NetOutRequest) GetNetwork() string {
	if m != nil && m.Network != nil {
		return *m.Network
	}
	return ""
}

func (m *NetOutRequest) GetPort() uint32 {
	if m != nil && m.Port != nil {
		return *m.Port
	}
	return 0
}

func (m *NetOutRequest) GetProtocol() NetOutRequest_Protocol {
	if m != nil && m.Protocol != nil {
		return *m.Protocol
	}
	return NetOutRequest_ALL
}

type NetOutResponse struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *NetOutResponse) Reset()         { *m = NetOutResponse{} }
func (m *NetOutResponse) String() string { return proto.CompactTextString(m) }
func (*NetOutResponse) ProtoMessage()    {}

func init() {
	proto.RegisterEnum("garden.NetOutRequest_Protocol", NetOutRequest_Protocol_name, NetOutRequest_Protocol_value)
}

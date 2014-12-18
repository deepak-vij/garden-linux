// Code generated by protoc-gen-gogo.
// source: message.proto
// DO NOT EDIT!

package garden

import proto "github.com/gogo/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type Message_Type int32

const (
	Message_Error          Message_Type = 1
	Message_Create         Message_Type = 11
	Message_Stop           Message_Type = 12
	Message_Destroy        Message_Type = 13
	Message_Info           Message_Type = 14
	Message_NetIn          Message_Type = 31
	Message_NetOut         Message_Type = 32
	Message_LimitMemory    Message_Type = 51
	Message_LimitDisk      Message_Type = 52
	Message_LimitBandwidth Message_Type = 53
	Message_LimitCpu       Message_Type = 54
	Message_Run            Message_Type = 71
	Message_Attach         Message_Type = 72
	Message_ProcessPayload Message_Type = 73
	Message_Ping           Message_Type = 91
	Message_List           Message_Type = 92
	Message_Capacity       Message_Type = 94
	Message_StreamIn       Message_Type = 95
	Message_StreamOut      Message_Type = 96
)

var Message_Type_name = map[int32]string{
	1:  "Error",
	11: "Create",
	12: "Stop",
	13: "Destroy",
	14: "Info",
	31: "NetIn",
	32: "NetOut",
	51: "LimitMemory",
	52: "LimitDisk",
	53: "LimitBandwidth",
	54: "LimitCpu",
	71: "Run",
	72: "Attach",
	73: "ProcessPayload",
	91: "Ping",
	92: "List",
	94: "Capacity",
	95: "StreamIn",
	96: "StreamOut",
}
var Message_Type_value = map[string]int32{
	"Error":          1,
	"Create":         11,
	"Stop":           12,
	"Destroy":        13,
	"Info":           14,
	"NetIn":          31,
	"NetOut":         32,
	"LimitMemory":    51,
	"LimitDisk":      52,
	"LimitBandwidth": 53,
	"LimitCpu":       54,
	"Run":            71,
	"Attach":         72,
	"ProcessPayload": 73,
	"Ping":           91,
	"List":           92,
	"Capacity":       94,
	"StreamIn":       95,
	"StreamOut":      96,
}

func (x Message_Type) Enum() *Message_Type {
	p := new(Message_Type)
	*p = x
	return p
}
func (x Message_Type) String() string {
	return proto.EnumName(Message_Type_name, int32(x))
}
func (x *Message_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Message_Type_value, data, "Message_Type")
	if err != nil {
		return err
	}
	*x = Message_Type(value)
	return nil
}

type Message struct {
	Type             *Message_Type `protobuf:"varint,1,req,name=type,enum=garden.Message_Type" json:"type,omitempty"`
	Payload          []byte        `protobuf:"bytes,2,req,name=payload" json:"payload,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}

func (m *Message) GetType() Message_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Message_Error
}

func (m *Message) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto.RegisterEnum("garden.Message_Type", Message_Type_name, Message_Type_value)
}

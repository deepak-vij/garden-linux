// Code generated by protoc-gen-gogo.
// source: create.proto
// DO NOT EDIT!

package garden

import proto "github.com/gogo/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type CreateRequest_BindMount_Mode int32

const (
	CreateRequest_BindMount_RO CreateRequest_BindMount_Mode = 0
	CreateRequest_BindMount_RW CreateRequest_BindMount_Mode = 1
)

var CreateRequest_BindMount_Mode_name = map[int32]string{
	0: "RO",
	1: "RW",
}
var CreateRequest_BindMount_Mode_value = map[string]int32{
	"RO": 0,
	"RW": 1,
}

func (x CreateRequest_BindMount_Mode) Enum() *CreateRequest_BindMount_Mode {
	p := new(CreateRequest_BindMount_Mode)
	*p = x
	return p
}
func (x CreateRequest_BindMount_Mode) String() string {
	return proto.EnumName(CreateRequest_BindMount_Mode_name, int32(x))
}
func (x *CreateRequest_BindMount_Mode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(CreateRequest_BindMount_Mode_value, data, "CreateRequest_BindMount_Mode")
	if err != nil {
		return err
	}
	*x = CreateRequest_BindMount_Mode(value)
	return nil
}

type CreateRequest_BindMount_Origin int32

const (
	CreateRequest_BindMount_Host      CreateRequest_BindMount_Origin = 0
	CreateRequest_BindMount_Container CreateRequest_BindMount_Origin = 1
)

var CreateRequest_BindMount_Origin_name = map[int32]string{
	0: "Host",
	1: "Container",
}
var CreateRequest_BindMount_Origin_value = map[string]int32{
	"Host":      0,
	"Container": 1,
}

func (x CreateRequest_BindMount_Origin) Enum() *CreateRequest_BindMount_Origin {
	p := new(CreateRequest_BindMount_Origin)
	*p = x
	return p
}
func (x CreateRequest_BindMount_Origin) String() string {
	return proto.EnumName(CreateRequest_BindMount_Origin_name, int32(x))
}
func (x *CreateRequest_BindMount_Origin) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(CreateRequest_BindMount_Origin_value, data, "CreateRequest_BindMount_Origin")
	if err != nil {
		return err
	}
	*x = CreateRequest_BindMount_Origin(value)
	return nil
}

type CreateRequest struct {
	BindMounts       []*CreateRequest_BindMount `protobuf:"bytes,1,rep,name=bind_mounts" json:"bind_mounts,omitempty"`
	GraceTime        *uint32                    `protobuf:"varint,2,opt,name=grace_time" json:"grace_time,omitempty"`
	Handle           *string                    `protobuf:"bytes,3,opt,name=handle" json:"handle,omitempty"`
	Network          *string                    `protobuf:"bytes,4,opt,name=network" json:"network,omitempty"`
	Rootfs           *string                    `protobuf:"bytes,5,opt,name=rootfs" json:"rootfs,omitempty"`
	Properties       []*Property                `protobuf:"bytes,6,rep,name=properties" json:"properties,omitempty"`
	Env              []*EnvironmentVariable     `protobuf:"bytes,7,rep,name=env" json:"env,omitempty"`
	Privileged       *bool                      `protobuf:"varint,8,opt,name=privileged" json:"privileged,omitempty"`
	XXX_unrecognized []byte                     `json:"-"`
}

func (m *CreateRequest) Reset()         { *m = CreateRequest{} }
func (m *CreateRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()    {}

func (m *CreateRequest) GetBindMounts() []*CreateRequest_BindMount {
	if m != nil {
		return m.BindMounts
	}
	return nil
}

func (m *CreateRequest) GetGraceTime() uint32 {
	if m != nil && m.GraceTime != nil {
		return *m.GraceTime
	}
	return 0
}

func (m *CreateRequest) GetHandle() string {
	if m != nil && m.Handle != nil {
		return *m.Handle
	}
	return ""
}

func (m *CreateRequest) GetNetwork() string {
	if m != nil && m.Network != nil {
		return *m.Network
	}
	return ""
}

func (m *CreateRequest) GetRootfs() string {
	if m != nil && m.Rootfs != nil {
		return *m.Rootfs
	}
	return ""
}

func (m *CreateRequest) GetProperties() []*Property {
	if m != nil {
		return m.Properties
	}
	return nil
}

func (m *CreateRequest) GetEnv() []*EnvironmentVariable {
	if m != nil {
		return m.Env
	}
	return nil
}

func (m *CreateRequest) GetPrivileged() bool {
	if m != nil && m.Privileged != nil {
		return *m.Privileged
	}
	return false
}

type CreateRequest_BindMount struct {
	SrcPath          *string                         `protobuf:"bytes,1,req,name=src_path" json:"src_path,omitempty"`
	DstPath          *string                         `protobuf:"bytes,2,req,name=dst_path" json:"dst_path,omitempty"`
	Mode             *CreateRequest_BindMount_Mode   `protobuf:"varint,3,req,name=mode,enum=garden.CreateRequest_BindMount_Mode" json:"mode,omitempty"`
	Origin           *CreateRequest_BindMount_Origin `protobuf:"varint,4,opt,name=origin,enum=garden.CreateRequest_BindMount_Origin" json:"origin,omitempty"`
	XXX_unrecognized []byte                          `json:"-"`
}

func (m *CreateRequest_BindMount) Reset()         { *m = CreateRequest_BindMount{} }
func (m *CreateRequest_BindMount) String() string { return proto.CompactTextString(m) }
func (*CreateRequest_BindMount) ProtoMessage()    {}

func (m *CreateRequest_BindMount) GetSrcPath() string {
	if m != nil && m.SrcPath != nil {
		return *m.SrcPath
	}
	return ""
}

func (m *CreateRequest_BindMount) GetDstPath() string {
	if m != nil && m.DstPath != nil {
		return *m.DstPath
	}
	return ""
}

func (m *CreateRequest_BindMount) GetMode() CreateRequest_BindMount_Mode {
	if m != nil && m.Mode != nil {
		return *m.Mode
	}
	return CreateRequest_BindMount_RO
}

func (m *CreateRequest_BindMount) GetOrigin() CreateRequest_BindMount_Origin {
	if m != nil && m.Origin != nil {
		return *m.Origin
	}
	return CreateRequest_BindMount_Host
}

type CreateResponse struct {
	Handle           *string `protobuf:"bytes,1,req,name=handle" json:"handle,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CreateResponse) Reset()         { *m = CreateResponse{} }
func (m *CreateResponse) String() string { return proto.CompactTextString(m) }
func (*CreateResponse) ProtoMessage()    {}

func (m *CreateResponse) GetHandle() string {
	if m != nil && m.Handle != nil {
		return *m.Handle
	}
	return ""
}

func init() {
	proto.RegisterEnum("garden.CreateRequest_BindMount_Mode", CreateRequest_BindMount_Mode_name, CreateRequest_BindMount_Mode_value)
	proto.RegisterEnum("garden.CreateRequest_BindMount_Origin", CreateRequest_BindMount_Origin_name, CreateRequest_BindMount_Origin_value)
}

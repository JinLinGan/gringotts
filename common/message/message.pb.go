// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package message

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

type MonitorItemType int32

const (
	MonitorItemType_SELF        MonitorItemType = 0
	MonitorItemType_THIRD_PARTY MonitorItemType = 1
)

var MonitorItemType_name = map[int32]string{
	0: "SELF",
	1: "THIRD_PARTY",
}

var MonitorItemType_value = map[string]int32{
	"SELF":        0,
	"THIRD_PARTY": 1,
}

func (x MonitorItemType) String() string {
	return proto.EnumName(MonitorItemType_name, int32(x))
}

func (MonitorItemType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

type SelfMonitorFunc int32

const (
	SelfMonitorFunc_CPU SelfMonitorFunc = 0
	SelfMonitorFunc_MEM SelfMonitorFunc = 1
)

var SelfMonitorFunc_name = map[int32]string{
	0: "CPU",
	1: "MEM",
}

var SelfMonitorFunc_value = map[string]int32{
	"CPU": 0,
	"MEM": 1,
}

func (x SelfMonitorFunc) String() string {
	return proto.EnumName(SelfMonitorFunc_name, int32(x))
}

func (SelfMonitorFunc) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

type File struct {
	FileName             string   `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	Sha1Hash             string   `protobuf:"bytes,2,opt,name=sha1_hash,json=sha1Hash,proto3" json:"sha1_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *File) Reset()         { *m = File{} }
func (m *File) String() string { return proto.CompactTextString(m) }
func (*File) ProtoMessage()    {}
func (*File) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *File) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_File.Unmarshal(m, b)
}
func (m *File) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_File.Marshal(b, m, deterministic)
}
func (m *File) XXX_Merge(src proto.Message) {
	xxx_messageInfo_File.Merge(m, src)
}
func (m *File) XXX_Size() int {
	return xxx_messageInfo_File.Size(m)
}
func (m *File) XXX_DiscardUnknown() {
	xxx_messageInfo_File.DiscardUnknown(m)
}

var xxx_messageInfo_File proto.InternalMessageInfo

func (m *File) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func (m *File) GetSha1Hash() string {
	if m != nil {
		return m.Sha1Hash
	}
	return ""
}

type FileChunk struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileChunk) Reset()         { *m = FileChunk{} }
func (m *FileChunk) String() string { return proto.CompactTextString(m) }
func (*FileChunk) ProtoMessage()    {}
func (*FileChunk) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

func (m *FileChunk) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileChunk.Unmarshal(m, b)
}
func (m *FileChunk) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileChunk.Marshal(b, m, deterministic)
}
func (m *FileChunk) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileChunk.Merge(m, src)
}
func (m *FileChunk) XXX_Size() int {
	return xxx_messageInfo_FileChunk.Size(m)
}
func (m *FileChunk) XXX_DiscardUnknown() {
	xxx_messageInfo_FileChunk.DiscardUnknown(m)
}

var xxx_messageInfo_FileChunk proto.InternalMessageInfo

func (m *FileChunk) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type HeartBeatRequest struct {
	AgentId              string   `protobuf:"bytes,1,opt,name=agent_id,json=agentId,proto3" json:"agent_id,omitempty"`
	HostName             string   `protobuf:"bytes,2,opt,name=host_name,json=hostName,proto3" json:"host_name,omitempty"`
	Time                 int64    `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartBeatRequest) Reset()         { *m = HeartBeatRequest{} }
func (m *HeartBeatRequest) String() string { return proto.CompactTextString(m) }
func (*HeartBeatRequest) ProtoMessage()    {}
func (*HeartBeatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
}

func (m *HeartBeatRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartBeatRequest.Unmarshal(m, b)
}
func (m *HeartBeatRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartBeatRequest.Marshal(b, m, deterministic)
}
func (m *HeartBeatRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartBeatRequest.Merge(m, src)
}
func (m *HeartBeatRequest) XXX_Size() int {
	return xxx_messageInfo_HeartBeatRequest.Size(m)
}
func (m *HeartBeatRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartBeatRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HeartBeatRequest proto.InternalMessageInfo

func (m *HeartBeatRequest) GetAgentId() string {
	if m != nil {
		return m.AgentId
	}
	return ""
}

func (m *HeartBeatRequest) GetHostName() string {
	if m != nil {
		return m.HostName
	}
	return ""
}

func (m *HeartBeatRequest) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

type MonitorInfo struct {
	Items                []*MonitorItem `protobuf:"bytes,1,rep,name=Items,proto3" json:"Items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *MonitorInfo) Reset()         { *m = MonitorInfo{} }
func (m *MonitorInfo) String() string { return proto.CompactTextString(m) }
func (*MonitorInfo) ProtoMessage()    {}
func (*MonitorInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{3}
}

func (m *MonitorInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorInfo.Unmarshal(m, b)
}
func (m *MonitorInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorInfo.Marshal(b, m, deterministic)
}
func (m *MonitorInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorInfo.Merge(m, src)
}
func (m *MonitorInfo) XXX_Size() int {
	return xxx_messageInfo_MonitorInfo.Size(m)
}
func (m *MonitorInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorInfo.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorInfo proto.InternalMessageInfo

func (m *MonitorInfo) GetItems() []*MonitorItem {
	if m != nil {
		return m.Items
	}
	return nil
}

type MonitorItem struct {
	TaskId                    int64           `protobuf:"varint,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	ExecIntervalSecond        int64           `protobuf:"varint,2,opt,name=exec_interval_second,json=execIntervalSecond,proto3" json:"exec_interval_second,omitempty"`
	Type                      MonitorItemType `protobuf:"varint,3,opt,name=type,proto3,enum=MonitorItemType" json:"type,omitempty"`
	SelfFunc                  SelfMonitorFunc `protobuf:"varint,4,opt,name=self_func,json=selfFunc,proto3,enum=SelfMonitorFunc" json:"self_func,omitempty"`
	SelfFuncConfig            string          `protobuf:"bytes,5,opt,name=self_func_config,json=selfFuncConfig,proto3" json:"self_func_config,omitempty"`
	ThirdPartyFuncName        string          `protobuf:"bytes,6,opt,name=third_party_func_name,json=thirdPartyFuncName,proto3" json:"third_party_func_name,omitempty"`
	ThirdPartyDependencyFiles string          `protobuf:"bytes,7,opt,name=third_party_dependency_files,json=thirdPartyDependencyFiles,proto3" json:"third_party_dependency_files,omitempty"`
	XXX_NoUnkeyedLiteral      struct{}        `json:"-"`
	XXX_unrecognized          []byte          `json:"-"`
	XXX_sizecache             int32           `json:"-"`
}

func (m *MonitorItem) Reset()         { *m = MonitorItem{} }
func (m *MonitorItem) String() string { return proto.CompactTextString(m) }
func (*MonitorItem) ProtoMessage()    {}
func (*MonitorItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{4}
}

func (m *MonitorItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorItem.Unmarshal(m, b)
}
func (m *MonitorItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorItem.Marshal(b, m, deterministic)
}
func (m *MonitorItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorItem.Merge(m, src)
}
func (m *MonitorItem) XXX_Size() int {
	return xxx_messageInfo_MonitorItem.Size(m)
}
func (m *MonitorItem) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorItem.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorItem proto.InternalMessageInfo

func (m *MonitorItem) GetTaskId() int64 {
	if m != nil {
		return m.TaskId
	}
	return 0
}

func (m *MonitorItem) GetExecIntervalSecond() int64 {
	if m != nil {
		return m.ExecIntervalSecond
	}
	return 0
}

func (m *MonitorItem) GetType() MonitorItemType {
	if m != nil {
		return m.Type
	}
	return MonitorItemType_SELF
}

func (m *MonitorItem) GetSelfFunc() SelfMonitorFunc {
	if m != nil {
		return m.SelfFunc
	}
	return SelfMonitorFunc_CPU
}

func (m *MonitorItem) GetSelfFuncConfig() string {
	if m != nil {
		return m.SelfFuncConfig
	}
	return ""
}

func (m *MonitorItem) GetThirdPartyFuncName() string {
	if m != nil {
		return m.ThirdPartyFuncName
	}
	return ""
}

func (m *MonitorItem) GetThirdPartyDependencyFiles() string {
	if m != nil {
		return m.ThirdPartyDependencyFiles
	}
	return ""
}

type HeartBeatResponse struct {
	ServerId             string       `protobuf:"bytes,1,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty"`
	ConfigVersion        int64        `protobuf:"varint,2,opt,name=config_version,json=configVersion,proto3" json:"config_version,omitempty"`
	ConfigChangeTime     int64        `protobuf:"varint,3,opt,name=config_change_time,json=configChangeTime,proto3" json:"config_change_time,omitempty"`
	MonitorInfo          *MonitorInfo `protobuf:"bytes,4,opt,name=monitor_info,json=monitorInfo,proto3" json:"monitor_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *HeartBeatResponse) Reset()         { *m = HeartBeatResponse{} }
func (m *HeartBeatResponse) String() string { return proto.CompactTextString(m) }
func (*HeartBeatResponse) ProtoMessage()    {}
func (*HeartBeatResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{5}
}

func (m *HeartBeatResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartBeatResponse.Unmarshal(m, b)
}
func (m *HeartBeatResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartBeatResponse.Marshal(b, m, deterministic)
}
func (m *HeartBeatResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartBeatResponse.Merge(m, src)
}
func (m *HeartBeatResponse) XXX_Size() int {
	return xxx_messageInfo_HeartBeatResponse.Size(m)
}
func (m *HeartBeatResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartBeatResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HeartBeatResponse proto.InternalMessageInfo

func (m *HeartBeatResponse) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
}

func (m *HeartBeatResponse) GetConfigVersion() int64 {
	if m != nil {
		return m.ConfigVersion
	}
	return 0
}

func (m *HeartBeatResponse) GetConfigChangeTime() int64 {
	if m != nil {
		return m.ConfigChangeTime
	}
	return 0
}

func (m *HeartBeatResponse) GetMonitorInfo() *MonitorInfo {
	if m != nil {
		return m.MonitorInfo
	}
	return nil
}

type RegisterRequest struct {
	HostId               int64                      `protobuf:"varint,1,opt,name=host_id,json=hostId,proto3" json:"host_id,omitempty"`
	HostName             string                     `protobuf:"bytes,2,opt,name=host_name,json=hostName,proto3" json:"host_name,omitempty"`
	NetInfo              []*RegisterRequest_NetInfo `protobuf:"bytes,3,rep,name=net_info,json=netInfo,proto3" json:"net_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *RegisterRequest) Reset()         { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterRequest) ProtoMessage()    {}
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{6}
}

func (m *RegisterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRequest.Unmarshal(m, b)
}
func (m *RegisterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRequest.Marshal(b, m, deterministic)
}
func (m *RegisterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRequest.Merge(m, src)
}
func (m *RegisterRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterRequest.Size(m)
}
func (m *RegisterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRequest proto.InternalMessageInfo

func (m *RegisterRequest) GetHostId() int64 {
	if m != nil {
		return m.HostId
	}
	return 0
}

func (m *RegisterRequest) GetHostName() string {
	if m != nil {
		return m.HostName
	}
	return ""
}

func (m *RegisterRequest) GetNetInfo() []*RegisterRequest_NetInfo {
	if m != nil {
		return m.NetInfo
	}
	return nil
}

type RegisterRequest_NetInfo struct {
	IpAddress            string   `protobuf:"bytes,1,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	MacAddress           string   `protobuf:"bytes,2,opt,name=mac_address,json=macAddress,proto3" json:"mac_address,omitempty"`
	InterfaceName        string   `protobuf:"bytes,3,opt,name=interface_name,json=interfaceName,proto3" json:"interface_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterRequest_NetInfo) Reset()         { *m = RegisterRequest_NetInfo{} }
func (m *RegisterRequest_NetInfo) String() string { return proto.CompactTextString(m) }
func (*RegisterRequest_NetInfo) ProtoMessage()    {}
func (*RegisterRequest_NetInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{6, 0}
}

func (m *RegisterRequest_NetInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRequest_NetInfo.Unmarshal(m, b)
}
func (m *RegisterRequest_NetInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRequest_NetInfo.Marshal(b, m, deterministic)
}
func (m *RegisterRequest_NetInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRequest_NetInfo.Merge(m, src)
}
func (m *RegisterRequest_NetInfo) XXX_Size() int {
	return xxx_messageInfo_RegisterRequest_NetInfo.Size(m)
}
func (m *RegisterRequest_NetInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRequest_NetInfo.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRequest_NetInfo proto.InternalMessageInfo

func (m *RegisterRequest_NetInfo) GetIpAddress() string {
	if m != nil {
		return m.IpAddress
	}
	return ""
}

func (m *RegisterRequest_NetInfo) GetMacAddress() string {
	if m != nil {
		return m.MacAddress
	}
	return ""
}

func (m *RegisterRequest_NetInfo) GetInterfaceName() string {
	if m != nil {
		return m.InterfaceName
	}
	return ""
}

type RegisterResponse struct {
	AgentId              int64    `protobuf:"varint,1,opt,name=agent_id,json=agentId,proto3" json:"agent_id,omitempty"`
	ConfigVersion        int64    `protobuf:"varint,2,opt,name=config_version,json=configVersion,proto3" json:"config_version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterResponse) ProtoMessage()    {}
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{7}
}

func (m *RegisterResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterResponse.Unmarshal(m, b)
}
func (m *RegisterResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterResponse.Marshal(b, m, deterministic)
}
func (m *RegisterResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterResponse.Merge(m, src)
}
func (m *RegisterResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterResponse.Size(m)
}
func (m *RegisterResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterResponse proto.InternalMessageInfo

func (m *RegisterResponse) GetAgentId() int64 {
	if m != nil {
		return m.AgentId
	}
	return 0
}

func (m *RegisterResponse) GetConfigVersion() int64 {
	if m != nil {
		return m.ConfigVersion
	}
	return 0
}

func init() {
	proto.RegisterEnum("MonitorItemType", MonitorItemType_name, MonitorItemType_value)
	proto.RegisterEnum("SelfMonitorFunc", SelfMonitorFunc_name, SelfMonitorFunc_value)
	proto.RegisterType((*File)(nil), "File")
	proto.RegisterType((*FileChunk)(nil), "FileChunk")
	proto.RegisterType((*HeartBeatRequest)(nil), "HeartBeatRequest")
	proto.RegisterType((*MonitorInfo)(nil), "MonitorInfo")
	proto.RegisterType((*MonitorItem)(nil), "MonitorItem")
	proto.RegisterType((*HeartBeatResponse)(nil), "HeartBeatResponse")
	proto.RegisterType((*RegisterRequest)(nil), "RegisterRequest")
	proto.RegisterType((*RegisterRequest_NetInfo)(nil), "RegisterRequest.NetInfo")
	proto.RegisterType((*RegisterResponse)(nil), "RegisterResponse")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 712 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x6f, 0x4e, 0xdb, 0x4e,
	0x10, 0xc5, 0x38, 0x90, 0x64, 0x12, 0x12, 0xb3, 0xfa, 0xfd, 0x84, 0x49, 0x5b, 0x81, 0x4c, 0x91,
	0x22, 0x44, 0x0d, 0x84, 0x03, 0xb4, 0x34, 0x80, 0x88, 0x54, 0x10, 0x32, 0x69, 0xa5, 0x7e, 0xa9,
	0xb5, 0xb5, 0xc7, 0xb1, 0x45, 0xbc, 0x76, 0xbd, 0x1b, 0xda, 0x1c, 0xa4, 0x57, 0xe9, 0x25, 0x7a,
	0x97, 0x9e, 0xa1, 0xda, 0x5d, 0xe7, 0x8f, 0xd2, 0x0f, 0xed, 0xa7, 0xcc, 0xbe, 0xf7, 0xc6, 0x33,
	0x99, 0x79, 0xbb, 0xb0, 0x95, 0x22, 0xe7, 0x74, 0x84, 0x6e, 0x5e, 0x64, 0x22, 0x73, 0xde, 0x40,
	0xe5, 0x3a, 0x19, 0x23, 0x79, 0x06, 0xf5, 0x28, 0x19, 0xa3, 0xcf, 0x68, 0x8a, 0xb6, 0xb1, 0x6f,
	0x74, 0xeb, 0x5e, 0x4d, 0x02, 0x77, 0x34, 0x55, 0x24, 0x8f, 0xe9, 0x99, 0x1f, 0x53, 0x1e, 0xdb,
	0xeb, 0x9a, 0x94, 0xc0, 0x0d, 0xe5, 0xb1, 0xb3, 0x07, 0x75, 0xf9, 0x85, 0x7e, 0x3c, 0x61, 0x8f,
	0x84, 0x40, 0x25, 0xa4, 0x82, 0xaa, 0x2f, 0x34, 0x3d, 0x15, 0x3b, 0x9f, 0xc0, 0xba, 0x41, 0x5a,
	0x88, 0xb7, 0x48, 0x85, 0x87, 0x5f, 0x26, 0xc8, 0x05, 0xd9, 0x85, 0x1a, 0x1d, 0x21, 0x13, 0x7e,
	0x12, 0x96, 0xd5, 0xaa, 0xea, 0x3c, 0x08, 0x65, 0xb1, 0x38, 0xe3, 0x42, 0x77, 0x52, 0x16, 0x93,
	0x80, 0xea, 0x84, 0x40, 0x45, 0x24, 0x29, 0xda, 0xe6, 0xbe, 0xd1, 0x35, 0x3d, 0x15, 0x3b, 0x67,
	0xd0, 0xb8, 0xcd, 0x58, 0x22, 0xb2, 0x62, 0xc0, 0xa2, 0x8c, 0x38, 0xb0, 0x31, 0x10, 0x98, 0x72,
	0xdb, 0xd8, 0x37, 0xbb, 0x8d, 0x5e, 0xd3, 0x9d, 0x91, 0x02, 0x53, 0x4f, 0x53, 0xce, 0xcf, 0xf5,
	0x45, 0x8e, 0xc0, 0x94, 0xec, 0x40, 0x55, 0x50, 0xfe, 0x38, 0xeb, 0xc6, 0xf4, 0x36, 0xe5, 0x71,
	0x10, 0x92, 0x53, 0xf8, 0x0f, 0xbf, 0x61, 0xe0, 0x27, 0x4c, 0x60, 0xf1, 0x44, 0xc7, 0x3e, 0xc7,
	0x20, 0x63, 0xa1, 0xea, 0xcb, 0xf4, 0x88, 0xe4, 0x06, 0x25, 0xf5, 0xa0, 0x18, 0xf2, 0x12, 0x2a,
	0x62, 0x9a, 0xeb, 0x0e, 0x5b, 0x3d, 0x6b, 0xb9, 0xfa, 0x70, 0x9a, 0xa3, 0xa7, 0x58, 0xf2, 0x0a,
	0xea, 0x1c, 0xc7, 0x91, 0x1f, 0x4d, 0x58, 0x60, 0x57, 0x4a, 0xe9, 0x03, 0x8e, 0xa3, 0x52, 0x7e,
	0x3d, 0x61, 0x81, 0x57, 0x93, 0x12, 0x19, 0x91, 0x2e, 0x58, 0x73, 0xb9, 0x1f, 0x64, 0x2c, 0x4a,
	0x46, 0xf6, 0x86, 0x1a, 0x4d, 0x6b, 0xa6, 0xe9, 0x2b, 0x94, 0x9c, 0xc1, 0xff, 0x22, 0x4e, 0x8a,
	0xd0, 0xcf, 0x69, 0x21, 0xa6, 0x3a, 0x41, 0x4d, 0x72, 0x53, 0xc9, 0x89, 0x22, 0xef, 0x25, 0x27,
	0x93, 0xd4, 0x4c, 0x5f, 0xc3, 0xf3, 0xe5, 0x94, 0x10, 0x73, 0x64, 0x21, 0xb2, 0x60, 0xea, 0x4b,
	0x03, 0x70, 0xbb, 0xaa, 0x32, 0x77, 0x17, 0x99, 0x97, 0x73, 0x85, 0x5c, 0x3c, 0x77, 0x7e, 0x18,
	0xb0, 0xbd, 0xb4, 0x61, 0x9e, 0x67, 0x8c, 0x6b, 0xd3, 0x60, 0xf1, 0x84, 0xc5, 0x62, 0xc7, 0x35,
	0x0d, 0x0c, 0x42, 0x72, 0x08, 0x2d, 0xfd, 0x37, 0xfc, 0x27, 0x2c, 0x78, 0x92, 0xb1, 0x72, 0xa2,
	0x5b, 0x1a, 0xfd, 0xa0, 0x41, 0x72, 0x0c, 0xa4, 0x94, 0x05, 0x31, 0x65, 0x23, 0xf4, 0x97, 0x96,
	0x6f, 0x69, 0xa6, 0xaf, 0x88, 0x61, 0x92, 0x22, 0x39, 0x81, 0x66, 0xaa, 0xc7, 0xe7, 0x27, 0x2c,
	0xca, 0xd4, 0x5c, 0x97, 0x0d, 0xc0, 0xa2, 0xcc, 0x6b, 0xa4, 0x8b, 0x83, 0xf3, 0xcb, 0x80, 0xb6,
	0x87, 0xa3, 0x84, 0x0b, 0x2c, 0x66, 0xce, 0xdc, 0x81, 0xaa, 0xb2, 0xdf, 0xc2, 0x0a, 0xf2, 0xf8,
	0x37, 0x5f, 0x9e, 0x43, 0x8d, 0xa1, 0xd0, 0x65, 0x4d, 0xe5, 0x3b, 0xdb, 0x5d, 0xf9, 0xb2, 0x7b,
	0x87, 0x42, 0xb5, 0x50, 0x65, 0x3a, 0xe8, 0xe4, 0x50, 0x2d, 0x31, 0xf2, 0x02, 0x20, 0xc9, 0x7d,
	0x1a, 0x86, 0x05, 0x72, 0x5e, 0x4e, 0xab, 0x9e, 0xe4, 0x17, 0x1a, 0x20, 0x7b, 0xd0, 0x48, 0x69,
	0x30, 0xe7, 0x75, 0x75, 0x48, 0x69, 0x30, 0x13, 0x1c, 0x42, 0x4b, 0x59, 0x34, 0xa2, 0x41, 0x79,
	0x87, 0x4d, 0xa5, 0xd9, 0x9a, 0xa3, 0xb2, 0x4d, 0x67, 0x08, 0xd6, 0xa2, 0xab, 0x72, 0x4f, 0xab,
	0x57, 0xd1, 0x5c, 0x5c, 0xc5, 0x7f, 0xdb, 0xd2, 0xd1, 0x31, 0xb4, 0x57, 0x5c, 0x4e, 0x6a, 0x50,
	0x79, 0xb8, 0x7a, 0x77, 0x6d, 0xad, 0x91, 0x36, 0x34, 0x86, 0x37, 0x03, 0xef, 0xd2, 0xbf, 0xbf,
	0xf0, 0x86, 0x1f, 0x2d, 0xe3, 0xe8, 0x00, 0xda, 0x2b, 0x46, 0x27, 0x55, 0x30, 0xfb, 0xf7, 0xef,
	0xad, 0x35, 0x19, 0xdc, 0x5e, 0xdd, 0x5a, 0x46, 0xef, 0xbb, 0x01, 0xf5, 0x51, 0x91, 0xb0, 0x51,
	0x26, 0x04, 0x27, 0x3d, 0xa8, 0xcf, 0xfd, 0x45, 0xb6, 0xdd, 0xd5, 0xd7, 0xa4, 0x43, 0xdc, 0x3f,
	0xed, 0x77, 0x00, 0xcd, 0xcb, 0xec, 0x2b, 0x1b, 0x67, 0x34, 0x54, 0x0f, 0xdc, 0x86, 0x2b, 0x7f,
	0x3a, 0xe0, 0xce, 0x1f, 0xab, 0x53, 0x83, 0x9c, 0x40, 0x6d, 0x36, 0x0f, 0x62, 0xad, 0x2e, 0xac,
	0xb3, 0xed, 0xae, 0x0e, 0xeb, 0xf3, 0xa6, 0x7a, 0x35, 0xcf, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff,
	0xe2, 0x97, 0x2f, 0xa0, 0x46, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GringottsClient is the client API for Gringotts service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GringottsClient interface {
	HeartBeat(ctx context.Context, in *HeartBeatRequest, opts ...grpc.CallOption) (*HeartBeatResponse, error)
	DownloadFile(ctx context.Context, in *File, opts ...grpc.CallOption) (Gringotts_DownloadFileClient, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
}

type gringottsClient struct {
	cc *grpc.ClientConn
}

func NewGringottsClient(cc *grpc.ClientConn) GringottsClient {
	return &gringottsClient{cc}
}

func (c *gringottsClient) HeartBeat(ctx context.Context, in *HeartBeatRequest, opts ...grpc.CallOption) (*HeartBeatResponse, error) {
	out := new(HeartBeatResponse)
	err := c.cc.Invoke(ctx, "/gringotts/HeartBeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gringottsClient) DownloadFile(ctx context.Context, in *File, opts ...grpc.CallOption) (Gringotts_DownloadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Gringotts_serviceDesc.Streams[0], "/gringotts/DownloadFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &gringottsDownloadFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Gringotts_DownloadFileClient interface {
	Recv() (*FileChunk, error)
	grpc.ClientStream
}

type gringottsDownloadFileClient struct {
	grpc.ClientStream
}

func (x *gringottsDownloadFileClient) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gringottsClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/gringotts/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GringottsServer is the server API for Gringotts service.
type GringottsServer interface {
	HeartBeat(context.Context, *HeartBeatRequest) (*HeartBeatResponse, error)
	DownloadFile(*File, Gringotts_DownloadFileServer) error
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
}

// UnimplementedGringottsServer can be embedded to have forward compatible implementations.
type UnimplementedGringottsServer struct {
}

func (*UnimplementedGringottsServer) HeartBeat(ctx context.Context, req *HeartBeatRequest) (*HeartBeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeartBeat not implemented")
}
func (*UnimplementedGringottsServer) DownloadFile(req *File, srv Gringotts_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (*UnimplementedGringottsServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}

func RegisterGringottsServer(s *grpc.Server, srv GringottsServer) {
	s.RegisterService(&_Gringotts_serviceDesc, srv)
}

func _Gringotts_HeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartBeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GringottsServer).HeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gringotts/HeartBeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GringottsServer).HeartBeat(ctx, req.(*HeartBeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gringotts_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(File)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GringottsServer).DownloadFile(m, &gringottsDownloadFileServer{stream})
}

type Gringotts_DownloadFileServer interface {
	Send(*FileChunk) error
	grpc.ServerStream
}

type gringottsDownloadFileServer struct {
	grpc.ServerStream
}

func (x *gringottsDownloadFileServer) Send(m *FileChunk) error {
	return x.ServerStream.SendMsg(m)
}

func _Gringotts_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GringottsServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gringotts/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GringottsServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Gringotts_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gringotts",
	HandlerType: (*GringottsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HeartBeat",
			Handler:    _Gringotts_HeartBeat_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _Gringotts_Register_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadFile",
			Handler:       _Gringotts_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "message.proto",
}

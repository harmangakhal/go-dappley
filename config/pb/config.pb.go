// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/dappley/go-dappley/config/pb/config.proto

package configpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Config struct {
	ConsensusConfig      *ConsensusConfig `protobuf:"bytes,1,opt,name=consensus_config,json=consensusConfig,proto3" json:"consensus_config,omitempty"`
	NodeConfig           *NodeConfig      `protobuf:"bytes,2,opt,name=node_config,json=nodeConfig,proto3" json:"node_config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bcc5bb4ac8bd8b6e, []int{0}
}
func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (dst *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(dst, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetConsensusConfig() *ConsensusConfig {
	if m != nil {
		return m.ConsensusConfig
	}
	return nil
}

func (m *Config) GetNodeConfig() *NodeConfig {
	if m != nil {
		return m.NodeConfig
	}
	return nil
}

type ConsensusConfig struct {
	MinerAddress         string   `protobuf:"bytes,1,opt,name=miner_address,json=minerAddress,proto3" json:"miner_address,omitempty"`
	PrivateKey           string   `protobuf:"bytes,2,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConsensusConfig) Reset()         { *m = ConsensusConfig{} }
func (m *ConsensusConfig) String() string { return proto.CompactTextString(m) }
func (*ConsensusConfig) ProtoMessage()    {}
func (*ConsensusConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bcc5bb4ac8bd8b6e, []int{1}
}
func (m *ConsensusConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsensusConfig.Unmarshal(m, b)
}
func (m *ConsensusConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsensusConfig.Marshal(b, m, deterministic)
}
func (dst *ConsensusConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsensusConfig.Merge(dst, src)
}
func (m *ConsensusConfig) XXX_Size() int {
	return xxx_messageInfo_ConsensusConfig.Size(m)
}
func (m *ConsensusConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsensusConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ConsensusConfig proto.InternalMessageInfo

func (m *ConsensusConfig) GetMinerAddress() string {
	if m != nil {
		return m.MinerAddress
	}
	return ""
}

func (m *ConsensusConfig) GetPrivateKey() string {
	if m != nil {
		return m.PrivateKey
	}
	return ""
}

type NodeConfig struct {
	Port                 uint32   `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
	Seed                 []string `protobuf:"bytes,2,rep,name=seed,proto3" json:"seed,omitempty"`
	DbPath               string   `protobuf:"bytes,3,opt,name=db_path,json=dbPath,proto3" json:"db_path,omitempty"`
	RpcPort              uint32   `protobuf:"varint,4,opt,name=rpc_port,json=rpcPort,proto3" json:"rpc_port,omitempty"`
	KeyPath              string   `protobuf:"bytes,5,opt,name=key_path,json=keyPath,proto3" json:"key_path,omitempty"`
	TxPoolLimit          uint32   `protobuf:"varint,6,opt,name=tx_pool_limit,json=txPoolLimit,proto3" json:"tx_pool_limit,omitempty"`
	NodeAddress          string   `protobuf:"bytes,7,opt,name=node_address,json=nodeAddress,proto3" json:"node_address,omitempty"`
	GenesisPath          string   `protobuf:"bytes,8,opt,name=genesis_path,json=genesisPath,proto3" json:"genesis_path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeConfig) Reset()         { *m = NodeConfig{} }
func (m *NodeConfig) String() string { return proto.CompactTextString(m) }
func (*NodeConfig) ProtoMessage()    {}
func (*NodeConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bcc5bb4ac8bd8b6e, []int{2}
}
func (m *NodeConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeConfig.Unmarshal(m, b)
}
func (m *NodeConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeConfig.Marshal(b, m, deterministic)
}
func (dst *NodeConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeConfig.Merge(dst, src)
}
func (m *NodeConfig) XXX_Size() int {
	return xxx_messageInfo_NodeConfig.Size(m)
}
func (m *NodeConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeConfig.DiscardUnknown(m)
}

var xxx_messageInfo_NodeConfig proto.InternalMessageInfo

func (m *NodeConfig) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *NodeConfig) GetSeed() []string {
	if m != nil {
		return m.Seed
	}
	return nil
}

func (m *NodeConfig) GetDbPath() string {
	if m != nil {
		return m.DbPath
	}
	return ""
}

func (m *NodeConfig) GetRpcPort() uint32 {
	if m != nil {
		return m.RpcPort
	}
	return 0
}

func (m *NodeConfig) GetKeyPath() string {
	if m != nil {
		return m.KeyPath
	}
	return ""
}

func (m *NodeConfig) GetTxPoolLimit() uint32 {
	if m != nil {
		return m.TxPoolLimit
	}
	return 0
}

func (m *NodeConfig) GetNodeAddress() string {
	if m != nil {
		return m.NodeAddress
	}
	return ""
}

func (m *NodeConfig) GetGenesisPath() string {
	if m != nil {
		return m.GenesisPath
	}
	return ""
}

type DynastyConfig struct {
	Producers            []string `protobuf:"bytes,1,rep,name=producers,proto3" json:"producers,omitempty"`
	MaxProducers         uint32   `protobuf:"varint,2,opt,name=max_producers,json=maxProducers,proto3" json:"max_producers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DynastyConfig) Reset()         { *m = DynastyConfig{} }
func (m *DynastyConfig) String() string { return proto.CompactTextString(m) }
func (*DynastyConfig) ProtoMessage()    {}
func (*DynastyConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bcc5bb4ac8bd8b6e, []int{3}
}
func (m *DynastyConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DynastyConfig.Unmarshal(m, b)
}
func (m *DynastyConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DynastyConfig.Marshal(b, m, deterministic)
}
func (dst *DynastyConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DynastyConfig.Merge(dst, src)
}
func (m *DynastyConfig) XXX_Size() int {
	return xxx_messageInfo_DynastyConfig.Size(m)
}
func (m *DynastyConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_DynastyConfig.DiscardUnknown(m)
}

var xxx_messageInfo_DynastyConfig proto.InternalMessageInfo

func (m *DynastyConfig) GetProducers() []string {
	if m != nil {
		return m.Producers
	}
	return nil
}

func (m *DynastyConfig) GetMaxProducers() uint32 {
	if m != nil {
		return m.MaxProducers
	}
	return 0
}

type CliConfig struct {
	Port                 uint32   `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CliConfig) Reset()         { *m = CliConfig{} }
func (m *CliConfig) String() string { return proto.CompactTextString(m) }
func (*CliConfig) ProtoMessage()    {}
func (*CliConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bcc5bb4ac8bd8b6e, []int{4}
}
func (m *CliConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CliConfig.Unmarshal(m, b)
}
func (m *CliConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CliConfig.Marshal(b, m, deterministic)
}
func (dst *CliConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CliConfig.Merge(dst, src)
}
func (m *CliConfig) XXX_Size() int {
	return xxx_messageInfo_CliConfig.Size(m)
}
func (m *CliConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_CliConfig.DiscardUnknown(m)
}

var xxx_messageInfo_CliConfig proto.InternalMessageInfo

func (m *CliConfig) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *CliConfig) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func init() {
	proto.RegisterType((*Config)(nil), "configpb.Config")
	proto.RegisterType((*ConsensusConfig)(nil), "configpb.ConsensusConfig")
	proto.RegisterType((*NodeConfig)(nil), "configpb.NodeConfig")
	proto.RegisterType((*DynastyConfig)(nil), "configpb.DynastyConfig")
	proto.RegisterType((*CliConfig)(nil), "configpb.CliConfig")
}

func init() {
	proto.RegisterFile("github.com/dappley/go-dappley/config/pb/config.proto", fileDescriptor_config_bcc5bb4ac8bd8b6e)
}

var fileDescriptor_config_bcc5bb4ac8bd8b6e = []byte{
	// 397 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x92, 0xcd, 0x8e, 0x9b, 0x30,
	0x10, 0xc7, 0x45, 0x92, 0xf2, 0x31, 0x04, 0xa5, 0xb2, 0x2a, 0x35, 0xa9, 0x2a, 0xb5, 0xa5, 0x97,
	0x5e, 0x4a, 0xa4, 0x7e, 0x9c, 0x7a, 0xaa, 0x92, 0x5b, 0xab, 0x0a, 0x71, 0xe9, 0x11, 0x01, 0x76,
	0x89, 0x55, 0x82, 0x2d, 0xdb, 0xd9, 0x0d, 0x0f, 0xb0, 0x8f, 0xbb, 0xef, 0xb0, 0x66, 0x80, 0x44,
	0xbb, 0x87, 0xbd, 0xcd, 0xfc, 0x98, 0xff, 0x7f, 0x3e, 0x0c, 0x7c, 0xab, 0xb9, 0x39, 0x9c, 0xca,
	0xa4, 0x12, 0xc7, 0x2d, 0x2d, 0xa4, 0x6c, 0x58, 0xb7, 0xad, 0xc5, 0xe7, 0x29, 0xac, 0x44, 0xfb,
	0x8f, 0xd7, 0x5b, 0x59, 0x8e, 0x51, 0x22, 0x95, 0x30, 0x82, 0xf8, 0x43, 0x26, 0xcb, 0xf8, 0xce,
	0x01, 0x77, 0x87, 0x09, 0xd9, 0xc3, 0x4b, 0x8b, 0x35, 0x6b, 0xf5, 0x49, 0xe7, 0x43, 0xc1, 0xda,
	0x79, 0xef, 0x7c, 0x0a, 0xbf, 0x6c, 0x92, 0xa9, 0x3e, 0xd9, 0x4d, 0x15, 0x83, 0x28, 0x5b, 0x55,
	0x8f, 0x01, 0xf9, 0x0e, 0x61, 0x2b, 0x28, 0x9b, 0x0c, 0x66, 0x68, 0xf0, 0xea, 0x6a, 0xf0, 0xc7,
	0x7e, 0x1c, 0xb5, 0xd0, 0x5e, 0xe2, 0xf8, 0x2f, 0xac, 0x9e, 0x58, 0x93, 0x8f, 0x10, 0x1d, 0x79,
	0xcb, 0x54, 0x5e, 0x50, 0xaa, 0x98, 0xd6, 0x38, 0x4c, 0x90, 0x2d, 0x11, 0xfe, 0x1c, 0x18, 0x79,
	0x07, 0xa1, 0x54, 0xfc, 0xa6, 0x30, 0x2c, 0xff, 0xcf, 0x3a, 0x6c, 0x17, 0x64, 0x30, 0xa2, 0x5f,
	0xac, 0x8b, 0xef, 0x1d, 0x80, 0x6b, 0x4f, 0x42, 0x60, 0x21, 0x85, 0x32, 0xe8, 0x15, 0x65, 0x18,
	0xf7, 0x4c, 0x33, 0x46, 0xad, 0x78, 0x6e, 0xc5, 0x18, 0x93, 0xd7, 0xe0, 0xd1, 0x32, 0x97, 0x85,
	0x39, 0xac, 0xe7, 0xe8, 0xe9, 0xd2, 0x32, 0xb5, 0x19, 0xd9, 0x80, 0xaf, 0x64, 0x95, 0xa3, 0xc9,
	0x02, 0x4d, 0x3c, 0x9b, 0xa7, 0xbd, 0x8f, 0xfd, 0x64, 0x67, 0x18, 0x44, 0x2f, 0x50, 0xe4, 0xd9,
	0x1c, 0x55, 0x31, 0x44, 0xe6, 0x6c, 0x45, 0xa2, 0xc9, 0x1b, 0x7e, 0xe4, 0x66, 0xed, 0xa2, 0x34,
	0x34, 0xe7, 0xd4, 0xb2, 0xdf, 0x3d, 0x22, 0x1f, 0x60, 0x89, 0x97, 0x9b, 0xd6, 0xf5, 0xd0, 0x02,
	0xaf, 0x39, 0x6d, 0x6b, 0x4b, 0x6a, 0xd6, 0x32, 0xcd, 0xf5, 0xd0, 0xc5, 0x1f, 0x4a, 0x46, 0xd6,
	0x77, 0x8a, 0x33, 0x88, 0xf6, 0x5d, 0x5b, 0x68, 0xd3, 0x8d, 0x1b, 0xbf, 0x85, 0xc0, 0x3e, 0x3a,
	0x3d, 0x55, 0x4c, 0xf5, 0x27, 0xec, 0x57, 0xbc, 0x02, 0x3c, 0x72, 0x61, 0x27, 0xbb, 0x54, 0xcc,
	0x70, 0xb0, 0xa5, 0x85, 0xe9, 0xc4, 0xe2, 0x1f, 0x10, 0xec, 0x1a, 0xfe, 0xcc, 0x05, 0xdf, 0x80,
	0x2f, 0x0b, 0xad, 0x6f, 0x85, 0xa2, 0xe3, 0x13, 0x5c, 0xf2, 0xd2, 0xc5, 0x5f, 0xee, 0xeb, 0x43,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xb7, 0x1c, 0x13, 0x35, 0xaa, 0x02, 0x00, 0x00,
}

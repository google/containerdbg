// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.12
// source: summary.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ContainerAnalysisSummary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MissingFiles       []string                                `protobuf:"bytes,1,rep,name=missing_files,json=missingFiles,proto3" json:"missing_files,omitempty"`
	MoveFailures       []*ContainerAnalysisSummary_MoveFailure `protobuf:"bytes,2,rep,name=move_failures,json=moveFailures,proto3" json:"move_failures,omitempty"`
	MissingLibraries   []string                                `protobuf:"bytes,3,rep,name=missing_libraries,json=missingLibraries,proto3" json:"missing_libraries,omitempty"`
	ConnectionFailures []*ContainerAnalysisSummary_Connection  `protobuf:"bytes,4,rep,name=connection_failures,json=connectionFailures,proto3" json:"connection_failures,omitempty"`
	DnsFailures        []*ContainerAnalysisSummary_DnsFailure  `protobuf:"bytes,5,rep,name=dns_failures,json=dnsFailures,proto3" json:"dns_failures,omitempty"`
	StaticIps          []*ContainerAnalysisSummary_Connection  `protobuf:"bytes,6,rep,name=static_ips,json=staticIps,proto3" json:"static_ips,omitempty"`
}

func (x *ContainerAnalysisSummary) Reset() {
	*x = ContainerAnalysisSummary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_summary_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContainerAnalysisSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerAnalysisSummary) ProtoMessage() {}

func (x *ContainerAnalysisSummary) ProtoReflect() protoreflect.Message {
	mi := &file_summary_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerAnalysisSummary.ProtoReflect.Descriptor instead.
func (*ContainerAnalysisSummary) Descriptor() ([]byte, []int) {
	return file_summary_proto_rawDescGZIP(), []int{0}
}

func (x *ContainerAnalysisSummary) GetMissingFiles() []string {
	if x != nil {
		return x.MissingFiles
	}
	return nil
}

func (x *ContainerAnalysisSummary) GetMoveFailures() []*ContainerAnalysisSummary_MoveFailure {
	if x != nil {
		return x.MoveFailures
	}
	return nil
}

func (x *ContainerAnalysisSummary) GetMissingLibraries() []string {
	if x != nil {
		return x.MissingLibraries
	}
	return nil
}

func (x *ContainerAnalysisSummary) GetConnectionFailures() []*ContainerAnalysisSummary_Connection {
	if x != nil {
		return x.ConnectionFailures
	}
	return nil
}

func (x *ContainerAnalysisSummary) GetDnsFailures() []*ContainerAnalysisSummary_DnsFailure {
	if x != nil {
		return x.DnsFailures
	}
	return nil
}

func (x *ContainerAnalysisSummary) GetStaticIps() []*ContainerAnalysisSummary_Connection {
	if x != nil {
		return x.StaticIps
	}
	return nil
}

type AnalysisSummary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContainerSummaries []*AnalysisSummary_ContainerSummaryTuple `protobuf:"bytes,1,rep,name=container_summaries,json=containerSummaries,proto3" json:"container_summaries,omitempty"`
}

func (x *AnalysisSummary) Reset() {
	*x = AnalysisSummary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_summary_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalysisSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalysisSummary) ProtoMessage() {}

func (x *AnalysisSummary) ProtoReflect() protoreflect.Message {
	mi := &file_summary_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalysisSummary.ProtoReflect.Descriptor instead.
func (*AnalysisSummary) Descriptor() ([]byte, []int) {
	return file_summary_proto_rawDescGZIP(), []int{1}
}

func (x *AnalysisSummary) GetContainerSummaries() []*AnalysisSummary_ContainerSummaryTuple {
	if x != nil {
		return x.ContainerSummaries
	}
	return nil
}

type ContainerAnalysisSummary_MoveFailure struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source      string `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Destination string `protobuf:"bytes,2,opt,name=destination,proto3" json:"destination,omitempty"`
}

func (x *ContainerAnalysisSummary_MoveFailure) Reset() {
	*x = ContainerAnalysisSummary_MoveFailure{}
	if protoimpl.UnsafeEnabled {
		mi := &file_summary_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContainerAnalysisSummary_MoveFailure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerAnalysisSummary_MoveFailure) ProtoMessage() {}

func (x *ContainerAnalysisSummary_MoveFailure) ProtoReflect() protoreflect.Message {
	mi := &file_summary_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerAnalysisSummary_MoveFailure.ProtoReflect.Descriptor instead.
func (*ContainerAnalysisSummary_MoveFailure) Descriptor() ([]byte, []int) {
	return file_summary_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ContainerAnalysisSummary_MoveFailure) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *ContainerAnalysisSummary_MoveFailure) GetDestination() string {
	if x != nil {
		return x.Destination
	}
	return ""
}

type ContainerAnalysisSummary_Connection struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetFqdn string `protobuf:"bytes,1,opt,name=target_fqdn,json=targetFqdn,proto3" json:"target_fqdn,omitempty"`
	TargetIp   string `protobuf:"bytes,2,opt,name=target_ip,json=targetIp,proto3" json:"target_ip,omitempty"`
	Port       int32  `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *ContainerAnalysisSummary_Connection) Reset() {
	*x = ContainerAnalysisSummary_Connection{}
	if protoimpl.UnsafeEnabled {
		mi := &file_summary_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContainerAnalysisSummary_Connection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerAnalysisSummary_Connection) ProtoMessage() {}

func (x *ContainerAnalysisSummary_Connection) ProtoReflect() protoreflect.Message {
	mi := &file_summary_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerAnalysisSummary_Connection.ProtoReflect.Descriptor instead.
func (*ContainerAnalysisSummary_Connection) Descriptor() ([]byte, []int) {
	return file_summary_proto_rawDescGZIP(), []int{0, 1}
}

func (x *ContainerAnalysisSummary_Connection) GetTargetFqdn() string {
	if x != nil {
		return x.TargetFqdn
	}
	return ""
}

func (x *ContainerAnalysisSummary_Connection) GetTargetIp() string {
	if x != nil {
		return x.TargetIp
	}
	return ""
}

func (x *ContainerAnalysisSummary_Connection) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type ContainerAnalysisSummary_DnsFailure struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Query string         `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	Error *DnsQueryError `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *ContainerAnalysisSummary_DnsFailure) Reset() {
	*x = ContainerAnalysisSummary_DnsFailure{}
	if protoimpl.UnsafeEnabled {
		mi := &file_summary_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContainerAnalysisSummary_DnsFailure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerAnalysisSummary_DnsFailure) ProtoMessage() {}

func (x *ContainerAnalysisSummary_DnsFailure) ProtoReflect() protoreflect.Message {
	mi := &file_summary_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerAnalysisSummary_DnsFailure.ProtoReflect.Descriptor instead.
func (*ContainerAnalysisSummary_DnsFailure) Descriptor() ([]byte, []int) {
	return file_summary_proto_rawDescGZIP(), []int{0, 2}
}

func (x *ContainerAnalysisSummary_DnsFailure) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *ContainerAnalysisSummary_DnsFailure) GetError() *DnsQueryError {
	if x != nil {
		return x.Error
	}
	return nil
}

type AnalysisSummary_ContainerSummaryTuple struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source  *SourceId                 `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Summary *ContainerAnalysisSummary `protobuf:"bytes,2,opt,name=summary,proto3" json:"summary,omitempty"`
}

func (x *AnalysisSummary_ContainerSummaryTuple) Reset() {
	*x = AnalysisSummary_ContainerSummaryTuple{}
	if protoimpl.UnsafeEnabled {
		mi := &file_summary_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalysisSummary_ContainerSummaryTuple) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalysisSummary_ContainerSummaryTuple) ProtoMessage() {}

func (x *AnalysisSummary_ContainerSummaryTuple) ProtoReflect() protoreflect.Message {
	mi := &file_summary_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalysisSummary_ContainerSummaryTuple.ProtoReflect.Descriptor instead.
func (*AnalysisSummary_ContainerSummaryTuple) Descriptor() ([]byte, []int) {
	return file_summary_proto_rawDescGZIP(), []int{1, 0}
}

func (x *AnalysisSummary_ContainerSummaryTuple) GetSource() *SourceId {
	if x != nil {
		return x.Source
	}
	return nil
}

func (x *AnalysisSummary_ContainerSummaryTuple) GetSummary() *ContainerAnalysisSummary {
	if x != nil {
		return x.Summary
	}
	return nil
}

var File_summary_proto protoreflect.FileDescriptor

var file_summary_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x90, 0x05, 0x0a,
	0x18, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73,
	0x69, 0x73, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x6d, 0x69, 0x73,
	0x73, 0x69, 0x6e, 0x67, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0c, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x4a,
	0x0a, 0x0d, 0x6d, 0x6f, 0x76, 0x65, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x2e, 0x4d, 0x6f, 0x76, 0x65, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x52, 0x0c, 0x6d, 0x6f,
	0x76, 0x65, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6e, 0x67, 0x5f, 0x6c, 0x69, 0x62, 0x72, 0x61, 0x72, 0x69, 0x65, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x10, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x4c, 0x69,
	0x62, 0x72, 0x61, 0x72, 0x69, 0x65, 0x73, 0x12, 0x55, 0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x2e,
	0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x12, 0x63, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x12, 0x47,
	0x0a, 0x0c, 0x64, 0x6e, 0x73, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x2e,
	0x44, 0x6e, 0x73, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x52, 0x0b, 0x64, 0x6e, 0x73, 0x46,
	0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x12, 0x43, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x69,
	0x63, 0x5f, 0x69, 0x70, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x43, 0x6f,
	0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x53,
	0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x49, 0x70, 0x73, 0x1a, 0x47, 0x0a, 0x0b,
	0x4d, 0x6f, 0x76, 0x65, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x5e, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x66, 0x71,
	0x64, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x46, 0x71, 0x64, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x69,
	0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49,
	0x70, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x70, 0x6f, 0x72, 0x74, 0x1a, 0x48, 0x0a, 0x0a, 0x44, 0x6e, 0x73, 0x46, 0x61, 0x69, 0x6c,
	0x75, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x12, 0x24, 0x0a, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x44, 0x6e, 0x73, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22,
	0xdb, 0x01, 0x0a, 0x0f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x53, 0x75, 0x6d, 0x6d,
	0x61, 0x72, 0x79, 0x12, 0x57, 0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x5f, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x26, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x53, 0x75, 0x6d, 0x6d, 0x61,
	0x72, 0x79, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x53, 0x75, 0x6d, 0x6d,
	0x61, 0x72, 0x79, 0x54, 0x75, 0x70, 0x6c, 0x65, 0x52, 0x12, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69,
	0x6e, 0x65, 0x72, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x69, 0x65, 0x73, 0x1a, 0x6f, 0x0a, 0x15,
	0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x54, 0x75, 0x70, 0x6c, 0x65, 0x12, 0x21, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64,
	0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x33, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d,
	0x61, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x69, 0x6e, 0x65, 0x72, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x53, 0x75, 0x6d,
	0x6d, 0x61, 0x72, 0x79, 0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x42, 0x26, 0x5a,
	0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x64, 0x62, 0x67, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_summary_proto_rawDescOnce sync.Once
	file_summary_proto_rawDescData = file_summary_proto_rawDesc
)

func file_summary_proto_rawDescGZIP() []byte {
	file_summary_proto_rawDescOnce.Do(func() {
		file_summary_proto_rawDescData = protoimpl.X.CompressGZIP(file_summary_proto_rawDescData)
	})
	return file_summary_proto_rawDescData
}

var file_summary_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_summary_proto_goTypes = []interface{}{
	(*ContainerAnalysisSummary)(nil),              // 0: ContainerAnalysisSummary
	(*AnalysisSummary)(nil),                       // 1: AnalysisSummary
	(*ContainerAnalysisSummary_MoveFailure)(nil),  // 2: ContainerAnalysisSummary.MoveFailure
	(*ContainerAnalysisSummary_Connection)(nil),   // 3: ContainerAnalysisSummary.Connection
	(*ContainerAnalysisSummary_DnsFailure)(nil),   // 4: ContainerAnalysisSummary.DnsFailure
	(*AnalysisSummary_ContainerSummaryTuple)(nil), // 5: AnalysisSummary.ContainerSummaryTuple
	(*DnsQueryError)(nil),                         // 6: DnsQueryError
	(*SourceId)(nil),                              // 7: SourceId
}
var file_summary_proto_depIdxs = []int32{
	2, // 0: ContainerAnalysisSummary.move_failures:type_name -> ContainerAnalysisSummary.MoveFailure
	3, // 1: ContainerAnalysisSummary.connection_failures:type_name -> ContainerAnalysisSummary.Connection
	4, // 2: ContainerAnalysisSummary.dns_failures:type_name -> ContainerAnalysisSummary.DnsFailure
	3, // 3: ContainerAnalysisSummary.static_ips:type_name -> ContainerAnalysisSummary.Connection
	5, // 4: AnalysisSummary.container_summaries:type_name -> AnalysisSummary.ContainerSummaryTuple
	6, // 5: ContainerAnalysisSummary.DnsFailure.error:type_name -> DnsQueryError
	7, // 6: AnalysisSummary.ContainerSummaryTuple.source:type_name -> SourceId
	0, // 7: AnalysisSummary.ContainerSummaryTuple.summary:type_name -> ContainerAnalysisSummary
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_summary_proto_init() }
func file_summary_proto_init() {
	if File_summary_proto != nil {
		return
	}
	file_event_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_summary_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContainerAnalysisSummary); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_summary_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalysisSummary); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_summary_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContainerAnalysisSummary_MoveFailure); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_summary_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContainerAnalysisSummary_Connection); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_summary_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContainerAnalysisSummary_DnsFailure); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_summary_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalysisSummary_ContainerSummaryTuple); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_summary_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_summary_proto_goTypes,
		DependencyIndexes: file_summary_proto_depIdxs,
		MessageInfos:      file_summary_proto_msgTypes,
	}.Build()
	File_summary_proto = out.File
	file_summary_proto_rawDesc = nil
	file_summary_proto_goTypes = nil
	file_summary_proto_depIdxs = nil
}

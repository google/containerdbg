// Copyright 2022 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
// source: node_api.proto

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

type MonitorPodResponse_ErrorCode int32

const (
	MonitorPodResponse_OK      MonitorPodResponse_ErrorCode = 0
	MonitorPodResponse_UNKNOWN MonitorPodResponse_ErrorCode = 1
)

// Enum value maps for MonitorPodResponse_ErrorCode.
var (
	MonitorPodResponse_ErrorCode_name = map[int32]string{
		0: "OK",
		1: "UNKNOWN",
	}
	MonitorPodResponse_ErrorCode_value = map[string]int32{
		"OK":      0,
		"UNKNOWN": 1,
	}
)

func (x MonitorPodResponse_ErrorCode) Enum() *MonitorPodResponse_ErrorCode {
	p := new(MonitorPodResponse_ErrorCode)
	*p = x
	return p
}

func (x MonitorPodResponse_ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MonitorPodResponse_ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_node_api_proto_enumTypes[0].Descriptor()
}

func (MonitorPodResponse_ErrorCode) Type() protoreflect.EnumType {
	return &file_node_api_proto_enumTypes[0]
}

func (x MonitorPodResponse_ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MonitorPodResponse_ErrorCode.Descriptor instead.
func (MonitorPodResponse_ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_node_api_proto_rawDescGZIP(), []int{1, 0}
}

type MonitorPodRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id *SourceId `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// 0 means it was not specified
	Netns uint64 `protobuf:"varint,2,opt,name=netns,proto3" json:"netns,omitempty"`
}

func (x *MonitorPodRequest) Reset() {
	*x = MonitorPodRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MonitorPodRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MonitorPodRequest) ProtoMessage() {}

func (x *MonitorPodRequest) ProtoReflect() protoreflect.Message {
	mi := &file_node_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MonitorPodRequest.ProtoReflect.Descriptor instead.
func (*MonitorPodRequest) Descriptor() ([]byte, []int) {
	return file_node_api_proto_rawDescGZIP(), []int{0}
}

func (x *MonitorPodRequest) GetId() *SourceId {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *MonitorPodRequest) GetNetns() uint64 {
	if x != nil {
		return x.Netns
	}
	return 0
}

type MonitorPodResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code MonitorPodResponse_ErrorCode `protobuf:"varint,1,opt,name=code,proto3,enum=MonitorPodResponse_ErrorCode" json:"code,omitempty"`
}

func (x *MonitorPodResponse) Reset() {
	*x = MonitorPodResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MonitorPodResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MonitorPodResponse) ProtoMessage() {}

func (x *MonitorPodResponse) ProtoReflect() protoreflect.Message {
	mi := &file_node_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MonitorPodResponse.ProtoReflect.Descriptor instead.
func (*MonitorPodResponse) Descriptor() ([]byte, []int) {
	return file_node_api_proto_rawDescGZIP(), []int{1}
}

func (x *MonitorPodResponse) GetCode() MonitorPodResponse_ErrorCode {
	if x != nil {
		return x.Code
	}
	return MonitorPodResponse_OK
}

type ReportDnsQueryResultRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       *SourceId `protobuf:"bytes,4,opt,name=id,proto3" json:"id,omitempty"`
	DnsQuery string    `protobuf:"bytes,1,opt,name=dns_query,json=dnsQuery,proto3" json:"dns_query,omitempty"`
	// Types that are assignable to Result:
	//
	//	*ReportDnsQueryResultRequest_ReturnedIp
	//	*ReportDnsQueryResultRequest_Error
	Result isReportDnsQueryResultRequest_Result `protobuf_oneof:"result"`
}

func (x *ReportDnsQueryResultRequest) Reset() {
	*x = ReportDnsQueryResultRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportDnsQueryResultRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportDnsQueryResultRequest) ProtoMessage() {}

func (x *ReportDnsQueryResultRequest) ProtoReflect() protoreflect.Message {
	mi := &file_node_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportDnsQueryResultRequest.ProtoReflect.Descriptor instead.
func (*ReportDnsQueryResultRequest) Descriptor() ([]byte, []int) {
	return file_node_api_proto_rawDescGZIP(), []int{2}
}

func (x *ReportDnsQueryResultRequest) GetId() *SourceId {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ReportDnsQueryResultRequest) GetDnsQuery() string {
	if x != nil {
		return x.DnsQuery
	}
	return ""
}

func (m *ReportDnsQueryResultRequest) GetResult() isReportDnsQueryResultRequest_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (x *ReportDnsQueryResultRequest) GetReturnedIp() string {
	if x, ok := x.GetResult().(*ReportDnsQueryResultRequest_ReturnedIp); ok {
		return x.ReturnedIp
	}
	return ""
}

func (x *ReportDnsQueryResultRequest) GetError() *DnsQueryError {
	if x, ok := x.GetResult().(*ReportDnsQueryResultRequest_Error); ok {
		return x.Error
	}
	return nil
}

type isReportDnsQueryResultRequest_Result interface {
	isReportDnsQueryResultRequest_Result()
}

type ReportDnsQueryResultRequest_ReturnedIp struct {
	ReturnedIp string `protobuf:"bytes,2,opt,name=returned_ip,json=returnedIp,proto3,oneof"`
}

type ReportDnsQueryResultRequest_Error struct {
	Error *DnsQueryError `protobuf:"bytes,3,opt,name=error,proto3,oneof"`
}

func (*ReportDnsQueryResultRequest_ReturnedIp) isReportDnsQueryResultRequest_Result() {}

func (*ReportDnsQueryResultRequest_Error) isReportDnsQueryResultRequest_Result() {}

type ReportDnsQueryResultResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReportDnsQueryResultResponse) Reset() {
	*x = ReportDnsQueryResultResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportDnsQueryResultResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportDnsQueryResultResponse) ProtoMessage() {}

func (x *ReportDnsQueryResultResponse) ProtoReflect() protoreflect.Message {
	mi := &file_node_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportDnsQueryResultResponse.ProtoReflect.Descriptor instead.
func (*ReportDnsQueryResultResponse) Descriptor() ([]byte, []int) {
	return file_node_api_proto_rawDescGZIP(), []int{3}
}

type ReportDnsSearchValuesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Search []string  `protobuf:"bytes,1,rep,name=search,proto3" json:"search,omitempty"`
	Id     *SourceId `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ReportDnsSearchValuesRequest) Reset() {
	*x = ReportDnsSearchValuesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportDnsSearchValuesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportDnsSearchValuesRequest) ProtoMessage() {}

func (x *ReportDnsSearchValuesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_node_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportDnsSearchValuesRequest.ProtoReflect.Descriptor instead.
func (*ReportDnsSearchValuesRequest) Descriptor() ([]byte, []int) {
	return file_node_api_proto_rawDescGZIP(), []int{4}
}

func (x *ReportDnsSearchValuesRequest) GetSearch() []string {
	if x != nil {
		return x.Search
	}
	return nil
}

func (x *ReportDnsSearchValuesRequest) GetId() *SourceId {
	if x != nil {
		return x.Id
	}
	return nil
}

type ReportDnsSearchValuesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReportDnsSearchValuesResponse) Reset() {
	*x = ReportDnsSearchValuesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportDnsSearchValuesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportDnsSearchValuesResponse) ProtoMessage() {}

func (x *ReportDnsSearchValuesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_node_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportDnsSearchValuesResponse.ProtoReflect.Descriptor instead.
func (*ReportDnsSearchValuesResponse) Descriptor() ([]byte, []int) {
	return file_node_api_proto_rawDescGZIP(), []int{5}
}

var File_node_api_proto protoreflect.FileDescriptor

var file_node_api_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x44, 0x0a,
	0x11, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x50, 0x6f, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x19, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09,
	0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x6e, 0x65, 0x74, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x6e, 0x65,
	0x74, 0x6e, 0x73, 0x22, 0x69, 0x0a, 0x12, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x50, 0x6f,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f,
	0x72, 0x50, 0x6f, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x20, 0x0a, 0x09,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10,
	0x00, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x01, 0x22, 0xaa,
	0x01, 0x0a, 0x1b, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x44, 0x6e, 0x73, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x53, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x49, 0x64, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x6e, 0x73,
	0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x6e,
	0x73, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x21, 0x0a, 0x0b, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e,
	0x65, 0x64, 0x5f, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0a, 0x72,
	0x65, 0x74, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x49, 0x70, 0x12, 0x26, 0x0a, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x44, 0x6e, 0x73, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x42, 0x08, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x1e, 0x0a, 0x1c, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x44, 0x6e, 0x73, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x51, 0x0a, 0x1c, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x44, 0x6e, 0x73, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x12, 0x19, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x09, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x52, 0x02, 0x69, 0x64, 0x22, 0x1f,
	0x0a, 0x1d, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x44, 0x6e, 0x73, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32,
	0xee, 0x01, 0x0a, 0x11, 0x4e, 0x6f, 0x64, 0x65, 0x44, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x07, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72,
	0x12, 0x12, 0x2e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x50, 0x6f, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x50, 0x6f,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0e, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x44, 0x6e, 0x73, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x1c, 0x2e, 0x52, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x44, 0x6e, 0x73, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x44, 0x6e, 0x73, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x56, 0x0a, 0x15, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x44, 0x6e, 0x73, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x73, 0x12, 0x1d, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x44, 0x6e, 0x73, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1e, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x44, 0x6e, 0x73, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x26, 0x5a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x64,
	0x62, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_node_api_proto_rawDescOnce sync.Once
	file_node_api_proto_rawDescData = file_node_api_proto_rawDesc
)

func file_node_api_proto_rawDescGZIP() []byte {
	file_node_api_proto_rawDescOnce.Do(func() {
		file_node_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_api_proto_rawDescData)
	})
	return file_node_api_proto_rawDescData
}

var file_node_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_node_api_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_node_api_proto_goTypes = []interface{}{
	(MonitorPodResponse_ErrorCode)(0),     // 0: MonitorPodResponse.ErrorCode
	(*MonitorPodRequest)(nil),             // 1: MonitorPodRequest
	(*MonitorPodResponse)(nil),            // 2: MonitorPodResponse
	(*ReportDnsQueryResultRequest)(nil),   // 3: ReportDnsQueryResultRequest
	(*ReportDnsQueryResultResponse)(nil),  // 4: ReportDnsQueryResultResponse
	(*ReportDnsSearchValuesRequest)(nil),  // 5: ReportDnsSearchValuesRequest
	(*ReportDnsSearchValuesResponse)(nil), // 6: ReportDnsSearchValuesResponse
	(*SourceId)(nil),                      // 7: SourceId
	(*DnsQueryError)(nil),                 // 8: DnsQueryError
}
var file_node_api_proto_depIdxs = []int32{
	7, // 0: MonitorPodRequest.id:type_name -> SourceId
	0, // 1: MonitorPodResponse.code:type_name -> MonitorPodResponse.ErrorCode
	7, // 2: ReportDnsQueryResultRequest.id:type_name -> SourceId
	8, // 3: ReportDnsQueryResultRequest.error:type_name -> DnsQueryError
	7, // 4: ReportDnsSearchValuesRequest.id:type_name -> SourceId
	1, // 5: NodeDaemonService.Monitor:input_type -> MonitorPodRequest
	3, // 6: NodeDaemonService.ReportDnsQuery:input_type -> ReportDnsQueryResultRequest
	5, // 7: NodeDaemonService.ReportDnsSearchValues:input_type -> ReportDnsSearchValuesRequest
	2, // 8: NodeDaemonService.Monitor:output_type -> MonitorPodResponse
	4, // 9: NodeDaemonService.ReportDnsQuery:output_type -> ReportDnsQueryResultResponse
	6, // 10: NodeDaemonService.ReportDnsSearchValues:output_type -> ReportDnsSearchValuesResponse
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_node_api_proto_init() }
func file_node_api_proto_init() {
	if File_node_api_proto != nil {
		return
	}
	file_event_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_node_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MonitorPodRequest); i {
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
		file_node_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MonitorPodResponse); i {
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
		file_node_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportDnsQueryResultRequest); i {
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
		file_node_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportDnsQueryResultResponse); i {
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
		file_node_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportDnsSearchValuesRequest); i {
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
		file_node_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportDnsSearchValuesResponse); i {
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
	file_node_api_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*ReportDnsQueryResultRequest_ReturnedIp)(nil),
		(*ReportDnsQueryResultRequest_Error)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_node_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_node_api_proto_goTypes,
		DependencyIndexes: file_node_api_proto_depIdxs,
		EnumInfos:         file_node_api_proto_enumTypes,
		MessageInfos:      file_node_api_proto_msgTypes,
	}.Build()
	File_node_api_proto = out.File
	file_node_api_proto_rawDesc = nil
	file_node_api_proto_goTypes = nil
	file_node_api_proto_depIdxs = nil
}

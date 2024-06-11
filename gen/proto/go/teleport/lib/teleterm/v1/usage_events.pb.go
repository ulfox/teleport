//
// Teleport
// Copyright (C) 2023  Gravitational, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: teleport/lib/teleterm/v1/usage_events.proto

package teletermv1

import (
	v1alpha "github.com/gravitational/teleport/gen/proto/go/prehog/v1alpha"
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

type ReportUsageEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthClusterId string                             `protobuf:"bytes,1,opt,name=auth_cluster_id,json=authClusterId,proto3" json:"auth_cluster_id,omitempty"`
	PrehogReq     *v1alpha.SubmitConnectEventRequest `protobuf:"bytes,2,opt,name=prehog_req,json=prehogReq,proto3" json:"prehog_req,omitempty"`
}

func (x *ReportUsageEventRequest) Reset() {
	*x = ReportUsageEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teleport_lib_teleterm_v1_usage_events_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportUsageEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportUsageEventRequest) ProtoMessage() {}

func (x *ReportUsageEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_lib_teleterm_v1_usage_events_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportUsageEventRequest.ProtoReflect.Descriptor instead.
func (*ReportUsageEventRequest) Descriptor() ([]byte, []int) {
	return file_teleport_lib_teleterm_v1_usage_events_proto_rawDescGZIP(), []int{0}
}

func (x *ReportUsageEventRequest) GetAuthClusterId() string {
	if x != nil {
		return x.AuthClusterId
	}
	return ""
}

func (x *ReportUsageEventRequest) GetPrehogReq() *v1alpha.SubmitConnectEventRequest {
	if x != nil {
		return x.PrehogReq
	}
	return nil
}

var File_teleport_lib_teleterm_v1_usage_events_proto protoreflect.FileDescriptor

var file_teleport_lib_teleterm_v1_usage_events_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x6c, 0x69, 0x62, 0x2f, 0x74,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x72, 0x6d, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x61, 0x67, 0x65,
	0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x74,
	0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x6c, 0x69, 0x62, 0x2e, 0x74, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x72, 0x6d, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x70, 0x72, 0x65, 0x68, 0x6f, 0x67, 0x2f,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8b, 0x01, 0x0a, 0x17, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x55, 0x73, 0x61, 0x67, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x75, 0x74, 0x68,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x48, 0x0a, 0x0a, 0x70, 0x72, 0x65,
	0x68, 0x6f, 0x67, 0x5f, 0x72, 0x65, 0x71, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e,
	0x70, 0x72, 0x65, 0x68, 0x6f, 0x67, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x53,
	0x75, 0x62, 0x6d, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x09, 0x70, 0x72, 0x65, 0x68, 0x6f, 0x67,
	0x52, 0x65, 0x71, 0x42, 0x54, 0x5a, 0x52, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x67, 0x72, 0x61, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x2f,
	0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x6c,
	0x69, 0x62, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x72, 0x6d, 0x2f, 0x76, 0x31, 0x3b, 0x74,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x72, 0x6d, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_teleport_lib_teleterm_v1_usage_events_proto_rawDescOnce sync.Once
	file_teleport_lib_teleterm_v1_usage_events_proto_rawDescData = file_teleport_lib_teleterm_v1_usage_events_proto_rawDesc
)

func file_teleport_lib_teleterm_v1_usage_events_proto_rawDescGZIP() []byte {
	file_teleport_lib_teleterm_v1_usage_events_proto_rawDescOnce.Do(func() {
		file_teleport_lib_teleterm_v1_usage_events_proto_rawDescData = protoimpl.X.CompressGZIP(file_teleport_lib_teleterm_v1_usage_events_proto_rawDescData)
	})
	return file_teleport_lib_teleterm_v1_usage_events_proto_rawDescData
}

var file_teleport_lib_teleterm_v1_usage_events_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_teleport_lib_teleterm_v1_usage_events_proto_goTypes = []any{
	(*ReportUsageEventRequest)(nil),           // 0: teleport.lib.teleterm.v1.ReportUsageEventRequest
	(*v1alpha.SubmitConnectEventRequest)(nil), // 1: prehog.v1alpha.SubmitConnectEventRequest
}
var file_teleport_lib_teleterm_v1_usage_events_proto_depIdxs = []int32{
	1, // 0: teleport.lib.teleterm.v1.ReportUsageEventRequest.prehog_req:type_name -> prehog.v1alpha.SubmitConnectEventRequest
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_teleport_lib_teleterm_v1_usage_events_proto_init() }
func file_teleport_lib_teleterm_v1_usage_events_proto_init() {
	if File_teleport_lib_teleterm_v1_usage_events_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_teleport_lib_teleterm_v1_usage_events_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ReportUsageEventRequest); i {
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
			RawDescriptor: file_teleport_lib_teleterm_v1_usage_events_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_teleport_lib_teleterm_v1_usage_events_proto_goTypes,
		DependencyIndexes: file_teleport_lib_teleterm_v1_usage_events_proto_depIdxs,
		MessageInfos:      file_teleport_lib_teleterm_v1_usage_events_proto_msgTypes,
	}.Build()
	File_teleport_lib_teleterm_v1_usage_events_proto = out.File
	file_teleport_lib_teleterm_v1_usage_events_proto_rawDesc = nil
	file_teleport_lib_teleterm_v1_usage_events_proto_goTypes = nil
	file_teleport_lib_teleterm_v1_usage_events_proto_depIdxs = nil
}

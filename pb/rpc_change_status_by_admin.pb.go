// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: rpc_change_status_by_admin.proto

package pb

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

type ChangeStatusJobByAdminRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ChangeStatusJobByAdminRequest) Reset() {
	*x = ChangeStatusJobByAdminRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_change_status_by_admin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeStatusJobByAdminRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeStatusJobByAdminRequest) ProtoMessage() {}

func (x *ChangeStatusJobByAdminRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_change_status_by_admin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeStatusJobByAdminRequest.ProtoReflect.Descriptor instead.
func (*ChangeStatusJobByAdminRequest) Descriptor() ([]byte, []int) {
	return file_rpc_change_status_by_admin_proto_rawDescGZIP(), []int{0}
}

func (x *ChangeStatusJobByAdminRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ChangeStatusJobByAdminRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type ChangeStatusJobByAdminResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ChangeStatusJobByAdminResponse) Reset() {
	*x = ChangeStatusJobByAdminResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_change_status_by_admin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeStatusJobByAdminResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeStatusJobByAdminResponse) ProtoMessage() {}

func (x *ChangeStatusJobByAdminResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_change_status_by_admin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeStatusJobByAdminResponse.ProtoReflect.Descriptor instead.
func (*ChangeStatusJobByAdminResponse) Descriptor() ([]byte, []int) {
	return file_rpc_change_status_by_admin_proto_rawDescGZIP(), []int{1}
}

func (x *ChangeStatusJobByAdminResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_rpc_change_status_by_admin_proto protoreflect.FileDescriptor

var file_rpc_change_status_by_admin_proto_rawDesc = []byte{
	0x0a, 0x20, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x5f, 0x62, 0x79, 0x5f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0d, 0x6a, 0x6f, 0x62, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x2e, 0x6a, 0x6f,
	0x62, 0x22, 0x47, 0x0a, 0x1d, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x4a, 0x6f, 0x62, 0x42, 0x79, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x38, 0x0a, 0x1e, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4a, 0x6f, 0x62, 0x42, 0x79, 0x41,
	0x64, 0x6d, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x53, 0x45, 0x43, 0x2d, 0x4a, 0x6f, 0x62, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74,
	0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2d, 0x6a, 0x6f, 0x62, 0x2d, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_change_status_by_admin_proto_rawDescOnce sync.Once
	file_rpc_change_status_by_admin_proto_rawDescData = file_rpc_change_status_by_admin_proto_rawDesc
)

func file_rpc_change_status_by_admin_proto_rawDescGZIP() []byte {
	file_rpc_change_status_by_admin_proto_rawDescOnce.Do(func() {
		file_rpc_change_status_by_admin_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_change_status_by_admin_proto_rawDescData)
	})
	return file_rpc_change_status_by_admin_proto_rawDescData
}

var file_rpc_change_status_by_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_change_status_by_admin_proto_goTypes = []interface{}{
	(*ChangeStatusJobByAdminRequest)(nil),  // 0: jobstreet.job.ChangeStatusJobByAdminRequest
	(*ChangeStatusJobByAdminResponse)(nil), // 1: jobstreet.job.ChangeStatusJobByAdminResponse
}
var file_rpc_change_status_by_admin_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_change_status_by_admin_proto_init() }
func file_rpc_change_status_by_admin_proto_init() {
	if File_rpc_change_status_by_admin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_change_status_by_admin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeStatusJobByAdminRequest); i {
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
		file_rpc_change_status_by_admin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeStatusJobByAdminResponse); i {
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
			RawDescriptor: file_rpc_change_status_by_admin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_change_status_by_admin_proto_goTypes,
		DependencyIndexes: file_rpc_change_status_by_admin_proto_depIdxs,
		MessageInfos:      file_rpc_change_status_by_admin_proto_msgTypes,
	}.Build()
	File_rpc_change_status_by_admin_proto = out.File
	file_rpc_change_status_by_admin_proto_rawDesc = nil
	file_rpc_change_status_by_admin_proto_goTypes = nil
	file_rpc_change_status_by_admin_proto_depIdxs = nil
}

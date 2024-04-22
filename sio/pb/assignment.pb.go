// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v5.26.1
// source: assignment.proto

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

type AssignmentInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Formula string         `protobuf:"bytes,1,opt,name=formula,proto3" json:"formula,omitempty"`
	Mapping *AssignmentMap `protobuf:"bytes,2,opt,name=mapping,proto3" json:"mapping,omitempty"`
}

func (x *AssignmentInput) Reset() {
	*x = AssignmentInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assignment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignmentInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignmentInput) ProtoMessage() {}

func (x *AssignmentInput) ProtoReflect() protoreflect.Message {
	mi := &file_assignment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignmentInput.ProtoReflect.Descriptor instead.
func (*AssignmentInput) Descriptor() ([]byte, []int) {
	return file_assignment_proto_rawDescGZIP(), []int{0}
}

func (x *AssignmentInput) GetFormula() string {
	if x != nil {
		return x.Formula
	}
	return ""
}

func (x *AssignmentInput) GetMapping() *AssignmentMap {
	if x != nil {
		return x.Mapping
	}
	return nil
}

type AssignmentMap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mapping []*Assignment `protobuf:"bytes,1,rep,name=mapping,proto3" json:"mapping,omitempty"`
}

func (x *AssignmentMap) Reset() {
	*x = AssignmentMap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assignment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignmentMap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignmentMap) ProtoMessage() {}

func (x *AssignmentMap) ProtoReflect() protoreflect.Message {
	mi := &file_assignment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignmentMap.ProtoReflect.Descriptor instead.
func (*AssignmentMap) Descriptor() ([]byte, []int) {
	return file_assignment_proto_rawDescGZIP(), []int{1}
}

func (x *AssignmentMap) GetMapping() []*Assignment {
	if x != nil {
		return x.Mapping
	}
	return nil
}

type Assignment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Variable string `protobuf:"bytes,1,opt,name=variable,proto3" json:"variable,omitempty"`
	Value    bool   `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Assignment) Reset() {
	*x = Assignment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assignment_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Assignment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Assignment) ProtoMessage() {}

func (x *Assignment) ProtoReflect() protoreflect.Message {
	mi := &file_assignment_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Assignment.ProtoReflect.Descriptor instead.
func (*Assignment) Descriptor() ([]byte, []int) {
	return file_assignment_proto_rawDescGZIP(), []int{2}
}

func (x *Assignment) GetVariable() string {
	if x != nil {
		return x.Variable
	}
	return ""
}

func (x *Assignment) GetValue() bool {
	if x != nil {
		return x.Value
	}
	return false
}

var File_assignment_proto protoreflect.FileDescriptor

var file_assignment_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x60,
	0x0a, 0x0f, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x66, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x66, 0x6f, 0x72, 0x6d, 0x75, 0x6c, 0x61, 0x12, 0x33, 0x0a, 0x07, 0x6d,
	0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61,
	0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e,
	0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x61, 0x70, 0x52, 0x07, 0x6d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67,
	0x22, 0x41, 0x0a, 0x0d, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x61,
	0x70, 0x12, 0x30, 0x0a, 0x07, 0x6d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x6d, 0x61, 0x70, 0x70,
	0x69, 0x6e, 0x67, 0x22, 0x3e, 0x0a, 0x0a, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x76, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x2f, 0x6c, 0x6f, 0x67,
	0x69, 0x63, 0x6e, 0x67, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x73, 0x69, 0x6f,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_assignment_proto_rawDescOnce sync.Once
	file_assignment_proto_rawDescData = file_assignment_proto_rawDesc
)

func file_assignment_proto_rawDescGZIP() []byte {
	file_assignment_proto_rawDescOnce.Do(func() {
		file_assignment_proto_rawDescData = protoimpl.X.CompressGZIP(file_assignment_proto_rawDescData)
	})
	return file_assignment_proto_rawDescData
}

var file_assignment_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_assignment_proto_goTypes = []interface{}{
	(*AssignmentInput)(nil), // 0: assignment.AssignmentInput
	(*AssignmentMap)(nil),   // 1: assignment.AssignmentMap
	(*Assignment)(nil),      // 2: assignment.Assignment
}
var file_assignment_proto_depIdxs = []int32{
	1, // 0: assignment.AssignmentInput.mapping:type_name -> assignment.AssignmentMap
	2, // 1: assignment.AssignmentMap.mapping:type_name -> assignment.Assignment
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_assignment_proto_init() }
func file_assignment_proto_init() {
	if File_assignment_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_assignment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignmentInput); i {
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
		file_assignment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignmentMap); i {
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
		file_assignment_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Assignment); i {
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
			RawDescriptor: file_assignment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_assignment_proto_goTypes,
		DependencyIndexes: file_assignment_proto_depIdxs,
		MessageInfos:      file_assignment_proto_msgTypes,
	}.Build()
	File_assignment_proto = out.File
	file_assignment_proto_rawDesc = nil
	file_assignment_proto_goTypes = nil
	file_assignment_proto_depIdxs = nil
}

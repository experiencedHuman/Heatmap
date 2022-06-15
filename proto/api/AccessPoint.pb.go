// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: api/AccessPoint.proto

package api

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type AccessPoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Lat       string `protobuf:"bytes,2,opt,name=lat,proto3" json:"lat,omitempty"`
	Long      string `protobuf:"bytes,3,opt,name=long,proto3" json:"long,omitempty"`
	Intensity int64  `protobuf:"varint,4,opt,name=intensity,proto3" json:"intensity,omitempty"`
	Max       int64  `protobuf:"varint,5,opt,name=max,proto3" json:"max,omitempty"`
	Min       int64  `protobuf:"varint,6,opt,name=min,proto3" json:"min,omitempty"`
}

func (x *AccessPoint) Reset() {
	*x = AccessPoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_AccessPoint_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessPoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessPoint) ProtoMessage() {}

func (x *AccessPoint) ProtoReflect() protoreflect.Message {
	mi := &file_api_AccessPoint_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessPoint.ProtoReflect.Descriptor instead.
func (*AccessPoint) Descriptor() ([]byte, []int) {
	return file_api_AccessPoint_proto_rawDescGZIP(), []int{0}
}

func (x *AccessPoint) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AccessPoint) GetLat() string {
	if x != nil {
		return x.Lat
	}
	return ""
}

func (x *AccessPoint) GetLong() string {
	if x != nil {
		return x.Long
	}
	return ""
}

func (x *AccessPoint) GetIntensity() int64 {
	if x != nil {
		return x.Intensity
	}
	return 0
}

func (x *AccessPoint) GetMax() int64 {
	if x != nil {
		return x.Max
	}
	return 0
}

func (x *AccessPoint) GetMin() int64 {
	if x != nil {
		return x.Min
	}
	return 0
}

type APRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Timestamp string `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *APRequest) Reset() {
	*x = APRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_AccessPoint_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *APRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*APRequest) ProtoMessage() {}

func (x *APRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_AccessPoint_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use APRequest.ProtoReflect.Descriptor instead.
func (*APRequest) Descriptor() ([]byte, []int) {
	return file_api_AccessPoint_proto_rawDescGZIP(), []int{1}
}

func (x *APRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *APRequest) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

type APResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Accesspoint *AccessPoint `protobuf:"bytes,1,opt,name=accesspoint,proto3" json:"accesspoint,omitempty"`
}

func (x *APResponse) Reset() {
	*x = APResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_AccessPoint_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *APResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*APResponse) ProtoMessage() {}

func (x *APResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_AccessPoint_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use APResponse.ProtoReflect.Descriptor instead.
func (*APResponse) Descriptor() ([]byte, []int) {
	return file_api_AccessPoint_proto_rawDescGZIP(), []int{2}
}

func (x *APResponse) GetAccesspoint() *AccessPoint {
	if x != nil {
		return x.Accesspoint
	}
	return nil
}

var File_api_AccessPoint_proto protoreflect.FileDescriptor

var file_api_AccessPoint_proto_rawDesc = []byte{
	0x0a, 0x15, 0x61, 0x70, 0x69, 0x2f, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x69, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x89, 0x01, 0x0a, 0x0b, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x6c, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6c, 0x61, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6c, 0x6f, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6c, 0x6f, 0x6e, 0x67, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x74,
	0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x74, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x03, 0x6d, 0x61, 0x78, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x22, 0x3d, 0x0a, 0x09, 0x41, 0x50, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x40, 0x0a, 0x0a, 0x41, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x32, 0xc8, 0x01, 0x0a, 0x09, 0x41, 0x50, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x58, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x50,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x1e, 0x12, 0x1c, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73,
	0x2f, 0x68, 0x65, 0x61, 0x74, 0x6d, 0x61, 0x70, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x12,
	0x61, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x69,
	0x6e, 0x74, 0x73, 0x12, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x50, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x50, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x12, 0x15, 0x2f, 0x61,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2f, 0x68, 0x65, 0x61, 0x74,
	0x6d, 0x61, 0x70, 0x62, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x30, 0x01, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6b, 0x76, 0x6f, 0x67, 0x6c, 0x69, 0x2f, 0x48, 0x65, 0x61, 0x74, 0x6d, 0x61, 0x70, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_api_AccessPoint_proto_rawDescOnce sync.Once
	file_api_AccessPoint_proto_rawDescData = file_api_AccessPoint_proto_rawDesc
)

func file_api_AccessPoint_proto_rawDescGZIP() []byte {
	file_api_AccessPoint_proto_rawDescOnce.Do(func() {
		file_api_AccessPoint_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_AccessPoint_proto_rawDescData)
	})
	return file_api_AccessPoint_proto_rawDescData
}

var file_api_AccessPoint_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_AccessPoint_proto_goTypes = []interface{}{
	(*AccessPoint)(nil), // 0: api.AccessPoint
	(*APRequest)(nil),   // 1: api.APRequest
	(*APResponse)(nil),  // 2: api.APResponse
}
var file_api_AccessPoint_proto_depIdxs = []int32{
	0, // 0: api.APResponse.accesspoint:type_name -> api.AccessPoint
	1, // 1: api.APService.GetAccessPoint:input_type -> api.APRequest
	1, // 2: api.APService.ListAccessPoints:input_type -> api.APRequest
	0, // 3: api.APService.GetAccessPoint:output_type -> api.AccessPoint
	2, // 4: api.APService.ListAccessPoints:output_type -> api.APResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_AccessPoint_proto_init() }
func file_api_AccessPoint_proto_init() {
	if File_api_AccessPoint_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_AccessPoint_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessPoint); i {
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
		file_api_AccessPoint_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*APRequest); i {
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
		file_api_AccessPoint_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*APResponse); i {
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
			RawDescriptor: file_api_AccessPoint_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_AccessPoint_proto_goTypes,
		DependencyIndexes: file_api_AccessPoint_proto_depIdxs,
		MessageInfos:      file_api_AccessPoint_proto_msgTypes,
	}.Build()
	File_api_AccessPoint_proto = out.File
	file_api_AccessPoint_proto_rawDesc = nil
	file_api_AccessPoint_proto_goTypes = nil
	file_api_AccessPoint_proto_depIdxs = nil
}

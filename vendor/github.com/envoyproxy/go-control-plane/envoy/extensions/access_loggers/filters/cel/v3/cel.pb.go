// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.4
// source: envoy/extensions/access_loggers/filters/cel/v3/cel.proto

package celv3

import (
	_ "github.com/cncf/xds/go/udpa/annotations"
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

// ExpressionFilter is an access logging filter that evaluates configured
// symbolic Common Expression Language expressions to inform the decision
// to generate an access log.
type ExpressionFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Expression that, when evaluated, will be used to filter access logs.
	// Expressions are based on the set of Envoy :ref:`attributes <arch_overview_attributes>`.
	// The provided expression must evaluate to true for logging (expression errors are considered false).
	// Examples:
	//
	// * “response.code >= 400“
	// * “(connection.mtls && request.headers['x-log-mtls'] == 'true') || request.url_path.contains('v1beta3')“
	Expression string `protobuf:"bytes,1,opt,name=expression,proto3" json:"expression,omitempty"`
}

func (x *ExpressionFilter) Reset() {
	*x = ExpressionFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExpressionFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExpressionFilter) ProtoMessage() {}

func (x *ExpressionFilter) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExpressionFilter.ProtoReflect.Descriptor instead.
func (*ExpressionFilter) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDescGZIP(), []int{0}
}

func (x *ExpressionFilter) GetExpression() string {
	if x != nil {
		return x.Expression
	}
	return ""
}

var File_envoy_extensions_access_loggers_filters_cel_v3_cel_proto protoreflect.FileDescriptor

var file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDesc = []byte{
	0x0a, 0x38, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72,
	0x73, 0x2f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2f, 0x63, 0x65, 0x6c, 0x2f, 0x76, 0x33,
	0x2f, 0x63, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x2e, 0x65, 0x6e, 0x76, 0x6f,
	0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x73, 0x2e, 0x63, 0x65, 0x6c, 0x2e, 0x76, 0x33, 0x1a, 0x1d, 0x75, 0x64, 0x70, 0x61,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x32, 0x0a, 0x10, 0x45, 0x78, 0x70,
	0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x1e, 0x0a,
	0x0a, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0xaf, 0x01,
	0xba, 0x80, 0xc8, 0xd1, 0x06, 0x02, 0x10, 0x02, 0x0a, 0x3c, 0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76,
	0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f,
	0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e,
	0x63, 0x65, 0x6c, 0x2e, 0x76, 0x33, 0x42, 0x08, 0x43, 0x65, 0x6c, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x5b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65,
	0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x67, 0x6f, 0x2d, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x2d, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79,
	0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x5f, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x73, 0x2f, 0x66, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x73, 0x2f, 0x63, 0x65, 0x6c, 0x2f, 0x76, 0x33, 0x3b, 0x63, 0x65, 0x6c, 0x76, 0x33, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDescOnce sync.Once
	file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDescData = file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDesc
)

func file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDescGZIP() []byte {
	file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDescOnce.Do(func() {
		file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDescData = protoimpl.X.CompressGZIP(file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDescData)
	})
	return file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDescData
}

var file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_goTypes = []interface{}{
	(*ExpressionFilter)(nil), // 0: envoy.extensions.access_loggers.filters.cel.v3.ExpressionFilter
}
var file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_init() }
func file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_init() {
	if File_envoy_extensions_access_loggers_filters_cel_v3_cel_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExpressionFilter); i {
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
			RawDescriptor: file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_goTypes,
		DependencyIndexes: file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_depIdxs,
		MessageInfos:      file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_msgTypes,
	}.Build()
	File_envoy_extensions_access_loggers_filters_cel_v3_cel_proto = out.File
	file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_rawDesc = nil
	file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_goTypes = nil
	file_envoy_extensions_access_loggers_filters_cel_v3_cel_proto_depIdxs = nil
}

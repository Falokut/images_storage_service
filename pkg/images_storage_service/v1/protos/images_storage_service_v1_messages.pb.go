// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: images_storage_service_v1_messages.proto

package protos

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

type UploadImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Category string `protobuf:"bytes,1,opt,name=category,proto3" json:"category,omitempty"`
	Image    []byte `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
}

func (x *UploadImageRequest) Reset() {
	*x = UploadImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_images_storage_service_v1_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadImageRequest) ProtoMessage() {}

func (x *UploadImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_images_storage_service_v1_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadImageRequest.ProtoReflect.Descriptor instead.
func (*UploadImageRequest) Descriptor() ([]byte, []int) {
	return file_images_storage_service_v1_messages_proto_rawDescGZIP(), []int{0}
}

func (x *UploadImageRequest) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *UploadImageRequest) GetImage() []byte {
	if x != nil {
		return x.Image
	}
	return nil
}

type StreamingUploadImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Can be empty
	Category string `protobuf:"bytes,1,opt,name=category,proto3" json:"category,omitempty"`
	Data     []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *StreamingUploadImageRequest) Reset() {
	*x = StreamingUploadImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_images_storage_service_v1_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamingUploadImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamingUploadImageRequest) ProtoMessage() {}

func (x *StreamingUploadImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_images_storage_service_v1_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamingUploadImageRequest.ProtoReflect.Descriptor instead.
func (*StreamingUploadImageRequest) Descriptor() ([]byte, []int) {
	return file_images_storage_service_v1_messages_proto_rawDescGZIP(), []int{1}
}

func (x *StreamingUploadImageRequest) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *StreamingUploadImageRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type UploadImageResponce struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImageId string `protobuf:"bytes,1,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty"`
}

func (x *UploadImageResponce) Reset() {
	*x = UploadImageResponce{}
	if protoimpl.UnsafeEnabled {
		mi := &file_images_storage_service_v1_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadImageResponce) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadImageResponce) ProtoMessage() {}

func (x *UploadImageResponce) ProtoReflect() protoreflect.Message {
	mi := &file_images_storage_service_v1_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadImageResponce.ProtoReflect.Descriptor instead.
func (*UploadImageResponce) Descriptor() ([]byte, []int) {
	return file_images_storage_service_v1_messages_proto_rawDescGZIP(), []int{2}
}

func (x *UploadImageResponce) GetImageId() string {
	if x != nil {
		return x.ImageId
	}
	return ""
}

type GetImageResponce struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImageData []byte `protobuf:"bytes,1,opt,name=ImageData,json=image_data,proto3" json:"ImageData,omitempty"`
}

func (x *GetImageResponce) Reset() {
	*x = GetImageResponce{}
	if protoimpl.UnsafeEnabled {
		mi := &file_images_storage_service_v1_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetImageResponce) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetImageResponce) ProtoMessage() {}

func (x *GetImageResponce) ProtoReflect() protoreflect.Message {
	mi := &file_images_storage_service_v1_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetImageResponce.ProtoReflect.Descriptor instead.
func (*GetImageResponce) Descriptor() ([]byte, []int) {
	return file_images_storage_service_v1_messages_proto_rawDescGZIP(), []int{3}
}

func (x *GetImageResponce) GetImageData() []byte {
	if x != nil {
		return x.ImageData
	}
	return nil
}

type ImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Category string `protobuf:"bytes,1,opt,name=category,proto3" json:"category,omitempty"`
	ImageId  string `protobuf:"bytes,2,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty"`
}

func (x *ImageRequest) Reset() {
	*x = ImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_images_storage_service_v1_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageRequest) ProtoMessage() {}

func (x *ImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_images_storage_service_v1_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageRequest.ProtoReflect.Descriptor instead.
func (*ImageRequest) Descriptor() ([]byte, []int) {
	return file_images_storage_service_v1_messages_proto_rawDescGZIP(), []int{4}
}

func (x *ImageRequest) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *ImageRequest) GetImageId() string {
	if x != nil {
		return x.ImageId
	}
	return ""
}

type ImageExistResponce struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImageExist bool `protobuf:"varint,1,opt,name=ImageExist,json=image_exist,proto3" json:"ImageExist,omitempty"`
}

func (x *ImageExistResponce) Reset() {
	*x = ImageExistResponce{}
	if protoimpl.UnsafeEnabled {
		mi := &file_images_storage_service_v1_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImageExistResponce) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageExistResponce) ProtoMessage() {}

func (x *ImageExistResponce) ProtoReflect() protoreflect.Message {
	mi := &file_images_storage_service_v1_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageExistResponce.ProtoReflect.Descriptor instead.
func (*ImageExistResponce) Descriptor() ([]byte, []int) {
	return file_images_storage_service_v1_messages_proto_rawDescGZIP(), []int{5}
}

func (x *ImageExistResponce) GetImageExist() bool {
	if x != nil {
		return x.ImageExist
	}
	return false
}

type ReplaceImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Category         string `protobuf:"bytes,1,opt,name=category,proto3" json:"category,omitempty"`
	ImageId          string `protobuf:"bytes,2,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty"`
	ImageData        []byte `protobuf:"bytes,3,opt,name=ImageData,json=image_data,proto3" json:"ImageData,omitempty"`
	CreateIfNotExist bool   `protobuf:"varint,4,opt,name=CreateIfNotExist,json=create_if_not_exist,proto3" json:"CreateIfNotExist,omitempty"`
}

func (x *ReplaceImageRequest) Reset() {
	*x = ReplaceImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_images_storage_service_v1_messages_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplaceImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplaceImageRequest) ProtoMessage() {}

func (x *ReplaceImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_images_storage_service_v1_messages_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplaceImageRequest.ProtoReflect.Descriptor instead.
func (*ReplaceImageRequest) Descriptor() ([]byte, []int) {
	return file_images_storage_service_v1_messages_proto_rawDescGZIP(), []int{6}
}

func (x *ReplaceImageRequest) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *ReplaceImageRequest) GetImageId() string {
	if x != nil {
		return x.ImageId
	}
	return ""
}

func (x *ReplaceImageRequest) GetImageData() []byte {
	if x != nil {
		return x.ImageData
	}
	return nil
}

func (x *ReplaceImageRequest) GetCreateIfNotExist() bool {
	if x != nil {
		return x.CreateIfNotExist
	}
	return false
}

type ReplaceImageResponce struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// returns only if in request CreateIfNotExist(create_if_not_exist) = true,
	// otherwise this field will be empty
	ImageId string `protobuf:"bytes,1,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty"`
}

func (x *ReplaceImageResponce) Reset() {
	*x = ReplaceImageResponce{}
	if protoimpl.UnsafeEnabled {
		mi := &file_images_storage_service_v1_messages_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplaceImageResponce) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplaceImageResponce) ProtoMessage() {}

func (x *ReplaceImageResponce) ProtoReflect() protoreflect.Message {
	mi := &file_images_storage_service_v1_messages_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplaceImageResponce.ProtoReflect.Descriptor instead.
func (*ReplaceImageResponce) Descriptor() ([]byte, []int) {
	return file_images_storage_service_v1_messages_proto_rawDescGZIP(), []int{7}
}

func (x *ReplaceImageResponce) GetImageId() string {
	if x != nil {
		return x.ImageId
	}
	return ""
}

type UserErrorMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *UserErrorMessage) Reset() {
	*x = UserErrorMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_images_storage_service_v1_messages_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserErrorMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserErrorMessage) ProtoMessage() {}

func (x *UserErrorMessage) ProtoReflect() protoreflect.Message {
	mi := &file_images_storage_service_v1_messages_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserErrorMessage.ProtoReflect.Descriptor instead.
func (*UserErrorMessage) Descriptor() ([]byte, []int) {
	return file_images_storage_service_v1_messages_proto_rawDescGZIP(), []int{8}
}

func (x *UserErrorMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_images_storage_service_v1_messages_proto protoreflect.FileDescriptor

var file_images_storage_service_v1_messages_proto_rawDesc = []byte{
	0x0a, 0x28, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x76, 0x31, 0x5f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x73, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x22, 0x46, 0x0a, 0x12, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x22, 0x4d, 0x0a, 0x1b, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x30, 0x0a, 0x13, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x63, 0x65,
	0x12, 0x19, 0x0a, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x64, 0x22, 0x31, 0x0a, 0x10, 0x47,
	0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x63, 0x65, 0x12,
	0x1d, 0x0a, 0x09, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0a, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x22, 0x45,
	0x0a, 0x0c, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x49, 0x64, 0x22, 0x35, 0x0a, 0x12, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x45, 0x78,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0a, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x45, 0x78, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0b, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x65, 0x78, 0x69, 0x73, 0x74, 0x22, 0x9a, 0x01, 0x0a,
	0x13, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x12, 0x19, 0x0a, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x09, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2d, 0x0a, 0x10, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x49, 0x66, 0x4e, 0x6f, 0x74, 0x45, 0x78, 0x69, 0x73, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x13, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x66, 0x5f,
	0x6e, 0x6f, 0x74, 0x5f, 0x65, 0x78, 0x69, 0x73, 0x74, 0x22, 0x31, 0x0a, 0x14, 0x52, 0x65, 0x70,
	0x6c, 0x61, 0x63, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x63,
	0x65, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x64, 0x22, 0x2c, 0x0a, 0x10,
	0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x22, 0x5a, 0x20, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x73, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_images_storage_service_v1_messages_proto_rawDescOnce sync.Once
	file_images_storage_service_v1_messages_proto_rawDescData = file_images_storage_service_v1_messages_proto_rawDesc
)

func file_images_storage_service_v1_messages_proto_rawDescGZIP() []byte {
	file_images_storage_service_v1_messages_proto_rawDescOnce.Do(func() {
		file_images_storage_service_v1_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_images_storage_service_v1_messages_proto_rawDescData)
	})
	return file_images_storage_service_v1_messages_proto_rawDescData
}

var file_images_storage_service_v1_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_images_storage_service_v1_messages_proto_goTypes = []interface{}{
	(*UploadImageRequest)(nil),          // 0: images_storage_service.UploadImageRequest
	(*StreamingUploadImageRequest)(nil), // 1: images_storage_service.StreamingUploadImageRequest
	(*UploadImageResponce)(nil),         // 2: images_storage_service.UploadImageResponce
	(*GetImageResponce)(nil),            // 3: images_storage_service.GetImageResponce
	(*ImageRequest)(nil),                // 4: images_storage_service.ImageRequest
	(*ImageExistResponce)(nil),          // 5: images_storage_service.ImageExistResponce
	(*ReplaceImageRequest)(nil),         // 6: images_storage_service.ReplaceImageRequest
	(*ReplaceImageResponce)(nil),        // 7: images_storage_service.ReplaceImageResponce
	(*UserErrorMessage)(nil),            // 8: images_storage_service.UserErrorMessage
}
var file_images_storage_service_v1_messages_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_images_storage_service_v1_messages_proto_init() }
func file_images_storage_service_v1_messages_proto_init() {
	if File_images_storage_service_v1_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_images_storage_service_v1_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadImageRequest); i {
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
		file_images_storage_service_v1_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamingUploadImageRequest); i {
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
		file_images_storage_service_v1_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadImageResponce); i {
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
		file_images_storage_service_v1_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetImageResponce); i {
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
		file_images_storage_service_v1_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImageRequest); i {
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
		file_images_storage_service_v1_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImageExistResponce); i {
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
		file_images_storage_service_v1_messages_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplaceImageRequest); i {
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
		file_images_storage_service_v1_messages_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplaceImageResponce); i {
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
		file_images_storage_service_v1_messages_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserErrorMessage); i {
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
			RawDescriptor: file_images_storage_service_v1_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_images_storage_service_v1_messages_proto_goTypes,
		DependencyIndexes: file_images_storage_service_v1_messages_proto_depIdxs,
		MessageInfos:      file_images_storage_service_v1_messages_proto_msgTypes,
	}.Build()
	File_images_storage_service_v1_messages_proto = out.File
	file_images_storage_service_v1_messages_proto_rawDesc = nil
	file_images_storage_service_v1_messages_proto_goTypes = nil
	file_images_storage_service_v1_messages_proto_depIdxs = nil
}

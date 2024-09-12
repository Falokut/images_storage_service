package domain

type UploadImageRequest struct {
	Category string `validate:"required"`
	Image    []byte `validate:"required"`
}

type UploadImageResponse struct {
	ImageId string
}

type ImageRequest struct {
	ImageId  string `validate:"required"`
	Category string `validate:"required"`
}

type ImageExistResponse struct {
	ImageExist bool
}

type Empty struct{}

type ReplaceImageRequest struct {
	ImageId          string `validate:"required"`
	Category         string `validate:"required"`
	ImageData        []byte `validate:"required"`
	CreateIfNotExist bool   `validate:"required"`
}

type ReplaceImageResponse struct {
	ImageId string
}

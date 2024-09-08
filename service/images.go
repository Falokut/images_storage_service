package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Falokut/images_storage_service/domain"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Config struct {
	MaxImageSize int
}

//go:generate mockgen -source=repository.go -destination=mocks/imageStorage.go
type ImageStorage interface {
	SaveImage(ctx context.Context, img []byte, filename string, category string) error
	GetImage(ctx context.Context, imageID string, category string) ([]byte, error)
	IsImageExist(ctx context.Context, imageID string, category string) (bool, error)
	DeleteImage(ctx context.Context, imageID string, category string) error
	ReplaceImage(ctx context.Context, img []byte, filename string, category string) error
}

//go:generate mockgen -source=service.go -destination=mocks/service.go
type Metrics interface {
	IncBytesUploaded(bytesUploaded int)
}

type Images struct {
	metrics Metrics
	storage ImageStorage
	cfg     Config
}

func NewImages(metrics Metrics, storage ImageStorage, cfg Config) Images {
	return Images{
		metrics: metrics,
		storage: storage,
		cfg:     cfg,
	}
}

func (s Images) detectFileType(fileData []byte) string {
	fileType := mimetype.Detect(fileData)
	Type := strings.Split(fileType.String(), "/")
	return Type[0]
}

func (s Images) SaveImage(ctx context.Context, img []byte, category string) (string, error) {
	if fileType := s.detectFileType(img); fileType != "image" {
		return "", domain.NewInvalidArgumentError(fmt.Sprintf("file type is not supported. file type: '%s'", fileType))
	}

	err := s.checkImage(img)
	if err != nil {
		return "", errors.WithMessage(err, "check image")
	}

	imageId := uuid.NewString()
	err = s.storage.SaveImage(ctx, img, imageId, category)
	if err != nil {
		return "", errors.WithMessage(err, "save image")
	}

	s.metrics.IncBytesUploaded(len(img))
	return imageId, nil
}

func (s Images) GetImage(ctx context.Context, imageId string, category string) ([]byte, error) {
	img, err := s.storage.GetImage(ctx, imageId, category)
	if err != nil {
		return nil, errors.WithMessage(err, "get image")
	}
	return img, nil
}

func (s Images) IsImageExist(ctx context.Context, imageId string, category string) (bool, error) {
	exists, err := s.storage.IsImageExist(ctx, imageId, category)
	if err != nil {
		return false, errors.WithMessage(err, "is image exist")
	}
	return exists, nil
}

func (s Images) DeleteImage(ctx context.Context, imageId string, category string) error {
	err := s.storage.DeleteImage(ctx, imageId, category)
	if err != nil {
		return errors.WithMessage(err, "delete image")
	}
	return nil
}

func (s Images) RewriteImage(ctx context.Context, img []byte, imageId string,
	category string, createImageIfNotExist bool) (string, error) {
	if fileType := s.detectFileType(img); fileType != "image" {
		return "", domain.NewInvalidArgumentError(fmt.Sprintf("file type is not supported. file type: '%s'", fileType))
	}

	err := s.checkImage(img)
	if err != nil {
		return "", errors.WithMessage(err, "check image")
	}

	imageExist, err := s.storage.IsImageExist(ctx, imageId, category)
	if err != nil {
		return "", errors.WithMessage(err, "is image exist")
	}
	if imageExist {
		err = s.storage.ReplaceImage(ctx, img, imageId, category)
		if err != nil {
			return "", errors.WithMessage(err, "replace image")
		}
		return imageId, nil
	}
	if !createImageIfNotExist {
		return "", domain.ErrImageNotFound
	}

	newImageId := uuid.NewString()
	err = s.storage.SaveImage(ctx, img, newImageId, category)
	if err != nil {
		return "", errors.WithMessage(err, "save image")
	}
	return newImageId, nil
}

func (s Images) checkImage(image []byte) error {
	if len(image) == 0 {
		return domain.NewInvalidArgumentError("file has zero size")
	}
	if len(image) > s.cfg.MaxImageSize {
		return domain.NewInvalidArgumentError(
			fmt.Sprintf("image is too large max image size: %d, file size: %d",
				s.cfg.MaxImageSize, len(image)),
		)
	}

	return nil
}

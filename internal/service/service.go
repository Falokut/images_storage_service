package service

import (
	"context"
	"strings"

	"github.com/Falokut/images_storage_service/internal/models"
	"github.com/Falokut/images_storage_service/internal/repository"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Config struct {
	MaxImageSize int
}

//go:generate mockgen -source=service.go -destination=mocks/service.go
type Metrics interface {
	IncBytesUploaded(bytesUploaded int)
}

//go:generate mockgen -source=service.go -destination=mocks/service.go
type ImagesStorageService interface {
	SaveImage(ctx context.Context, img []byte, category string) (string, error)
	GetImage(ctx context.Context, imageId string, category string) ([]byte, error)
	IsImageExist(ctx context.Context, imageId string, category string) (bool, error)
	DeleteImage(ctx context.Context, imageId string, category string) error
	RewriteImage(ctx context.Context, img []byte, imageId string, category string, createImageIfNotExist bool) (string, error)
}

type imagesStorageService struct {
	logger  *logrus.Logger
	metrics Metrics
	storage repository.ImageStorage
	cfg     Config
}

func NewImagesStorageService(logger *logrus.Logger,
	metrics Metrics,
	storage repository.ImageStorage,
	cfg Config) ImagesStorageService {
	return &imagesStorageService{
		logger:  logger,
		metrics: metrics,
		storage: storage,
		cfg:     cfg,
	}
}

func (s *imagesStorageService) detectFileType(fileData []byte) string {
	fileType := mimetype.Detect(fileData)
	Type := strings.Split(fileType.String(), "/")
	return Type[0]
}

func (s *imagesStorageService) SaveImage(ctx context.Context, img []byte, category string) (imageId string, err error) {
	if err = s.checkImage(img); err != nil {
		return
	}

	imageId = uuid.NewString()
	err = s.storage.SaveImage(ctx, img, imageId, category)
	if err != nil {
		return
	}

	s.metrics.IncBytesUploaded(len(img))
	return
}

func (s *imagesStorageService) GetImage(ctx context.Context, imageId string, category string) ([]byte, error) {
	return s.storage.GetImage(ctx, imageId, category)
}

func (s *imagesStorageService) IsImageExist(ctx context.Context, imageId string, category string) (bool, error) {
	return s.storage.IsImageExist(ctx, imageId, category)
}

func (s *imagesStorageService) DeleteImage(ctx context.Context, imageId string, category string) error {
	return s.storage.DeleteImage(ctx, imageId, category)
}

func (s *imagesStorageService) RewriteImage(ctx context.Context, img []byte, imageId string,
	category string, createImageIfNotExist bool) (newImageId string, err error) {
	if err = s.checkImage(img); err != nil {
		return
	}

	imageExist, err := s.storage.IsImageExist(ctx, imageId, category)
	if err != nil {
		return
	}
	if !createImageIfNotExist && !imageExist {
		err = models.Error(models.NotFound, "can't find image with specified id")
		return
	}

	if !imageExist && createImageIfNotExist {
		newImageId = uuid.NewString()
		err = s.storage.SaveImage(ctx, img, newImageId, category)
		return
	}

	err = s.storage.RewriteImage(ctx, img, imageId, category)
	return
}

func (s *imagesStorageService) checkImage(image []byte) error {
	if len(image) == 0 {
		return models.Error(models.InvalidArgument, "the received file has zero size")
	}
	if len(image) > s.cfg.MaxImageSize {
		return models.Errorf(models.InvalidArgument,
			"image is too large max image size: %d, file size: %d",
			s.cfg.MaxImageSize, len(image))
	}

	s.logger.Info("Checking filetype")
	if fileType := s.detectFileType(image); fileType != "image" {
		return models.Error(models.InvalidArgument, "the received file type is not supported")
	}

	return nil
}

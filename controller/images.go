package controller

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Falokut/go-kit/http/apierrors"
	"github.com/Falokut/images_storage_service/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go
type ImagesStorageService interface {
	SaveImage(ctx context.Context, img []byte, category string) (string, error)
	GetImage(ctx context.Context, imageId string, category string) ([]byte, error)
	IsImageExist(ctx context.Context, imageId string, category string) (bool, error)
	DeleteImage(ctx context.Context, imageId string, category string) error
	RewriteImage(ctx context.Context, img []byte, imageId string, category string, createImageIfNotExist bool) (string, error)
}

type Images struct {
	service      ImagesStorageService
	maxImageSize int
}

func NewImages(service ImagesStorageService, maxImageSize int) Images {
	return Images{
		service:      service,
		maxImageSize: maxImageSize,
	}
}

// UploadImage
//
//	@Tags			image
//	@Summary		Upload image
//	@Description	Загрузить изображение в хранилище
//	@Accept			json
//	@Produce		json
//
//	@Param			category	path		string						true	"Категория изображения"
//
//	@Param			body		body		domain.UploadImageRequest	true	"request body"
//	@Success		200			{object}	domain.UploadImageResponse
//	@Failure		400			{object}	apierrors.Error
//	@Failure		500			{object}	apierrors.Error
//	@Router			/image/:category [POST]
func (c Images) UploadImage(ctx context.Context, req domain.UploadImageRequest) (*domain.UploadImageResponse, error) {
	imageId, err := c.service.SaveImage(ctx, req.Image, req.Category)
	invalidArgError := domain.InvalidArgumentError{}
	switch {
	case errors.As(err, &invalidArgError):
		return nil, apierrors.NewBusinessError(invalidArgError.ErrCode, invalidArgError.Reason, err)
	case err != nil:
		return nil, err
	default:
		return &domain.UploadImageResponse{ImageId: imageId}, nil
	}
}

// GetImage
//
//	@Tags			image
//	@Summary		Get image
//	@Description	Получить изображение из хранилища
//	@Accept			json
//	@Produce		image/*
//
//	@Param			category	path		string	true	"Категория изображения"
//	@Param			imageId		path		string	true	"Идентификатор изображения"
//
//	@Success		200			{array}		byte
//	@Failure		400			{object}	apierrors.Error
//	@Failure		404			{object}	apierrors.Error
//	@Failure		500			{object}	apierrors.Error
//	@Router			/image/:category/:imageId [GET]
func (c Images) GetImage(ctx context.Context, w http.ResponseWriter, req domain.ImageRequest) error {
	image, err := c.service.GetImage(ctx, req.ImageId, req.Category)
	switch {
	case errors.Is(err, domain.ErrImageNotFound):
		return apierrors.New(http.StatusNotFound, domain.ErrCodeImageNotFound, domain.ErrImageNotFound.Error(), err)
	case err != nil:
		return err
	}
	w.Header().Set("Content-Type", http.DetectContentType(image))
	_, err = w.Write(image)
	return err
}

// IsImageExist
//
//	@Tags			image
//	@Summary		Is image exist
//	@Description	Проверить наличие изображения в хранилище
//	@Accept			json
//	@Produce		json
//
//	@Param			category	path		string	true	"Категория изображения"
//	@Param			imageId		path		string	true	"Идентификатор изображения"
//
//	@Success		200			{object}	domain.ImageExistResponse
//	@Failure		400			{object}	apierrors.Error
//	@Failure		404			{object}	apierrors.Error
//	@Failure		500			{object}	apierrors.Error
//	@Router			/image/:category/:imageId/exist [GET]
func (c Images) IsImageExist(ctx context.Context, req domain.ImageRequest) (*domain.ImageExistResponse, error) {
	imageExist, err := c.service.IsImageExist(ctx, req.ImageId, req.Category)
	if err != nil {
		return nil, err
	}

	return &domain.ImageExistResponse{ImageExist: imageExist}, nil
}

// DeleteImage
//
//	@Tags			image
//	@Summary		Delete image
//	@Description	Удалить изображение из хранилища
//	@Accept			json
//	@Produce		json
//
//	@Param			category	path		string	true	"Категория изображения"
//	@Param			imageId		path		string	true	"Идентификатор изображения"
//
//	@Success		200			{object}	domain.Empty
//	@Failure		400			{object}	apierrors.Error
//	@Failure		404			{object}	apierrors.Error
//	@Failure		500			{object}	apierrors.Error
//	@Router			/image/:category/:imageId [DELETE]
func (c Images) DeleteImage(ctx context.Context, req domain.ImageRequest) error {
	err := c.service.DeleteImage(ctx, req.ImageId, req.Category)
	switch {
	case errors.Is(err, domain.ErrImageNotFound):
		return apierrors.New(http.StatusNotFound, domain.ErrCodeImageNotFound, domain.ErrImageNotFound.Error(), err)
	case err != nil:
		return status.Error(codes.Internal, "internal server error")
	default:
		return nil
	}
}

// ReplaceImage
//
//	@Tags			image
//	@Summary		Replace image
//	@Description	Заменить изображение в хранилище
//	@Accept			json
//	@Produce		json
//
//	@Param			category	path	string						true	"Категория изображения"
//	@Param			imageId		path	string						true	"Идентификатор изображения"
//	@Param			body		body	domain.ReplaceImageRequest	true	"request body"

// @Success	200	{object}	domain.ReplaceImageResponse
// @Failure	400	{object}	apierrors.Error
// @Failure	404	{object}	apierrors.Error
// @Failure	500	{object}	apierrors.Error
// @Router		/image/:category/:imageId/replace [POST]
func (c Images) ReplaceImage(ctx context.Context, in domain.ReplaceImageRequest) (*domain.ReplaceImageResponse, error) {
	imageId, err := c.service.RewriteImage(ctx, in.ImageData, in.ImageId, in.Category, in.CreateIfNotExist)
	invalidArgError := domain.InvalidArgumentError{}
	switch {
	case errors.As(err, &invalidArgError):
		return nil, apierrors.NewBusinessError(invalidArgError.ErrCode, invalidArgError.Reason, err)
	case errors.Is(err, domain.ErrImageNotFound):
		return nil, apierrors.New(http.StatusNotFound, domain.ErrCodeImageNotFound, domain.ErrImageNotFound.Error(), err)
	case err != nil:
		return nil, err
	default:
		return &domain.ReplaceImageResponse{ImageId: imageId}, nil
	}
}

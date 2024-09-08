package controller

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Falokut/go-kit/log"
	"github.com/Falokut/images_storage_service/domain"
	img_storage_serv "github.com/Falokut/images_storage_service/pkg/images_storage_service/v1/protos"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
	img_storage_serv.UnimplementedImagesStorageServiceV1Server
	logger       log.Logger
	service      ImagesStorageService
	maxImageSize int
}

func NewImagesStorageServiceHandler(
	logger log.Logger,
	service ImagesStorageService,
	maxImageSize int,
) Images {
	return Images{
		logger:       logger,
		service:      service,
		maxImageSize: maxImageSize,
	}
}

func (c Images) UploadImage(ctx context.Context,
	in *img_storage_serv.UploadImageRequest) (res *img_storage_serv.UploadImageResponse, err error) {
	imageId, err := c.service.SaveImage(ctx, in.Image, in.Category)
	invalidArgError := &domain.InvalidArgumentError{}
	switch {
	case errors.As(err, &invalidArgError):
		c.logger.Warn(ctx, "invalid argument", log.Any("error", err))
		return nil, status.Error(codes.InvalidArgument, invalidArgError.Reason)
	case err != nil:
		c.logger.Error(ctx, "internal error", log.Any("error", err))
		return nil, status.Error(codes.Internal, "internal service error")
	default:
		return &img_storage_serv.UploadImageResponse{ImageId: imageId}, nil
	}
}

func (c Images) StreamingUploadImage(
	stream img_storage_serv.ImagesStorageServiceV1_StreamingUploadImageServer,
) error {
	ctx := stream.Context()
	req, imageData, err := c.receiveUploadImage(stream)
	if err != nil {
		c.logger.Error(ctx, "internal error", log.Any("error", err))
		return status.Error(codes.Internal, "internal service error")
	}
	if req == nil {
		c.logger.Warn(ctx, "invalid argument", log.Any("error", "empty request"))
		return status.Error(codes.InvalidArgument, "the received request is nil")
	}

	imageId, err := c.service.SaveImage(ctx, imageData, req.Category)
	invalidArgError := &domain.InvalidArgumentError{}
	switch {
	case errors.As(err, &invalidArgError):
		c.logger.Warn(ctx, "invalid argument", log.Any("error", err))
		return status.Error(codes.InvalidArgument, invalidArgError.Reason)
	case err != nil:
		c.logger.Error(ctx, "internal error", log.Any("error", err))
		return status.Error(codes.Internal, "internal service error")
	}

	err = stream.SendAndClose(&img_storage_serv.UploadImageResponse{ImageId: imageId})
	if err != nil {
		return status.Error(codes.Internal, "internal service error")
	}

	return nil
}

func (c Images) receiveUploadImage(
	stream img_storage_serv.ImagesStorageServiceV1_StreamingUploadImageServer,
) (*img_storage_serv.StreamingUploadImageRequest, []byte, error) {

	firstReq := &img_storage_serv.StreamingUploadImageRequest{}
	imageData := bytes.Buffer{}
	for {
		if stream.Context().Err() != nil {
			return nil, nil, status.Error(codes.Canceled, "")
		}

		req, rerr := stream.Recv()
		if firstReq == nil && req != nil {
			firstReq = req
		}

		if rerr == io.EOF {
			return firstReq, imageData.Bytes(), nil
		}
		if rerr != nil {
			return nil, nil, rerr
		}

		chunkSize := len(req.Data)
		imageSize := imageData.Len() + chunkSize
		if imageSize > c.maxImageSize {
			return nil, nil, status.Errorf(codes.InvalidArgument,
				"image is too large max image size: %d, file size: %d",
				c.maxImageSize, imageSize)
		}
		imageData.Write(req.Data)
	}
}

func (c Images) GetImage(ctx context.Context,
	in *img_storage_serv.ImageRequest) (res *httpbody.HttpBody, err error) {

	image, err := c.service.GetImage(ctx, in.ImageId, in.Category)
	switch {
	case errors.Is(err, domain.ErrImageNotFound):
		c.logger.Warn(ctx, "image not found", log.Any("error", err))
		return nil, status.Error(codes.NotFound, domain.ErrImageNotFound.Error())
	case err != nil:
		c.logger.Error(ctx, "internal error", log.Any("error", err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &httpbody.HttpBody{ContentType: http.DetectContentType(image), Data: image}, nil
}

func (c Images) IsImageExist(ctx context.Context,
	in *img_storage_serv.ImageRequest) (res *img_storage_serv.ImageExistResponse, err error) {
	imageExist, err := c.service.IsImageExist(ctx, in.ImageId, in.Category)
	if err != nil {
		c.logger.Error(ctx, "internal error", log.Any("error", err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &img_storage_serv.ImageExistResponse{ImageExist: imageExist}, nil
}

func (c Images) DeleteImage(ctx context.Context,
	in *img_storage_serv.ImageRequest) (res *emptypb.Empty, err error) {
	err = c.service.DeleteImage(ctx, in.ImageId, in.Category)
	switch {
	case errors.Is(err, domain.ErrImageNotFound):
		c.logger.Warn(ctx, "image not found", log.Any("error", err))
		return nil, status.Error(codes.NotFound, domain.ErrImageNotFound.Error())
	case err != nil:
		c.logger.Error(ctx, "internal error", log.Any("error", err))
		return nil, status.Error(codes.Internal, "internal server error")
	default:
		return &emptypb.Empty{}, nil
	}
}

func (c Images) ReplaceImage(ctx context.Context, in *img_storage_serv.ReplaceImageRequest) (*img_storage_serv.ReplaceImageResponse, error) {
	imageId, err := c.service.RewriteImage(ctx, in.ImageData, in.ImageId, in.Category, in.CreateIfNotExist)
	switch {
	case errors.Is(err, domain.ErrImageNotFound):
		c.logger.Warn(ctx, "image not found", log.Any("error", err))
		return nil, status.Error(codes.NotFound, domain.ErrImageNotFound.Error())
	case err != nil:
		c.logger.Error(ctx, "internal error", log.Any("error", err))
		return nil, status.Error(codes.Internal, "internal service error")
	default:
		return &img_storage_serv.ReplaceImageResponse{ImageId: imageId}, nil
	}
}

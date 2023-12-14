package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/Falokut/images_storage_service/internal/repository"
	img_storage_serv "github.com/Falokut/images_storage_service/pkg/images_storage_service/v1/protos"
	"github.com/Falokut/images_storage_service/pkg/metrics"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Config struct {
	MaxImageSize int
}

type ImagesStorageService struct {
	img_storage_serv.UnimplementedImagesStorageServiceV1Server
	logger       *logrus.Logger
	cfg          Config
	imageStorage repository.ImageStorage
	errHandler   errorHandler
	metrics      metrics.Metrics
}

func NewImagesStorageService(
	logger *logrus.Logger,
	cfg Config,
	imageStorage repository.ImageStorage,
	metrics metrics.Metrics,
) *ImagesStorageService {
	errHandler := newErrorHandler(logger)
	return &ImagesStorageService{
		logger:       logger,
		cfg:          cfg,
		imageStorage: imageStorage,
		errHandler:   errHandler,
		metrics:      metrics,
	}
}

func (s *ImagesStorageService) UploadImage(ctx context.Context,
	in *img_storage_serv.UploadImageRequest) (*img_storage_serv.UploadImageResponce, error) {
	s.logger.Info("Start uploading image")
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesStorageService.UploadImage")
	defer span.Finish()

	if err := s.checkImage(ctx, in.Image); err != nil {
		span.SetTag("grpc.status", status.Code(err))
		ext.LogError(span, err)
		return nil, err
	}

	imageId, err := s.saveImage(ctx, in.Image, in.Category)
	if err != nil {
		span.SetTag("grpc.status", status.Code(err))
		ext.LogError(span, err)
		return nil, err
	}

	s.logger.Info("Image uploaded")

	span.SetTag("grpc.status", codes.OK)
	return &img_storage_serv.UploadImageResponce{ImageId: imageId}, nil
}

func (s *ImagesStorageService) checkImage(ctx context.Context, Image []byte) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ImagesStorageService.checkImage")
	defer span.Finish()

	if len(Image) == 0 {
		return s.errHandler.createErrorResponceWithSpan(span, ErrZeroSizeFile, "")
	}
	if len(Image) > int(s.cfg.MaxImageSize) {
		return s.errHandler.createExtendedErrorResponceWithSpan(span,
			ErrImageTooLarge,
			"",
			fmt.Sprintf("max image size: %d, file size: %d",
				s.cfg.MaxImageSize, s.cfg.MaxImageSize),
		)
	}

	s.logger.Info("Checking filetype")
	if fileType := s.detectFileType(&Image); fileType != "image" {
		return s.errHandler.createExtendedErrorResponceWithSpan(span, ErrUnsupportedFileType, "", "unsupported file type")
	}

	span.SetTag("grpc.status", codes.OK)
	return nil
}

func (s *ImagesStorageService) StreamingUploadImage(
	stream img_storage_serv.ImagesStorageServiceV1_StreamingUploadImageServer,
) error {
	span, ctx := opentracing.StartSpanFromContext(
		stream.Context(),
		"ImagesStorageService.StreamingUploadImage",
	)
	defer span.Finish()
	s.logger.Info("Start receiving image data")

	req, imageData, err := s.receiveUploadImage(ctx, stream)
	if err != nil {
		span.SetTag("grpc.status", status.Code(err))
		ext.LogError(span, err)
		return err
	}
	if req == nil {
		return s.errHandler.createErrorResponceWithSpan(span, ErrReceivedNilRequest, "")
	}

	s.logger.Info("Image data received. Calling upload method")
	res, err := s.UploadImage(ctx, &img_storage_serv.UploadImageRequest{
		Image:    imageData,
		Category: req.Category,
	})
	if err != nil {
		span.SetTag("grpc.status", status.Code(err))
		ext.LogError(span, err)
		return err // error alredy logged
	}
	if err = stream.SendAndClose(&img_storage_serv.UploadImageResponce{ImageId: res.ImageId}); err != nil {
		return s.errHandler.createErrorResponceWithSpan(span, err, "can't send response")
	}

	s.logger.Info("Responce successfully send")
	span.SetTag("grpc.status", codes.OK)
	return nil
}

func (s *ImagesStorageService) receiveUploadImage(ctx context.Context,
	stream img_storage_serv.ImagesStorageServiceV1_StreamingUploadImageServer) (*img_storage_serv.StreamingUploadImageRequest,
	[]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx,
		"ImagesStorageService.receiveUploadImage")
	defer span.Finish()

	var firstReq *img_storage_serv.StreamingUploadImageRequest
	imageData := bytes.Buffer{}
	for {
		err := stream.Context().Err()
		if err != nil {
			return nil, []byte{}, s.errHandler.createErrorResponceWithSpan(span, err, "")
		}

		s.logger.Info("Waiting to receive more data")

		req, err := stream.Recv()
		if firstReq == nil && req != nil {
			firstReq = req
		}
		if err == io.EOF {
			s.logger.Info("No more data")
			err = nil
			span.SetTag("grpc.status", codes.OK)
			return firstReq, imageData.Bytes(), nil
		}

		chunkSize := len(req.Data)
		imageSize := imageData.Len() + chunkSize
		s.logger.Debugf("Received a chunk with size: %d", chunkSize)
		if imageSize > s.cfg.MaxImageSize {
			return nil, []byte{}, s.errHandler.createErrorResponceWithSpan(span,
				ErrImageTooLarge,
				fmt.Sprintf("image size: %d, max supported size: %d",
					imageSize, s.cfg.MaxImageSize),
			)
		}
		imageData.Write(req.Data)
	}
}

func (s *ImagesStorageService) GetImage(ctx context.Context,
	in *img_storage_serv.ImageRequest) (*httpbody.HttpBody, error) {
	s.logger.Info("Start getting image")
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImagesStorageService.GetImage")
	defer span.Finish()

	s.logger.Info("Calling storage to get image")
	image, err := s.imageStorage.GetImage(ctx, in.ImageId, in.Category)
	if errors.Is(err, repository.ErrNotExist) {
		return nil, s.errHandler.createExtendedErrorResponceWithSpan(span, ErrCantFindImageByID, "", "image not found")
	}
	if err != nil {
		return nil, s.errHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}
	s.logger.Info("Writting responce")
	span.SetTag("grpc.status", codes.OK)
	return &httpbody.HttpBody{ContentType: http.DetectContentType(image), Data: image}, nil
}

func (s *ImagesStorageService) IsImageExist(ctx context.Context,
	in *img_storage_serv.ImageRequest) (*img_storage_serv.ImageExistResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesStorageService.IsImageExist")
	defer span.Finish()

	imageExist := s.imageStorage.IsImageExist(ctx, in.ImageId, in.Category)
	span.SetTag("grpc.status", codes.OK)
	return &img_storage_serv.ImageExistResponce{ImageExist: imageExist}, nil
}

func (s *ImagesStorageService) DeleteImage(ctx context.Context,
	in *img_storage_serv.ImageRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImagesStorageService.DeleteImage")
	defer span.Finish()

	err := s.imageStorage.DeleteImage(ctx, in.ImageId, in.Category)
	if errors.Is(err, repository.ErrNotExist) {
		return nil, s.errHandler.createExtendedErrorResponceWithSpan(span, ErrCantFindImageByID, "", "image not found")
	}
	if err != nil {
		return nil, s.errHandler.createErrorResponceWithSpan(span, ErrCantDeleteImage, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return &emptypb.Empty{}, nil
}

func (s *ImagesStorageService) saveImage(ctx context.Context, Image []byte, Category string) (string, error) {
	s.logger.Info("Start saving image")
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImagesStorageService.saveImage")
	defer span.Finish()

	s.logger.Info("Getting file extension")
	ext, _ := mime.ExtensionsByType(http.DetectContentType(Image))
	ImageId := uuid.NewString() + ext[0]

	s.metrics.IncBytesUploaded(len(Image))
	s.logger.Info("Calling storage to save image")
	if err := s.imageStorage.SaveImage(ctx, Image, ImageId, Category); err != nil {
		return "", s.errHandler.createErrorResponceWithSpan(span, ErrCantSaveImage, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return ImageId, nil
}

func (s *ImagesStorageService) ReplaceImage(
	ctx context.Context,
	in *img_storage_serv.ReplaceImageRequest,
) (*img_storage_serv.ReplaceImageResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesStorageService.ReplaceImage")
	defer span.Finish()

	if err := s.checkImage(ctx, in.ImageData); err != nil {
		span.SetTag("grpc.status", status.Code(err))
		ext.LogError(span, err)
		return nil, err
	}

	imageExist := s.imageStorage.IsImageExist(ctx, in.ImageId, in.Category)

	if in.CreateIfNotExist && !imageExist {
		imageID, err := s.saveImage(ctx, in.ImageData, in.Category)
		if err != nil {
			span.SetTag("grpc.status", status.Code(err))
			ext.LogError(span, err)
			return nil, err
		}
		return &img_storage_serv.ReplaceImageResponce{ImageId: imageID}, nil
	} else if !imageExist {
		return nil, s.errHandler.createErrorResponceWithSpan(span, ErrCantFindImageByID, "")
	}
	if err := s.imageStorage.RewriteImage(ctx, in.ImageData, in.ImageId, in.Category); err != nil {
		return nil, s.errHandler.createErrorResponceWithSpan(span, ErrCantReplaceImage, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return &img_storage_serv.ReplaceImageResponce{ImageId: in.ImageId}, nil
}

func (s *ImagesStorageService) detectFileType(fileData *[]byte) string {
	fileType := http.DetectContentType(*fileData)
	Type := strings.Split(fileType, "/")
	return Type[0]
}

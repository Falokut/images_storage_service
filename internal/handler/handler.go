package handler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/Falokut/images_storage_service/internal/models"
	"github.com/Falokut/images_storage_service/internal/service"
	img_storage_serv "github.com/Falokut/images_storage_service/pkg/images_storage_service/v1/protos"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Config struct {
	MaxImageSize int
}

type ImagesStorageServiceHandler struct {
	img_storage_serv.UnimplementedImagesStorageServiceV1Server
	logger  *logrus.Logger
	cfg     Config
	service service.ImagesStorageService
}

func NewImagesStorageServiceHandler(
	logger *logrus.Logger,
	cfg Config,
	service service.ImagesStorageService,
) *ImagesStorageServiceHandler {
	return &ImagesStorageServiceHandler{
		logger:  logger,
		cfg:     cfg,
		service: service,
	}
}

func (h *ImagesStorageServiceHandler) UploadImage(ctx context.Context,
	in *img_storage_serv.UploadImageRequest) (res *img_storage_serv.UploadImageResponse, err error) {
	defer h.handleError(&err)

	h.logger.Info("Image uploading starts")
	imageId, err := h.service.SaveImage(ctx, in.Image, in.Category)
	if err != nil {
		return
	}

	h.logger.Info("Image uploaded")
	return &img_storage_serv.UploadImageResponse{ImageId: imageId}, nil
}

func (h *ImagesStorageServiceHandler) StreamingUploadImage(
	stream img_storage_serv.ImagesStorageServiceV1_StreamingUploadImageServer,
) (err error) {
	ctx := stream.Context()
	defer h.handleError(&err)

	h.logger.Info("Start receiving image data")
	req, imageData, err := h.receiveUploadImage(stream)
	if err != nil {
		return
	}
	if req == nil {
		return status.Error(codes.InvalidArgument, "the received request is nil")
	}

	h.logger.Info("Image data received. Calling upload method")
	imageId, err := h.service.SaveImage(ctx, imageData, req.Category)
	if err != nil {
		return
	}

	if err = stream.SendAndClose(&img_storage_serv.UploadImageResponse{ImageId: imageId}); err != nil {
		return
	}

	h.logger.Info("Response successfully send")
	return nil
}

func (h *ImagesStorageServiceHandler) receiveUploadImage(
	stream img_storage_serv.ImagesStorageServiceV1_StreamingUploadImageServer) (firstReq *img_storage_serv.StreamingUploadImageRequest,
	image []byte, err error) {

	imageData := bytes.Buffer{}
	for {
		err = stream.Context().Err()
		if err != nil {
			err = status.Error(codes.Canceled, "")
			return
		}

		h.logger.Info("Waiting to receive more data")

		req, rerr := stream.Recv()
		if firstReq == nil && req != nil {
			firstReq = req
		}

		if rerr == io.EOF {
			h.logger.Info("No more data")
			return firstReq, imageData.Bytes(), nil
		}
		if rerr != nil {
			err = rerr
			return
		}

		chunkSize := len(req.Data)
		imageSize := imageData.Len() + chunkSize
		h.logger.Debugf("Received a chunk with size: %d", chunkSize)
		if imageSize > h.cfg.MaxImageSize {
			err = status.Errorf(codes.InvalidArgument,
				"image is too large max image size: %d, file size: %d",
				h.cfg.MaxImageSize, imageSize)
			return
		}
		imageData.Write(req.Data)
	}
}

func (h *ImagesStorageServiceHandler) GetImage(ctx context.Context,
	in *img_storage_serv.ImageRequest) (res *httpbody.HttpBody, err error) {
	defer h.handleError(&err)

	h.logger.Info("Start getting image")
	h.logger.Info("Calling storage to get image")
	image, err := h.service.GetImage(ctx, in.ImageId, in.Category)
	if err != nil {
		return
	}

	return &httpbody.HttpBody{ContentType: http.DetectContentType(image), Data: image}, nil
}

func (h *ImagesStorageServiceHandler) IsImageExist(ctx context.Context,
	in *img_storage_serv.ImageRequest) (res *img_storage_serv.ImageExistResponse, err error) {
	defer h.handleError(&err)

	imageExist, err := h.service.IsImageExist(ctx, in.ImageId, in.Category)
	if err != nil {
		return
	}

	return &img_storage_serv.ImageExistResponse{ImageExist: imageExist}, nil
}

func (h *ImagesStorageServiceHandler) DeleteImage(ctx context.Context,
	in *img_storage_serv.ImageRequest) (res *emptypb.Empty, err error) {
	defer h.handleError(&err)

	err = h.service.DeleteImage(ctx, in.ImageId, in.Category)
	if err != nil {
		return
	}
	res = &emptypb.Empty{}

	return
}

func (h *ImagesStorageServiceHandler) ReplaceImage(
	ctx context.Context,
	in *img_storage_serv.ReplaceImageRequest,
) (res *img_storage_serv.ReplaceImageResponse, err error) {
	defer h.handleError(&err)

	imageId, err := h.service.RewriteImage(ctx, in.ImageData, in.ImageId, in.Category, in.CreateIfNotExist)
	if err != nil {
		return
	}

	res = &img_storage_serv.ReplaceImageResponse{ImageId: imageId}
	return
}

func (h *ImagesStorageServiceHandler) handleError(err *error) {
	if err == nil || *err == nil {
		return
	}

	serviceErr := &models.ServiceError{}
	if errors.As(*err, &serviceErr) {
		*err = status.Error(convertServiceErrCodeToGrpc(serviceErr.Code), serviceErr.Msg)
	} else if _, ok := status.FromError(*err); !ok {
		e := *err
		*err = status.Error(codes.Unknown, e.Error())
	}
}

func convertServiceErrCodeToGrpc(code models.ErrorCode) codes.Code {
	switch code {
	case models.Internal:
		return codes.Internal
	case models.InvalidArgument:
		return codes.InvalidArgument
	case models.Unauthenticated:
		return codes.Unauthenticated
	case models.Conflict:
		return codes.AlreadyExists
	case models.NotFound:
		return codes.NotFound
	case models.Canceled:
		return codes.Canceled
	case models.DeadlineExceeded:
		return codes.DeadlineExceeded
	case models.PermissionDenied:
		return codes.PermissionDenied
	default:
		return codes.Unknown
	}
}

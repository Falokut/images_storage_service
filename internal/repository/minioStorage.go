package repository

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/Falokut/images_storage_service/internal/models"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type MinioStorage struct {
	logger  *logrus.Logger
	storage *minio.Client
}

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Secure          bool
}

func NewMinio(cfg MinioConfig) (*minio.Client, error) {
	return minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.Secure,
	})
}

func NewMinioStorage(logger *logrus.Logger, storage *minio.Client) *MinioStorage {
	return &MinioStorage{logger: logger, storage: storage}
}

const baseCategory = "image"

func (s *MinioStorage) SaveImage(ctx context.Context, img []byte, filename string, category string) (err error) {
	defer s.handleError(ctx, &err, "GetImage")

	category = baseCategory + category

	s.logger.Info("Start saving")
	exists, err := s.storage.BucketExists(ctx, category)
	if err != nil || !exists {
		s.logger.Warnf("no bucket %s. creating new one...", category)
		err := s.storage.MakeBucket(ctx, category, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	reader := bytes.NewReader(img)
	s.logger.Debugf("put new object %s to bucket %s", filename, category)
	_, err = s.storage.PutObject(ctx, category, filename, reader, int64(len(img)),
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": filename,
			},
			ContentType: "image/jpeg",
		})
	if err != nil {
		return
	}

	return
}

func (s *MinioStorage) GetImage(ctx context.Context, filename string, category string) (image []byte, err error) {
	defer s.handleError(ctx, &err, "GetImage")

	s.logger.Info("Start getting image")
	category = baseCategory + category

	obj, err := s.storage.GetObject(ctx, category, filename, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	defer obj.Close()
	objectInfo, err := obj.Stat()
	if err != nil {
		return
	}

	image = make([]byte, objectInfo.Size)
	_, err = obj.Read(image)
	if err == io.EOF {
		err = nil
	}

	return
}

func (s *MinioStorage) IsImageExist(ctx context.Context, filename string, category string) (exist bool, err error) {
	category = baseCategory + category

	_, err = s.storage.StatObject(ctx, category, filename, minio.StatObjectOptions{})
	s.handleError(ctx, &err, "IsImageExist")
	if models.Code(err) == models.NotFound {
		return false, nil
	}

	exist = true
	return
}

func (s *MinioStorage) DeleteImage(ctx context.Context, filename string, category string) (err error) {
	defer s.handleError(ctx, &err, "DeleteImage")

	category = baseCategory + category

	err = s.storage.RemoveObject(ctx, category, filename, minio.RemoveObjectOptions{ForceDelete: true})
	return
}

func (s *MinioStorage) RewriteImage(ctx context.Context, img []byte, filename string, category string) (err error) {
	defer s.handleError(ctx, &err, "RewriteImage")
	category = baseCategory + category

	s.logger.Info("Start saving")
	exists, err := s.storage.BucketExists(ctx, category)
	if err != nil || !exists {
		s.logger.Warnf("no bucket %s. creating new one...", category)
		err = s.storage.MakeBucket(ctx, category, minio.MakeBucketOptions{})
		if err != nil {
			return
		}
	}

	reader := bytes.NewReader(img)
	s.logger.Debugf("put new object %s to bucket %s", filename, category)
	_, err = s.storage.PutObject(ctx, category, filename, reader, int64(len(img)),
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": filename,
			},
			ContentType: "image/jpeg",
		})

	return
}

func (s *MinioStorage) handleError(ctx context.Context, err *error, functionName string) {
	if ctx.Err() != nil {
		var code models.ErrorCode
		switch {
		case errors.Is(ctx.Err(), context.Canceled):
			code = models.Canceled
		case errors.Is(ctx.Err(), context.DeadlineExceeded):
			code = models.DeadlineExceeded
		}
		*err = models.Error(code, ctx.Err().Error())
		s.logError(*err, functionName)
		return
	}

	if err == nil || *err == nil {
		return
	}

	s.logError(*err, functionName)
	var repoErr = &models.ServiceError{}
	if !errors.As(*err, &repoErr) {
		errResp := minio.ToErrorResponse(*err)
		switch errResp.StatusCode {
		case http.StatusNotFound:
			*err = models.Error(models.NotFound, "image not found")
		case http.StatusBadRequest:
			*err = models.Error(models.InvalidArgument, errResp.Message)
		default:
			*err = models.Error(models.Internal, errResp.Message)
		}

	}
}

func (s *MinioStorage) logError(err error, functionName string) {
	if err == nil {
		return
	}

	var repoErr = &models.ServiceError{}
	if errors.As(err, &repoErr) {
		s.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           repoErr.Msg,
				"error.code":          repoErr.Code,
			},
		).Error("minio images storage error occurred")
	} else {
		s.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("minio images storage error occurred")
	}
}

package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"

	"github.com/Falokut/go-kit/log"
	"github.com/Falokut/images_storage_service/domain"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	logger  log.Logger
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

func NewMinioStorage(logger log.Logger, storage *minio.Client) *MinioStorage {
	return &MinioStorage{logger: logger, storage: storage}
}

func (s *MinioStorage) SaveImage(ctx context.Context, img []byte, filename string, category string) error {
	bucketName := getBucketName(category)
	exists, err := s.storage.BucketExists(ctx, bucketName)
	if err != nil {
		return errors.WithMessage(err, "check is bucket exists")
	}
	if !exists {
		s.logger.Info(ctx, "creating bucket", log.Any("bucketName", bucketName))
		err = s.storage.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return errors.WithMessage(err, "make bucket")
		}
	}

	reader := bytes.NewReader(img)
	s.logger.Info(ctx, "save file", log.Any("bucketName", bucketName), log.Any("filename", filename))
	_, err = s.storage.PutObject(ctx, bucketName, filename, reader, int64(len(img)),
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": filename,
			},
			ContentType: "image/jpeg",
		})
	if err != nil {
		return errors.WithMessage(err, "put object")
	}

	return nil
}

func (s *MinioStorage) GetImage(ctx context.Context, filename string, category string) ([]byte, error) {
	bucketName := getBucketName(category)
	obj, err := s.storage.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
	errResp := minio.ToErrorResponse(err)
	switch {
	case errResp.StatusCode == http.StatusNotFound:
		return nil, domain.ErrImageNotFound
	case err != nil:
		return nil, errors.WithMessage(err, "get object")
	}
	defer obj.Close()

	objectInfo, err := obj.Stat()
	if err != nil {
		return nil, errors.WithMessage(err, "object stat")
	}

	image := make([]byte, objectInfo.Size)
	_, err = obj.Read(image)
	if err != nil && err != io.EOF {
		return nil, errors.WithMessage(err, "read object")
	}

	return image, nil
}

func (s *MinioStorage) IsImageExist(ctx context.Context, filename string, category string) (exist bool, err error) {
	bucketName := getBucketName(category)
	_, err = s.storage.StatObject(ctx, bucketName, filename, minio.StatObjectOptions{})
	switch {
	case minio.ToErrorResponse(err).StatusCode == http.StatusNotFound:
		return false, nil
	case err != nil:
		return false, errors.WithMessage(err, "stat object")
	default:
		return true, nil
	}
}

func (s *MinioStorage) DeleteImage(ctx context.Context, filename string, category string) error {
	bucketName := getBucketName(category)
	err := s.storage.RemoveObject(ctx, bucketName, filename, minio.RemoveObjectOptions{ForceDelete: true})
	switch {
	case minio.ToErrorResponse(err).StatusCode == http.StatusNotFound:
		return domain.ErrImageNotFound
	case err != nil:
		return errors.WithMessage(err, "stat object")
	default:
		return nil
	}
}

func (s *MinioStorage) ReplaceImage(ctx context.Context, img []byte, filename string, category string) error {
	bucketName := getBucketName(category)
	exists, err := s.storage.BucketExists(ctx, bucketName)
	if err != nil {
		return errors.WithMessage(err, "check is bucket exists")
	}
	if !exists {
		return domain.ErrImageNotFound
	}

	reader := bytes.NewReader(img)
	s.logger.Info(ctx, "save file", log.Any("bucketName", bucketName), log.Any("filename", filename))
	_, err = s.storage.PutObject(ctx, bucketName, filename, reader, int64(len(img)),
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": filename,
			},
			ContentType: "image/jpeg",
		})
	if err != nil {
		return errors.WithMessage(err, "put object")
	}
	return nil
}

func getBucketName(category string) string {
	return fmt.Sprintf("image-%s", category)
}

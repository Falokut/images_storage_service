package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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

const baseCategory = "image-"

func (s *MinioStorage) SaveImage(ctx context.Context, img []byte, filename string, category string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MinioStorage.SaveImage")
	defer span.Finish()
	category = baseCategory + category

	s.logger.Info("Start saving")
	exists, err := s.storage.BucketExists(ctx, category)
	if err != nil || !exists {
		s.logger.Warnf("no bucket %s. creating new one...", category)
		err := s.storage.MakeBucket(ctx, category, minio.MakeBucketOptions{})
		if err != nil {
			ext.LogError(span, err)
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
		ext.LogError(span, err)
		return err
	}

	span.SetTag("error", false)
	return nil
}

func (s *MinioStorage) GetImage(ctx context.Context, filename string, category string) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MinioStorage.GetImage")
	defer span.Finish()
	s.logger.Info("Start getting image")
	category = baseCategory + category

	obj, err := s.storage.GetObject(ctx, category, filename, minio.GetObjectOptions{})
	if err != nil {
		ext.LogError(span, err)
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" || errResponse.Code == "NoSuchBucket" {
			return []byte{}, ErrNotExist
		}
		return []byte{}, err
	}
	defer obj.Close()
	objectInfo, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file. err: %w", err)
	}

	buffer := make([]byte, objectInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		ext.LogError(span, err)
		return []byte{}, err
	}

	span.SetTag("error", false)
	return buffer, nil
}

func (s *MinioStorage) IsImageExist(ctx context.Context, filename string, category string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MinioStorage.IsImageExist")
	defer span.Finish()
	category = baseCategory + category

	_, err := s.storage.StatObject(ctx, category, filename, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			span.SetTag("error", false)
			return false
		}
		ext.LogError(span, err)
		return false
	}

	span.SetTag("error", false)
	return true
}

func (s *MinioStorage) DeleteImage(ctx context.Context, filename string, category string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MinioStorage.DeleteImage")
	defer span.Finish()
	category = baseCategory + category

	err := s.storage.RemoveObject(ctx, category, filename, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		ext.LogError(span, err)

		if errResponse.Code == "NoSuchKey" {
			return ErrNotExist
		}
		return err
	}

	span.SetTag("error", false)
	return nil
}

func (s *MinioStorage) RewriteImage(ctx context.Context, img []byte, filename string, category string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MinioStorage.RewriteImage")
	defer span.Finish()
	category = baseCategory + category

	s.logger.Info("Start saving")
	exists, err := s.storage.BucketExists(ctx, category)
	if err != nil || !exists {
		s.logger.Warnf("no bucket %s. creating new one...", category)
		err := s.storage.MakeBucket(ctx, category, minio.MakeBucketOptions{})
		if err != nil {
			ext.LogError(span, err)
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
		ext.LogError(span, err)
		return err
	}

	span.SetTag("error", false)
	return nil
}

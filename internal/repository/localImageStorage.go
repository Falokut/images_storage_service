package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Falokut/images_storage_service/internal/models"
	"github.com/sirupsen/logrus"
)

type LocalImageStorage struct {
	logger   *logrus.Logger
	basePath string
}

var wg sync.WaitGroup

func NewLocalStorage(logger *logrus.Logger, baseStoragePath string) *LocalImageStorage {
	return &LocalImageStorage{logger: logger, basePath: baseStoragePath}
}

func (s *LocalImageStorage) Shutdown() {
	s.logger.Info("Shutting down local image storage")
	wg.Wait()
}

func (s *LocalImageStorage) SaveImage(ctx context.Context, img []byte, filename string, relativePath string) (err error) {
	defer s.handleError(ctx, &err, "SaveImage")

	s.logger.Info("Start saving image")

	wg.Add(1)
	defer wg.Done()
	s.logger.Info("Creating a file")
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))

	s.logger.Debugf("Saving relativePath: %s", relativePath)

	err = os.MkdirAll(filepath.Dir(relativePath), 0755)
	if err != nil || os.IsExist(err) {
		err = errors.New("can't create dir for file")
		return
	}

	f, err := os.OpenFile(relativePath, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0660)
	if err != nil {
		return
	}

	s.logger.Info("Writing data into file")
	_, err = f.Write(img)
	s.logger.Info("Image saving is completed")
	return
}

func (s *LocalImageStorage) GetImage(ctx context.Context, filename string, relativePath string) (image []byte, err error) {
	defer s.handleError(ctx, &err, "GetImage")

	s.logger.Info("Start getting image")

	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))

	image, err = os.ReadFile(relativePath)
	return
}

func (s *LocalImageStorage) IsImageExist(ctx context.Context, filename string, relativePath string) (exist bool, err error) {

	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))
	_, err = os.Stat(relativePath)

	s.handleError(ctx, &err, "IsImageExist")
	if models.Code(err) == models.NotFound {
		return false, nil
	}

	exist = true
	return
}

func (s *LocalImageStorage) DeleteImage(ctx context.Context, filename string, relativePath string) (err error) {
	defer s.handleError(ctx, &err, "DeleteImage")
	wg.Add(1)
	defer wg.Done()
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))
	err = os.Remove(relativePath)

	return nil
}

func (s *LocalImageStorage) RewriteImage(ctx context.Context, img []byte, filename string, relativePath string) (err error) {
	defer s.handleError(ctx, &err, "RewriteImage")

	s.logger.Info("Start getting image file")
	wg.Add(1)
	defer wg.Done()
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))

	f, err := os.OpenFile(relativePath, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return
	}

	s.logger.Info("Truncate file")
	f.Truncate(0)
	f.Seek(0, 0)
	s.logger.Info("Writing data into file")
	_, err = f.Write(img)
	s.logger.Info("Image saving is completed")
	return 
}

func (s *LocalImageStorage) handleError(ctx context.Context, err *error, functionName string) {
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
		switch {
		case errors.Is(*err, os.ErrNotExist):
			*err = models.Error(models.NotFound, "image not found")
		default:
			*err = models.Error(models.Internal, "images storage internal error")
		}
	}
}

func (s *LocalImageStorage) logError(err error, functionName string) {
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
		).Error("local images storage error occurred")
	} else {
		s.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("local images storage error occurred")
	}
}

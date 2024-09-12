package assembly

import (
	"context"
	"strings"

	"github.com/Falokut/go-kit/http/endpoint"
	"github.com/Falokut/go-kit/http/router"
	"github.com/Falokut/go-kit/log"
	"github.com/Falokut/images_storage_service/conf"
	"github.com/Falokut/images_storage_service/controller"
	"github.com/Falokut/images_storage_service/repository"
	"github.com/Falokut/images_storage_service/routes"
	"github.com/Falokut/images_storage_service/service"
	"github.com/pkg/errors"
)

type Config struct {
	Mux *router.Router
}


func Locator(_ context.Context, logger log.Logger, cfg conf.LocalConfig) (Config, error) {
	var imagesStorage service.ImageStorage
	storageMode := strings.ToUpper(cfg.StorageMode)
	switch storageMode {
	case "MINIO":
		minioStorage, err := repository.NewMinio(repository.MinioConfig{
			Endpoint:        cfg.MinioConfig.Endpoint,
			AccessKeyID:     cfg.MinioConfig.AccessKeyID,
			SecretAccessKey: cfg.MinioConfig.SecretAccessKey,
			Secure:          cfg.MinioConfig.Secure,
		})
		imagesStorage = repository.NewMinioStorage(logger, minioStorage)
		if err != nil {
			return Config{}, errors.WithMessage(err, "new minio client")
		}
	default:
		imagesStorage = repository.NewLocalStorage(cfg.BaseLocalStoragePath)
	}

	imagesService := service.NewImages(imagesStorage, cfg.MaxImageSize*mb)
	imagesController := controller.NewImages(imagesService, cfg.MaxImageSize*mb)
	c := routes.Router{
		Images: imagesController,
	}
	return Config{Mux: c.InitRoutes(endpoint.DefaultWrapper(logger))}, nil
}

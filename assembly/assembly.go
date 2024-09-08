package assembly

import (
	"context"
	"io"
	"strings"

	server "github.com/Falokut/grpc_rest_server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"

	"github.com/Falokut/go-kit/log"
	"github.com/Falokut/images_storage_service/app"
	"github.com/Falokut/images_storage_service/conf"
	"github.com/Falokut/images_storage_service/controller"
	img_storage_serv "github.com/Falokut/images_storage_service/pkg/images_storage_service/v1/protos"
	jaegerTracer "github.com/Falokut/images_storage_service/pkg/jaeger"
	"github.com/Falokut/images_storage_service/pkg/metrics"
	"github.com/Falokut/images_storage_service/repository"
	"github.com/Falokut/images_storage_service/service"
	"github.com/pkg/errors"
)

const (
	kb = 8 << 10
	mb = kb << 10
)

type Assembly struct {
	logger       log.Logger
	server       server.Server
	cfg          *conf.LocalConfig
	metric       metrics.Metrics
	jaegerCloser io.Closer
}

func New(ctx context.Context, logger log.Logger, cfg *conf.LocalConfig) (*Assembly, error) {
	metric, closer, err := getMetrics(cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "get metrics")
	}
	locatorCfg, err := Locator(ctx, logger, cfg, metric)
	if err != nil {
		return nil, errors.WithMessage(err, "init locator")
	}
	server := server.NewServer(logger, locatorCfg.Mux)
	return &Assembly{
		logger:       logger,
		server:       server,
		jaegerCloser: closer,
		metric:       metric,
		cfg:          cfg,
	}, nil
}

type Config struct {
	Mux any
}

func Locator(_ context.Context, logger log.Logger, cfg *conf.LocalConfig, metric metrics.Metrics) (Config, error) {
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

	imagesService := service.NewImages(metric, imagesStorage, service.Config{
		MaxImageSize: cfg.MaxImageSize,
	})
	imagesController := controller.NewImagesStorageServiceHandler(logger, imagesService, cfg.MaxImageSize)
	return Config{
		Mux: imagesController,
	}, nil
}

func (a *Assembly) Runners() []app.RunnerFunc {
	return []app.RunnerFunc{
		func(ctx context.Context) error {
			err := a.server.Run(ctx, getListenServerConfig(a.cfg), a.metric, nil, nil)
			if err != nil {
				return errors.WithMessage(err, "run server")
			}
			return nil
		},
		func(ctx context.Context) error {
			if !a.cfg.EnableMetrics {
				return nil
			}
			err := metrics.RunMetricServer(a.cfg.PrometheusConfig.ServerConfig)
			if err != nil {
				return errors.WithMessage(err, "run metrics server")
			}
			return nil
		},
	}
}

func (a *Assembly) Closers() []app.CloserFunc {
	return []app.CloserFunc{
		func(ctx context.Context) error {
			a.server.Shutdown(ctx)
			return nil
		},
		func(context.Context) error {
			return a.jaegerCloser.Close()
		},
	}
}

func getMetrics(cfg *conf.LocalConfig) (metrics.Metrics, io.Closer, error) {
	if !cfg.EnableMetrics {
		return metrics.EmptyMetrics{}, nil, nil
	}

	tracer, closer, err := jaegerTracer.InitJaeger(cfg.JaegerConfig)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "init jaeger tracer")
	}
	opentracing.SetGlobalTracer(tracer)

	metric, err := metrics.CreateMetrics(cfg.PrometheusConfig.Name)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "error while creating metrics")
	}

	return metric, closer, nil
}

func getListenServerConfig(cfg *conf.LocalConfig) server.Config {
	return server.Config{
		Mode:           cfg.Listen.Mode,
		Host:           cfg.Listen.Host,
		Port:           cfg.Listen.Port,
		AllowedHeaders: cfg.Listen.AllowedHeaders,
		ServiceDesc:    &img_storage_serv.ImagesStorageServiceV1_ServiceDesc,
		RegisterRestHandlerServer: func(ctx context.Context,
			mux *runtime.ServeMux, service any) error {
			serv, ok := service.(img_storage_serv.ImagesStorageServiceV1Server)
			if !ok {
				return errors.New("can't convert service")
			}
			return img_storage_serv.RegisterImagesStorageServiceV1HandlerServer(context.Background(),
				mux, serv)
		},
		MaxRequestSize:  cfg.Listen.MaxRequestSize * mb,
		MaxResponceSize: cfg.Listen.MaxResponseSize * mb,
	}
}

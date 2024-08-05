package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	server "github.com/Falokut/grpc_rest_server"
	"github.com/Falokut/healthcheck"
	"github.com/Falokut/images_storage_service/internal/config"
	"github.com/Falokut/images_storage_service/internal/handler"
	"github.com/Falokut/images_storage_service/internal/repository"
	"github.com/Falokut/images_storage_service/internal/service"
	img_storage_serv "github.com/Falokut/images_storage_service/pkg/images_storage_service/v1/protos"
	jaegerTracer "github.com/Falokut/images_storage_service/pkg/jaeger"
	"github.com/Falokut/images_storage_service/pkg/logging"
	"github.com/Falokut/images_storage_service/pkg/metrics"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	kb = 8 << 10
	mb = kb << 10
)

func main() {
	logging.NewEntry(logging.ConsoleOutput)
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Logger.SetLevel(logLevel)

	var metric metrics.Metrics
	shutdown := make(chan error, 1)
	if cfg.EnableMetrics {
		tracer, closer, err := jaegerTracer.InitJaeger(cfg.JaegerConfig)
		if err != nil {
			logger.Errorf("cannot create tracer %v", err)
			return
		}
		logger.Info("Jaeger connected")
		defer closer.Close()

		opentracing.SetGlobalTracer(tracer)

		logger.Info("Metrics initializing")
		metric, err = metrics.CreateMetrics(cfg.PrometheusConfig.Name)
		if err != nil {
			logger.Errorf("error while creating metrics %v", err)
			return
		}

		go func() {
			logger.Info("Metrics server running")
			if err = metrics.RunMetricServer(cfg.PrometheusConfig.ServerConfig); err != nil {
				logger.Errorf("Shutting down, error while running metrics server %v", err)
				shutdown <- err
				return
			}
		}()
	} else {
		metric = metrics.EmptyMetrics{}
	}

	var storage repository.ImageStorage
	cfg.StorageMode = strings.ToUpper(cfg.StorageMode)
	switch cfg.StorageMode {
	case "MINIO":
		minioStorage, err := repository.NewMinio(repository.MinioConfig{
			Endpoint:        cfg.MinioConfig.Endpoint,
			AccessKeyID:     cfg.MinioConfig.AccessKeyID,
			SecretAccessKey: cfg.MinioConfig.SecretAccessKey,
			Secure:          cfg.MinioConfig.Secure,
		})
		if err != nil {
			logger.Errorf("Shutting down, error while connecting to the minio storage: %v", err)
			return
		}
		storage = repository.NewMinioStorage(logger.Logger, minioStorage)
		go func() {
			logger.Info("Healthcheck initializing")
			healthcheckManager := healthcheck.NewHealthManager(logger.Logger,
				[]healthcheck.HealthcheckResource{}, cfg.HealthcheckPort, func(ctx context.Context) error {
					healthcheckTime, ok := ctx.Deadline()
					if !ok {
						healthcheckTime = time.Now().Add(time.Second * 5)
					}
					dur := time.Until(healthcheckTime)
					cancel, err := minioStorage.HealthCheck(dur)
					if err != nil {
						return err
					}
					defer cancel()

					if minioStorage.IsOffline() {
						return errors.New("minio storage offline")
					}

					return nil
				})
			if err := healthcheckManager.RunHealthcheckEndpoint(); err != nil {
				logger.Errorf("Shutting down, error while running healthcheck endpoint %s", err.Error())
				shutdown <- err
				return
			}
		}()
	default:
		logger.Info("Local storage initializing")
		storage = repository.NewLocalStorage(logger.Logger, cfg.BaseLocalStoragePath)

		healthcheckManager := healthcheck.NewHealthManager(logger.Logger,
			[]healthcheck.HealthcheckResource{}, cfg.HealthcheckPort, func(ctx context.Context) error { return nil })
		if err := healthcheckManager.RunHealthcheckEndpoint(); err != nil {
			logger.Errorf("Shutting down, error while running healthcheck endpoint %s", err.Error())
			shutdown <- err
			return
		}
	}

	logger.Info("Service initializing")
	service := service.NewImagesStorageService(logger.Logger, metric, storage,
		service.Config{MaxImageSize: cfg.MaxImageSize * mb})
		
	logger.Info("Handler initializing")
	handler := handler.NewImagesStorageServiceHandler(logger.Logger,
		handler.Config{MaxImageSize: cfg.MaxImageSize * mb}, service)

	logger.Info("Server initializing")
	s := server.NewServer(logger.Logger, handler)
	go func() {
		if err := s.Run(getListenServerConfig(cfg), metric, nil, nil); err != nil {
			logger.Errorf("Shutting down, error while running server %s", err.Error())
			shutdown <- err
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	select {
	case <-quit:
		break
	case <-shutdown:
		break
	}

	s.Shutdown()
}

func getListenServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Mode:           cfg.Listen.Mode,
		Host:           cfg.Listen.Host,
		Port:           cfg.Listen.Port,
		AllowedHeaders: cfg.Listen.AllowedHeaders,
		ServiceDesc:    &img_storage_serv.ImagesStorageServiceV1_ServiceDesc,
		RegisterRestHandlerServer: func(ctx context.Context, mux *runtime.ServeMux, service any) error {
			serv, ok := service.(img_storage_serv.ImagesStorageServiceV1Server)
			if !ok {
				return errors.New("can't convert")
			}
			return img_storage_serv.RegisterImagesStorageServiceV1HandlerServer(context.Background(),
				mux, serv)
		},
		MaxRequestSize:  cfg.Listen.MaxRequestSize * mb,
		MaxResponceSize: cfg.Listen.MaxResponseSize * mb,
	}
}

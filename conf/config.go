package conf

import (
	"github.com/Falokut/images_storage_service/pkg/jaeger"
	"github.com/Falokut/images_storage_service/pkg/metrics"
)

type LocalConfig struct {
	BaseLocalStoragePath string `yaml:"base_local_storage_path" env:"BASE_LOCAL_STORAGE_PATH"`
	StorageMode string `yaml:"storage_mode" env:"STORAGE_MODE"` // MINIO or LOCAL
	MinioConfig struct {
		Endpoint        string `yaml:"endpoint" env:"MINIO_ENDPOINT"`
		AccessKeyID     string `yaml:"access_key_id" env:"MINIO_ACCESS_KEY_ID"`
		SecretAccessKey string `yaml:"secret_access_key" env:"MINIO_SECRET_ACCESS_KEY"`
		Secure          bool   `yaml:"secure" env:"MINIO_SECURE"`
	} `yaml:"minio"`
	MaxImageSize int `yaml:"max_image_size" env:"MAX_IMAGE_SIZE"` // in mb
	Listen       struct {
		Host            string   `yaml:"host" env:"HOST"`
		Port            string   `yaml:"port" env:"PORT"`
		AllowedHeaders  []string `yaml:"allowed_headers"`
		Mode            string   `yaml:"server_mode" env:"SERVER_MODE"`
		MaxRequestSize  int      `yaml:"max_request_size" env:"MAX_REQUEST_SIZE"`
		MaxResponseSize int      `yaml:"max_response_size" env:"MAX_RESPONSE_SIZE"`
	} `yaml:"listen"`

	EnableMetrics    bool `yaml:"enable_metrics" env:"ENABLE_METRICS"`
	PrometheusConfig struct {
		Name         string                      `yaml:"service_name" env:"PROMETHEUS_SERVICE_NAME"`
		ServerConfig metrics.MetricsServerConfig `yaml:"server_config"`
	} `yaml:"prometheus"`
	JaegerConfig jaeger.Config `yaml:"jaeger"`
}
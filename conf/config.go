package conf

import (
	"context"
	"fmt"
	"sync"

	"github.com/Falokut/images_storage_service/pkg/jaeger"

	"github.com/Falokut/go-kit/log"
	"github.com/Falokut/images_storage_service/pkg/metrics"
	"github.com/ilyakaznacheev/cleanenv"
)

type LocalConfig struct {
	Log struct {
		LogLevel      string `yaml:"level" env:"LOG_LEVEL"`
		ConsoleOutput bool   `yaml:"console_output" env:"LOG_CONSOLE_OUTPUT"`
		Filepath      string `yaml:"filepath" env:"LOG_FILEPATH"`
	} `yaml:"log"`
	App struct {
		Id      string `yaml:"id" env:"APP_ID"`
		Version string `yaml:"version" env:"APP_VERSION"`
	} `yaml:"app"`
	BaseLocalStoragePath string `yaml:"base_local_storage_path" env:"BASE_LOCAL_STORAGE_PATH"`

	StorageMode string `yaml:"storage_mode" env:"STORAGE_MODE"` // MINIO or LOCAL
	MinioConfig struct {
		Endpoint        string `yaml:"endpoint" env:"MINIO_ENDPOINT"`
		AccessKeyID     string `yaml:"access_key_id" env:"MINIO_ACCESS_KEY_ID"`
		SecretAccessKey string `yaml:"secret_access_key" env:"MINIO_SECRET_ACCESS_KEY"`
		Secure          bool   `yaml:"secure" env:"MINIO_SECURE"`
	} `yaml:"minio"`
	MaxImageSize int `yaml:"max_image_size" env:"MAX_IMAGE_SIZE"`
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

var instance *LocalConfig
var once sync.Once

const configsPath = "conf/"

func GetLocalConfig() *LocalConfig {
	once.Do(func() {
		instance = &LocalConfig{}

		if err := cleanenv.ReadConfig(configsPath+"config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger, _ := log.NewFromConfig(log.Config{
				Loglevel: log.FatalLevel,
				Output: log.OutputConfig{
					Console: true,
				},
			})
			logger.Fatal(context.Background(), fmt.Sprint(help, err))
		}
	})

	return instance
}

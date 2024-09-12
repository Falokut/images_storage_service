package conf

import "github.com/Falokut/go-kit/config"

type LocalConfig struct {
	BaseLocalStoragePath string `yaml:"base_local_storage_path" env:"BASE_LOCAL_STORAGE_PATH"`
	StorageMode          string `yaml:"storage_mode" env:"STORAGE_MODE"` // MINIO or LOCAL
	MinioConfig          struct {
		Endpoint        string `yaml:"endpoint" env:"MINIO_ENDPOINT"`
		AccessKeyID     string `yaml:"access_key_id" env:"MINIO_ACCESS_KEY_ID"`
		SecretAccessKey string `yaml:"secret_access_key" env:"MINIO_SECRET_ACCESS_KEY"`
		Secure          bool   `yaml:"secure" env:"MINIO_SECURE"`
	} `yaml:"minio"`
	MaxImageSize int           `yaml:"max_image_size" env:"MAX_IMAGE_SIZE"` // in mb
	Listen       config.Listen `yaml:"listen"`
}

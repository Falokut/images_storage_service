package repository

import (
	"context"
)

//go:generate mockgen -source=repository.go -destination=mocks/imageStorage.go
type ImageStorage interface {
	SaveImage(ctx context.Context, img []byte, filename string, category string) error
	GetImage(ctx context.Context, imageID string, category string) ([]byte, error)
	IsImageExist(ctx context.Context, imageID string, category string) (bool, error)
	DeleteImage(ctx context.Context, imageID string, category string) error
	RewriteImage(ctx context.Context, img []byte, filename string, category string) error
}

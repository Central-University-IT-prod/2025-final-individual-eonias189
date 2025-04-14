package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Connect(cfg Config) (*minio.Client, error) {
	return minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.User, cfg.Password, ""),
	})
}

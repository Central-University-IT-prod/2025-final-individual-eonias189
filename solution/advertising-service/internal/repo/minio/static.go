package minio

import (
	"advertising/advertising-service/internal/models"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

type StaticRepo struct {
	cli        *minio.Client
	bucketName string
}

func NewStaticRepo(cli *minio.Client, bucketName string) *StaticRepo {
	return &StaticRepo{
		cli:        cli,
		bucketName: bucketName,
	}
}

func (sr *StaticRepo) SaveStatic(ctx context.Context, name string, static models.Static) error {
	op := "StaticRepo.SaveStatic"

	_, err := sr.cli.PutObject(ctx, sr.bucketName, name, static.Data, static.Size, minio.PutObjectOptions{
		ContentType: static.ContentType,
	})

	if err != nil {
		return fmt.Errorf("%s: cli.PutObject: %w", op, err)
	}

	return nil
}

func (sr *StaticRepo) LoadStatic(ctx context.Context, name string) (models.Static, error) {
	op := "StaticRepo.LoadStatic"

	obj, err := sr.cli.GetObject(ctx, sr.bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return models.Static{}, fmt.Errorf("%s: cli.GetObject: %w", op, err)
	}

	stat, err := obj.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return models.Static{}, models.ErrStaticNotFound
		}
		return models.Static{}, fmt.Errorf("%s: obj.Stat: %w", op, err)
	}

	return models.Static{
		Data:        obj,
		Size:        stat.Size,
		ContentType: stat.ContentType,
	}, nil
}

func (sr StaticRepo) DeleteStatic(ctx context.Context, name string) error {
	op := "StaticRepo.DeleteStatic"

	err := sr.cli.RemoveObject(ctx, sr.bucketName, name, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("%s: cli.RemoveObject: %w", op, err)
	}

	return nil
}

package repository

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
}

type S3Repository struct {
	client     *minio.Client
	bucketName string
}

func NewS3Repository(cfg S3Config) (*S3Repository, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &S3Repository{client: client, bucketName: cfg.BucketName}, nil
}

func (r *S3Repository) Upload(key string, file multipart.File, contentType string) error {
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	_, err = r.client.PutObject(
		context.Background(),
		r.bucketName,
		key,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	return err
}

func (r *S3Repository) Download(key string) ([]byte, error) {
	obj, err := r.client.GetObject(context.Background(), r.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	return io.ReadAll(obj)
}

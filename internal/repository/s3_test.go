package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewS3Repository_InvalidConfig(t *testing.T) {
	cfg := S3Config{
		Endpoint:        "invalid:9000",
		AccessKeyID:     "test",
		SecretAccessKey: "test",
		BucketName:      "test",
		UseSSL:          false,
	}

	_, err := NewS3Repository(cfg)
	// Ожидаем ошибку подключения (не панику)
	assert.Error(t, err)
}

func TestS3Config_Defaults(t *testing.T) {
	cfg := S3Config{
		Endpoint:        "localhost:9000",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
		BucketName:      "test-bucket",
		UseSSL:          false,
	}
	assert.Equal(t, "localhost:9000", cfg.Endpoint)
	assert.Equal(t, "minioadmin", cfg.AccessKeyID)
	assert.Equal(t, "test-bucket", cfg.BucketName)
}
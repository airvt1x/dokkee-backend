package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresDB_InvalidConfig(t *testing.T) {
	cfg := Config{
		Host:     "invalid-host",
		Port:     "9999",
		Username: "user",
		Password: "pass",
		DBName:   "db",
		SSLMode:  "disable",
	}

	db, err := NewPostgresDB(cfg)
	assert.Error(t, err)
	assert.Nil(t, db)
}

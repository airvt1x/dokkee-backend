package service

import (
	"testing"

	"github.com/airvt1x/dokkee-backend/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	repos := &repository.Repository{}
	svc := NewService(repos)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.Authorization)
	assert.NotNil(t, svc.Document)
	assert.NotNil(t, svc.Result)
	assert.NotNil(t, svc.Audit)
}

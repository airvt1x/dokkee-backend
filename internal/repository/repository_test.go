package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository(nil, nil)
	assert.NotNil(t, repo)
}

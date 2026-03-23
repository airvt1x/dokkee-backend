package service

import (
	"testing"

	"github.com/airvt1x/dokkee-backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthorization struct {
	mock.Mock
}

func (m *MockAuthorization) CreateUser(user dokkee.User) (int, error) {
	args := m.Called(user)
	return args.Int(0), args.Error(1)
}

func (m *MockAuthorization) GetUser(email, password string) (dokkee.User, error) {
	args := m.Called(email, password)
	return args.Get(0).(dokkee.User), args.Error(1)
}

func TestAuthService_CreateUser(t *testing.T) {
	mockRepo := new(MockAuthorization)
	service := NewAuthService(mockRepo)

	user := dokkee.User{
		Username:  "testuser",
		Password:  "password",
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Phone:     "+1234567890",
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("dokkee.User")).Return(1, nil)

	id, err := service.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GenerateToken(t *testing.T) {
	mockRepo := new(MockAuthorization)
	service := NewAuthService(mockRepo)

	email := "test@example.com"
	password := "password"

	user := dokkee.User{
		Id:    1,
		Email: email,
	}

	mockRepo.On("GetUser", email, mock.AnythingOfType("string")).Return(user, nil)

	token, err := service.GenerateToken(email, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_ParseToken(t *testing.T) {
	mockRepo := new(MockAuthorization)
	service := NewAuthService(mockRepo)

	user := dokkee.User{Id: 1}
	mockRepo.On("GetUser", "test@example.com", mock.AnythingOfType("string")).Return(user, nil)

	token, err := service.GenerateToken("test@example.com", "password")
	assert.NoError(t, err)

	userId, err := service.ParseToken(token)

	assert.NoError(t, err)
	assert.Equal(t, 1, userId)
}

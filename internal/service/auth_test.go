package service

import (
	"testing"

	dokkee "github.com/airvt1x/dokkee-backend"
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

func (m *MockAuthorization) GetUser(username string) (dokkee.User, error) {
	args := m.Called(username)
	return args.Get(0).(dokkee.User), args.Error(1)
}

func (m *MockAuthorization) GetProfile(userID int) (dokkee.User, error) {
	args := m.Called(userID)
	return args.Get(0).(dokkee.User), args.Error(1)
}

func (m *MockAuthorization) UpdateProfile(userID int, input dokkee.UpdateProfileInput) error {
	args := m.Called(userID, input)
	return args.Error(0)
}

func TestAuthService_CreateUser(t *testing.T) {
	mockRepo := new(MockAuthorization)
	svc := NewAuthService(mockRepo)

	user := dokkee.User{
		Username:  "testuser",
		Password:  "password",
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Phone:     "+1234567890",
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("dokkee.User")).Return(1, nil)

	id, err := svc.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GenerateToken(t *testing.T) {
	mockRepo := new(MockAuthorization)
	svc := NewAuthService(mockRepo)

	username := "testuser"
	password := "password"

	// bcrypt hash for "password"
	hash, _ := generatePasswordHash(password)
	user := dokkee.User{
		Id:       1,
		Username: username,
		Password: hash,
	}

	mockRepo.On("GetUser", username).Return(user, nil)

	token, err := svc.GenerateToken(username, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_ParseToken(t *testing.T) {
	mockRepo := new(MockAuthorization)
	svc := NewAuthService(mockRepo)

	username := "testuser"
	password := "password"
	hash, _ := generatePasswordHash(password)
	user := dokkee.User{Id: 1, Username: username, Password: hash}

	mockRepo.On("GetUser", username).Return(user, nil)

	token, err := svc.GenerateToken(username, password)
	assert.NoError(t, err)

	userID, err := svc.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, 1, userID)
}

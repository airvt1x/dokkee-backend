package service

import (
	"testing"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuditRepository struct {
	mock.Mock
}

func (m *MockAuditRepository) Insert(record dokkee.AuditRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func TestAuditService_Log(t *testing.T) {
	mockRepo := new(MockAuditRepository)
	svc := NewAuditService(mockRepo)

	event := AuditEvent{
		Type:    "DOCUMENT_UPLOADED",
		UserID:  123,
		IP:      "192.168.1.1",
		Success: true,
	}

	mockRepo.On("Insert", mock.MatchedBy(func(record dokkee.AuditRecord) bool {
		return record.EventType == event.Type &&
			len(record.UserIDHash) > 0 &&
			len(record.IPHash) > 0 &&
			record.Success == event.Success
	})).Return(nil)

	err := svc.Log(event)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuditService_Log_WithDocID(t *testing.T) {
	mockRepo := new(MockAuditRepository)
	svc := NewAuditService(mockRepo)

	docID := 456
	event := AuditEvent{
		Type:    "RESULT_ACCESSED",
		UserID:  123,
		DocID:   &docID,
		IP:      "10.0.0.1",
		Success: true,
	}

	mockRepo.On("Insert", mock.MatchedBy(func(record dokkee.AuditRecord) bool {
		return record.EventType == event.Type &&
			record.DocIDHash != nil &&
			len(*record.DocIDHash) > 0
	})).Return(nil)

	err := svc.Log(event)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuditService_Log_Error(t *testing.T) {
	mockRepo := new(MockAuditRepository)
	svc := NewAuditService(mockRepo)

	event := AuditEvent{Type: "API_REQUEST", UserID: 1, IP: "127.0.0.1"}
	mockRepo.On("Insert", mock.Anything).Return(assert.AnError)

	err := svc.Log(event)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHashValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"simple text", "test", "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"},
		{"with spaces", "hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hashValue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHashID(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"zero", 0, "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9"},
		{"one", 1, "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"},
		{"large", 999999, "937377f056160fc4b15e0b770c67136a5f03c15205b4d3bf918268fefa2c6d0a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hashID(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

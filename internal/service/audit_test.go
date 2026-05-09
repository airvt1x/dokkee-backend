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

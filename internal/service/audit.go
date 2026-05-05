package service

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/airvt1x/dokkee-backend/internal/repository"
)

type AuditEvent struct {
	Type      string
	UserID    int
	DocID     *int
	IP        string
	Success   bool
	ErrorCode string
}

type AuditService struct {
	repo repository.Audit
}

func NewAuditService(repo repository.Audit) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Log(event AuditEvent) error {
	record := dokkee.AuditRecord{
		EventType:  event.Type,
		UserIDHash: hashID(event.UserID),
		IPHash:     hashValue(event.IP),
		Success:    event.Success,
		ErrorCode:  event.ErrorCode,
	}

	if event.DocID != nil {
		h := hashID(*event.DocID)
		record.DocIDHash = &h
	}

	return s.repo.Insert(record)
}

func hashValue(v string) string {
	h := sha256.New()
	h.Write([]byte(v))
	return hex.EncodeToString(h.Sum(nil))
}

func hashID(id int) string {
	return hashValue(strconv.Itoa(id))
}

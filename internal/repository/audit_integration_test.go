//go:build integration

package repository

import (
	"testing"

	dokkee "github.com/airvt1x/dokkee-backend"
	"github.com/stretchr/testify/assert"
)

func TestAuditPostgres_Insert_Integration(t *testing.T) {
	repo := &AuditPostgres{db: testDB}

	record := dokkee.AuditRecord{
		EventType:  "API_REQUEST",
		UserIDHash: "hash123",
		DocIDHash:  nil,
		IPHash:     "iphash",
		Success:    true,
		ErrorCode:  "",
	}
	err := repo.Insert(record)
	assert.NoError(t, err)

	// Cleanup: удаляем последнюю вставленную запись (по UserIDHash)
	t.Cleanup(func() {
		testDB.Exec("DELETE FROM audit_log WHERE user_id_hash='hash123'")
	})
}

func TestAuditPostgres_Insert_Multiple_Integration(t *testing.T) {
	repo := &AuditPostgres{db: testDB}

	for i := 0; i < 3; i++ {
		record := dokkee.AuditRecord{
			EventType:  "API_REQUEST",
			UserIDHash: "user",
			IPHash:     "ip",
			Success:    true,
		}
		err := repo.Insert(record)
		assert.NoError(t, err)
	}

	var count int
	err := testDB.Get(&count, "SELECT COUNT(*) FROM audit_log WHERE user_id_hash='user'")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, count, 3)

	t.Cleanup(func() {
		testDB.Exec("DELETE FROM audit_log WHERE user_id_hash='user'")
	})
}

func TestAuditPostgres_Insert_ErrorCode_Integration(t *testing.T) {
	repo := &AuditPostgres{db: testDB}

	record := dokkee.AuditRecord{
		EventType:  "DOCUMENT_ANALYSIS_FAILED",
		UserIDHash: "user_fail",
		IPHash:     "ip",
		Success:    false,
		ErrorCode:  "ERR_AI_TIMEOUT",
	}
	err := repo.Insert(record)
	assert.NoError(t, err)

	t.Cleanup(func() {
		testDB.Exec("DELETE FROM audit_log WHERE user_id_hash='user_fail'")
	})
}
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

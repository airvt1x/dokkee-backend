//go:build integration

package repository

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var testDB *sqlx.DB

func TestMain(m *testing.M) {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5433" // порт, который пробрасывает port-forwarder
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "dokkee"
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		password = "dokkee_secret"
	}
	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		dbname = "dokkee"
	}

	dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	var err error
	testDB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}
	defer testDB.Close()

	code := m.Run()
	os.Exit(code)
}
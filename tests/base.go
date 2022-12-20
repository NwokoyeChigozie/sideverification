package tests

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/vesicash/verification-ms/internal/config"
	"github.com/vesicash/verification-ms/internal/models/migrations"
	"github.com/vesicash/verification-ms/pkg/repository/storage/postgresql"
)

func Setup() {
	config := config.Setup("../../app")
	db := postgresql.ConnectToDatabases(config.TestDatabases)
	if config.TestDatabases.Migrate {
		migrations.RunAllMigrations(db)
	}
}

func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	return res
}

func AssertStatusCode(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("handler returned wrong status code: got status %d expected status %d", got, expected)
	}
}

func AssertResponseMessage(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("handler returned wrong message: got message: %q expected: %q", got, expected)
	}
}

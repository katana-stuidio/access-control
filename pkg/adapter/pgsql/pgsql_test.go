package pgsql

import (
	"database/sql"
	"os"
	"testing"

	"github.com/katana-stuidio/access-control/internal/config"
)

func TestNew_DefaultConfig(t *testing.T) {
	os.Setenv("SRV_DB_HOST", "localhost")
	os.Setenv("SRV_DB_USER", "postgres")
	os.Setenv("SRV_DB_PASS", "password")
	os.Setenv("SRV_DB_NAME", "testdb")

	conf := &config.Config{
		PGSQLConfig: &config.PGSQLConfig{},
	}
	db := New(conf)

	if db == nil {
		t.Fatal("Esperado um objeto dbpool, mas obteve nil")
	}
}

func TestGetDB(t *testing.T) {
	conf := &config.Config{
		PGSQLConfig: &config.PGSQLConfig{
			DB_DRIVE: "postgres",
			DB_DSN:   "host=localhost user=postgres password=password dbname=testdb sslmode=disable",
		},
	}
	db := New(conf)

	if db.GetDB() == nil {
		t.Error("Esperado um objeto *sql.DB, mas obteve nil")
	}
}

func TestCloseConnection(t *testing.T) {
	conf := &config.Config{
		PGSQLConfig: &config.PGSQLConfig{
			DB_DRIVE: "postgres",
			DB_DSN:   "host=localhost user=postgres password=password dbname=testdb sslmode=disable",
		},
	}
	db := New(conf)

	if err := db.CloseConnection(); err != nil && err != sql.ErrConnDone {
		t.Errorf("Erro ao fechar a conexão: %v", err)
	}
}

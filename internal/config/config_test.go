package config

import (
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestNewConfig_DefaultValues(t *testing.T) {
	conf := NewConfig()

	if conf.PORT != "8080" {
		t.Errorf("PORT esperado '8080', mas obteve '%s'", conf.PORT)
	}

	if conf.Mode != DEVELOPER {
		t.Errorf("Mode esperado '%s', mas obteve '%s'", DEVELOPER, conf.Mode)
	}

	if conf.JWTSecretKey != "LinuxRust162!" {
		t.Errorf("JWTSecretKey esperado 'LinuxRust162!', mas obteve '%s'", conf.JWTSecretKey)
	}

	if conf.JWTTokenExp != 480 {
		t.Errorf("JWTTokenExp esperado 480, mas obteve %d", conf.JWTTokenExp)
	}

	if conf.JWTRefreshExp != 720 {
		t.Errorf("JWTRefreshExp esperado 720, mas obteve %d", conf.JWTRefreshExp)
	}
}

func TestNewConfig_EnvVariables(t *testing.T) {
	os.Setenv("SRV_PORT", "9090")
	os.Setenv("SRV_MODE", PRODUCTION)
	os.Setenv("SRV_JWT_SECRET_KEY", "NewSecret")
	os.Setenv("SRV_JWT_TOKEN_EXP", "600")
	os.Setenv("SRV_JWT_REFRESH_EXP", "900")

	conf := NewConfig()

	if conf.PORT != "9090" {
		t.Errorf("PORT esperado '9090', mas obteve '%s'", conf.PORT)
	}

	if conf.Mode != PRODUCTION {
		t.Errorf("Mode esperado '%s', mas obteve '%s'", PRODUCTION, conf.Mode)
	}

	if conf.JWTSecretKey != "NewSecret" {
		t.Errorf("JWTSecretKey esperado 'NewSecret', mas obteve '%s'", conf.JWTSecretKey)
	}

	expToken, _ := strconv.Atoi("600")
	if conf.JWTTokenExp != expToken {
		t.Errorf("JWTTokenExp esperado %d, mas obteve %d", expToken, conf.JWTTokenExp)
	}

	expRefresh, _ := strconv.Atoi("900")
	if conf.JWTRefreshExp != expRefresh {
		t.Errorf("JWTRefreshExp esperado %d, mas obteve %d", expRefresh, conf.JWTRefreshExp)
	}
}

func TestNewConfig_PGSQLDefaults(t *testing.T) {
	conf := NewConfig()
	expectedDB := &PGSQLConfig{
		DB_DRIVE: "postgres",
		DB_HOST:  "localhost",
		DB_PORT:  "5432",
		DB_USER:  "postgres",
		DB_PASS:  "supersenha",
		DB_NAME:  "drive_db_dev",
	}

	if !reflect.DeepEqual(conf.PGSQLConfig, expectedDB) {
		t.Errorf("Configuração do banco de dados não corresponde ao esperado")
	}
}

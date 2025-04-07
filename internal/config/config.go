package config

import (
	"os"
	"strconv"
)

const (
	DEVELOPER    = "developer"
	HOMOLOGATION = "homologation"
	PRODUCTION   = "production"
)

type Config struct {
	PORT          string `json:"port"`
	Mode          string `json:"mode"`
	JWTSecretKey  string `json:"jwt_secret_key"`
	JWTTokenExp   int    `json:"jwt_token_exp"`
	JWTRefreshExp int    `json:"jwt_refresh_exp"`
	*PGSQLConfig
}

type PGSQLConfig struct {
	DB_DRIVE                  string `json:"db_drive"`
	DB_HOST                   string `json:"db_host"`
	DB_PORT                   string `json:"db_port"`
	DB_USER                   string `json:"db_user"`
	DB_PASS                   string `json:"db_pass"`
	DB_NAME                   string `json:"db_name"`
	DB_DSN                    string `json:"-"`
	DB_SET_MAX_OPEN_CONNS     int    `json:"db_set_max_open_conns"`
	DB_SET_MAX_IDLE_CONNS     int    `json:"db_set_max_idle_conns"`
	DB_SET_CONN_MAX_LIFE_TIME int    `json:"db_set_conn_max_life_time"`
	SRV_DB_SSL_MODE           bool   `json:"srv_db_ssl_mode"`
}

func NewConfig() *Config {
	conf := defaultConf()

	SRV_PORT := os.Getenv("SRV_PORT")
	if SRV_PORT != "" {
		conf.PORT = SRV_PORT
	}

	SRV_MODE := os.Getenv("SRV_MODE")
	if SRV_MODE != "" {
		conf.Mode = SRV_MODE
	}

	SRV_JWT_SECRET_KEY := os.Getenv("SRV_JWT_SECRET_KEY")
	if SRV_JWT_SECRET_KEY != "" {
		conf.JWTSecretKey = SRV_JWT_SECRET_KEY
	}

	SRV_JWT_TOKEN_EXP := os.Getenv("SRV_JWT_TOKEN_EXP")
	if SRV_JWT_SECRET_KEY != "" {
		conf.JWTTokenExp, _ = strconv.Atoi(SRV_JWT_TOKEN_EXP)
	}

	SRV_JWT_REFRESH_EXP := os.Getenv("SRV_JWT_REFRESH_EXP")
	if SRV_JWT_SECRET_KEY != "" {
		conf.JWTRefreshExp, _ = strconv.Atoi(SRV_JWT_REFRESH_EXP)
	}

	return conf
}

func defaultConf() *Config {
	default_conf := Config{
		PORT:          "8080",
		Mode:          DEVELOPER,
		JWTSecretKey:  "LinuxRust162!", // "----your-256-bit-secret-here----" length 32
		JWTTokenExp:   480,             // 8 horas
		JWTRefreshExp: 720,             // 12 horas                               // 240 = 4h

		PGSQLConfig: &PGSQLConfig{
			DB_DRIVE: "postgres",
			DB_PORT:  "5432",
		},
	}

	return &default_conf
}

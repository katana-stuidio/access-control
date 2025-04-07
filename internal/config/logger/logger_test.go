package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestGetOutputLogs_Default(t *testing.T) {
	os.Unsetenv(LOG_OUTPUT)
	if output := getOutputLogs(); output != "stdout" {
		t.Errorf("Esperado 'stdout', mas obteve '%s'", output)
	}
}

func TestGetOutputLogs_FromEnv(t *testing.T) {
	os.Setenv(LOG_OUTPUT, "file.log")
	if output := getOutputLogs(); output != "file.log" {
		t.Errorf("Esperado 'file.log', mas obteve '%s'", output)
	}
}

func TestGetLevelLogs_Default(t *testing.T) {
	os.Unsetenv(LOG_LEVEL)
	if level := getLevelLogs(); level != zapcore.InfoLevel {
		t.Errorf("Esperado zapcore.InfoLevel, mas obteve %v", level)
	}
}

func TestGetLevelLogs_FromEnv(t *testing.T) {
	cases := map[string]zapcore.Level{
		"info":  zapcore.InfoLevel,
		"error": zapcore.ErrorLevel,
		"debug": zapcore.DebugLevel,
		"other": zapcore.InfoLevel,
	}
	for env, expected := range cases {
		os.Setenv(LOG_LEVEL, env)
		if level := getLevelLogs(); level != expected {
			t.Errorf("Para '%s', esperado %v, mas obteve %v", env, expected, level)
		}
	}
}

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logConfig := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{MessageKey: "message"}),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)
	log := zap.New(logConfig)
	log.Info("Teste de log")

	if !strings.Contains(buf.String(), "Teste de log") {
		t.Errorf("Log esperado contendo 'Teste de log', mas obteve '%s'", buf.String())
	}
}

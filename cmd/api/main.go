package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katana-stuidio/access-control/internal/config"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	hand_usr "github.com/katana-stuidio/access-control/internal/handler/user"
	"github.com/katana-stuidio/access-control/pkg/adapter/pgsql"
	"github.com/katana-stuidio/access-control/pkg/server"
	service_usr "github.com/katana-stuidio/access-control/pkg/service/user"
)

var (
	VERSION = "0.1.0-dev"
	COMMIT  = "ABCDEFG-dev"
)

func main() {
	logger.Info("start Drive Auth application")

	conf := config.NewConfig()
	conn_pg := pgsql.New(conf)

	usr_service := service_usr.NewUserService(conn_pg)

	// Criação do router com Gin
	router := gin.Default()

	// Healthcheck básico
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"MSG":    "Server Ok",
			"codigo": 200,
		})
	})

	// Registra handlers do módulo user
	hand_usr.RegisterUserAPIHandlers(router, usr_service, conf)

	// Cria servidor HTTP
	srv := server.NewHTTPServer(router, conf)

	// Inicia o servidor em goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Printf("Server Run on [Port: %s], [Mode: %s], [Version: %s], [Commit: %s]", conf.PORT, conf.Mode, VERSION, COMMIT)

	// Impede o encerramento da aplicação
	select {}
}

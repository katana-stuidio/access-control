package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katana-stuidio/access-control/internal/config"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	hand_ten "github.com/katana-stuidio/access-control/internal/handler/tenant"
	hand_usr "github.com/katana-stuidio/access-control/internal/handler/user"
	"github.com/katana-stuidio/access-control/pkg/adapter/pgsql"
	"github.com/katana-stuidio/access-control/pkg/server"
	service_ten "github.com/katana-stuidio/access-control/pkg/service/tenant"
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
	tenat_service := service_ten.NewTenantService(conn_pg)

	// Criação do router com Gin
	router := gin.Default()

	// Configure CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Healthcheck básico
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"MSG":    "Server Ok",
			"codigo": 200,
		})
	})

	// Registra handlers do módulo user
	hand_usr.RegisterUserAPIHandlers(router, usr_service, conf)
	hand_ten.RegisterTenantAPIHandlers(router, tenat_service)

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

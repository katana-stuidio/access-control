package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
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

	// Carrega configurações e inicializa conexão com PostgreSQL
	conf := config.NewConfig()
	conn_pg := pgsql.New(conf)

	// Instancia serviços de usuário e tenant
	usr_service := service_usr.NewUserService(conn_pg)
	tenat_service := service_ten.NewTenantService(conn_pg)

	// Criação do router com Gin
	router := gin.Default()

	// Middleware CORS configurado corretamente
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token", "X-Requested-With", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "Cookie", "Host", "Pragma", "Referer", "User-Agent"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,        // Deve ser false quando AllowOrigins é "*"
		MaxAge:           12 * 60 * 60, // 12 horas
	}))

	// Rota básica de healthcheck
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"MSG":    "Server Ok",
			"codigo": 200,
		})
	})

	// Registra handlers dos módulos user e tenant
	hand_usr.RegisterUserAPIHandlers(router, usr_service, conf)
	hand_ten.RegisterTenantAPIHandlers(router, tenat_service)

	// Cria e inicia servidor HTTP
	srv := server.NewHTTPServer(router, conf)

	// Roda o servidor em uma goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Log de inicialização
	log.Printf("Server Run on [Port: %s], [Mode: %s], [Version: %s], [Commit: %s]",
		conf.PORT, conf.Mode, VERSION, COMMIT)

	// Aguarda indefinidamente
	select {}
}

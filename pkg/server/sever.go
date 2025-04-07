package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/katana-stuidio/access-control/internal/config"
)

func NewHTTPServer(router *gin.Engine, conf *config.Config) *http.Server {
	srv := &http.Server{
		Addr:         ":" + conf.PORT,
		Handler:      router,
		ErrorLog:     log.New(os.Stderr, "logger: ", log.Lshortfile),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv
}

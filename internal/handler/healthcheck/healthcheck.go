package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/pkg/service/healthcheck"
)

func HealthcheckHandler(service healthcheck.HealthcheckServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := service.CheckDB()

		if err != nil || !ok {
			logger.Error("Erro ao verificar o banco de dados", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"MSG":    "Erro ao verificar o banco de dados",
				"codigo": 500,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"MSG":    "IAMGEM MAIS",
			"codigo": 200,
		})
	}
}

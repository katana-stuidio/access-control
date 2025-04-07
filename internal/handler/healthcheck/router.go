package healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/katana-stuidio/access-control/pkg/service/healthcheck"
)

func RegisterHealthcheckAPIHandlers(router *gin.Engine, service healthcheck.HealthcheckServiceInterface) {
	v1 := router.Group("/int/v1")
	{
		v1.GET("/healthcheck", HealthcheckHandler(service))
	}
}

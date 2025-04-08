package tenant

import (
	"github.com/gin-gonic/gin"
	"github.com/katana-stuidio/access-control/pkg/service/tenant"
)

func RegisterTenantAPIHandlers(r *gin.Engine, service tenant.TenantServiceInterface) {
	TenantGroup := r.Group("/api/v1/Tenant")
	{
		TenantGroup.POST("/", createTenant(service))
		TenantGroup.GET("/:id", getTenant(service))

		TenantGroup.PATCH("/:id", updateTenant(service))
		TenantGroup.DELETE("/:id", deleteTenant(service))
		TenantGroup.GET("/", getAllTenant(service))
	}
}

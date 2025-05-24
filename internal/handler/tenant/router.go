package tenant

import (
	"github.com/gin-gonic/gin"
	"github.com/katana-stuidio/access-control/pkg/service/tenant"
)

func RegisterTenantAPIHandlers(r *gin.Engine, service tenant.TenantServiceInterface) {
	tenantGroup := r.Group("/api/v1/Tenant")
	{
		tenantGroup.POST("/", gin.WrapH(createTenant(service)))
		tenantGroup.GET("/:id", gin.WrapH(getTenant(service)))
		tenantGroup.PATCH("/:id", gin.WrapH(updateTenant(service)))
		tenantGroup.DELETE("/:id", gin.WrapH(deleteTenant(service)))
		tenantGroup.GET("/", gin.WrapH(getAllTenant(service)))
	}
}

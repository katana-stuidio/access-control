package tenant_group

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures the tenant group routes
func SetupRoutes(router *gin.Engine, handler *TenantGroupHandler) {
	// Tenant Group routes
	tenantGroupRoutes := router.Group("/api/v1/tenant-groups")
	{
		tenantGroupRoutes.POST("/", handler.Create)
		tenantGroupRoutes.GET("/", handler.GetAll)
		tenantGroupRoutes.GET("/:id", handler.GetByID)
		tenantGroupRoutes.PUT("/:id", handler.Update)
		tenantGroupRoutes.DELETE("/:id", handler.Delete)
	}
}

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/katana-stuidio/access-control/internal/config"
	service_ten "github.com/katana-stuidio/access-control/pkg/service/tenant"
	service_ten_group "github.com/katana-stuidio/access-control/pkg/service/tenant_group"
	"github.com/katana-stuidio/access-control/pkg/service/token"
	"github.com/katana-stuidio/access-control/pkg/service/user"
)

func RegisterUserAPIHandlers(r *gin.Engine, service user.UserServiceInterface, tenantService service_ten.TenantServiceInterface, tenantGroupService service_ten_group.TenantGroupServiceInterface, conf *config.Config, tokenService token.TokenServiceInterface) {
	userGroup := r.Group("/api/v1/user")
	{
		userGroup.POST("/", createUser(service))
		userGroup.GET("/:id", getUser(service))
		userGroup.POST("/getjwt", getJWT(service, tenantService, tenantGroupService, conf, tokenService))
		userGroup.POST("/refreshjwt", refreshToken(conf, tokenService))
		userGroup.POST("/validatejwt", validateToken(conf))
		userGroup.POST("/logout", logout(tokenService))
		userGroup.PATCH("/:id", updateUser(service))
		userGroup.DELETE("/:id", deleteUser(service))
		userGroup.GET("/", getAllUser(service))
		userGroup.PATCH("/changepassword", gin.WrapH(changePassword(service)))
	}
}

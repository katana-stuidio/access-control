package tenant

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/internal/dto"
	"github.com/katana-stuidio/access-control/pkg/model"
	"github.com/katana-stuidio/access-control/pkg/service/tenant"
)

func getAllTenant(service tenant.TenantServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		Tenants := service.GetAll(c.Request.Context())
		c.JSON(http.StatusOK, Tenants)
	}
}

func getTenant(service tenant.TenantServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		externalID := c.Param("id")
		id, err := uuid.Parse(externalID)
		if err != nil || id == uuid.Nil {
			ErroHttpMsgTenantIdIsRequired.Write(c.Writer)
			return
		}

		Tenant := service.GetByID(c.Request.Context(), id)
		if Tenant.ID == uuid.Nil {
			ErroHttpMsgTenantNotFound.Write(c.Writer)
			return
		}

		c.JSON(http.StatusOK, Tenant)
	}
}

func createTenant(service tenant.TenantServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var TenantDto dto.TenantRequestDtoInput

		if err := c.ShouldBindJSON(&TenantDto); err != nil {
			logger.Error("Invalid Tenant request: ", err)
			ErroHttpMsgToParseRequestTenantToJson.Write(c.Writer)
			return
		}

		TenantModel := model.Tenant{

			Name:       TenantDto.Name,
			SchemaName: TenantDto.SchemaName,
		}

		tn, err := model.NewTenant(&TenantModel)
		if err != nil {
			logger.Error("Invalid login request: ", err)
			ErroHttpMsgToParseRequestTenantToJson.Write(c.Writer)
			return
		}

		if strings.TrimSpace(tn.Name) == "" {
			ErroHttpMsgTenantNameIsRequired.Write(c.Writer)
			return
		}

		if strings.TrimSpace(tn.SchemaName) == "" {
			ErroHttpMsgTenantSchemaNameIsRequired.Write(c.Writer)
			return
		}

		TenantExist, err := service.GetExistTenantName(c.Request.Context(), tn.Name)
		if err != nil {
			ErroHttpMsgToInsertTenant.Write(c.Writer)
			return
		}

		if TenantExist {
			ErroHttpMsgTenantAlreadyExist.Write(c.Writer)
			return
		}

		result, err := service.Create(c.Request.Context(), tn)
		if err != nil {
			ErroHttpMsgToInsertTenant.Write(c.Writer)
			return
		}

		resultOut := dto.TenantRequestDtoOutPut{
			ID:         result.ID,
			Name:       result.Name,
			SchemaName: result.SchemaName,
			IsActive:   result.IsActive,
			CreatedAt:  result.CreatedAt,
			UpdatedAt:  result.UpdatedAt,
		}

		c.JSON(http.StatusCreated, resultOut)
	}
}

func updateTenant(service tenant.TenantServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		externalID := c.Param("id")
		id, err := uuid.Parse(externalID)
		if err != nil || id == uuid.Nil {
			ErroHttpMsgTenantIdIsRequired.Write(c.Writer)
			return
		}

		var requestToUpdate model.Tenant
		if err := c.ShouldBindJSON(&requestToUpdate); err != nil {
			ErroHttpMsgToParseRequestTenantToJson.Write(c.Writer)
			return
		}

		if strings.TrimSpace(requestToUpdate.Name) == "" {
			ErroHttpMsgTenantNameIsRequired.Write(c.Writer)
			return
		}

		if strings.TrimSpace(requestToUpdate.Name) == "" {
			ErroHttpMsgTenantNameIsRequired.Write(c.Writer)
			return
		}

		Tenant := service.GetByID(c.Request.Context(), id)
		if Tenant.ID == uuid.Nil {
			ErroHttpMsgTenantNotFound.Write(c.Writer)
			return
		}

		rowsAffected := service.Update(c.Request.Context(), id, &requestToUpdate)
		if rowsAffected == 0 {
			ErroHttpMsgToUpdateTenant.Write(c.Writer)
			return
		}

		c.JSON(http.StatusOK, requestToUpdate)
	}
}

func deleteTenant(service tenant.TenantServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		externalID := c.Param("id")
		id, err := uuid.Parse(externalID)
		if err != nil || id == uuid.Nil {
			ErroHttpMsgTenantIdIsRequired.Write(c.Writer)
			return
		}

		Tenant := service.GetByID(c.Request.Context(), id)
		if Tenant.ID == uuid.Nil {
			ErroHttpMsgTenantNotFound.Write(c.Writer)
			return
		}

		rowsAffected := service.Delete(c.Request.Context(), id)
		if rowsAffected == 0 {
			ErroHttpMsgToDeleteTenant.Write(c.Writer)
			return
		}

		SuccessHttpMsgToDeleteTenant.Write(c.Writer)
	}
}

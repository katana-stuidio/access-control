package tenant_group

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/internal/dto"
	"github.com/katana-stuidio/access-control/pkg/model"
	"github.com/katana-stuidio/access-control/pkg/service/tenant_group"
)

type TenantGroupHandler struct {
	tenantGroupService tenant_group.TenantGroupServiceInterface
}

func NewTenantGroupHandler(tenantGroupService tenant_group.TenantGroupServiceInterface) *TenantGroupHandler {
	return &TenantGroupHandler{
		tenantGroupService: tenantGroupService,
	}
}

// Create handles the creation of a new tenant group
func (h *TenantGroupHandler) Create(c *gin.Context) {
	var request dto.TenantGroupCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Error binding JSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate CNPJ
	tenantGroup := &model.TenantGroup{
		Name: request.Name,
		CNPJ: request.CNPJ,
	}

	if !tenantGroup.CheckCNPJ(request.CNPJ) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CNPJ format"})
		return
	}

	// Check if CNPJ already exists
	exists, err := h.tenantGroupService.GetExistCNPJ(c.Request.Context(), request.CNPJ)
	if err != nil {
		logger.Error("Error checking existing CNPJ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "CNPJ already exists"})
		return
	}

	// Create tenant group
	newTenantGroup, err := model.NewTenantGroup(tenantGroup)
	if err != nil {
		logger.Error("Error creating tenant group", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	createdTenantGroup, err := h.tenantGroupService.Create(c.Request.Context(), newTenantGroup)
	if err != nil {
		logger.Error("Error saving tenant group", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := dto.TenantGroupResponse{
		ID:        createdTenantGroup.ID,
		Name:      createdTenantGroup.Name,
		CNPJ:      createdTenantGroup.CNPJ,
		IsActive:  createdTenantGroup.IsActive,
		CreatedAt: createdTenantGroup.CreatedAt,
		UpdatedAt: createdTenantGroup.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetAll handles the retrieval of all tenant groups with pagination
func (h *TenantGroupHandler) GetAll(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	paginate, err := h.tenantGroupService.GetAll(c.Request.Context(), limit, page)
	if err != nil {
		logger.Error("Error getting tenant groups", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convert to response format
	var tenantGroups []dto.TenantGroupResponse
	if tenantGroupList, ok := paginate.Data.(*model.TenantGroupList); ok {
		for _, tg := range tenantGroupList.List {
			tenantGroups = append(tenantGroups, dto.TenantGroupResponse{
				ID:        tg.ID,
				Name:      tg.Name,
				CNPJ:      tg.CNPJ,
				IsActive:  tg.IsActive,
				CreatedAt: tg.CreatedAt,
				UpdatedAt: tg.UpdatedAt,
			})
		}
	}

	response := dto.TenantGroupListResponse{
		Total: paginate.Total,
		Page:  paginate.Currente,
		Last:  paginate.Last,
		Data:  tenantGroups,
	}

	c.JSON(http.StatusOK, response)
}

// GetByID handles the retrieval of a tenant group by ID
func (h *TenantGroupHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	tenantGroup := h.tenantGroupService.GetByID(c.Request.Context(), id)
	if tenantGroup.ID == uuid.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant group not found"})
		return
	}

	response := dto.TenantGroupResponse{
		ID:        tenantGroup.ID,
		Name:      tenantGroup.Name,
		CNPJ:      tenantGroup.CNPJ,
		IsActive:  tenantGroup.IsActive,
		CreatedAt: tenantGroup.CreatedAt,
		UpdatedAt: tenantGroup.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Update handles the update of a tenant group
func (h *TenantGroupHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var request dto.TenantGroupUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Error binding JSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate CNPJ
	tenantGroup := &model.TenantGroup{
		Name:     request.Name,
		CNPJ:     request.CNPJ,
		IsActive: request.IsActive,
	}

	if !tenantGroup.CheckCNPJ(request.CNPJ) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CNPJ format"})
		return
	}

	// Check if CNPJ already exists for other tenant groups
	existingTenantGroup, err := h.tenantGroupService.GetByCNPJ(c.Request.Context(), request.CNPJ)
	if err == nil && existingTenantGroup.ID != id {
		c.JSON(http.StatusConflict, gin.H{"error": "CNPJ already exists"})
		return
	}

	rowsAffected := h.tenantGroupService.Update(c.Request.Context(), id, tenantGroup)
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant group not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant group updated successfully"})
}

// Delete handles the deletion of a tenant group
func (h *TenantGroupHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	rowsAffected := h.tenantGroupService.Delete(c.Request.Context(), id)
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant group not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant group deleted successfully"})
}

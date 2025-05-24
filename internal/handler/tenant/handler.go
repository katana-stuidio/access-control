package tenant

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/internal/dto"
	"github.com/katana-stuidio/access-control/pkg/model"
	"github.com/katana-stuidio/access-control/pkg/service/tenant"
)

// @Summary Get all tenants
// @Description Get a paginated list of all tenants
// @Tags tenants
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page (default: 10)"
// @Param page query int false "Page number (default: 1)"
// @Success 200 {object} model.Paginate
// @Failure 500 {object} handler.HttpMsg
// @Router /api/v1/tenant/ [get]
func getAllTenant(service tenant.TenantServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get pagination parameters from query string
		limitStr := r.URL.Query().Get("limit")
		pageStr := r.URL.Query().Get("page")

		// Convert to int64 with defaults
		limit := int64(10) // default limit
		page := int64(1)   // default page

		if limitStr != "" {
			if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil && l > 0 {
				limit = l
			}
		}
		if pageStr != "" {
			if p, err := strconv.ParseInt(pageStr, 10, 64); err == nil && p > 0 {
				page = p
			}
		}

		result, err := service.GetAll(r.Context(), limit, page)
		if err != nil {
			ErroHttpMsgToConvertingResponseTenantListToJson.Write(w)
			return
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			ErroHttpMsgToConvertingResponseTenantListToJson.Write(w)
			return
		}
	})
}

// @Summary Get tenant by ID
// @Description Get a tenant by their ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} dto.TenantRequestDtoOutPut
// @Failure 400 {object} handler.HttpMsg
// @Failure 404 {object} handler.HttpMsg
// @Router /api/v1/tenant/{id} [get]
func getTenant(service tenant.TenantServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		externalID := r.URL.Query().Get("id")
		id, err := uuid.Parse(externalID)
		if err != nil {
			ErroHttpMsgTenantIdIsRequired.Write(w)
			return
		}

		if id == uuid.Nil {
			ErroHttpMsgTenantIdIsRequired.Write(w)
			return
		}

		tenant := service.GetByID(r.Context(), id)
		if tenant.ID == uuid.Nil {
			ErroHttpMsgTenantNotFound.Write(w)
			return
		}

		err = json.NewEncoder(w).Encode(tenant)
		if err != nil {
			ErroHttpMsgToParseResponseTenantToJson.Write(w)
			return
		}
	}
}

// @Summary Create a new tenant
// @Description Create a new tenant with the provided details
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenant body dto.TenantRequestDtoInput true "Tenant details"
// @Success 201 {object} dto.TenantRequestDtoOutPut
// @Failure 400 {object} handler.HttpMsg
// @Failure 500 {object} handler.HttpMsg
// @Router /api/v1/tenant/ [post]
func createTenant(service tenant.TenantServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantDto := dto.TenantRequestDtoInput{}

		err := json.NewDecoder(r.Body).Decode(&tenantDto)
		if err != nil {
			logger.Error("Invalid tenant request: ", err)
			ErroHttpMsgToParseRequestTenantToJson.Write(w)
			return
		}

		var tenantModel model.Tenant
		tenantModel.Name = tenantDto.Name
		tenantModel.CNPJ = tenantDto.CNPJ

		tenantCad, err := model.NewTenant(&tenantModel)
		if err != nil {
			logger.Error("Invalid tenant request: ", err)
			ErroHttpMsgToParseRequestTenantToJson.Write(w)
			return
		}

		if tenantCad.Name == " " || tenantCad.Name == "" {
			ErroHttpMsgTenantNameIsRequired.Write(w)
			return
		}

		if tenantCad.CNPJ == " " || tenantCad.CNPJ == "" {
			ErroHttpMsgTenantCNPJIsRequired.Write(w)
			return
		}

		if !tenantCad.CheckCNPJ(tenantCad.CNPJ) {
			ErroHttpMsgTenantCNPJIsInvalid.Write(w)
			return
		}

		tenantExist, err := service.GetExistCNPJ(r.Context(), tenantCad.CNPJ)
		if err != nil {
			ErroHttpMsgToInsertTenant.Write(w)
			return
		}

		if tenantExist {
			ErroHttpMsgTenantAlreadyExist.Write(w)
			return
		}
		tenantCad.ID = uuid.New()
		tenantCad.SchemaName = tenantCad.ID.String()
		result, err := service.Create(r.Context(), tenantCad)
		if err != nil {
			ErroHttpMsgToInsertTenant.Write(w)
			return
		}

		var resultOut dto.TenantRequestDtoOutPut
		resultOut.ID = result.ID
		resultOut.Name = result.Name
		resultOut.CNPJ = result.CNPJ
		resultOut.SchemaName = result.SchemaName
		resultOut.IsActive = result.IsActive
		resultOut.CreatedAt = result.CreatedAt
		resultOut.UpdatedAt = result.UpdatedAt

		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(resultOut)
		if err != nil {
			ErroHttpMsgToParseResponseTenantToJson.Write(w)
			return
		}
	}
}

// @Summary Update tenant
// @Description Update an existing tenant's details
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Param tenant body model.Tenant true "Tenant details"
// @Success 200 {object} model.Tenant
// @Failure 400 {object} handler.HttpMsg
// @Failure 404 {object} handler.HttpMsg
// @Failure 500 {object} handler.HttpMsg
// @Router /api/v1/tenant/{id} [patch]
func updateTenant(service tenant.TenantServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		externalID := r.URL.Query().Get("id")
		id, err := uuid.Parse(externalID)
		if err != nil {
			ErroHttpMsgTenantIdIsRequired.Write(w)
			return
		}

		if id == uuid.Nil {
			ErroHttpMsgTenantIdIsRequired.Write(w)
			return
		}

		request_to_update_tenant := model.Tenant{}
		err = json.NewDecoder(r.Body).Decode(&request_to_update_tenant)
		if err != nil {
			ErroHttpMsgToParseRequestTenantToJson.Write(w)
			return
		}

		if request_to_update_tenant.Name == " " || request_to_update_tenant.Name == "" {
			ErroHttpMsgTenantNameIsRequired.Write(w)
			return
		}

		if request_to_update_tenant.CNPJ == " " || request_to_update_tenant.CNPJ == "" {
			ErroHttpMsgTenantCNPJIsRequired.Write(w)
			return
		}

		if !request_to_update_tenant.CheckCNPJ(request_to_update_tenant.CNPJ) {
			ErroHttpMsgTenantCNPJIsInvalid.Write(w)
			return
		}

		rowsAffected := service.Update(r.Context(), id, &request_to_update_tenant)
		if rowsAffected == 0 {
			ErroHttpMsgToUpdateTenant.Write(w)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// @Summary Delete tenant
// @Description Delete a tenant by their ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} handler.HttpMsg
// @Failure 400 {object} handler.HttpMsg
// @Failure 404 {object} handler.HttpMsg
// @Failure 500 {object} handler.HttpMsg
// @Router /api/v1/tenant/{id} [delete]
func deleteTenant(service tenant.TenantServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		externalID := r.URL.Query().Get("id")
		id, err := uuid.Parse(externalID)
		if err != nil {
			ErroHttpMsgTenantIdIsRequired.Write(w)
			return
		}

		if id == uuid.Nil {
			ErroHttpMsgTenantIdIsRequired.Write(w)
			return
		}

		rowsAffected := service.Delete(r.Context(), id)
		if rowsAffected == 0 {
			ErroHttpMsgToDeleteTenant.Write(w)
			return
		}

		SuccessHttpMsgToDeleteTenant.Write(w)
	}
}

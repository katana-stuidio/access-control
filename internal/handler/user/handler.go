package user

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/internal/dto"
	"github.com/katana-stuidio/access-control/pkg/jwt"
	"github.com/katana-stuidio/access-control/pkg/model"
	service_ten "github.com/katana-stuidio/access-control/pkg/service/tenant"
	service_ten_group "github.com/katana-stuidio/access-control/pkg/service/tenant_group"
	"github.com/katana-stuidio/access-control/pkg/service/token"
	"github.com/katana-stuidio/access-control/pkg/service/user"
	"github.com/potatowski/brazilcode"
)

type HttpMsg struct {
	Message string
	Code    int
}

func (h HttpMsg) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h.Code)
	json.NewEncoder(w).Encode(h)
}

var SuccessHttpMsgToChangePassword = HttpMsg{
	Message: "Password changed successfully",
	Code:    http.StatusOK,
}

var ErroHttpMsgUserEmailAlreadyExists = HttpMsg{
	Message: "Email already exists",
	Code:    http.StatusBadRequest,
}

var ErroHttpMsgUserEmailIsRequired = HttpMsg{
	Message: "Email is required",
	Code:    http.StatusBadRequest,
}

var ErroHttpMsgUserRoleIsRequired = HttpMsg{
	Message: "Role is required",
	Code:    http.StatusBadRequest,
}

// @Summary Get all users
// @Description Get a paginated list of all users
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page (default: 10)"
// @Param page query int false "Page number (default: 1)"
// @Success 200 {object} model.Paginate
// @Failure 500 {object} handler.HttpMsg
// @Router /api/v1/user/ [get]
func getAllUser(service user.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := int64(10) // default limit
		page := int64(1)   // default page

		users, err := service.GetAll(c.Request.Context(), limit, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} handler.HttpMsg
// @Failure 404 {object} handler.HttpMsg
// @Router /api/v1/user/{id} [get]
func getUser(service user.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		externalID := c.Param("id")
		id, err := uuid.Parse(externalID)
		if err != nil || id == uuid.Nil {
			ErroHttpMsgUserIdIsRequired.Write(c.Writer)
			return
		}

		user := service.GetByID(c.Request.Context(), id)
		if user.ID == uuid.Nil {
			ErroHttpMsgUserNotFound.Write(c.Writer)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// @Summary Create a new user
// @Description Create a new user with the provided details. Role must be one of: Professor, Estudante, Instituicao, Admin
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserRequestDtoInput true "User details"
// @Success 201 {object} dto.UserRequestDtoOutPut
// @Failure 400 {object} handler.HttpMsg
// @Failure 500 {object} handler.HttpMsg
// @Router /api/v1/user/ [post]
func createUser(service user.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userDto dto.UserRequestDtoInput

		if err := c.ShouldBindJSON(&userDto); err != nil {
			logger.Error("Invalid login request: ", err)
			ErroHttpMsgToParseRequestUserToJson.Write(c.Writer)
			return
		}

		userModel := model.User{
			Name:     userDto.Name,
			Username: userDto.Username,
			Password: userDto.Password,
			CNPJ:     userDto.CNPJ,
			Email:    userDto.Email,
			Role:     userDto.Role,
		}

		logger.Info("Received role from request: " + userDto.Role)

		if userDto.CNPJ == "" {
			ErroHttpMsgToInsertUser.Write(c.Writer)
			return
		}
		if err := brazilcode.CNPJIsValid(userDto.CNPJ); err != nil {
			ErroHttpMsgUserCnpjIsInvalid.Write(c.Writer)
			return
		}

		if userDto.Email == "" {
			c.JSON(ErroHttpMsgUserEmailIsRequired.Code, ErroHttpMsgUserEmailIsRequired)
			return
		}

		if userDto.Role == "" {
			ErroHttpMsgUserRoleIsRequired.Write(c.Writer)
			return
		}

		validRoles := map[string]bool{
			"Professor":   true,
			"Estudante":   true,
			"Instituicao": true,
			"Admin":       true,
		}

		if !validRoles[userDto.Role] {
			ErroHttpMsgInvalidRole.Write(c.Writer)
			return
		}

		logger.Info("Email: " + userDto.Email)
		emailExist, err := service.EmailExists(c.Request.Context(), userDto.Email)
		if err != nil {
			ErroHttpMsgToInsertUser.Write(c.Writer)
			return
		}

		if emailExist {
			c.JSON(ErroHttpMsgUserEmailAlreadyExists.Code, ErroHttpMsgUserEmailAlreadyExists)
			return
		}

		tenant_id, err := service.GetByCNPJ(c.Request.Context(), userDto.CNPJ)
		if err != nil {
			ErroHttpMsgCNPJNotFound.Write(c.Writer)
			return
		}

		if tenant_id == "" {
			ErroHttpMsgToInsertUser.Write(c.Writer)
			return
		}

		tenantUUID, err := uuid.Parse(tenant_id)
		if err != nil {
			ErroHttpMsgToInsertUser.Write(c.Writer)
			return
		}

		usrCad, err := model.NewUser(&userModel)
		if err != nil {
			logger.Error("Invalid login request: ", err)
			ErroHttpMsgToParseRequestUserToJson.Write(c.Writer)
			return
		}

		if strings.TrimSpace(usrCad.Username) == "" {
			ErroHttpMsgUserNameIsRequired.Write(c.Writer)
			return
		}

		if strings.TrimSpace(usrCad.Name) == "" {
			ErroHttpMsgUserNameIsRequired.Write(c.Writer)
			return
		}

		if strings.TrimSpace(usrCad.Password) == "" {
			ErroHttpMsgUserPasswordIsRequired.Write(c.Writer)
			return
		}

		usrCad.TenantID = tenantUUID

		userExist, err := service.GetExistUserName(c.Request.Context(), usrCad.Username)
		if err != nil {
			ErroHttpMsgToInsertUser.Write(c.Writer)
			return
		}

		if userExist {
			ErroHttpMsgUserAlreadyExist.Write(c.Writer)
			return
		}

		result, err := service.Create(c.Request.Context(), usrCad)
		if err != nil {
			ErroHttpMsgToInsertUser.Write(c.Writer)
			return
		}

		resultOut := dto.UserRequestDtoOutPut{
			ID:        result.ID,
			Username:  result.Username,
			Name:      result.Name,
			Enable:    result.Enable,
			Role:      result.Role,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		}

		c.JSON(http.StatusCreated, resultOut)
	}
}

// @Summary Update user
// @Description Update an existing user's details
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body model.User true "User details"
// @Success 200 {object} model.User
// @Failure 400 {object} handler.HttpMsg
// @Failure 404 {object} handler.HttpMsg
// @Failure 500 {object} handler.HttpMsg
// @Router /api/v1/user/{id} [patch]
func updateUser(service user.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		externalID := c.Param("id")
		id, err := uuid.Parse(externalID)
		if err != nil || id == uuid.Nil {
			ErroHttpMsgUserIdIsRequired.Write(c.Writer)
			return
		}

		var requestToUpdate model.User
		if err := c.ShouldBindJSON(&requestToUpdate); err != nil {
			ErroHttpMsgToParseRequestUserToJson.Write(c.Writer)
			return
		}

		if strings.TrimSpace(requestToUpdate.Username) == "" {
			ErroHttpMsgUserNameIsRequired.Write(c.Writer)
			return
		}

		if strings.TrimSpace(requestToUpdate.Name) == "" {
			ErroHttpMsgUserNameIsRequired.Write(c.Writer)
			return
		}

		if strings.TrimSpace(requestToUpdate.Password) == "" {
			ErroHttpMsgUserPasswordIsRequired.Write(c.Writer)
			return
		}

		user := service.GetByID(c.Request.Context(), id)
		if user.ID == uuid.Nil {
			ErroHttpMsgUserNotFound.Write(c.Writer)
			return
		}

		rowsAffected := service.Update(c.Request.Context(), id, &requestToUpdate)
		if rowsAffected == 0 {
			ErroHttpMsgToUpdateUser.Write(c.Writer)
			return
		}

		c.JSON(http.StatusOK, requestToUpdate)
	}
}

// @Summary Delete user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} handler.HttpMsg
// @Failure 400 {object} handler.HttpMsg
// @Failure 404 {object} handler.HttpMsg
// @Failure 500 {object} handler.HttpMsg
// @Router /api/v1/user/{id} [delete]
func deleteUser(service user.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		externalID := c.Param("id")
		id, err := uuid.Parse(externalID)
		if err != nil || id == uuid.Nil {
			ErroHttpMsgUserIdIsRequired.Write(c.Writer)
			return
		}

		user := service.GetByID(c.Request.Context(), id)
		if user.ID == uuid.Nil {
			ErroHttpMsgUserNotFound.Write(c.Writer)
			return
		}

		rowsAffected := service.Delete(c.Request.Context(), id)
		if rowsAffected == 0 {
			ErroHttpMsgToDeleteUser.Write(c.Writer)
			return
		}

		SuccessHttpMsgToDeleteUser.Write(c.Writer)
	}
}

// @Summary Get JWT token
// @Description Authenticate user and get JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} jwt.TokenDetails
// @Failure 400 {object} handler.HttpMsg
// @Failure 401 {object} handler.HttpMsg
// @Failure 500 {object} handler.HttpMsg
// @Router /api/v1/user/getjwt [post]
func getJWT(service user.UserServiceInterface, tenantService service_ten.TenantServiceInterface, tenantGroupService service_ten_group.TenantGroupServiceInterface, conf *config.Config, tokenService token.TokenServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest dto.LoginRequest

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			logger.Error("Invalid login request: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		user, err := service.Authenticate(loginRequest.Username, loginRequest.Password)
		if err != nil {
			logger.Error("Authentication failed: ", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Fetch tenant information
		tenant := tenantService.GetByID(c.Request.Context(), user.TenantID)
		if tenant.ID == uuid.Nil {
			logger.Error("Tenant not found for user: "+user.ID.String(), nil)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Tenant information not found"})
			return
		}

		// Fetch tenant group information (now mandatory)
		tenantGroup := tenantGroupService.GetByID(c.Request.Context(), tenant.GroupID)

		tokenDetails, err := jwt.GenerateToken(user, tenant, tenantGroup, conf, tokenService)
		if err != nil {
			logger.Error("Failed to generate JWT: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, tokenDetails)
	}
}

// @Summary Validate JWT token
// @Description Validate a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} jwt.Claims
// @Failure 401 {object} handler.HttpMsg
// @Router /api/v1/user/validatejwt [post]
func validateToken(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		claims, err := jwt.ValidateToken(tokenStr, conf)
		if err != nil {
			logger.Error("Token validation failed: ", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.JSON(http.StatusOK, claims)
	}
}

// @Summary Refresh JWT token
// @Description Refresh an expired JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {refresh_token}"
// @Success 200 {object} jwt.TokenDetails
// @Failure 401 {object} handler.HttpMsg
// @Router /api/v1/user/refreshjwt [post]
func refreshToken(conf *config.Config, tokenService token.TokenServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		refreshToken := strings.TrimPrefix(authHeader, "Bearer ")
		if refreshToken == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		tokenDetails, ok := jwt.RefreshJWT(refreshToken, conf, tokenService)
		if !ok {
			logger.Error("Token refresh failed: ", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		c.JSON(http.StatusOK, tokenDetails)
	}
}

// @Summary Logout user
// @Description Logout user by revoking refresh token
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} handler.HttpMsg
// @Failure 401 {object} handler.HttpMsg
// @Router /api/v1/user/logout [post]
func logout(tokenService token.TokenServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		tokenID, err := jwt.ExtractTokenID(authHeader)
		if err != nil {
			logger.Error("Failed to extract token ID: ", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		err = jwt.RevokeToken(tokenID, tokenService)
		if err != nil {
			logger.Error("Failed to revoke token: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
			"code":    200,
		})
	}
}

// @Summary Change user password
// @Description Change a user's password
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.UserChangePasswordOutPut true "Password change details"
// @Success 200 {object} handler.HttpMsg
// @Failure 400 {object} handler.HttpMsg
// @Failure 500 {object} handler.HttpMsg
// @Router /api/v1/user/changepassword [patch]
func changePassword(service user.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userChange dto.UserChangePasswordOutPut
		err := json.NewDecoder(r.Body).Decode(&userChange)
		if err != nil {
			logger.Error("Error decoding user change password: ", err)
			http.Error(w, "Error decoding user change password", http.StatusBadRequest)
			return
		}

		if userChange.Username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		if userChange.NewPassowrd == "" || userChange.OldPassowrd == "" {
			http.Error(w, "New password and confirm password are required", http.StatusBadRequest)
			return
		}

		err = service.ChangePassword(r.Context(), userChange.Username, userChange.OldPassowrd, userChange.NewPassowrd)
		if err != nil {
			errorMsg := err.Error()
			if errorMsg == "new password length is too short" ||
				errorMsg == "new password must contain at least one uppercase letter" ||
				errorMsg == "new password must contain at least one number" ||
				errorMsg == "new password must contain at least one symbol" {

				requirements := `Password must meet the following requirements:
1. At least 8 characters long
2. At least one uppercase letter
3. At least one number
4. At least one special symbol (!@#$%^&*()\-_+=)`

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(HttpMsg{
					Message: requirements,
					Code:    http.StatusBadRequest,
				})
				return
			}
			http.Error(w, errorMsg, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(SuccessHttpMsgToChangePassword.Code)
		json.NewEncoder(w).Encode(SuccessHttpMsgToChangePassword)
	}
}

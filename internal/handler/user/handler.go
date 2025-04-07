package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/internal/dto"
	"github.com/katana-stuidio/access-control/pkg/jwt"
	"github.com/katana-stuidio/access-control/pkg/model"
	"github.com/katana-stuidio/access-control/pkg/service/user"
)

func getAllUser(service user.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		users := service.GetAll(c.Request.Context())
		c.JSON(http.StatusOK, users)
	}
}

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

func createUser(service user.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userDto dto.UserRequestDtoInput

		if err := c.ShouldBindJSON(&userDto); err != nil {
			logger.Error("Invalid login request: ", err)
			ErroHttpMsgToParseRequestUserToJson.Write(c.Writer)
			return
		}

		userModel := model.User{
			Username: userDto.Username,
			Name:     userDto.Name,
			Password: userDto.Password,
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

		if !usrCad.CheckCpf(usrCad.Username) {
			ErroHttpMsgUserCpfIsInvalid.Write(c.Writer)
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
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		}

		c.JSON(http.StatusCreated, resultOut)
	}
}

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

func getJWT(service user.UserServiceInterface, conf *config.Config) gin.HandlerFunc {
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

		tokenDetails, err := jwt.GenerateToken(user, conf)
		if err != nil {
			logger.Error("Failed to generate JWT: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, tokenDetails)
	}
}

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

func refreshToken(conf *config.Config) gin.HandlerFunc {
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

		tokenDetails, ok := jwt.RefreshJWT(refreshToken, conf)
		if !ok {
			logger.Error("Token refresh failed: ", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		c.JSON(http.StatusOK, tokenDetails)
	}
}

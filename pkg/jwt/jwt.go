package jwt

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/katana-stuidio/access-control/internal/config"
	"github.com/katana-stuidio/access-control/pkg/model"
	"github.com/katana-stuidio/access-control/pkg/service/token"
)

type TokenDetails struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	TokenID      string `json:"tokenId,omitempty"`
}

type Claims struct {
	Username    string `json:"username"`
	UserID      string `json:"user_id"`
	TenantID    string `json:"tenant_id"`
	TenantName  string `json:"tenant_name,omitempty"`
	GroupID     string `json:"group_id,omitempty"`
	GroupName   string `json:"group_name,omitempty"`
	Role        string `json:"role"`
	FirstAccess bool   `json:"first_access"`
	Renew       bool   `json:"renew,omitempty"`
	TokenID     string `json:"token_id,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken generates both access and refresh tokens with Redis integration
func GenerateToken(user *model.User, tenant *model.Tenant, tenantGroup *model.TenantGroup, conf *config.Config, tokenService token.TokenServiceInterface) (*TokenDetails, error) {
	jwtKey := []byte(conf.JWTSecretKey)

	// Generate a unique token ID for Redis storage
	tokenID, err := generateTokenID()
	if err != nil {
		log.Println("Error generating token ID:", err)
		return nil, err
	}

	// Generate Access Token (short-lived)
	accessExpiration := time.Now().Add(time.Duration(conf.JWTTokenExp) * time.Minute)
	accessClaims := &Claims{
		Username:   user.Username,
		UserID:     user.ID.String(),
		TenantID:   user.TenantID.String(),
		TenantName: tenant.Name,
		Role:       user.Role,
		Renew:      false,
		TokenID:    tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Add group information (now mandatory)
	if tenantGroup != nil {
		accessClaims.GroupID = tenantGroup.ID.String()
		accessClaims.GroupName = tenantGroup.Name
	}
	accessToken, err := createToken(accessClaims, jwtKey)
	if err != nil {
		log.Println("Error generating access token:", err)
		return nil, err
	}

	// Generate Refresh Token (long-lived)
	refreshExpiration := time.Now().Add(time.Duration(conf.JWTRefreshExp) * time.Minute)
	refreshClaims := &Claims{
		Username:   user.Username,
		UserID:     user.ID.String(),
		TenantID:   user.TenantID.String(),
		TenantName: tenant.Name,
		Role:       user.Role,
		Renew:      true,
		TokenID:    tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Add group information (now mandatory)
	if tenantGroup != nil {
		refreshClaims.GroupID = tenantGroup.ID.String()
		refreshClaims.GroupName = tenantGroup.Name
	}
	refreshToken, err := createToken(refreshClaims, jwtKey)
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return nil, err
	}

	// Save refresh token to Redis
	ctx := context.Background()
	err = tokenService.SaveRefreshToken(
		ctx,
		tokenID,
		user.ID.String(),
		user.Username,
		user.TenantID.String(),
		user.Role,
		time.Now(),
		refreshExpiration,
	)
	if err != nil {
		log.Println("Error saving refresh token to Redis:", err)
		return nil, err
	}

	return &TokenDetails{
		AccessToken:  fmt.Sprintf("Bearer %s", accessToken),
		RefreshToken: fmt.Sprintf("Bearer %s", refreshToken),
		TokenID:      tokenID,
	}, nil
}

// ValidateToken validates a token string and returns the claims
func ValidateToken(tokenStr string, conf *config.Config) (*Claims, error) {
	jwtKey := []byte(conf.JWTSecretKey)
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			log.Println("Invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		log.Println("Invalid token")
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshJWT creates a new access token using a valid refresh token with Redis validation
func RefreshJWT(tknStr string, conf *config.Config, tokenService token.TokenServiceInterface) (token *TokenDetails, ok bool) {
	jwtKey := []byte(conf.JWTSecretKey)

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		log.Println(err.Error())
		if err == jwt.ErrSignatureInvalid {
			return token, false
		}
		return token, false
	}

	if !claims.Renew {
		log.Println("Error: this is not a valid refresh token")
		return token, false
	}

	if !tkn.Valid {
		log.Println("Error: the token is not valid")
		return token, false
	}

	// Validate refresh token against Redis
	ctx := context.Background()
	isValid, err := tokenService.IsTokenValid(ctx, claims.TokenID)
	if err != nil || !isValid {
		log.Println("Error: refresh token not found in Redis or invalid")
		return token, false
	}

	// Generate new access token
	claims.Renew = false
	expirationTime := time.Now().Add(time.Duration(conf.JWTTokenExp) * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	newThk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	thkString, err := newThk.SignedString(jwtKey)
	if err != nil {
		log.Println("Error creating new token")
		return token, false
	}

	token = &TokenDetails{
		AccessToken: fmt.Sprintf("Bearer %s", thkString),
		TokenID:     claims.TokenID,
	}

	return token, true
}

// RevokeToken removes a refresh token from Redis
func RevokeToken(tokenID string, tokenService token.TokenServiceInterface) error {
	ctx := context.Background()
	return tokenService.DeleteRefreshToken(ctx, tokenID)
}

// RevokeAllUserTokens removes all refresh tokens for a user
func RevokeAllUserTokens(userID string, tokenService token.TokenServiceInterface) error {
	ctx := context.Background()
	return tokenService.DeleteAllUserTokens(ctx, userID)
}

// RevokeAllTenantTokens removes all refresh tokens for a tenant
func RevokeAllTenantTokens(tenantID string, tokenService token.TokenServiceInterface) error {
	ctx := context.Background()
	return tokenService.DeleteAllTenantTokens(ctx, tenantID)
}

// Helper function to create a new token
func createToken(claims *Claims, jwtKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Helper function to generate a unique token ID
func generateTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ExtractTokenID extracts token ID from Authorization header
func ExtractTokenID(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return "", errors.New("invalid token format")
	}

	// Parse the token to extract claims
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(""), nil // We don't need to validate here, just extract claims
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.TokenID, nil
}

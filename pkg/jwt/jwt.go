package jwt

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/katana-stuidio/access-control/internal/config"
	"github.com/katana-stuidio/access-control/pkg/model"
)

type TokenDetails struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type Claims struct {
	Username    string `json:"username"`
	UserID      string `json:"user_id"`
	TenantID    string `json:"tenant_id"`
	Role        string `json:"role"`
	FirstAccess bool   `json:"first_access"`
	Renew       bool   `json:"renew,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken generates both access and refresh tokens
func GenerateToken(user *model.User, conf *config.Config) (*TokenDetails, error) {
	jwtKey := []byte(conf.JWTSecretKey)

	// Generate Access Token
	accessExpiration := time.Now().Add(time.Duration(conf.JWTTokenExp) * time.Minute)
	accessClaims := &Claims{
		Username: user.Username,
		UserID:   user.ID.String(),
		TenantID: user.TenantID.String(),
		Role:     user.Role,
		Renew:    false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiration),
		},
	}
	accessToken, err := createToken(accessClaims, jwtKey)
	if err != nil {
		log.Println("Error generating access token:", err)
		return nil, err
	}

	// Generate Refresh Token
	refreshExpiration := time.Now().Add(time.Duration(conf.JWTRefreshExp) * time.Minute)
	refreshClaims := &Claims{
		Username: user.Username,
		UserID:   user.ID.String(),
		TenantID: user.TenantID.String(),
		Role:     user.Role,
		Renew:    true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
		},
	}
	refreshToken, err := createToken(refreshClaims, jwtKey)
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return nil, err
	}

	return &TokenDetails{
		AccessToken:  fmt.Sprintf("Bearer %s", accessToken),
		RefreshToken: fmt.Sprintf("Bearer %s", refreshToken),
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

// RefreshToken creates a new access token using a valid refresh token
func RefreshJWT(tknStr string, conf *config.Config) (token *TokenDetails, ok bool) {

	jwtKey := []byte(conf.JWTSecretKey)

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) { return jwtKey, nil })

	if err != nil {
		log.Println(err.Error())
		if err == jwt.ErrSignatureInvalid {
			return token, false
		}
		return token, false
	}

	if !claims.Renew {
		log.Println("Error this not a valid refreshToken")
		return token, false
	}

	if !tkn.Valid {
		log.Println("Error the token is not valid")
		return token, false
	}

	claims.Renew = false // set false to create a simple Token when true creates a RefreshToken

	expirationTime := time.Now().Add(time.Duration(conf.JWTTokenExp) * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	newThk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	thkString, err := newThk.SignedString(jwtKey)
	if err != nil {
		log.Println("Error the make a new token")
		return token, false
	}

	token = &TokenDetails{
		AccessToken: fmt.Sprintf("Bearer %s", thkString),
	}

	return token, true
}

// Helper function to create a new token
func createToken(claims *Claims, jwtKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

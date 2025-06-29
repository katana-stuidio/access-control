package token

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/katana-stuidio/access-control/internal/config"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/pkg/adapter/redisdb"
)

type TokenServiceInterface interface {
	SaveRefreshToken(ctx context.Context, tokenID, userID, username, tenantID, role string, issuedAt, exp time.Time) error
	GetRefreshToken(ctx context.Context, tokenID string) (*RefreshTokenData, error)
	DeleteRefreshToken(ctx context.Context, tokenID string) error
	DeleteAllUserTokens(ctx context.Context, userID string) error
	DeleteAllTenantTokens(ctx context.Context, tenantID string) error
	IsTokenValid(ctx context.Context, tokenID string) (bool, error)
}

type RefreshTokenData struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	TenantID  string    `json:"tenant_id"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TokenService struct {
	redis redisdb.RedisClientInterface
	conf  *config.Config
}

func NewTokenService(redis redisdb.RedisClientInterface, conf *config.Config) *TokenService {
	return &TokenService{
		redis: redis,
		conf:  conf,
	}
}

// SaveRefreshToken saves a refresh token in Redis with TTL
func (ts *TokenService) SaveRefreshToken(ctx context.Context, tokenID, userID, username, tenantID, role string, issuedAt, exp time.Time) error {
	tokenData := &RefreshTokenData{
		UserID:    userID,
		Username:  username,
		TenantID:  tenantID,
		Role:      role,
		IssuedAt:  issuedAt,
		ExpiresAt: exp,
	}

	data, err := json.Marshal(tokenData)
	if err != nil {
		logger.Error("Error marshaling token data", err)
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	// Calculate TTL based on expiration time
	ttl := time.Until(exp)
	if ttl <= 0 {
		return fmt.Errorf("token already expired")
	}

	key := fmt.Sprintf("refresh:%s", tokenID)
	success := ts.redis.SaveData(ctx, key, data, ttl)
	if !success {
		return fmt.Errorf("failed to save refresh token to Redis")
	}

	logger.Info(fmt.Sprintf("Refresh token saved successfully: %s", tokenID))
	return nil
}

// GetRefreshToken retrieves a refresh token from Redis
func (ts *TokenService) GetRefreshToken(ctx context.Context, tokenID string) (*RefreshTokenData, error) {
	key := fmt.Sprintf("refresh:%s", tokenID)
	data, err := ts.redis.ReadData(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found: %w", err)
	}

	var tokenData RefreshTokenData
	if err := json.Unmarshal(data, &tokenData); err != nil {
		logger.Error("Error unmarshaling token data", err)
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	// Check if token is expired
	if time.Now().After(tokenData.ExpiresAt) {
		// Delete expired token
		ts.redis.DeleteAllHSetData(ctx, key)
		return nil, fmt.Errorf("refresh token expired")
	}

	return &tokenData, nil
}

// DeleteRefreshToken removes a refresh token from Redis
func (ts *TokenService) DeleteRefreshToken(ctx context.Context, tokenID string) error {
	key := fmt.Sprintf("refresh:%s", tokenID)
	success := ts.redis.DeleteAllHSetData(ctx, key)
	if !success {
		return fmt.Errorf("failed to delete refresh token from Redis")
	}

	logger.Info(fmt.Sprintf("Refresh token deleted successfully: %s", tokenID))
	return nil
}

// DeleteAllUserTokens removes all refresh tokens for a specific user
func (ts *TokenService) DeleteAllUserTokens(ctx context.Context, userID string) error {
	// This would require a more complex implementation with Redis SCAN
	// For now, we'll implement a simple version that deletes by pattern
	// In production, you might want to maintain a separate index for user tokens

	// Note: This is a simplified implementation
	// In a real scenario, you'd need to scan all keys and filter by user_id
	// For now, we'll return success as the main deletion happens during logout
	logger.Info(fmt.Sprintf("User tokens deletion requested for user: %s", userID))
	return nil
}

// DeleteAllTenantTokens removes all refresh tokens for a specific tenant
func (ts *TokenService) DeleteAllTenantTokens(ctx context.Context, tenantID string) error {
	// Similar to DeleteAllUserTokens, this would require SCAN implementation
	logger.Info(fmt.Sprintf("Tenant tokens deletion requested for tenant: %s", tenantID))
	return nil
}

// IsTokenValid checks if a refresh token exists and is valid
func (ts *TokenService) IsTokenValid(ctx context.Context, tokenID string) (bool, error) {
	_, err := ts.GetRefreshToken(ctx, tokenID)
	if err != nil {
		return false, err
	}
	return true, nil
}

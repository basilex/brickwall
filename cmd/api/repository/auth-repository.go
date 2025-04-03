package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"brickwall/internal/common"
	"brickwall/internal/provider"
)

const (
	AuthPrefix = "auth:"

	RefreshTokenPrefix   = AuthPrefix + "refresh:"
	BlacklistTokenPrefix = AuthPrefix + "blacklist:"
	PasswordResetPrefix  = AuthPrefix + "reset:"
)

type IAuthRepository interface {
	StoreToken(prefix, userID, tokenID string, ttl time.Duration) error
	IsTokenExist(prefix, tokenID string) (bool, error)
	DeleteToken(prefix, userID, tokenID string) error

	StoreRefreshToken(userID, tokenID string, ttl time.Duration) error
	IsTokenBlacklisted(tokenID string) (bool, error)
	BlacklistToken(tokenID string, ttl time.Duration) error
	DeleteRefreshToken(userID, tokenID string) error

	StorePasswordResetToken(userID, token string, ttl time.Duration) error
	GetPasswordResetToken(userID string) (string, error)
	DeletePasswordResetToken(userID string) error
}

type AuthRepository struct {
	ctx   context.Context
	redis *redis.Client
}

func NewAuthRepository(ctx context.Context) IAuthRepository {
	redis := ctx.Value(common.KeyRedisProvider).(provider.IRedisProvider)

	return &AuthRepository{
		ctx:   ctx,
		redis: redis.Client(),
	}
}

// Universal code for token saving
func (rcv *AuthRepository) StoreToken(prefix, userID, tokenID string, ttl time.Duration) error {
	key := prefix + userID + ":" + tokenID
	return rcv.redis.Set(rcv.ctx, key, "valid", ttl).Err()
}

// Universal method for checking of existing token
func (rcv *AuthRepository) IsTokenExist(prefix, tokenID string) (bool, error) {
	key := prefix + tokenID
	exists, err := rcv.redis.Exists(rcv.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// Universal method for token deletions
func (rcv *AuthRepository) DeleteToken(prefix, userID, tokenID string) error {
	key := prefix + userID + ":" + tokenID
	return rcv.redis.Del(rcv.ctx, key).Err()
}

// Save refresh token
func (rcv *AuthRepository) StoreRefreshToken(userID, tokenID string, ttl time.Duration) error {
	return rcv.StoreToken(RefreshTokenPrefix, userID, tokenID, ttl)
}

// Check if the token is int the blacklist
func (rcv *AuthRepository) IsTokenBlacklisted(tokenID string) (bool, error) {
	return rcv.IsTokenExist(BlacklistTokenPrefix, tokenID)
}

// Add token to blacklist with TTL
func (rcv *AuthRepository) BlacklistToken(tokenID string, ttl time.Duration) error {
	key := BlacklistTokenPrefix + tokenID
	return rcv.redis.SetEx(rcv.ctx, key, "revoked", ttl).Err()
}

// Delete the refresh token
func (rcv *AuthRepository) DeleteRefreshToken(userID, tokenID string) error {
	return rcv.DeleteToken(RefreshTokenPrefix, userID, tokenID)
}

// Save the token of the password reset (SetNX - not existing data)
func (rcv *AuthRepository) StorePasswordResetToken(userID, token string, ttl time.Duration) error {
	key := PasswordResetPrefix + userID
	_, err := rcv.redis.SetNX(rcv.ctx, key, token, ttl).Result()
	return err
}

// Get the token of the password reset
func (rcv *AuthRepository) GetPasswordResetToken(userID string) (string, error) {
	key := PasswordResetPrefix + userID
	return rcv.redis.Get(rcv.ctx, key).Result()
}

// Delete the token of the password reset
func (rcv *AuthRepository) DeletePasswordResetToken(userID string) error {
	key := PasswordResetPrefix + userID
	return rcv.redis.Del(rcv.ctx, key).Err()
}

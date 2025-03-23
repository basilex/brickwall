package adapter

import (
	"context"
	"fmt"
	"time"

	"brickwall/internal/common"
	"brickwall/internal/provider"
)

type IJwtAdapter interface {
	ValidateToken(string) (*provider.Claims, error)
	InvalidateToken(string, time.Duration) error
	IsTokenInvalidated(string) bool
	StoreToken(string, time.Duration) error
	DeleteToken(string) error
}

type JwtAdapter struct {
	ctx           context.Context
	redisProvider provider.IRedisProvider
	jwtProvider   provider.IJwtProvider
}

func NewJwtAdapter(ctx context.Context) IJwtAdapter {
	return &JwtAdapter{
		ctx:           ctx,
		jwtProvider:   ctx.Value(common.KeyJwtProvider).(provider.IJwtProvider),
		redisProvider: ctx.Value(common.KeyRedisProvider).(provider.IRedisProvider),
	}
}

func (rcv *JwtAdapter) ValidateToken(tokenString string) (*provider.Claims, error) {
	val, err := rcv.redisProvider.Client().Get(rcv.ctx, tokenString).Result()
	if err == nil && val == common.JwtTokenInvalid {
		return nil, fmt.Errorf("%w: %v", common.ErrJwtTokenInvalidated, "marked as invalid")
	}
	return rcv.jwtProvider.ValidateToken(tokenString)
}

func (rcv *JwtAdapter) InvalidateToken(tokenString string, accessExpiration time.Duration) error {
	return rcv.redisProvider.Client().Set(rcv.ctx, tokenString, common.JwtTokenInvalid, accessExpiration).Err()
}

func (rcv *JwtAdapter) IsTokenInvalidated(tokenString string) bool {
	val, err := rcv.redisProvider.Client().Get(rcv.ctx, tokenString).Result()
	return err == nil && val == common.JwtTokenInvalid
}

func (rcv *JwtAdapter) StoreToken(tokenString string, accessExpiration time.Duration) error {
	return rcv.redisProvider.Client().Set(rcv.ctx, tokenString, common.JwtTokenValid, accessExpiration).Err()
}

func (rcv *JwtAdapter) DeleteToken(tokenString string) error {
	return rcv.redisProvider.Client().Del(rcv.ctx, tokenString).Err()
}

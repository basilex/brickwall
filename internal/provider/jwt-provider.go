package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v3"

	"brickwall/internal/common"
)

const (
	TokenValid   = "valid"
	TokenInvalid = "invalid"
)

type IJwtProvider interface {
	GenerateTokens(string) (string, string, error)
	RefreshTokens(string) (string, string, error)
	ValidateToken(string) (*Claims, error)
	InvalidateToken(string) error
	IsTokenInvalidated(string) bool
	StoreToken(string) error
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

type JwtProvider struct {
	ctx               context.Context
	redis             *redis.Client
	secret            string
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func NewJwtProvider(ctx context.Context) IJwtProvider {
	cli := ctx.Value(common.KeyCommand).(*cli.Command)
	redis := ctx.Value(common.KeyRedisProvider).(IRedisProvider)

	return &JwtProvider{
		ctx:               ctx,
		redis:             redis.Client(),
		secret:            cli.String("jwt-secret"),
		accessExpiration:  cli.Duration("jwt-access-expiration"),
		refreshExpiration: cli.Duration("jwt-refresh-expiration"),
	}
}

func (rcv *JwtProvider) GenerateTokens(userID string) (string, string, error) {
	AccessExpiration := time.Now().Add(rcv.accessExpiration)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(AccessExpiration),
		},
	})
	signedAccessToken, err := accessToken.SignedString([]byte(rcv.secret))
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", common.ErrJwtTokenSigning, err)
	}
	refreshExpiration := time.Now().Add(rcv.refreshExpiration)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshExpiration),
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte(rcv.secret))
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", common.ErrJwtTokenSigning, err)
	}
	return signedAccessToken, signedRefreshToken, nil
}

func (rcv *JwtProvider) RefreshTokens(tokenString string) (string, string, error) {
	claims, err := rcv.ValidateToken(tokenString)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", common.ErrJwtTokenInvalid, err)
	}
	return rcv.GenerateTokens(claims.UserID)
}

func (rcv *JwtProvider) ValidateToken(tokenString string) (*Claims, error) {
	val, err := rcv.redis.Get(context.Background(), tokenString).Result()
	if err == nil && val == TokenInvalid {
		return nil, fmt.Errorf("%w: %v", common.ErrJwtTokenInvalidated, "marked as invalid")
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return rcv.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrJwtTokenClaims, "failed to parse claims")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("%w: %v", common.ErrJwtTokenClaims, "failed to bind claims")
	}
	return claims, nil
}

func (rcv *JwtProvider) InvalidateToken(tokenString string) error {
	ctx := context.Background()
	return rcv.redis.Set(ctx, tokenString, TokenInvalid, rcv.accessExpiration).Err()
}

func (rcv *JwtProvider) IsTokenInvalidated(tokenString string) bool {
	ctx := context.Background()
	val, err := rcv.redis.Get(ctx, tokenString).Result()
	return err == nil && val == TokenInvalid
}

func (rcv *JwtProvider) StoreToken(tokenString string) error {
	ctx := context.Background()
	return rcv.redis.Set(ctx, tokenString, TokenValid, rcv.accessExpiration).Err()
}

func (rcv *JwtProvider) DeleteToken(tokenString string) error {
	ctx := context.Background()
	return rcv.redis.Del(ctx, tokenString).Err()
}

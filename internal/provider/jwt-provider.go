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

const TokenIvalid = "invalid"

type IJwtProvider interface {
	GenerateTokens(string) (string, string, error)
	RefreshTokens(string) (string, string, error)
	ValidateToken(string) (*Claims, error)
	InvalidateToken(string) error
	IsTokenInvalidated(string) bool
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
	redis := ctx.Value(common.KeyRedisProvider).(*redis.Client)

	return &JwtProvider{
		ctx:               ctx,
		redis:             redis,
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
	signedAccessToken, err := accessToken.SignedString(rcv.secret)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", common.ErrJwtTokenSigning, err)
	}
	refreshExpiration := time.Now().Add(rcv.refreshExpiration)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshExpiration),
	})
	signedRefreshToken, err := refreshToken.SignedString(rcv.secret)
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
	if err == nil && val == TokenIvalid {
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
	return rcv.redis.Set(ctx, tokenString, TokenIvalid, rcv.accessExpiration).Err()
}

func (rcv *JwtProvider) IsTokenInvalidated(tokenString string) bool {
	ctx := context.Background()
	val, err := rcv.redis.Get(ctx, tokenString).Result()
	return err == nil && val == TokenIvalid
}

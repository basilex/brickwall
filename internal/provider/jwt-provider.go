package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/xid"

	"brickwall/internal/common"
)

type IJwtProvider interface {
	GenerateAllTokens(string) (string, string, error)
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)

	RefreshTokens(string) (string, string, error)
	ValidateToken(string) (*Claims, error)
	InvalidateToken(string) error
	IsTokenInvalidated(string) bool
	StoreToken(string) error
	DeleteToken(string) error
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
	env := ctx.Value(common.KeyEnvProvider).(IEnvProvider)
	redis := ctx.Value(common.KeyRedisProvider).(IRedisProvider)

	return &JwtProvider{
		ctx:               ctx,
		redis:             redis.Client(),
		secret:            env.GetString("JWT_SECRET", DefJwtSecret),
		accessExpiration:  env.GetDuration("JWT_ACCESS_EXPIRATION", DefJwtAccessExpiration),
		refreshExpiration: env.GetDuration("JWT_REFRESH_EXPIRATION", DefJwtRefreshExpiration),
	}
}

/////////////////////////////////////////////////////////////////////////////// old

// func (rcv *JwtProvider) GenerateTokens(userID string) (string, string, error) {
// 	// access token
// 	accessExpiration := time.Now().Add(rcv.accessExpiration)
// 	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
// 		UserID: userID,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ID:        xid.New().String(),
// 			ExpiresAt: jwt.NewNumericDate(accessExpiration),
// 		},
// 	})
// 	signedAccessToken, err := accessToken.SignedString([]byte(rcv.secret))
// 	if err != nil {
// 		return "", "", fmt.Errorf("%w: %v", common.ErrJwtTokenSigning, err)
// 	}

// 	// refresh token
// 	refreshExpiration := time.Now().Add(rcv.refreshExpiration)
// 	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
// 		ID:        xid.New().String(),
// 		ExpiresAt: jwt.NewNumericDate(refreshExpiration),
// 	})
// 	signedRefreshToken, err := refreshToken.SignedString([]byte(rcv.secret))
// 	if err != nil {
// 		return "", "", fmt.Errorf("%w: %v", common.ErrJwtTokenSigning, err)
// 	}
// 	return signedAccessToken, signedRefreshToken, nil
// }

// ///////////////////////////////////////////////////////////////////////////// new start
func (rcv *JwtProvider) GenerateAllTokens(userID string) (string, string, error) {
	access, err := rcv.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}
	refresh, err := rcv.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (rcv *JwtProvider) GenerateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"jti": xid.New().String(),
		"exp": time.Now().Add(rcv.accessExpiration),
	})
	return token.SignedString([]byte(rcv.secret))
}

func (rcv *JwtProvider) GenerateRefreshToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"jti": xid.New().String(),
		"exp": time.Now().Add(rcv.refreshExpiration),
	})
	return token.SignedString([]byte(rcv.secret))
}

/////////////////////////////////////////////////////////////////////////////// new end

func (rcv *JwtProvider) RefreshTokens(tokenString string) (string, string, error) {
	claims, err := rcv.ValidateToken(tokenString)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", common.ErrJwtTokenInvalid, err)
	}
	return rcv.GenerateAllTokens(claims.UserID)
}

func (rcv *JwtProvider) ValidateToken(tokenString string) (*Claims, error) {
	val, err := rcv.redis.Get(rcv.ctx, tokenString).Result()
	if err == nil && val == common.JwtTokenInvalid {
		return nil, fmt.Errorf("%w: %v", common.ErrJwtTokenInvalidated, "marked as invalid")
	}
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
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
	return rcv.redis.Set(rcv.ctx, tokenString, common.JwtTokenInvalid, rcv.refreshExpiration).Err()
}

func (rcv *JwtProvider) IsTokenInvalidated(tokenString string) bool {
	val, err := rcv.redis.Get(rcv.ctx, tokenString).Result()
	return err == nil && val == common.JwtTokenInvalid
}

func (rcv *JwtProvider) StoreToken(tokenString string) error {
	return rcv.redis.Set(rcv.ctx, tokenString, common.JwtTokenValid, rcv.accessExpiration).Err()
}

func (rcv *JwtProvider) DeleteToken(tokenString string) error {
	return rcv.redis.Del(rcv.ctx, tokenString).Err()
}

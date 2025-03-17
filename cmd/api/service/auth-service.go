package service

import (
	"context"
	"errors"

	"brickwall/internal/common"
	"brickwall/internal/provider"
	"brickwall/internal/storage/dbs"
)

type IAuthService interface {
	LoginWith2FA(string, string, string) (bool, error)
}

type AuthService struct {
	ctx     context.Context
	queries *dbs.Queries

	twoFA provider.I2FAProvider
}

func NewAuthService(ctx context.Context, queries *dbs.Queries) IAuthService {
	twoFA := ctx.Value(common.Key2FAProvider).(provider.I2FAProvider)
	return &AuthService{ctx: ctx, twoFA: twoFA}
}

func (rcv *AuthService) LoginWith2FA(userEmail, password, code string) (bool, error) {
	// Заглушка проверки пароля (заменить на реальную логику)
	if password != "correct_password" {
		return false, errors.New("invalid credentials")
	}

	// Проверяем, включен ли 2FA (секрет должен быть в БД)
	secret := "user_saved_secret" // Достать из БД по userEmail
	if secret != "" {
		if !rcv.twoFA.VerifyCode(secret, code) {
			return false, errors.New("invalid 2FA code")
		}
	}

	return true, nil
}

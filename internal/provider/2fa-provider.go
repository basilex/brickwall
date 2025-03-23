package provider

import (
	"context"
	"fmt"

	"github.com/pquerna/otp/totp"

	"brickwall/internal/common"
)

const Issuer = "BSP"

type I2FAProvider interface {
	GenerateSecretKey(string) (string, string, error)
	VerifyCode(string, string) bool
}

type TwoFAProvider struct {
	ctx context.Context
}

func New2FAProvider(ctx context.Context) I2FAProvider {
	return &TwoFAProvider{ctx: ctx}
}

func (rcv *TwoFAProvider) GenerateSecretKey(userEmail string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      Issuer,
		AccountName: userEmail,
	})
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", common.Err2FAKeyGeneration, err)
	}
	return key.Secret(), key.URL(), nil
}

func (s *TwoFAProvider) VerifyCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

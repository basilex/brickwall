package handler

import (
	"log/slog"

	"brickwall/cmd/api/exchange"
	"brickwall/internal/provider"
)

func UserResetPasswordEmail(data []byte, encoder provider.IEncoder) error {
	var user *exchange.AuthUserResetRes

	if err := encoder.Decode(data, &user); err != nil {
		return err
	}
	slog.Info("topic: user-reset-password-email received", "data", user)
	return nil
}

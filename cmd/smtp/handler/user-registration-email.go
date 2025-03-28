package handler

import (
	"log/slog"

	"brickwall/internal/provider"
	"brickwall/internal/storage/dbs"
)

func UserRegistrationEmail(data []byte, encoder provider.IEncoder) error {
	var user *dbs.UserNewRow

	if err := encoder.Decode(data, &user); err != nil {
		return err
	}
	slog.Info("topic: user-registration-email received", "user-new-row", user)
	return nil
}

package service

import (
	"context"
	"errors"
	"fmt"

	"brickwall/cmd/api/exchange"
	"brickwall/internal/common"
	"brickwall/internal/provider"
	"brickwall/internal/storage/dbs"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Signup(*exchange.AuthSignupReq) (*dbs.UserNewRow, error)
	Signin(*exchange.AuthSigninReq) (*dbs.UserSelectRow, error)
	Signout() (bool, error)
}
type AuthService struct {
	ctx     context.Context
	queries *dbs.Queries

	pgxProvider   provider.IPgxProvider
	twoFAProvider provider.I2FAProvider
}

func NewAuthService(ctx context.Context, queries *dbs.Queries) IAuthService {
	return &AuthService{
		ctx:           ctx,
		pgxProvider:   ctx.Value(common.KeyPgxProvider).(provider.IPgxProvider),
		twoFAProvider: ctx.Value(common.Key2FAProvider).(provider.I2FAProvider),
	}
}

func (rcv *AuthService) Signup(req *exchange.AuthSignupReq) (*dbs.UserNewRow, error) {
	ctx := context.Background()

	// begin new transaction
	trx, err := rcv.pgxProvider.Pool().BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrDBTrxError, err)
	}
	defer func() {
		if err != nil {
			trx.Rollback(ctx)
		}
	}()
	qtx := rcv.queries.WithTx(trx)

	// create user
	passwordCrypted, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	user, err := qtx.UserNew(context.Background(), &dbs.UserNewParams{
		Username: req.Email,
		Password: string(passwordCrypted),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrDBRecordInsert, err)
	}

	// create contact email
	_, err = qtx.ContactNew(ctx, &dbs.ContactNewParams{
		UserID: user.ID, Class: "email", Content: req.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrDBRecordInsert, err)
	}
	// create default profile
	_, err = qtx.ProfileNew(ctx, &dbs.ProfileNewParams{
		UserID:    user.ID,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrDBRecordInsert, err)
	}

	// send an verification email
	// TODO: should be implemented

	// commit transaction
	if err := trx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrDBTrxError, err)
	}
	return user, nil
}

func (rcv *AuthService) Signin(*exchange.AuthSigninReq) (*dbs.UserSelectRow, error) {
	return nil, fmt.Errorf("%w: %v", common.ErrNotImplemented, errors.New("Auth.Signin()"))
}

func (rcv *AuthService) Signout() (bool, error) {
	return true, fmt.Errorf("%w: %v", common.ErrNotImplemented, errors.New("Auth.Signout()"))
}

// --------------------------------------------------------------------------------------
// func (rcv *AuthService) LoginWith2FA(userEmail, password, code string) (bool, error) {
// 	// Заглушка проверки пароля (заменить на реальную логику)
// 	if password != "correct_password" {
// 		return false, errors.New("invalid credentials")
// 	}

// 	// Проверяем, включен ли 2FA (секрет должен быть в БД)
// 	secret := "user_saved_secret" // Достать из БД по userEmail
// 	if secret != "" {
// 		if !rcv.twoFA.VerifyCode(secret, code) {
// 			return false, errors.New("invalid 2FA code")
// 		}
// 	}

// 	return true, nil
// }

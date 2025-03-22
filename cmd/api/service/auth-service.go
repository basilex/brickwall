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
	Signin(*exchange.AuthSigninReq) (*exchange.AuthSigninRes, error)
	Signout() (bool, error)
}
type AuthService struct {
	ctx     context.Context
	queries *dbs.Queries

	pgxProvider   provider.IPgxProvider
	jwtProvider   provider.IJwtProvider
	twoFAProvider provider.I2FAProvider
}

func NewAuthService(ctx context.Context, queries *dbs.Queries) IAuthService {
	return &AuthService{
		ctx:           ctx,
		queries:       queries,
		pgxProvider:   ctx.Value(common.KeyPgxProvider).(provider.IPgxProvider),
		jwtProvider:   ctx.Value(common.KeyJwtProvider).(provider.IJwtProvider),
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

func (rcv *AuthService) Signin(req *exchange.AuthSigninReq) (*exchange.AuthSigninRes, error) {
	ctx := context.Background()

	// check user password
	user, err := rcv.queries.AuthSelectUserCredentials(ctx, req.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%w: %v", common.ErrDBNotFound, err)
		} else {
			return nil, fmt.Errorf("%w: %v", common.ErrDBRecordSelect, err)
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrAuthInvalidPassword, err)
	}

	// check for blocked or checked statuses
	if user.IsBlocked {
		return nil, fmt.Errorf("%w: %v", common.ErrAuthUserBlocked, errors.New("block status detected"))
	}
	if !user.IsChecked {
		return nil, fmt.Errorf("%w: %v", common.ErrAuthUserNotChecked, errors.New("email check required"))
	}

	// generate tokens
	accessToken, refreshToken, err := rcv.jwtProvider.GenerateTokens(user.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrAuthGenerateTokens, err)
	}

	// update user visited_at
	updated, err := rcv.queries.AuthUpdateVisitedAt(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrDBRecordUpdate, err)
	}

	// store refresh token in redis
	if err := rcv.jwtProvider.StoreToken(refreshToken); err != nil {
		return nil, fmt.Errorf("%w: %v", common.ErrDBRecordInsert, err)
	}

	// response about logged in user
	res := &exchange.AuthSigninRes{
		User: &exchange.AuthUser{
			ID:        updated.ID,
			Username:  updated.Username,
			CheckedAt: updated.CheckedAt,
			VisitedAt: updated.VisitedAt,
			CreatedAt: updated.CreatedAt,
		},
		Tokens: &exchange.AuthTokens{
			Access:  accessToken,
			Refresh: refreshToken,
		},
	}
	return res, nil
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

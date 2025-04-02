package exchange

import "github.com/jackc/pgx/v5/pgtype"

type AuthSignupReq struct {
	Email     string `json:"email" binding:"required,email,max=255"`
	Password  string `json:"password" binding:"required,min=4,max=72"`
	Firstname string `json:"firstname" binding:"required,min=1,max=255"`
	Lastname  string `json:"lastname" binding:"required,min=1,max=255"`
}
type AuthSigninReq struct {
	Username string `json:"username" binding:"required,min=1,max=64"`
	Password string `json:"password" binding:"required,min=4,max=72"`
}
type AuthTokenRefreshReq struct {
	Token string `json:"token" binding:"required"`
}
type AuthTokenInvalidateReq struct {
	Token string `json:"token" binding:"required"`
}
type AuthPasswordResetReq struct {
	Email string `json:"email" binding:"required,email,max=255"`
}
type AuthPasswordChangeReq struct {
	Token    string `json:"token" binding:"required,max=255"`
	Password string `json:"password" binding:"required,min=4,max=72"`
}

// responses
type AuthTokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

// signup response
type AuthUserSignup struct {
	ID        string           `json:"id"`
	Username  string           `json:"username"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type AuthUserContact struct {
	ID        string           `json:"id"`
	Class     string           `json:"class"`
	Content   string           `json:"content"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}
type AuthUserProfile struct {
	ID        string           `json:"id"`
	Firstname string           `json:"firstname"`
	Lastname  string           `json:"lastname"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type AuthUserSignupRes struct {
	User    *AuthUserSignup  `json:"user"`
	Contact *AuthUserContact `json:"contact"`
	Profile *AuthUserProfile `json:"profile"`
}

// signin response
type AuthUserSignin struct {
	ID        string           `json:"id"`
	Username  string           `json:"username"`
	CheckedAt pgtype.Timestamp `json:"checked_at"`
	VisitedAt pgtype.Timestamp `json:"visited_at"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type AuthUserSigninRes struct {
	User   *AuthUserSignin `json:"user"`
	Tokens *AuthTokens     `json:"tokens"`
}

// reset response
type AuthUserReset struct {
	ID        string           `json:"id"`
	Email     string           `json:"email"`
	Username  string           `json:"username"`
	CheckedAt pgtype.Timestamp `json:"checked_at"`
	VisitedAt pgtype.Timestamp `json:"visited_at"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type AuthUserResetRes struct {
	User   *AuthUserReset `json:"user"`
	Tokens *AuthTokens    `json:"tokens"`
}

// change response
type AuthUserChangeRes struct {
	ID        string           `json:"id"`
	Username  string           `json:"username"`
	CheckedAt pgtype.Timestamp `json:"checked_at"`
	VisitedAt pgtype.Timestamp `json:"visited_at"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

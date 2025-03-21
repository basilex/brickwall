package exchange

type AuthSignupReq struct {
	Email     string `json:"email" binding:"required,email,max=255"`
	Password  string `json:"password" binding:"required,min=4,max=72"`
	Firstname string `json:"firstname" binding:"required,min=1,max=255"`
	Lastname  string `json:"lastname" binding:"required,min=1,max=255"`
}

type AuthSigninReq struct {
	Username string `json:"name" binding:"required,min=1,max=64"`
	Password string `json:"password" binding:"required,min=4,max=72"`
}

type AuthTokenRefreshReq struct {
	Token string `json:"token" binding:"required"`
}

type AuthTokenInvalidateReq struct {
	Token string `json:"token" binding:"required"`
}

type AuthPasswordResetReq struct {
	Password string `json:"password" binding:"required,min=4,max=72"`
}

type AuthPasswordChangeReq struct {
	Password string `json:"password" binding:"required,min=4,max=72"`
}

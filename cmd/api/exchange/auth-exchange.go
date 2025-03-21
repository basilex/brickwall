package exchange

type AuthSignupReq struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Username string `json:"username" binding:"required,min=1,max=64"`
	Password string `json:"password" binding:"required,min=4,max=72"`
}

type AuthSigninReq struct {
	ID       string `json:"id" binding:"required,max=32"`
	Username string `json:"name" binding:"required,min=1,max=64"`
	Password string `json:"password" binding:"required,min=4,max=72"`
}

type AuthSignoutReq struct {
	ID string `uri:"id" binding:"required,max=32,alphanum"`
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

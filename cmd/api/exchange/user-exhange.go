package exchange

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UserUriID struct {
	ID string `uri:"id" binding:"required,max=32,alphanum"`
}

type UserNewReq struct {
	Username string `json:"name" binding:"required,max=64"`
	Password string `json:"password" binding:"required,max=255"`
}

type UserUpdateCredentialsReq struct {
	ID       string `json:"id" binding:"required,max=32"`
	Username string `json:"username" binding:"required,max=64"`
	Password string `json:"password" binding:"required,max=72"`
}

type UserUpdateIsBlockedByIDReq struct {
	ID        string `json:"id" binding:"required,max=32"`
	IsBlocked bool   `json:"is_blocked" binding:"required,boolean"`
}

type UserUpdateIsCheckedByIDReq struct {
	ID        string `json:"id" binding:"required,max=32"`
	IsChecked bool   `json:"is_checked" binding:"required,boolean"`
}

type UserUpdateVisitedAtByIDReq struct {
	ID        string           `json:"id" binding:"required,max=32"`
	VisitedAt pgtype.Timestamp `json:"visited_at" binding:"required" time_format:"2006-01-02T15:04:05.999999999Z07:00"`
}

type UserQuery struct {
	Page  int    `form:"page" binding:"required,min=1,numeric"`
	Size  int    `form:"size" binding:"required,min=5,max=100,numeric"`
	Order string `form:"order" binding:"omitempty,oneof=id name iso2 iso3 num_code"`
}

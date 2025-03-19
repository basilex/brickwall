package exchange

type RoleUriID struct {
	ID string `uri:"id" binding:"required,max=32,alphanum"`
}

type RoleNewReq struct {
	Name string `json:"name" binding:"required,max=255"`
}

type RoleUpdateReq struct {
	ID   string `json:"id" binding:"required,max=32"`
	Name string `json:"name" binding:"required,max=255"`
}

type RoleQuery struct {
	Page  int    `form:"page" binding:"required,min=1,numeric"`
	Size  int    `form:"size" binding:"required,min=5,max=100,numeric"`
	Order string `form:"order" binding:"omitempty,oneof=id name iso2 iso3 num_code"`
}

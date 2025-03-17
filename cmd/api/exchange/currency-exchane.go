package exchange

type CurrencyUriID struct {
	ID string `uri:"id" binding:"required,max=32,alphanum"`
}

type CurrencyNewReq struct {
	Name    string `json:"name" binding:"required,max=255"`
	Code    string `json:"code" binding:"required,len=3,uppercase"`
	NumCode int16  `json:"num_code" binding:"required,min=0,max=1000,numeric"`
	Symbol  string `json:"symbol" binding:"required,min=1,max=8"`
}

type CurrencyUpdateReq struct {
	ID      string `json:"id" binding:"required,max=32"`
	Name    string `json:"name" binding:"required,max=255"`
	Code    string `json:"code" binding:"required,len=3,uppercase"`
	NumCode int16  `json:"num_code" binding:"required,min=0,max=1000,numeric"`
	Symbol  string `json:"symbol" binding:"required,min=1,max=8"`
}

type CurrencyQuery struct {
	Page  int    `form:"page" binding:"required,min=1,numeric"`
	Size  int    `form:"size" binding:"required,min=5,max=100,numeric"`
	Order string `form:"order" binding:"omitempty,oneof=id name code num_code symbol"`
}

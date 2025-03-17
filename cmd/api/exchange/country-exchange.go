package exchange

type CountryUriID struct {
	ID string `uri:"id" binding:"required,max=32,alphanum"`
}

type CountryNewReq struct {
	Name    string `json:"name" binding:"required,max=255"`
	Iso2    string `json:"iso2" binding:"required,len=2,uppercase"`
	Iso3    string `json:"iso3" binding:"required,len=3,uppercase"`
	NumCode int16  `json:"num_code" binding:"required,min=0,max=1000,numeric"`
}

type CountryUpdateReq struct {
	ID      string `json:"id" binding:"required,max=32"`
	Name    string `json:"name" binding:"required,max=255"`
	Iso2    string `json:"iso2" binding:"required,len=2,uppercase"`
	Iso3    string `json:"iso3" binding:"required,len=3,uppercase"`
	NumCode int16  `json:"num_code" binding:"required,min=0,max=1000,numeric"`
}

type CountryQuery struct {
	Page  int    `form:"page" binding:"required,min=1,numeric"`
	Size  int    `form:"size" binding:"required,min=5,max=100,numeric"`
	Order string `form:"order" binding:"omitempty,oneof=id name iso2 iso3 num_code"`
}

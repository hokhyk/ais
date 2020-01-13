package controller

//PR 請購單資料結構
type PR struct {
	prList   prList
	prDetail []prDetail
}

//prList 請購單單頭
type prList struct {
	ID               int    `json:"id"`
	OrganizationName string `json:"organization_name"`
	OrganizationID   int    `json:"organization_id"`
	PayTo            int    `json:"pay_to"`
	VendorName       string `json:"vendor_name"`
	PayType          int    `json:"pay_type"`
	ListType         int    `json:"list_type"`
	UsersID          int    `json:"users_id"`
	UsersName        string `json:"users_name"`
	PayMethod        int    `json:"pay_method"`
	BankAccount      string `json:"bank_account"`
	CreateAt         string `json:"create_at"`
}

//prDetail 請購單單身
type prDetail struct {
	ID           int     `json:"id"`
	Currency     int     `json:"currency"`
	UnitPrice    float32 `json:"unit_float"`
	Quantity     int     `json:"quantity"`
	ExchangeRate float32 `json:"exchange_rate"`
	Tax          float32 `json:"tax"`
	TotalPrice   float32 `json:"total_price"`
}

//New 建構式
func (pr PR) New() *PR {
	return &pr
}

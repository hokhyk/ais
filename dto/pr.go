package dto

import "time"

//PR 請購單資料結構
type PR struct {
	List   PrList     `json:"prList"`
	Detail []PrDetail `json:"prDetail"`
}

//PrList 請購單單頭
type PrList struct {
	ID             int       `json:"id" gorm:"primary_key:yes"`
	OrganizationID int       `json:"organization_id"`
	PayTo          int       `json:"pay_to"`
	VendorName     string    `json:"vendor_name"`
	PayType        int       `json:"pay_type"`
	ListType       int       `json:"list_type"`
	UsersID        int       `json:"users_id"`
	PayMethod      int       `json:"pay_method"`
	BankAccount    string    `json:"bank_account"`
	Proof          string    `json:"proof"`
	Status         int       `json:"status"`
	SignAt         time.Time `json:"sign_at"`
	CreateAt       time.Time `json:"create_at"`
}

//PrDetail 請購單單身
type PrDetail struct {
	ID           int     `json:"id" gorm:"primary_key:yes"`
	PRListID     int     `json:"pr_list_id"`
	Currency     string  `json:"currency"`
	UnitPrice    float32 `json:"unit_price"`
	Quantity     int     `json:"quantity"`
	ExchangeRate float32 `json:"exchange_rate"`
	Tax          float32 `json:"tax"`
	TotalPrice   float32 `json:"total_price"`
}

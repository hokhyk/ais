package dto

import "time"

//PR 請購單資料結構
type PR struct {
	List   PrList     `json:"prList"`
	Detail []PrDetail `json:"prDetail"`
}

//PrList 請購單單頭
type PrList struct {
	ID              int       `json:"id" gorm:"primary_key:yes"`
	Serial          string    `json:"serial"`
	PayTo           int       `json:"pay_to"`
	Company         int       `json:"company"`
	PayType         int       `json:"pay_type"`
	ListType        int       `json:"list_type"`
	UsersID         int       `json:"users_id"`
	PayMethod       int       `json:"pay_method"`
	BankAccount     string    `json:"bank_account"`
	Proof           string    `json:"proof"`
	InstallmentPlan int       `json:"installment_plan"`
	PayBy           int       `json:"pay_by"`
	Memo            string    `json:"memo"`
	PrItem          int       `json:"pr_item"`
	Status          int       `json:"status"`
	SignAt          time.Time `json:"sign_at"`
	PayDate         time.Time `json:"pay_date"`
	CreateAt        time.Time `json:"create_at"`
}

//PrDetail 請購單單身
type PrDetail struct {
	ID           int     `json:"id" gorm:"primary_key:yes"`
	Name         string  `json:"name"`
	PRListID     int     `json:"pr_list_id"`
	Currency     string  `json:"currency"`
	UnitPrice    float32 `json:"unit_price"`
	Quantity     int     `json:"quantity"`
	ExchangeRate float32 `json:"exchange_rate"`
	Tax          float32 `json:"tax"`
	TotalPrice   float32 `json:"total_price"`
}

//GetList 取得請購單列表
type GetList struct {
	Serial       string  `json:"serial"`
	Name         string  `json:"name"`
	PRListID     int     `json:"pr_list_id"`
	Currency     string  `json:"currency"`
	UnitPrice    float32 `json:"unit_price"`
	Quantity     int     `json:"quantity"`
	ExchangeRate float32 `json:"exchange_rate"`
	Tax          float32 `json:"tax"`
	TotalPrice   float32 `json:"total_price"`
	Proof        string  `json:"proof"`
}

//PrSearch 搜尋條件
type PrSearch struct {
	ID             int       `json:"id"`
	Serial         string    `json:"serial"`
	OrganizationID int       `json:"organization_id"`
	PrItem         int       `json:"pr_item"`
	Begin          time.Time `json:"begin"`
	End            time.Time `json:"end"`
	Num            int       `json:"num"`
	Page           int       `json:"page"`
}

//PrListResult 請款單列表取得結果
type PrListResult struct {
	ID              int       `json:"id"`
	Serial          string    `json:"serial"`
	PayTo           int       `json:"pay_to"`
	Company         int       `json:"company"`
	PayType         int       `json:"pay_type"`
	ListType        int       `json:"list_type"`
	UsersID         int       `json:"users_id"`
	Email           string    `json:"email"`
	Identifier      string    `json:"identifier"`
	Lastname        string    `json:"lastname"`
	Firstname       string    `json:"firstname"`
	PayMethod       int       `json:"pay_method"`
	BankAccount     string    `json:"bank_account"`
	Proof           string    `json:"proof"`
	InstallmentPlan int       `json:"installment_plan"`
	PayBy           int       `json:"pay_by"`
	Memo            string    `json:"memo"`
	PrItem          int       `json:"pr_item"`
	Status          int       `json:"status"`
	SignAt          time.Time `json:"sign_at"`
	PayDate         time.Time `json:"pay_date"`
	CreateAt        time.Time `json:"create_at"`
}

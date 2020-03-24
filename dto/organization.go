package dto

//Organization 單位資料表結構
type Organization struct {
	ID  int    `json:"id"`
	Key string `json:"key" gorm:"column:name"`
}

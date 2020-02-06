package dto

//ResultObject 使用者相關資料結構
type ResultObject struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

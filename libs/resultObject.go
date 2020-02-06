package libs

import (
	"encoding/json"

	"github.com/teed7334-restore/ais/dto"
)

//ResultObject 物件參數
type ResultObject struct{}

//New 建構式
func (ro ResultObject) New() *ResultObject {
	return &ro
}

//BuildJSON 設定RO物件並回傳JSON
func (ro *ResultObject) BuildJSON(status int, message string) string {
	dtoRO := &dto.ResultObject{}
	dtoRO.Status = status
	dtoRO.Message = message
	json, _ := json.Marshal(dtoRO)
	content := string(json)
	return content
}

//Build 製作成RO物件
func (ro *ResultObject) Build(status int, message string) *dto.ResultObject {
	dtoRO := &dto.ResultObject{}
	dtoRO.Status = status
	dtoRO.Message = message
	return dtoRO
}

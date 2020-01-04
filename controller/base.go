package controller

import (
	"encoding/json"
)

//resultObject 回傳物件
type resultObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//buildRO 設定RO物件並回傳JSON
func buildRO(code int, message string) string {
	ro := &resultObject{}
	ro.Code = code
	ro.Message = message
	json, _ := json.Marshal(ro)
	content := string(json)
	return content
}

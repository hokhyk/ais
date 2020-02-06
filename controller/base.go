package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
)

//RO 資料回傳物件
var RO = libs.ResultObject{}.New()

//TimeFormat 時間格式
var TimeFormat = "2006/01/02 15:04:05"

//PrintRO 輸出RO物件到網頁
func PrintRO(w http.ResponseWriter, dtoRO *dto.ResultObject, message string) {
	if dtoRO.Status != 1 {
		ro, _ := json.Marshal(dtoRO)
		content := string(ro)
		fmt.Fprintf(w, content)
		return
	}
	content := RO.BuildJSON(1, message)
	fmt.Fprintf(w, content)
	return
}

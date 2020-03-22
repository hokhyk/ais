package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/teed7334-restore/ais/service"
)

//CreditType 入帳類別資料結構
type CreditType struct{}

var sCreditType = service.CreditType{}.New()

//New 建構式
func (ct CreditType) New() *CreditType {
	return &ct
}

//GetCreditType 取得入帳類別
func (ct *CreditType) GetCreditType(w http.ResponseWriter, r *http.Request) {
	dtoRO, dtoCreditType := sCreditType.GetCreditType()
	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	jsonByte, _ := json.Marshal(dtoCreditType)
	result := string(jsonByte)
	fmt.Fprintf(w, result)
}

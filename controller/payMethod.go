package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/teed7334-restore/ais/service"
)

//PayMethod 支付方式資料結構
type PayMethod struct{}

var sPayMethod = service.PayMethod{}.New()

//New 建構式
func (pm PayMethod) New() *PayMethod {
	return &pm
}

//GetPayMethod 取得支付方式
func (pm *PayMethod) GetPayMethod(w http.ResponseWriter, r *http.Request) {
	dtoRO, dtoPayMethod := sPayMethod.GetPayMethod()
	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	jsonByte, _ := json.Marshal(dtoPayMethod)
	result := string(jsonByte)
	fmt.Fprintf(w, result)
}

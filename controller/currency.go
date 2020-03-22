package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/teed7334-restore/ais/service"
)

//Currency 貨幣檔案資料結構
type Currency struct{}

var sCurrency = service.Currency{}.New()

//New 建構式
func (c Currency) New() *Currency {
	return &c
}

//GetCurrency 取得貨幣資料
func (c *Currency) GetCurrency(w http.ResponseWriter, r *http.Request) {
	dtoRO, dtoCurrency := sCurrency.GetCurrency()
	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	jsonByte, _ := json.Marshal(dtoCurrency)
	result := string(jsonByte)
	fmt.Fprintf(w, result)
}

package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/teed7334-restore/ais/service"
)

//PrItem 應收項目資料結構
type PrItem struct{}

var sPrItem = service.PrItem{}.New()

//New 建構式
func (pi PrItem) New() *PrItem {
	return &pi
}

//GetPrItem 取得應付項目
func (pi *PrItem) GetPrItem(w http.ResponseWriter, r *http.Request) {
	dtoRO, dtoPrItem := sPrItem.GetPrItem()
	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	jsonByte, _ := json.Marshal(dtoPrItem)
	result := string(jsonByte)
	fmt.Fprintf(w, result)
}

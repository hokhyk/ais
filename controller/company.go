package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/teed7334-restore/ais/service"
)

//Company 廠商資料結構
type Company struct{}

var sCompany = service.Company{}.New()

//New 建構式
func (c Company) New() *Company {
	return &c
}

//GetCompany 取得廠商列表
func (c *Company) GetCompany(w http.ResponseWriter, r *http.Request) {
	dtoRO, dtoCompany := sCompany.GetCompany()
	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	jsonByte, _ := json.Marshal(dtoCompany)
	result := string(jsonByte)
	fmt.Fprintf(w, result)
}

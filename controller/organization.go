package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/teed7334-restore/ais/service"
)

//Organization 單位資料結構
type Organization struct{}

var sOrganization = service.Organization{}.New()

//New 建構式
func (o Organization) New() *Organization {
	return &o
}

//GetOrganization 取得部門項目
func (o *Organization) GetOrganization(w http.ResponseWriter, r *http.Request) {
	dtoRO, dtoOrganization := sOrganization.GetOrganization()
	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	jsonByte, _ := json.Marshal(dtoOrganization)
	result := string(jsonByte)
	fmt.Fprintf(w, result)
}

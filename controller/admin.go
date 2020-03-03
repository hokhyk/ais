package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
	"github.com/teed7334-restore/ais/service"
)

//Admin 管理員頁面資料結構
type Admin struct{}

//adminPrItemResult 用來存放請購單資訊
type adminPrItemResult struct {
	Status int               `json:"status"`
	List   *dto.PrListResult `json:"list"`
	Detail *[]dto.PrDetail   `json:"detail"`
}

var sa = service.Admin{}.New()

//New 建構式
func (a Admin) New() *Admin {
	return &a
}

//GetItem 取得請購單資訊
func (a *Admin) GetItem(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	_, dtoUsers := users.GetUser(token)

	if !a.isAdmin(dtoUsers) {
		content := RO.BuildJSON(0, "此會員並非管理員")
		fmt.Fprintf(w, content)
		return
	}

	val := r.FormValue("id")
	if val == "" {
		content := RO.BuildJSON(0, "請購單號為空白")
		fmt.Fprintf(w, content)
		return
	}

	id, err := strconv.Atoi(val)
	if err != nil {
		content := RO.BuildJSON(0, "請購單號需為數字")
		fmt.Fprintf(w, content)
		return
	}

	dtoPrSearch := &dto.PrSearch{}
	dtoPrSearch.ID = id
	dtoRO, dtoPrListResult, dtoPrDetail := sa.GetItem(dtoPrSearch)

	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	result := &adminPrItemResult{}
	result.Status = 1
	result.List = dtoPrListResult
	result.Detail = dtoPrDetail

	jsonByte, _ := json.Marshal(result)
	content := string(jsonByte)
	fmt.Fprintf(w, content)
}

//GetList 取得請購單列表
func (a *Admin) GetList(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	_, dtoUsers := users.GetUser(token)

	if !a.isAdmin(dtoUsers) {
		content := RO.BuildJSON(0, "此會員並非管理員")
		fmt.Fprintf(w, content)
		return
	}

	dtoPrSearch := &dto.PrSearch{}

	if r.FormValue("begin") != "" {
		t := r.FormValue("begin") + " 00:00:00"
		dtoPrSearch.Begin, _ = time.ParseInLocation(TimeFormat, t, time.Local)
	}
	if r.FormValue("end") != "" {
		t := r.FormValue("end") + " 23:59:59"
		dtoPrSearch.End, _ = time.ParseInLocation(TimeFormat, t, time.Local)
	}
	dtoPrSearch.Num = 10
	if r.FormValue("num") != "" {
		dtoPrSearch.Num, _ = strconv.Atoi(r.FormValue("num"))
	}
	dtoPrSearch.Page = 1
	if r.FormValue("page") != "" {
		dtoPrSearch.Page, _ = strconv.Atoi(r.FormValue("page"))
	}

	dtoRO, dtoGetList := sa.GetList(dtoPrSearch)

	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	result := &getListResult{}
	result.Status = 1
	result.List = dtoGetList

	jsonByte, _ := json.Marshal(result)
	content := string(jsonByte)
	fmt.Fprintf(w, content)
}

//isAdmin 是否為管理員
func (a *Admin) isAdmin(user *dto.Users) bool {
	helper := libs.Helper{}.New()
	admins := []string{"1", "12", "50"}
	id := strconv.Itoa(user.ID)
	return helper.InArray(admins, id)
}

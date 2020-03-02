package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/service"
)

//PR 請購單資料結構
type PR struct{}

//prListResult 用來存放列表頁資訊
type prListResult struct {
	Status int             `json:"status"`
	List   *[]dto.PrDetail `json:"list"`
}

//prItemResult 用來存放請購單資訊
type prItemResult struct {
	Status int               `json:"status"`
	List   *dto.PrListResult `json:"list"`
	Detail *[]dto.PrDetail   `json:"detail"`
}

var spr = service.PR{}.New()

//New 建構式
func (pr PR) New() *PR {
	return &pr
}

//GetItem 取得請購單資訊
func (pr *PR) GetItem(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	_, dtoUsers := users.GetUser(token)

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

	dtoSearch := &dto.PrSearch{}
	dtoSearch.ID = id
	dtoRO, dtoPrList, dtoPrDetail := spr.GetItem(dtoSearch, dtoUsers)

	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	result := &prItemResult{}
	result.Status = 1
	result.List = dtoPrList
	result.Detail = dtoPrDetail

	jsonByte, _ := json.Marshal(result)
	content := string(jsonByte)
	fmt.Fprintf(w, content)
}

//GetList 取得請購單列表
func (pr *PR) GetList(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	_, dtoUsers := users.GetUser(token)
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

	dtoRO, dtoPrDetails := spr.GetList(dtoPrSearch, dtoUsers)

	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	result := &prListResult{}
	result.Status = 1
	result.List = dtoPrDetails

	jsonByte, _ := json.Marshal(result)
	content := string(jsonByte)
	fmt.Fprintf(w, content)
}

//SetCancel 作廢請購單
func (pr *PR) SetCancel(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	dtoRO, dtoUsers := users.GetUser(token)

	id := r.FormValue("id")
	if id == "" {
		content := RO.BuildJSON(0, "請購單號為空白")
		fmt.Fprintf(w, content)
		return
	}

	dtoRO = spr.SetCancel(dtoUsers, id)

	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	PrintRO(w, dtoRO, "true")
}

//Add 新增請購單
func (pr *PR) Add(w http.ResponseWriter, r *http.Request) {
	dtoPR := new(dto.PR)
	dtoRO, dtoPR := pr.filterList(r, dtoPR)
	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	dtoRO, dtoPR = pr.filterDetail(r, dtoPR)
	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	r.ParseMultipartForm(5000000)
	proof := r.MultipartForm
	dtoRO = spr.Add(dtoPR, proof)
	PrintRO(w, dtoRO, "true")
}

func (pr *PR) filterDetail(r *http.Request, dtoPR *dto.PR) (*dto.ResultObject, *dto.PR) {
	name, ok := r.MultipartForm.Value["name[]"]
	if !ok {
		dtoRO := RO.Build(0, "項目不得為空")
		return dtoRO, dtoPR
	}
	currency, ok := r.MultipartForm.Value["currency[]"]
	if !ok {
		dtoRO := RO.Build(0, "幣值不得為空")
		return dtoRO, dtoPR
	}
	unitPrice, ok := r.MultipartForm.Value["unit_price[]"]
	if !ok {
		dtoRO := RO.Build(0, "單價不得為空")
		return dtoRO, dtoPR
	}
	quantity, ok := r.MultipartForm.Value["quantity[]"]
	if !ok {
		dtoRO := RO.Build(0, "數量不得為空")
		return dtoRO, dtoPR
	}
	exchangeRate, ok := r.MultipartForm.Value["exchange_rate[]"]
	if !ok {
		dtoRO := RO.Build(0, "匯率不得為空")
		return dtoRO, dtoPR
	}
	tax, ok := r.MultipartForm.Value["tax[]"]
	if !ok {
		dtoRO := RO.Build(0, "稅額不得為空")
		return dtoRO, dtoPR
	}
	for k := range currency {
		dtoDetail := new(dto.PrDetail)
		if len(name) < k {
			dtoRO := RO.Build(0, "項目不得為空")
			return dtoRO, dtoPR
		}
		dtoDetail.Name = name[k]
		if len(currency) < k {
			dtoRO := RO.Build(0, "幣值不得為空")
			return dtoRO, dtoPR
		}
		dtoDetail.Currency = currency[k]
		if len(unitPrice) < k {
			dtoRO := RO.Build(0, "單價不得為空")
			return dtoRO, dtoPR
		}
		up, err := strconv.ParseFloat(unitPrice[k], 32)
		if err != nil {
			dtoRO := RO.Build(0, "單價不得為空")
			return dtoRO, dtoPR
		}
		dtoDetail.UnitPrice = float32(up)
		qty, err := strconv.Atoi(quantity[k])
		if err != nil {
			dtoRO := RO.Build(0, "數量不得為空")
			return dtoRO, dtoPR
		}
		dtoDetail.Quantity = qty
		er, err := strconv.ParseFloat(exchangeRate[k], 32)
		if err != nil {
			dtoRO := RO.Build(0, "匯率不得為空")
			return dtoRO, dtoPR
		}
		dtoDetail.ExchangeRate = float32(er)
		tx, err := strconv.ParseFloat(tax[k], 32)
		if err != nil {
			dtoRO := RO.Build(0, "稅額不得為空")
			return dtoRO, dtoPR
		}
		dtoDetail.Tax = float32(tx)
		dtoDetail.TotalPrice = dtoDetail.UnitPrice * float32(dtoDetail.Quantity) * dtoDetail.ExchangeRate * (1.0 + dtoDetail.Tax)
		dtoPR.Detail = append(dtoPR.Detail, *dtoDetail)
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoPR
}

//filterList 過濾單頭必要參數為空值之參數
func (pr *PR) filterList(r *http.Request, dtoPR *dto.PR) (*dto.ResultObject, *dto.PR) {
	dtoUsers := new(dto.Users)
	r.ParseMultipartForm(128)

	token, ok := r.MultipartForm.Value["token"]
	dtoRO, dtoUsers := users.GetUser(token[0])

	val, ok := r.MultipartForm.Value["sign_at"]
	if !ok {
		dtoRO := RO.Build(0, "請輸入日期")
		return dtoRO, dtoPR
	}
	signAt := fmt.Sprintf("%s 00:00:00", val[0])

	val, ok = r.MultipartForm.Value["pay_date"]
	payDate := "0001/01/01 00:00:00"
	if ok {
		payDate = fmt.Sprintf("%s 00:00:00", val[0])
	}

	now := time.Now()

	dtoPR.List.OrganizationID = dtoUsers.OrganizationID
	dtoPR.List.UsersID = dtoUsers.ID
	dtoPR.List.SignAt, _ = time.ParseInLocation(TimeFormat, signAt, time.Local)
	dtoPR.List.PayDate, _ = time.ParseInLocation(TimeFormat, payDate, time.Local)

	dtoPR.List.CreateAt, _ = time.ParseInLocation(TimeFormat, now.Format(TimeFormat), time.Local)
	dtoPR.List.Status = 1

	val, ok = r.MultipartForm.Value["pay_to"]
	if !ok {
		dtoRO := RO.Build(0, "請勾選支付對象")
		return dtoRO, dtoPR
	}
	if val[0] == "" {
		dtoRO := RO.Build(0, "請勾選支付對象")
		return dtoRO, dtoPR
	}
	payTo, err := strconv.Atoi(val[0])
	if err != nil {
		dtoRO := RO.Build(0, "請勾選支付對象")
		return dtoRO, dtoPR
	}

	vendorName := ""
	val, ok = r.MultipartForm.Value["vendor_name"]
	if ok && val[0] == "" && payTo == 1 {
		dtoRO := RO.Build(0, "請填寫支付廠商")
		return dtoRO, dtoPR
	}
	vendorName = val[0]

	val, ok = r.MultipartForm.Value["pay_type"]
	if !ok || val[0] == "" {
		dtoRO := RO.Build(0, "請選擇入賬類別")
		return dtoRO, dtoPR
	}
	payType, err := strconv.Atoi(val[0])
	if err != nil {
		dtoRO := RO.Build(0, "請選擇入賬類別")
		return dtoRO, dtoPR
	}

	val, ok = r.MultipartForm.Value["list_type"]
	if !ok || val[0] == "" {
		dtoRO := RO.Build(0, "請選擇類別")
		return dtoRO, dtoPR
	}
	listType, err := strconv.Atoi(val[0])
	if err != nil {
		dtoRO := RO.Build(0, "請選擇類別")
		return dtoRO, dtoPR
	}

	val, ok = r.MultipartForm.Value["pay_method"]
	if !ok || val[0] == "" {
		dtoRO := RO.Build(0, "請選擇支付方式")
		return dtoRO, dtoPR
	}
	payMethod, err := strconv.Atoi(val[0])
	if err != nil {
		dtoRO := RO.Build(0, "請選擇支付方式")
		return dtoRO, dtoPR
	}

	val, ok = r.MultipartForm.Value["bank_account"]
	bankAccount := ""
	if ok && val[0] == "" && payMethod == 3 {
		dtoRO := RO.Build(0, "請輸入銀行帳號")
		return dtoRO, dtoPR
	}
	bankAccount = val[0]

	dtoPR.List.PayTo = payTo
	dtoPR.List.VendorName = vendorName
	dtoPR.List.PayType = payType
	dtoPR.List.ListType = listType
	dtoPR.List.PayMethod = payMethod
	dtoPR.List.BankAccount = bankAccount

	dtoRO = RO.Build(1, "")
	return dtoRO, dtoPR
}

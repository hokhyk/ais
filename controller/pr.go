package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/service"
)

//PR 請購單資料結構
type PR struct{}

var spr = service.PR{}.New()

//New 建構式
func (pr PR) New() *PR {
	return &pr
}

//Add 新增請購單
func (pr *PR) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		content := buildRO(201, "不支持此種HTTP Method")
		fmt.Fprintf(w, content)
		return
	}

	dtoPR := new(dto.PR)
	status, dtoPR := pr.filterList(r, dtoPR)
	switch status {
	case -1:
		content := buildRO(202, "使用者Token有誤")
		fmt.Fprintf(w, content)
	case -2:
		content := buildRO(203, "請選擇支付對象")
		fmt.Fprintf(w, content)
	case -3:
		content := buildRO(204, "請輸入支付廠商")
		fmt.Fprintf(w, content)
	case -4:
		content := buildRO(205, "請選擇入賬類別")
		fmt.Fprintf(w, content)
	case -5:
		content := buildRO(206, "請選擇類別")
		fmt.Fprintf(w, content)
	case -6:
		content := buildRO(207, "請選擇支付方式")
		fmt.Fprintf(w, content)
	case -7:
		content := buildRO(208, "請輸入銀行帳號")
		fmt.Fprintf(w, content)
	}
	if status != 1 {
		return
	}

	r.ParseMultipartForm(5000000)
	proof := r.MultipartForm
	status = spr.Add(dtoPR, proof)
	switch status {
	case -1:
		content := buildRO(209, "未上傳佐証資料")
		fmt.Fprintf(w, content)
	case -2:
		content := buildRO(210, "無法新增檔案")
		fmt.Fprintf(w, content)
	case -3:
		content := buildRO(211, "上傳資料夾權限不符")
		fmt.Fprintf(w, content)
	case -4:
		content := buildRO(212, "無法新增Zip壓縮檔")
		fmt.Fprintf(w, content)
	case -5:
		content := buildRO(213, "伺服器中無對應之Zip壓縮檔")
		fmt.Fprintf(w, content)
	case -6:
		content := buildRO(214, "無法解析Zip壓縮檔")
		fmt.Fprintf(w, content)
	case -7:
		content := buildRO(215, "無法取得Zip壓縮檔檔頭資訊")
		fmt.Fprintf(w, content)
	case -8:
		content := buildRO(216, "無法建立Zip壓縮檔檔頭")
		fmt.Fprintf(w, content)
	case -9:
		content := buildRO(217, "無法寫入檔到到Zip壓縮檔")
		fmt.Fprintf(w, content)
	}
	if status != 1 {
		return
	}

	content := buildRO(200, "true")
	fmt.Fprintf(w, content)
}

func (pr *PR) filterDetail(r *http.Request, dtoPR *dto.PR) (int, *dto.PR) {
	currency, ok := r.MultipartForm.Value["currency"]
	if !ok {
		return -1, dtoPR
	}
	unitPrice, ok := r.MultipartForm.Value["unit_price"]
	if !ok {
		return -2, dtoPR
	}
	quantity, ok := r.MultipartForm.Value["quantity"]
	if !ok {
		return -3, dtoPR
	}
	exchangeRate, ok := r.MultipartForm.Value["exchange_rate"]
	if !ok {
		return -4, dtoPR
	}
	tax, ok := r.MultipartForm.Value["tax"]
	if !ok {
		return -5, dtoPR
	}
	for k := range currency {
		dtoDetail := new(dto.PrDetail)
		if len(currency) < k {
			return -6, dtoPR
		}
		dtoDetail.Currency = currency[k]
		if len(unitPrice) < k {
			return -7, dtoPR
		}
		up, err := strconv.ParseFloat(unitPrice[k], 32)
		if err != nil {
			return -7, dtoPR
		}
		dtoDetail.UnitPrice = float32(up)
		qty, err := strconv.Atoi(quantity[k])
		if err != nil {
			return -8, dtoPR
		}
		dtoDetail.Quantity = qty
		er, err := strconv.ParseFloat(exchangeRate[k], 32)
		if err != nil {
			return -9, dtoPR
		}
		dtoDetail.ExchangeRate = float32(er)
		tx, err := strconv.ParseFloat(tax[k], 32)
		if err != nil {
			return -10, dtoPR
		}
		dtoDetail.Tax = float32(tx)
	}
	return 1, dtoPR
}

//filterList 過濾單頭必要參數為空值之參數
func (pr *PR) filterList(r *http.Request, dtoPR *dto.PR) (int, *dto.PR) {
	now := time.Now()
	dtoUsers := new(dto.Users)
	r.ParseMultipartForm(128)

	token, ok := r.MultipartForm.Value["token"]
	if !ok {
		return -1, dtoPR
	}
	status, dtoUsers := users.GetUser(token[0])
	if status != 1 {
		return -1, dtoPR
	}

	dtoPR.List.OrganizationName = dtoUsers.OrganizationName
	dtoPR.List.OrganizationID = dtoUsers.OrganizationID
	dtoPR.List.UsersID = dtoUsers.ID
	dtoPR.List.UsersName = dtoUsers.FirstName + dtoUsers.LastName
	dtoPR.List.CreateAt = now.Format("2006-01-02 15:04:05")
	dtoPR.List.Status = 1

	val, ok := r.MultipartForm.Value["pay_to"]
	if !ok {
		return -2, dtoPR
	}
	if val[0] == "" {
		return -2, dtoPR
	}
	payTo, err := strconv.Atoi(val[0])
	if err != nil {
		return -2, dtoPR
	}

	val, ok = r.MultipartForm.Value["vendor_name"]
	if !ok {
		return -3, dtoPR
	}
	if val[0] == "" && payTo == 1 {
		return -3, dtoPR
	}
	vendorName := val[0]

	val, ok = r.MultipartForm.Value["pay_type"]
	if !ok {
		return -4, dtoPR
	}
	if val[0] == "" {
		return -4, dtoPR
	}
	payType, err := strconv.Atoi(val[0])
	if err != nil {
		return -4, dtoPR
	}

	val, ok = r.MultipartForm.Value["list_type"]
	if !ok {
		return -5, dtoPR
	}
	if val[0] == "" {
		return -5, dtoPR
	}
	listType, err := strconv.Atoi(val[0])
	if err != nil {
		return -5, dtoPR
	}

	val, ok = r.MultipartForm.Value["pay_method"]
	if !ok {
		return -6, dtoPR
	}
	if val[0] == "" {
		return -6, dtoPR
	}
	payMethod, err := strconv.Atoi(val[0])
	if err != nil {
		return -6, dtoPR
	}

	val, ok = r.MultipartForm.Value["bank_account"]
	if !ok {
		return -7, dtoPR
	}
	if val[0] == "" && payMethod == 3 {
		return -7, dtoPR
	}
	bankAccount := val[0]

	dtoPR.List.PayTo = payTo
	dtoPR.List.VendorName = vendorName
	dtoPR.List.PayType = payType
	dtoPR.List.ListType = listType
	dtoPR.List.PayMethod = payMethod
	dtoPR.List.BankAccount = bankAccount

	return 1, dtoPR
}

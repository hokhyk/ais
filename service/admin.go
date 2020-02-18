package service

import (
	"fmt"
	"strings"

	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
)

//Admin 管理員頁面資料結構
type Admin struct{}

//New 建構式
func (a Admin) New() *Admin {
	return &a
}

//GetItem 取得請購單資料
func (a *Admin) GetItem(search *dto.PrSearch) (*dto.ResultObject, *dto.PrListResult, *[]dto.PrDetail) {
	mysql := libs.MySQL{}.New()
	search.Num = 1
	search.Page = 1
	dtoPrListResults := a.getListFromDB(search, mysql)
	dtoPrDetail := &[]dto.PrDetail{}
	if len(*dtoPrListResults) == 0 {
		dtoRO := RO.Build(0, "查無任何資料")
		return dtoRO, &dto.PrListResult{}, dtoPrDetail
	}
	dtoPrListResults = a.setProofURL(dtoPrListResults)
	dtoPrList := &(*dtoPrListResults)[0]
	dtoPrDetail = a.getDetailFromDB(dtoPrList, mysql)
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoPrList, dtoPrDetail
}

//GetList 取得請購單列表
func (a *Admin) GetList(search *dto.PrSearch) (*dto.ResultObject, *[]dto.PrListResult) {
	mysql := libs.MySQL{}.New()
	dtoPrListResults := a.getListFromDB(search, mysql)
	if len(*dtoPrListResults) == 0 {
		dtoRO := RO.Build(0, "查無任何資料")
		return dtoRO, dtoPrListResults
	}
	dtoPrListResults = a.setProofURL(dtoPrListResults)
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoPrListResults
}

//setProofURL 將佐証資資路徑轉成相對路徑
func (a *Admin) setProofURL(results *[]dto.PrListResult) *[]dto.PrListResult {
	for k, v := range *results {
		arr := strings.Split(v.Proof, "/")
		fileName := arr[len(arr)-1]
		(*results)[k].Proof = "/download/getFile?proof=" + fileName
	}
	return results
}

//getDetailFromDB 從資料庫取得請購單單身
func (a *Admin) getDetailFromDB(list *dto.PrListResult, m MySQL) *[]dto.PrDetail {
	db := m.GetAdater()
	dtoPrDetail := &[]dto.PrDetail{}
	db.Where("pr_list_id = ?", list.ID).Order("id ASC").Find(dtoPrDetail)
	return dtoPrDetail
}

//getListFromDB 從資料庫取得請購單列表
func (a *Admin) getListFromDB(search *dto.PrSearch, m MySQL) *[]dto.PrListResult {
	db := m.GetAdater()
	sql := `
		SELECT 
			pl.id, 
			pl.organization_id, 
			o.name AS organization_name, 
			pl.pay_to, 
			pl.vendor_name, 
			pl.pay_type, 
			pl.list_type, 
			pl.users_id, 
			u.email, 
			u.identifier, 
			u.lastname, 
			u.firstname, 
			pl.pay_method, 
			pl.bank_account, 
			pl.proof, 
			pl.status, 
			pl.sign_at, 
			pl.create_at
		FROM 
			pr_lists pl
		INNER JOIN 
			users u ON pl.users_id = u.id
		INNER JOIN
			organization o ON pl.organization_id = o.id
		WHERE
			pl.status = 1 %s
		ORDER BY
			pl.sign_at DESC
		LIMIT %d, %d
	`
	where := ""
	if !search.Begin.IsZero() && !search.End.IsZero() {
		where = where + " %s"
		begin := search.Begin.Format(TimeFormat)
		end := search.End.Format(TimeFormat)
		where = fmt.Sprintf(" AND pl.sign_at >= '%s' AND pl.sign_at <= '%s'", begin, end)
	}
	if search.ID != 0 {
		where = where + " %s"
		where = fmt.Sprintf(" AND pl.id = %d", search.ID)
	}
	offset := (search.Page - 1) * search.Num
	sql = fmt.Sprintf(sql, where, offset, search.Num)
	PrListResults := &[]dto.PrListResult{}
	db.Raw(sql).Scan(PrListResults)
	return PrListResults
}

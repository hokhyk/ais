package service

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
)

//PR 請購單資料結構
type PR struct{}

const (
	uploadPath = "./resources/uploads"
	proofPath  = "./resources/proof"
)

// New 建構式
func (pr PR) New() *PR {
	return &pr
}

//GetItem 取得請購單資料
func (pr *PR) GetItem(search *dto.PrSearch, user *dto.Users) (*dto.ResultObject, *dto.PrListResult, *[]dto.PrDetail) {
	mysql := libs.MySQL{}.New()
	search.Page = 1
	search.Num = 1
	dtoPrListResults := pr.getHeaderFromDB(search, user, mysql)
	dtoPrDetail := &[]dto.PrDetail{}
	if len(*dtoPrListResults) == 0 {
		dtoRO := RO.Build(0, "查無任何資料")
		return dtoRO, &dto.PrListResult{}, dtoPrDetail
	}
	dtoPrListResults = pr.setProofURL(dtoPrListResults)
	dtoPrListResult := &(*dtoPrListResults)[0]
	dtoPrDetail = pr.getDetailFromDB(dtoPrListResult, mysql)
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoPrListResult, dtoPrDetail
}

//GetList 取得請購單列表
func (pr *PR) GetList(search *dto.PrSearch, user *dto.Users) (*dto.ResultObject, *[]dto.GetList) {
	mysql := libs.MySQL{}.New()
	dtoGetList := pr.getListFromDB(search, user, mysql)
	if len(*dtoGetList) == 0 {
		dtoRO := RO.Build(0, "查無任何資料")
		return dtoRO, dtoGetList
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoGetList
}

//SetCancel 作廢請購單
func (pr *PR) SetCancel(u *dto.Users, id string) *dto.ResultObject {
	mysql := libs.MySQL{}.New()
	return pr.doSetCancelToDB(u, id, mysql)
}

//Add 新增請購單
func (pr *PR) Add(params *dto.PR, files *multipart.Form) *dto.ResultObject {
	dtoRO, uploads := pr.doUpload(files)
	if dtoRO.Status != 1 {
		return dtoRO
	}
	ziper := libs.Zip{}.New()
	dtoRO, proof := pr.doCompress(uploads, ziper)
	if dtoRO.Status != 1 {
		return dtoRO
	}
	params.List.Proof = proof
	pr.doRemoveTempFiles(uploads)
	mysql := libs.MySQL{}.New()
	if params.List.Serial == "" {
		params.List.Serial = pr.getSerial(params)
	}
	dtoRO = pr.doInsertToDB(params, mysql)
	return dtoRO
}

//doUpload 處理上傳檔案
func (pr *PR) doUpload(files *multipart.Form) (*dto.ResultObject, []string) {
	var dst []string
	for _, fs := range files.File {
		for i := range fs {
			file, err := fs[i].Open()
			defer file.Close()
			if err != nil {
				log.Println(err)
				dtoRO := RO.Build(0, "未上傳佐証資料")
				return dtoRO, dst
			}
			ext := strings.Split(fs[i].Filename, ".")[1]
			now := time.Now()
			filename := fmt.Sprintf("%d%d.%s", now.UnixNano(), i, ext)
			fp := filepath.Join(uploadPath, filename)
			filename = fmt.Sprintf("%s/%s", uploadPath, filename)
			dst = append(dst, filename)
			out, err := os.Create(fp)
			defer out.Close()
			if err != nil {
				log.Println(err)
				dtoRO := RO.Build(0, "無法新增檔案")
				return dtoRO, dst
			}
			_, err = io.Copy(out, file)
			if err != nil {
				log.Println(err)
				dtoRO := RO.Build(0, "上傳資料夾權限不符")
				return dtoRO, dst
			}
		}
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, dst
}

//doCompress 壓縮檔案
func (pr *PR) doCompress(files []string, ziper Zip) (*dto.ResultObject, string) {
	now := time.Now()
	filename := fmt.Sprintf("%d.zip", now.UnixNano())
	dst := fmt.Sprintf("%s/%s", proofPath, filename)
	status := ziper.Compress(files, dst)
	switch status {
	case -1:
		dtoRO := RO.Build(0, "無法新增Zip壓縮檔")
		return dtoRO, ""
	case -2:
		dtoRO := RO.Build(0, "伺服器中無對應之Zip壓縮檔")
		return dtoRO, ""
	case -3:
		dtoRO := RO.Build(0, "無法解析Zip壓縮檔")
		return dtoRO, ""
	case -4:
		dtoRO := RO.Build(0, "無法取得Zip壓縮檔檔頭資訊")
		return dtoRO, ""
	case -5:
		dtoRO := RO.Build(0, "無法建立Zip壓縮檔檔頭")
		return dtoRO, ""
	case -6:
		dtoRO := RO.Build(0, "無法寫入檔到到Zip壓縮檔")
		return dtoRO, ""
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, dst
}

//doRemoveTempFiles 刪除暫存檔案
func (pr *PR) doRemoveTempFiles(files []string) {
	for _, file := range files {
		os.Remove(file)
	}
}

//setProofURL 將佐証資資路徑轉成相對路徑
func (pr *PR) setProofURL(results *[]dto.PrListResult) *[]dto.PrListResult {
	for k, v := range *results {
		arr := strings.Split(v.Proof, "/")
		fileName := arr[len(arr)-1]
		(*results)[k].Proof = "/download/getFile?proof=" + fileName
	}
	return results
}

//getSerial 取得單號
func (pr *PR) getSerial(params *dto.PR) string {
	now := time.Now().Format("20060102")
	userID := strconv.Itoa(params.List.UsersID)
	num := len(userID)
	for i := num; i <= 3; i++ {
		userID = "0" + userID
	}
	times := strings.Split(time.Now().Format("15:04:05"), ":")
	serial := "P" + now + userID + times[0] + times[1] + times[2]
	return serial
}

//getHeaderFromDB 從資料庫取得請購單單頭
func (pr *PR) getHeaderFromDB(search *dto.PrSearch, user *dto.Users, m MySQL) *[]dto.PrListResult {
	db := m.GetAdater()
	sql := `
		SELECT 
			pl.id, 
			pl.pay_to, 
			pl.company, 
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
			pl.pay_date,
			pl.serial,
			pl.pr_item,
			pl.installment_plan,
			pl.pay_by,
			pl.memo,
			pl.create_at
		FROM 
			pr_lists pl
		INNER JOIN 
			users u ON pl.users_id = u.id
		WHERE
			pl.status = 1 %s
		ORDER BY
			pl.sign_at DESC
		LIMIT %d, %d
	`
	where := ""
	if user.ID != 0 {
		where = where + " %s"
		where = fmt.Sprintf(" AND pl.users_id = %d", user.ID)
	}
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
	dtoPrListResults := &[]dto.PrListResult{}
	db.Raw(sql).Scan(dtoPrListResults)
	return dtoPrListResults
}

//getDetailFromDB 從資料庫取得請購單單身
func (pr *PR) getDetailFromDB(list *dto.PrListResult, m MySQL) *[]dto.PrDetail {
	db := m.GetAdater()
	dtoPrDetail := &[]dto.PrDetail{}
	db.Where("pr_list_id = ?", list.ID).Order("id ASC").Find(dtoPrDetail)
	return dtoPrDetail
}

//getListFromDB 從資料庫取得請購單列表
func (pr *PR) getListFromDB(search *dto.PrSearch, user *dto.Users, m MySQL) *[]dto.GetList {
	db := m.GetAdater()
	sql := `
		SELECT 
			pd.*,
			pl.proof
		FROM 
			pr_details pd 
		INNER JOIN 
			pr_lists pl ON pl.id = pd.pr_list_id
		WHERE 
			pd.id IN (
				SELECT 
					MIN(pd.id)
				FROM 
					pr_lists pl
				INNER JOIN
					pr_details pd ON pl.id = pd.pr_list_id
				INNER JOIN 
					users u ON pl.users_id = u.id
				WHERE
					pl.status = 1 %s
				GROUP BY
					pd.pr_list_id
			)
		ORDER BY
			pl.sign_at DESC
		LIMIT %d, %d
	`
	where := ""
	if user.ID != 0 {
		where = where + " %s"
		where = fmt.Sprintf(" AND pl.users_id = %d", user.ID)
	}
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
	dtoGetList := &[]dto.GetList{}
	db.Raw(sql).Scan(dtoGetList)
	return dtoGetList
}

//doSetCancelToDB 將作廢資訊寫入資料蟀
func (pr *PR) doSetCancelToDB(u *dto.Users, id string, m MySQL) *dto.ResultObject {
	db := m.GetAdater()
	dtoPrList := &dto.PrList{}
	count := db.Model(&dtoPrList).Where("id = ? AND users_id = ? AND status != 0", id, u.ID).Update("status", 0).RowsAffected
	if count == 0 {
		dtoRO := RO.Build(0, "請購單作廢失敗")
		return dtoRO
	}
	dtoRO := RO.Build(1, "")
	return dtoRO
}

//doInsertToDB 將資料寫入資料庫
func (pr *PR) doInsertToDB(dtoPR *dto.PR, m MySQL) *dto.ResultObject {
	db := m.GetAdater()
	tx := db.Begin()
	err := tx.Create(&dtoPR.List).Error
	if err != nil {
		tx.Rollback()
		log.Println(err)
		dtoRO := RO.Build(0, "請購單單頭寫入失敗")
		return dtoRO
	}
	for _, v := range dtoPR.Detail {
		v.PRListID = dtoPR.List.ID
		err = tx.Create(&v).Error
		if err != nil {
			tx.Rollback()
			log.Println(err)
			dtoRO := RO.Build(0, "請購單單身寫入失敗")
			return dtoRO
		}
	}
	tx.Commit()
	dtoRO := RO.Build(1, "")
	return dtoRO
}

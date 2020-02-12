package service

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
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

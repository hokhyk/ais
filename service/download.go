package service

import (
	"strconv"

	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
)

//Download 下載檔案資料結構
type Download struct{}

//New 建構式
func (d Download) New() *Download {
	return &d
}

//GetFile 取得佐証檔
func (d *Download) GetFile(token string, fileName string) *dto.ResultObject {
	users := Users{}.New()

	dtoRO, dtoUsers := users.GetUser(token)
	if dtoRO.Status != 1 {
		return dtoRO
	}

	mysql := libs.MySQL{}.New()
	dtoRO = d.checkFile(dtoUsers, fileName, mysql)
	return dtoRO
}

//checkFile 檢查是否有權限下載檔案
func (d *Download) checkFile(u *dto.Users, fileName string, m MySQL) *dto.ResultObject {
	db := m.GetAdater()
	helper := libs.Helper{}.New()
	usersID := strconv.Itoa(u.ID)
	dtoPrList := &dto.PrList{}
	skiper := []string{"1", "12", "50"}
	db = db.Where("status = 1")
	if helper.InArray(skiper, usersID) {
		db.Where("proof = ?", fileName).Find(dtoPrList)
	} else {
		db.Where("users_id = ? AND proof = ?", usersID, fileName).Find(dtoPrList)
	}
	if dtoPrList.ID == 0 {
		dtoRO := RO.Build(0, "無此權限下載此檔案")
		return dtoRO
	}
	dtoRO := RO.Build(1, "")
	return dtoRO
}

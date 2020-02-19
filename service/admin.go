package service

import (
	"github.com/teed7334-restore/ais/dto"
)

//Admin 管理員頁面資料結構
type Admin struct{}

//New 建構式
func (a Admin) New() *Admin {
	return &a
}

//GetItem 取得請購單資料
func (a *Admin) GetItem(search *dto.PrSearch) (*dto.ResultObject, *dto.PrListResult, *[]dto.PrDetail) {
	user := &dto.Users{}
	pr := &PR{}
	return pr.GetItem(search, user)
}

//GetList 取得請購單列表
func (a *Admin) GetList(search *dto.PrSearch) (*dto.ResultObject, *[]dto.PrDetail) {
	user := &dto.Users{}
	pr := &PR{}
	return pr.GetList(search, user)
}

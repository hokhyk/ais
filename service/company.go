package service

import (
	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
)

//Company 廠商資料結構
type Company struct{}

//New 建構式
func (c Company) New() *Company {
	return &c
}

//GetCompany 取得廠商資料
func (c *Company) GetCompany() (*dto.ResultObject, *[]dto.Company) {
	mysql := libs.MySQL{}.New()
	dtoCompany := c.getCompanyFromDB(mysql)
	if len(*dtoCompany) == 0 {
		dtoRO := RO.Build(0, "查無任何資料")
		return dtoRO, dtoCompany
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoCompany
}

//getCompanyFromDB 從資料庫取得廠商列表
func (c *Company) getCompanyFromDB(m MySQL) *[]dto.Company {
	db := m.GetAdater()
	dtoCompany := &[]dto.Company{}
	db.Find(dtoCompany)
	return dtoCompany
}

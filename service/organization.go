package service

import (
	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
)

//Organization 單位資料結構
type Organization struct{}

//New 建構式
func (o Organization) New() *Organization {
	return &o
}

//GetOrganization 取得單位資料
func (o *Organization) GetOrganization() (*dto.ResultObject, *[]dto.Organization) {
	mysql := libs.MySQL{}.New()
	dtoOrganization := o.getOrganizationFromDB(mysql)
	if len(*dtoOrganization) == 0 {
		dtoRO := RO.Build(0, "查無任何資料")
		return dtoRO, dtoOrganization
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoOrganization
}

//getOrganizationFromDB 從資料庫取得單位列表
func (o *Organization) getOrganizationFromDB(m MySQL) *[]dto.Organization {
	db := m.GetAdater()
	dtoOrganization := &[]dto.Organization{}
	db.Table("organization").Find(dtoOrganization)
	return dtoOrganization
}

package service

import (
	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
)

//PayMethod 支付方式資料結構
type PayMethod struct{}

//New 建構式
func (pm PayMethod) New() *PayMethod {
	return &pm
}

//GetPayMethod 取得支付方式
func (pm *PayMethod) GetPayMethod() (*dto.ResultObject, *[]dto.PayMethod) {
	mysql := libs.MySQL{}.New()
	dtoPayMethod := pm.getPayMethodFromDB(mysql)
	if len(*dtoPayMethod) == 0 {
		dtoRO := RO.Build(0, "查無任何資料")
		return dtoRO, dtoPayMethod
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoPayMethod
}

//getPayMethodFromDB 從資料庫取得支付方式
func (pm *PayMethod) getPayMethodFromDB(m MySQL) *[]dto.PayMethod {
	db := m.GetAdater()
	dtoPayMethod := &[]dto.PayMethod{}
	db.Find(dtoPayMethod)
	return dtoPayMethod
}

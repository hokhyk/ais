package service

import (
	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
)

//CreditType 入帳類別資料結構
type CreditType struct{}

//New 建構式
func (ct CreditType) New() *CreditType {
	return &ct
}

//GetCreditType 取得入帳類別
func (ct *CreditType) GetCreditType() (*dto.ResultObject, *[]dto.CreditType) {
	mysql := libs.MySQL{}.New()
	dtoCreditType := ct.getCreditTypeFromDB(mysql)
	if len(*dtoCreditType) == 0 {
		dtoRO := RO.Build(0, "查無任何資料")
		return dtoRO, dtoCreditType
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoCreditType
}

//getCreditTypeFromDB 從資料庫取得支付方式
func (ct *CreditType) getCreditTypeFromDB(m MySQL) *[]dto.CreditType {
	db := m.GetAdater()
	dtoCreditType := &[]dto.CreditType{}
	db.Find(dtoCreditType)
	return dtoCreditType
}

package service

import (
	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
)

//PrItem 應付項目資料結構
type PrItem struct{}

//New 建構式
func (pi PrItem) New() *PrItem {
	return &pi
}

//GetPrItem 取得應付項目資料
func (pi *PrItem) GetPrItem() (*dto.ResultObject, *[]dto.PrItem) {
	mysql := libs.MySQL{}.New()
	dtoPrItem := pi.getPrItemFromDB(mysql)
	if len(*dtoPrItem) == 0 {
		dtoRO := RO.Build(0, "查無任何資料")
		return dtoRO, dtoPrItem
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoPrItem
}

//getPrItemFromDB 從資料庫取得貨幣列表
func (pi *PrItem) getPrItemFromDB(m MySQL) *[]dto.PrItem {
	db := m.GetAdater()
	dtoPrItem := &[]dto.PrItem{}
	db.Find(dtoPrItem)
	return dtoPrItem
}

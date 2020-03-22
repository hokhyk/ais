package service

import (
	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
)

//Currency 貨幣檔案資料結構
type Currency struct{}

//New 建構式
func (c Currency) New() *Currency {
	return &c
}

//GetCurrency 取得貨幣資料
func (c *Currency) GetCurrency() (*dto.ResultObject, *[]dto.Currency) {
	mysql := libs.MySQL{}.New()
	dtoCurrency := c.getCurrencyFromDB(mysql)
	if len(*dtoCurrency) == 0 {
		dtoRO := RO.Build(0, "查無任何資料")
		return dtoRO, dtoCurrency
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoCurrency
}

//getCurrencyFromDB 從資料庫取得貨幣列表
func (c *Currency) getCurrencyFromDB(m MySQL) *[]dto.Currency {
	db := m.GetAdater()
	dtoCurrency := &[]dto.Currency{}
	db.Find(dtoCurrency)
	return dtoCurrency
}

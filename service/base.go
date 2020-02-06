package service

import (
	"github.com/jinzhu/gorm"
	"github.com/teed7334-restore/ais/libs"
)

//TimeFormat 時間格式
var TimeFormat = "2006-01-02 15:04:05"

//RO 資料回傳物件
var RO = libs.ResultObject{}.New()

//Curl 物件介面
type Curl interface {
	Post(url string, params []byte, header map[string]string) []byte
}

//Redis 物件介面
type Redis interface {
	Get(key string) string
	Set(key, value string) bool
}

//Zip 物件介面
type Zip interface {
	Compress(files []string, dst string) int
}

//MySQL 物件介面
type MySQL interface {
	GetAdater() *gorm.DB
}

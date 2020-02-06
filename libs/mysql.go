package libs

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

//MySQL 物件參數
type MySQL struct {
	Conn *gorm.DB
}

//New 建構式
func (m MySQL) New() *MySQL {
	user := os.Getenv("mysql.user")
	password := os.Getenv("mysql.password")
	host := os.Getenv("mysql.host")
	database := os.Getenv("mysql.database")
	charset := os.Getenv("mysql.charset")
	parseTime := os.Getenv("mysql.parsetime")
	loc := os.Getenv("mysql.loc")
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=%s&loc=%s", user, password, host, database, charset, parseTime, loc)
	conn, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
	}
	m.Conn = conn
	return &m
}

//GetAdater 取得MySQL元件
func (m *MySQL) GetAdater() *gorm.DB {
	return m.Conn
}

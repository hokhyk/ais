package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
	"golang.org/x/crypto/scrypt"
)

//Users 使用者相關資料結構
type Users struct{}

type token struct {
	Expired string `json:"expired"`
	User    string `json:"user"`
	Token   string `json:"token"`
}

var addr = os.Getenv("redis.addr")
var pass = os.Getenv("redis.password")
var db, _ = strconv.Atoi(os.Getenv("redis.db"))
var redis = libs.Redis{}.New(addr, pass, db)
var curl = libs.Curl{}.New()

//New 建構式
func (u Users) New() *Users {
	return &u
}

//Login 登入
func (u *Users) Login(login, password string) (*dto.ResultObject, string) {
	valid := u.validateUser(login, password, curl)
	if valid != true {
		dtoRO := RO.Build(0, "使用者帳號或密碼有誤")
		return dtoRO, ""
	}
	hash := u.generateHash(login)
	token := u.generateToken(login, hash, redis)
	if token != true {
		dtoRO := RO.Build(0, "寫入Token失敗")
		return dtoRO, ""
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, hash
}

//GetUser 取得使用者資訊
func (u *Users) GetUser(hash string) (*dto.ResultObject, *dto.Users) {
	data := redis.Get(hash)
	dtoUsers := new(dto.Users)
	if data == "" {
		dtoRO := RO.Build(0, "請登入會員")
		return dtoRO, dtoUsers
	}
	token := &token{}
	json.Unmarshal([]byte(data), token)
	q := url.Values{}
	q.Add("login", token.User)
	params := []byte(q.Encode())
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	api := fmt.Sprintf("%s/hack/getUserByAjax", os.Getenv("hrm.url"))
	result := curl.Post(api, params, header)
	err := json.Unmarshal(result, dtoUsers)
	if err != nil {
		log.Println(err)
		dtoRO := RO.Build(0, "查無此會員")
		return dtoRO, dtoUsers
	}
	dtoRO := RO.Build(1, "")
	return dtoRO, dtoUsers
}

//Logout 登出
func (u *Users) Logout(hash string) {
	redis.Remove(hash)
}

//CheckLogin 檢查是否處於登入狀態
func (u *Users) CheckLogin(hash string) *dto.ResultObject {
	data := redis.Get(hash)
	if data == "" {
		dtoRO := RO.Build(0, "請登入會員")
		return dtoRO
	}
	token := &token{}
	json.Unmarshal([]byte(data), token)
	expired, _ := time.ParseInLocation(TimeFormat, token.Expired, time.Local)
	now := time.Now()
	if !now.Before(expired) {
		redis.Remove(hash)
		dtoRO := RO.Build(0, "此令牌已過期")
		return dtoRO
	}
	dtoRO := RO.Build(1, "")
	return dtoRO
}

//validateUser 驗証使用者
func (u *Users) validateUser(login, password string, curl Curl) bool {
	api := fmt.Sprintf("%s/hack/checkLoginByAjax", os.Getenv("hrm.url"))
	q := url.Values{}
	q.Add("login", login)
	q.Add("password", password)
	params := []byte(q.Encode())
	header := make(map[string]string)
	header["Content-Type"] = "application/x-www-form-urlencoded"
	result := string(curl.Post(api, params, header))
	status := true

	if result != "true" {
		status = false
	}

	return status
}

//generateHash 生成Hash值
func (u *Users) generateHash(login string) string {
	salt := os.Getenv("salt")
	dk, _ := scrypt.Key([]byte(login), []byte(salt), 32768, 8, 1, 32)
	txt := base58.Encode(dk)
	log.Println(txt)
	return txt
}

//generateToken 生成登入令牌
func (u *Users) generateToken(login, hash string, redis Redis) bool {
	token := &token{}
	now := time.Now()
	expired, _ := time.ParseDuration("3h")
	token.Expired = now.Add(expired).Format(TimeFormat)
	token.User = login
	token.Token = hash
	json, _ := json.Marshal(token)
	jsonStr := string(json)
	status := redis.Set(hash, jsonStr)
	return status
}

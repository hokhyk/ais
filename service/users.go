package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/teed7334-restore/ais/libs"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

//Users 使用者相關資料結構
type Users struct {
}

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
func (u *Users) Login(login, password string) (status int, hash string) {
	valid := u.validateUser(login, password, curl)
	if valid != true {
		return -1, ""
	}
	hash = u.generateHash(login)
	token := u.generateToken(login, hash, redis)
	if token != true {
		return -2, ""
	}
	return 1, hash
}

//Logout 登出
func (u *Users) Logout(hash string) {
	redis.Remove(hash)
}

//CheckLogin 檢查是否處於登入狀態
func (u *Users) CheckLogin(hash string) int {
	data := redis.Get(hash)
	if data == "" {
		return -1
	}
	token := &token{}
	json.Unmarshal([]byte(data), token)
	expired, _ := time.ParseInLocation(TimeFormat, token.Expired, time.Local)
	now := time.Now()
	if !now.Before(expired) {
		redis.Remove(hash)
		return -2
	}
	return 1
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
	data := []byte(login + salt)
	s := sha256.New()
	s.Write(data)
	hash := s.Sum(nil)
	txt := hex.EncodeToString(hash)
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
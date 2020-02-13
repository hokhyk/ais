package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/teed7334-restore/ais/service"
)

//Users 使用者相關資料結構
type Users struct{}

var users = service.Users{}.New()

//New 建構式
func (u Users) New() *Users {
	return &u
}

//GetUser 取得使用者
func (u *Users) GetUser(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	dtoRO, dtoUsers := users.GetUser(token)
	data, _ := json.Marshal(dtoUsers)
	message := string(data)
	PrintRO(w, dtoRO, message)
	return
}

//Login 登入
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	if login == "" || password == "" {
		content := RO.BuildJSON(0, "使用者帳號或密碼不可為空白")
		fmt.Fprintf(w, content)
		return
	}

	dtoRO, hash := users.Login(login, password)
	PrintRO(w, dtoRO, hash)
	return
}

//Logout 登出
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	users.Logout(token)
	content := RO.BuildJSON(1, "true")
	fmt.Fprintf(w, content)
}

//CheckLogin 檢查是否處於登入狀態
func (u *Users) CheckLogin(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	dtoRO := users.CheckLogin(token)
	PrintRO(w, dtoRO, "true")
	return
}

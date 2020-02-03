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

	if r.Method != "POST" {
		content := buildRO(201, "不支持此種HTTP Method")
		fmt.Fprintf(w, content)
		return
	}

	token := r.FormValue("token")

	if token == "" {
		content := buildRO(202, "使用者令牌為空白")
		fmt.Fprintf(w, content)
		return
	}

	status, message := users.GetUser(token)

	switch status {
	case -1:
		content := buildRO(203, "請登入會員")
		fmt.Fprintf(w, content)
	case -2:
		content := buildRO(204, "查無此會員")
		fmt.Fprintf(w, content)
	case 1:
		result, _ := json.Marshal(message)
		response := string(result)
		content := buildRO(200, response)
		fmt.Fprintf(w, content)
	}
}

//Login 登入
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		content := buildRO(201, "不支持此種HTTP Method")
		fmt.Fprintf(w, content)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	if login == "" || password == "" {
		content := buildRO(202, "使用者帳號或密碼不可為空白")
		fmt.Fprintf(w, content)
		return
	}

	status, hash := users.Login(login, password)

	switch status {
	case -1:
		content := buildRO(203, "使用者帳號或密碼有誤")
		fmt.Fprintf(w, content)
	case -2:
		content := buildRO(204, "寫入Token失敗")
		fmt.Fprintf(w, content)
	case 1:
		content := buildRO(200, hash)
		fmt.Fprintf(w, content)
	}

	return
}

//Logout 登出
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		content := buildRO(201, "不支持此種HTTP Method")
		fmt.Fprintf(w, content)
		return
	}

	token := r.FormValue("token")

	users.Logout(token)

	content := buildRO(200, "true")
	fmt.Fprintf(w, content)
}

//CheckLogin 檢查是否處於登入狀態
func (u *Users) CheckLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		content := buildRO(201, "不支持此種HTTP Method")
		fmt.Fprintf(w, content)
		return
	}

	token := r.FormValue("token")

	if token == "" {
		content := buildRO(202, "使用者令牌為空白")
		fmt.Fprintf(w, content)
		return
	}

	status := users.CheckLogin(token)

	switch status {
	case -1:
		content := buildRO(203, "此帳號已登出")
		fmt.Fprintf(w, content)
	case -2:
		content := buildRO(204, "此令牌已過期")
		fmt.Fprintf(w, content)
	case 1:
		content := buildRO(200, "true")
		fmt.Fprintf(w, content)
	}
}

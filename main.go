package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/teed7334-restore/ais/controller"
	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
	"github.com/teed7334-restore/ais/service"
)

var ro = libs.ResultObject{}.New()

var user = controller.Users{}.New()
var pr = controller.PR{}.New()
var download = controller.Download{}.New()

var suser = service.Users{}.New()

var helper = libs.Helper{}.New()

//middleware 中間件
func middleware(next http.Handler, method []string, needLogin int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = setCORS(w)
		dtoRO := checkHTTPMethod(r, method)
		if dtoRO.Status != 1 {
			content := ro.BuildJSON(0, dtoRO.Message)
			fmt.Fprintf(w, content)
			return
		}
		if needLogin == 1 {
			dtoRO = doAuth(r)
			if dtoRO.Status != 1 {
				content := ro.BuildJSON(0, dtoRO.Message)
				fmt.Fprintf(w, content)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

//setCORS 設定CORS
func setCORS(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return w
}

//checkHTTPMethod 檢查HTTP Method
func checkHTTPMethod(r *http.Request, method []string) *dto.ResultObject {
	accept := helper.InArray(method, r.Method)
	dtoRO := ro.Build(1, "")
	if !accept {
		dtoRO = ro.Build(0, "不支持此種HTTP Method")
		return dtoRO
	}
	return dtoRO
}

//doAuth 使用者驗証
func doAuth(r *http.Request) *dto.ResultObject {
	token := r.FormValue("token")
	dtoRO := ro.Build(1, "")
	if token == "" {
		dtoRO = ro.Build(0, "使用者令牌為空白")
		return dtoRO
	}
	dtoRO = suser.CheckLogin(token)
	return dtoRO
}

//main 主程式
func main() {
	port := fmt.Sprintf(":%s", os.Getenv("port"))
	handler := http.HandlerFunc(user.GetUser)
	http.Handle("/users/getUser", middleware(handler, []string{"POST"}, 1))
	handler = http.HandlerFunc(user.Login)
	http.Handle("/users/login", middleware(handler, []string{"POST"}, 0))
	handler = http.HandlerFunc(user.CheckLogin)
	http.Handle("/users/checkLogin", middleware(handler, []string{"POST"}, 1))
	handler = http.HandlerFunc(user.Logout)
	http.Handle("/users/logout", middleware(handler, []string{"POST"}, 0))
	handler = http.HandlerFunc(pr.Add)
	http.Handle("/pr/add", middleware(handler, []string{"POST"}, 1))
	handler = http.HandlerFunc(pr.SetCancel)
	http.Handle("/pr/setCancel", middleware(handler, []string{"POST"}, 1))
	handler = http.HandlerFunc(download.GetFile)
	http.Handle("/download/getFile", middleware(handler, []string{"GET", "POST"}, 1))
	http.ListenAndServe(port, nil)
}

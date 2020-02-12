package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/teed7334-restore/ais/controller"
	"github.com/teed7334-restore/ais/libs"
)

var ro = libs.ResultObject{}.New()

var user = controller.Users{}.New()
var pr = controller.PR{}.New()
var download = controller.Download{}.New()

var helper = libs.Helper{}.New()

//middleware 中間層
func middleware(next http.Handler, method []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := helper.InArray(method, r.Method)
		if !accept {
			content := ro.BuildJSON(0, "不支持此種HTTP Method")
			fmt.Fprintf(w, content)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

//main 主程式
func main() {
	port := fmt.Sprintf(":%s", os.Getenv("port"))
	handler := http.HandlerFunc(user.GetUser)
	http.Handle("/users/getUser", middleware(handler, []string{"POST"}))
	handler = http.HandlerFunc(user.Login)
	http.Handle("/users/login", middleware(handler, []string{"POST"}))
	handler = http.HandlerFunc(user.CheckLogin)
	http.Handle("/users/checkLogin", middleware(handler, []string{"POST"}))
	handler = http.HandlerFunc(user.Logout)
	http.Handle("/users/logout", middleware(handler, []string{"POST"}))
	handler = http.HandlerFunc(pr.Add)
	http.Handle("/pr/add", middleware(handler, []string{"POST"}))
	handler = http.HandlerFunc(pr.SetCancel)
	http.Handle("/pr/setCancel", middleware(handler, []string{"POST"}))
	handler = http.HandlerFunc(download.GetFile)
	http.Handle("/download/getFile", middleware(handler, []string{"GET", "POST"}))
	http.ListenAndServe(port, nil)
}

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

//middleware 中間層
func middleware(next http.Handler, method []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := inArray(method, r.Method)
		if !accept {
			content := ro.BuildJSON(0, "不支持此種HTTP Method")
			fmt.Fprintf(w, content)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

//InArray 於陣列中有對應之值
func inArray(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
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
	http.ListenAndServe(port, nil)
}

package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/teed7334-restore/ais/controller"
)

var user = controller.Users{}.New()
var pr = controller.PR{}.New()

//middleware 中間層
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

//main 主程式
func main() {
	port := fmt.Sprintf(":%s", os.Getenv("port"))
	handler := http.HandlerFunc(user.GetUser)
	http.Handle("/users/getUser", middleware(handler))
	handler = http.HandlerFunc(user.Login)
	http.Handle("/users/login", middleware(handler))
	handler = http.HandlerFunc(user.CheckLogin)
	http.Handle("/users/checkLogin", middleware(handler))
	handler = http.HandlerFunc(user.Logout)
	http.Handle("/users/logout", middleware(handler))
	handler = http.HandlerFunc(pr.Add)
	http.Handle("/pr/add", middleware(handler))
	http.ListenAndServe(port, nil)
}

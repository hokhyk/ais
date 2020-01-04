package main

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/teed7334-restore/ais/controller"
)

var user = controller.Users{}.New()

func main() {
	http.HandleFunc("/users/login", user.Login)
	http.HandleFunc("/users/checkLogin", user.CheckLogin)
	http.HandleFunc("/users/logout", user.Logout)
	http.ListenAndServe(":8080", nil)
}

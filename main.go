package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/teed7334-restore/ais/controller"
)

var user = controller.Users{}.New()

func main() {
	port := fmt.Sprintf(":%s", os.Getenv("port"))
	http.HandleFunc("/users/login", user.Login)
	http.HandleFunc("/users/checkLogin", user.CheckLogin)
	http.HandleFunc("/users/logout", user.Logout)
	http.ListenAndServe(port, nil)
}

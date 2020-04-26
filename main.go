package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"todoapp/app"
	"todoapp/controllers"
	"todoapp/model"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.JWTAuth)
	router.Use(app.TranceRequests)

	port := os.Getenv("APP_PORT")
	fmt.Println("APP_PORT: " + port)
	model.GetDB()

	router.HandleFunc("/api/users", controllers.GetAccounts).Methods(http.MethodGet)
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods(http.MethodPost)
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods(http.MethodPost)
	router.HandleFunc("/api/user", controllers.GetAccount).Methods(http.MethodGet)

	err := http.ListenAndServe(":"+port, router) //Запустите приложение, посетите localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}

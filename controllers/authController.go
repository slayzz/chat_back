package controllers

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"todoapp/model"
	u "todoapp/utils"
)

func GetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := model.GetUsers()
	if err != nil {
		log.Println("error: in GetAccounts", err)
		u.Respond(w, u.Message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	u.Respond(w, map[string]interface{}{"data": accounts}, http.StatusOK)
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user")

	account := &model.Account{
		ID: userID.(primitive.ObjectID),
	}
	err := account.GetUser()
	if err != nil {
		log.Println("error: in GetAccount", err)
		u.Respond(w, u.Message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	u.Respond(w, map[string]interface{}{"data": account}, http.StatusOK)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := &model.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		log.Println("error: in CreateAccount", err)
		u.Respond(w, u.Message(false, err.Error()), http.StatusInternalServerError)
		return
	}

	err = account.Create()
	if err != nil {
		log.Println("error: in CreateAccount", err)
		u.Respond(w, u.Message(false, err.Error()), http.StatusUnauthorized)
		return
	}

	u.Respond(w, map[string]interface{}{"data": account}, http.StatusCreated)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	account := &model.Account{}

	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		log.Println("error: in Authenticate", err)
		u.Respond(w, u.Message(false, err.Error()), http.StatusUnauthorized)
		return
	}

	err = account.Login()
	if err != nil {
		log.Println("error: in Authenticate", err)
		u.Respond(w, u.Message(false, err.Error()), http.StatusUnauthorized)
		return
	}
	u.Respond(w, map[string]interface{}{"data": account}, http.StatusAccepted)
}

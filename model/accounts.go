package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

type Token struct {
	UserID primitive.ObjectID
	jwt.StandardClaims
}

type Account struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string             `json:"email"`
	Password string             `json:"password,omitempty"`
	Token    string             `json:"token"`
}

var collectionName = "accounts"

func (account *Account) Validate() error {

	if !strings.Contains(account.Email, "@") {
		return errors.New("email address is required")
	}

	if len(account.Password) < 6 {
		return errors.New("password address is required")
	}

	temp := Account{}

	mdb := GetDB()
	mdb.Collection(collectionName).FindOne(context.Background(), bson.D{{Key: "email", Value: account.Email}}).Decode(&temp)

	if temp.Email != "" {
		return errors.New("email must be unique")
	}

	return nil
}

func (account *Account) Create() error {
	if err := account.Validate(); err != nil {
		return err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	mdb := GetDB()
	result, err := mdb.Collection(collectionName).InsertOne(context.Background(), account)
	if err != nil {
		return err
	}

	id := result.InsertedID.(primitive.ObjectID)
	account.ID = id

	tk := &Token{UserID: id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	account.Token = tokenString
	return nil
}

func (account *Account) Login() error {
	mdb := GetDB()

	dbAccount := &Account{}
	err := mdb.Collection(collectionName).FindOne(context.Background(), bson.D{{Key: "email", Value: account.Email}}).Decode(dbAccount)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(dbAccount.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("invalid login credentials. Please try again")
	}


	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))

	account.Token = tokenString
	account.ID = dbAccount.ID
	account.Password = dbAccount.Password
	account.Email = dbAccount.Email

	return nil
}

func (account *Account) GetUser() error {
	mdb := GetDB()

	err := mdb.Collection(collectionName).FindOne(context.Background(), bson.D{{Key: "_id", Value: account.ID}}).Decode(&account)
	if err != nil {
		return err
	}

	return nil
}

func GetUsers() ([]Account, error) {
	mdb := GetDB()
	cur, err := mdb.Collection("accounts").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var accounts []Account
	for cur.Next(context.Background()) {
		account := Account{}

		err := cur.Decode(&account)
		if err != nil {
			fmt.Println(err)
		}

		accounts = append(accounts, account)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

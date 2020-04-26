package model

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client //база данных
var err error
var dbName = os.Getenv("DB_NAME")

func init() {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")

	dbUri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin", username, password, dbHost, dbName)
	fmt.Println(dbUri)

	client, err = mongo.NewClient(options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal("can't create client to mongodb", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("can't open connection to mongodb", err)
	}

	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

// возвращает дескриптор объекта DB
func GetDB() *mongo.Database {
	return client.Database(dbName)
}

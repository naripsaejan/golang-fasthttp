package main

import (
	//lib
	"context"
	"fmt"
	"log"

	//utils
	"examp/hello-fast-http/service"

	// BookGetAll

	//lib

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// -------------------- variable -----------------------//
const (
	mongoURI   = "mongodb://test:password@10.138.41.195:27017,10.138.41.196:27017,10.138.41.197:27017/?authSource=test&replicaSet=nmgw"
)

//-------------------- function -----------------------//

func main() {
	r := router.New()
		connectToMongo()
	defer service.MongoClient.Disconnect(context.Background())

	r.GET("/books/read", service.BookGetAll)
	r.GET("/books/readConsumer", service.BookGetAllConsumer)
	r.GET("/books/read/{id}", service.BookGetByID)
	r.POST("/books/create", service.BookPost)
	r.PATCH("/books/update/{id}", service.BookPatch)
	r.DELETE("/books/delete/{id}", service.BookDeleteById)

	fasthttp.ListenAndServe(":3000", r.Handler)
}

func connectToMongo() {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	service.MongoClient = client
	fmt.Println("Connected to MongoDB")
}
package utils

// import (
// 	"context"
// 	"examp/hello-fast-http/service"
// 	"fmt"
// 	"log"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func connectToMongo() {
// 	clientOptions := options.Client().ApplyURI(service.mongoURI)
// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	if err != nil {
// 		log.Fatal("Failed to connect to MongoDB:", err)
// 	}
// 	service.MongoClient = client
// 	fmt.Println("Connected to MongoDB")
// }
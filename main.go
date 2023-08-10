package main

import (
	//lib
	"context"
	"encoding/json"
	"fmt"
	"log"

	//utils
	"examp/hello-fast-http/service"
	"examp/hello-fast-http/utils"

	// BookGetAll

	//lib
	"github.com/IBM/sarama"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// -------------------- variable -----------------------//
// Constants for error messages and status codes
const (
	invalidRequest  = `{"status":"40000","error":"Invalid request"}`
	internalError   = `{"status":"50000","error":"Internal server error"}`
	notFoundError   = `{"status":"40004","error":"Book not found"}`
	successMessage = `{"status":"20000","message":"Operation successful"}`
)

// MongoDB and Kafka configuration constants
const (
	dbName     = "test"
	kafkaName  = "rip-test"
	mongoURI   = "mongodb://test:password@10.138.41.195:27017,10.138.41.196:27017,10.138.41.197:27017/?authSource=test&replicaSet=nmgw"
)


var (
	mongoClient *mongo.Client
	kafkaBrokers = []string{"10.138.41.195:9092", "10.138.41.196:9092"}

)

//-------------------- function -----------------------//

func main() {
	r := router.New()
	connectToMongo()
	defer mongoClient.Disconnect(context.Background())

	r.GET("/books/read", service.BookGetAll)
	r.GET("/books/read/{id}", service.BookGetByID)
	r.POST("/books/create", BookPost)
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

func isBookIDUnique(id string) bool {
	collection := mongoClient.Database(dbName).Collection("books")

	filter := bson.M{"id": id}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Println("Error checking book ID uniqueness:", err)
		return false
	}

	return count == 0
}

// Post
func BookPost(ctx *fasthttp.RequestCtx) {
    // Set response content type
    ctx.SetContentType("application/json")

    // Parse request body
    var newBook service.Book
    err := json.Unmarshal(ctx.Request.Body(), &newBook)
    if err != nil {
        ctx.SetStatusCode(fasthttp.StatusBadRequest)
        ctx.Write([]byte(invalidRequest))
        return
    }

    // Check if the book ID is unique :set unique in database
    if !isBookIDUnique(newBook.ID) || newBook.ID == "" {
        ctx.SetStatusCode(fasthttp.StatusBadRequest)
        ctx.Write([]byte(invalidRequest))
        return
    }

    // Insert the new book into MongoDB
    collection := mongoClient.Database(dbName).Collection("books")
    _, err = collection.InsertOne(context.Background(), newBook)
    if err != nil {
        ctx.SetStatusCode(fasthttp.StatusInternalServerError)
        ctx.Write([]byte(internalError))
        return
    }

    // Marshal the newBook to JSON
    responseJSON, err := json.Marshal(newBook)
    if err != nil {
        ctx.SetStatusCode(fasthttp.StatusInternalServerError)
        ctx.Write([]byte(internalError))
        return
    }

    // Create a Kafka producer config
    kafkaConfig := sarama.NewConfig()
    kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
    kafkaConfig.Producer.Return.Successes = true
    kafkaConfig.Producer.Return.Errors = true

    // Create a Kafka producer instance

	kafkaProducer, err := utils.InitializeKafkaProducer(kafkaBrokers)
	if err != nil {
		log.Fatal("Failed to initialize Kafka producer:", err)
	}
	defer kafkaProducer.Close()
	
        // Send a Kafka message for the new book
		err = utils.SendMessageToKafka(kafkaProducer, kafkaName, newBook.ID, responseJSON)
		if err != nil {
			log.Println("Failed to send Kafka message:", err)
		} else {
			log.Println("Kafka message sent successfully")
		}

    // Respond with the added book
    ctx.Write(responseJSON)
}

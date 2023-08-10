package service

import (
	"context"
	"encoding/json"
	"log"

	"examp/hello-fast-http/utils"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
)


func isBookIDUnique(id string) bool {
	collection := MongoClient.Database(dbName).Collection("books")

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
    var newBook Book
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
    collection := MongoClient.Database(dbName).Collection("books")
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

    // Create a Kafka producer instance
	kafkaProducer, err := utils.InitializeKafkaProducer(kafkaBrokers)
	if err != nil {
		log.Fatal("Failed to initialize Kafka producer:", err)
	}
	defer kafkaProducer.Close()
	
        // Send a Kafka message for the new book
		err = utils.SendMessageToKafka(kafkaProducer, kafkaTopics, newBook.ID, responseJSON)
		if err != nil {
			log.Println("Failed to send Kafka message:", err)
		} else {
			log.Println("Kafka message sent successfully")
		}

    // Respond with the added book
    ctx.Write(responseJSON)
}

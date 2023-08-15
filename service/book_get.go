package service

import (
	"context"
	"encoding/json"
	"examp/hello-fast-http/utils"

	"log"

	"github.com/valyala/fasthttp"

	"go.mongodb.org/mongo-driver/bson"
)

// BookGetAllHandler handles retrieving all books from the MongoDB collection.
func BookGetAll(ctx *fasthttp.RequestCtx) {
	// Set response content type
	ctx.SetContentType("application/json")

	// Retrieve data from MongoDB
	collection := utils.MongoClient.Database(utils.DbName).Collection("books")
	filter := bson.D{} // You can add filtering criteria here

	var results []Book
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var book Book
		if err := cur.Decode(&book); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte(internalError))
			return
		}
		results = append(results, book)
	}

	// Respond with data
	jsonBytes, _ := json.Marshal(results)
	ctx.Write(jsonBytes)

}

// GetById
func BookGetByID(ctx *fasthttp.RequestCtx) {
	// Set response content type
	ctx.SetContentType("application/json")

	idStr := ctx.UserValue("id").(string)
	log.Println("bookID", idStr)
	// Retrieve book by ID from MongoDB
	collection := utils.MongoClient.Database(utils.DbName).Collection("books")

	// filter := bson.M{"_id": idStr}
	filter := bson.M{"id": idStr}

	var book Book
	err := collection.FindOne(context.Background(), filter).Decode(&book)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write([]byte(notFoundError))
		return
	}

	// Respond with the retrieved book
	jsonBytes, _ := json.Marshal(book)
	ctx.Write(jsonBytes)
}

//----------------get consumer -------------------//
// ConsumeKafkaMessages consumes Kafka messages.
func ConsumeKafkaMessages() {
	err := utils.ConsumeMessagesFromKafka(utils.KafkaBrokers, utils.KafkaTopics)
	if err != nil {
		log.Println("Error consuming Kafka messages:", err)
	}
}

func BookGetAllConsumer(ctx *fasthttp.RequestCtx) {
	// Set response content type
	ctx.SetContentType("application/json")

	// Retrieve data from MongoDB
	collection := utils.MongoClient.Database(utils.DbName).Collection("books")
	filter := bson.D{} // You can add filtering criteria here

	var results []Book
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var book Book
		if err := cur.Decode(&book); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte(internalError))
			return
		}
		results = append(results, book)
	}

	// Respond with data
	jsonBytes, _ := json.Marshal(results)
	ctx.Write(jsonBytes)
    
	// Consume Kafka messages in the background
	go func() {
log.Println("get=> ConsumeMessagesFromKafka")

		err := utils.ConsumeMessagesFromKafka(utils.KafkaBrokers, utils.KafkaTopics)
		if err != nil {
			log.Println("Error consuming Kafka messages:", err)
		}
	}()

}

func BookGetAllConsumerGroup(ctx *fasthttp.RequestCtx) {
	// Set response content type
	ctx.SetContentType("application/json")

	// Retrieve data from MongoDB
	collection := utils.MongoClient.Database(utils.DbName).Collection("books")
	filter := bson.D{} // You can add filtering criteria here

	var results []Book
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var book Book
		if err := cur.Decode(&book); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte(internalError))
			return
		}
		results = append(results, book)
	}

	// Respond with data
	jsonBytes, _ := json.Marshal(results)
	ctx.Write(jsonBytes)

	// Consume Kafka messages in the background using a consumer group
	go func() {
		log.Println("get=> ConsumeMessagesWithConsumerGroup")

		err := utils.ConsumeMessagesWithConsumerGroup(utils.KafkaBrokers, utils.KafkaTopics, "rip-group1")
		if err != nil {
			log.Println("Error consuming Kafka messages:", err)
		}
	}()
}

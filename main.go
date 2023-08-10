package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// -------------------- struct -----------------------//
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// -------------------- variable -----------------------//
const (
	dbName          = "test"
	
	invalidRequest  = `{"status":"40000","error":"Invalid request"}`
	internalError   = `{"status":"50000","error":"Internal server error"}`
	notFoundError   = `{"status":"40004","error":"Book not found"}`
	successMessage = `{"status":"20000","message": "Operation successful"}`
)

var (
	mongoClient   *mongo.Client
)

//-------------------- function -----------------------//

// connectDB
func connectToMongo() error {
	mongoURI := "mongodb://test:password@10.138.41.195:27017,10.138.41.196:27017,10.138.41.197:27017/?authSource=test&replicaSet=nmgw"
	// mongoURI := "mongodb://localhost:27017/"
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	mongoClient = client
	fmt.Println("Connected to MongoDB")
	return nil
}


// check duplicate ID
func isBookIDUnique(id string) bool {
	collection := mongoClient.Database(dbName).Collection("books")

	filter := bson.M{"id": id}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		// Handle error (e.g., log, return false)
		return false
	}

	return count == 0
}

// Get
func BookGetAll(ctx *fasthttp.RequestCtx) {
	// Set response content type
	ctx.SetContentType("application/json")

	// Retrieve data from MongoDB
	collection := mongoClient.Database(dbName).Collection("books")
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
	collection := mongoClient.Database(dbName).Collection("books")

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

// Post
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
    kafkaProducer, err := sarama.NewSyncProducer([]string{"10.138.41.195:9092","10.138.41.196:9092"}, kafkaConfig)
    if err != nil {
        log.Println("Failed to create Kafka producer:", err)
    } else {
        defer kafkaProducer.Close()

        // Send a Kafka message for the new book
        kafkaMessage := &sarama.ProducerMessage{
            Topic: "rip-test",
            Key:   sarama.StringEncoder(newBook.ID),         // Use book ID as the key
            Value: sarama.StringEncoder(string(responseJSON)), // Use responseJSON as the value
        }

        partition, offset, err := kafkaProducer.SendMessage(kafkaMessage)
        if err != nil {
            log.Println("Failed to send Kafka message:", err)
        } else {
            log.Printf("Kafka message sent. Partition: %d, Offset: %d\n", partition, offset)
        }
    }

    // Respond with the added book
    ctx.Write(responseJSON)
}



// Patch
func BookPatch(ctx *fasthttp.RequestCtx) {
	// Set response content type
	ctx.SetContentType("application/json")

	// Extract the book ID from the URL path parameters
	idStr := ctx.UserValue("id").(string)

	// Parse the request body
	var updatedBook Book
	err := json.Unmarshal(ctx.PostBody(), &updatedBook)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(invalidRequest))
		return
	}

	// Check the request id in body
	if len(updatedBook.ID) != 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(invalidRequest))
		return
	}

	// Update the book in MongoDB
	collection := mongoClient.Database(dbName).Collection("books")
	filter := bson.M{"id": idStr}

	updateFields := bson.M{}

	updateFields["title"] = updatedBook.Title

	updateFields["author"] = updatedBook.Author

	update := bson.M{"$set": updateFields}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	log.Println("filter", filter)
	log.Println("update", update)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}

	// Retrieve the updated book from MongoDB to include the id field
	var updatedBookWithID Book
	err = collection.FindOne(context.Background(), filter).Decode(&updatedBookWithID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}

	// Respond with the updated book
	responseJSON, err := json.Marshal(updatedBookWithID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}
	ctx.Write(responseJSON)
}

// Delete ById
func BookDeleteById(ctx *fasthttp.RequestCtx) {
	// Set response content type
	ctx.SetContentType("application/json")

	idStr, ok := ctx.UserValue("id").(string)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(invalidRequest))
		return
	}

	// Delete the book from MongoDB
	collection := mongoClient.Database(dbName).Collection("books")
	filter := bson.M{"id": idStr}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}

	ctx.Write([]byte(successMessage))
}

func main() {
	r := router.New()
	//run database
	connectToMongo()
	defer mongoClient.Disconnect(context.Background())


	//------EndPoint-----------//
	r.GET("/books/read", BookGetAll)
	r.GET("/books/read/{id}", BookGetByID)
	r.POST("/books/create", BookPost)
	r.PATCH("/books/update/{id}", BookPatch)
	r.DELETE("/books/delete/{id}", BookDeleteById)

	fasthttp.ListenAndServe(":3000", r.Handler)

}

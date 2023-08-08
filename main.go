package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
var (
	mongoClient *mongo.Client
	dbName      = "test"
)

//-------------------- function -----------------------//

// connectDB
func connectToMongo() {
	mongoURI := "mongodb://test:password@10.138.41.195:27017,10.138.41.196:27017,10.138.41.197:27017/?authSource=test&replicaSet=nmgw"
	// mongoURI := "mongodb://localhost:27017/"
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	mongoClient = client
	fmt.Println("Connected to MongoDB")
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
	// Retrieve data from MongoDB
	collection := mongoClient.Database(dbName).Collection("books")
	filter := bson.D{} // You can add filtering criteria here

	var results []Book
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString("Error retrieving data from MongoDB")
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var book Book
		if err := cur.Decode(&book); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.WriteString("Error decoding data from MongoDB")
			return
		}
		results = append(results, book)
	}

	// Respond with the retrieved data
	ctx.SetContentType("application/json")
	jsonBytes, _ := json.Marshal(results)
	ctx.Write(jsonBytes)

}

// GetById
func BookGetByID(ctx *fasthttp.RequestCtx) {
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
		ctx.WriteString("Book not found")
		return
	}

	// Respond with the retrieved book
	ctx.SetContentType("application/json")
	jsonBytes, _ := json.Marshal(book)
	ctx.Write(jsonBytes)
}

// Post
func BookPost(ctx *fasthttp.RequestCtx) {

	// Parse request body
	var newBook Book
	err := json.Unmarshal(ctx.Request.Body(), &newBook)
	if err != nil {
		ctx.Error("Invalid JSON", fasthttp.StatusBadRequest)
		return
	}

	log.Println("check new id", newBook.ID)
	// Check if the book ID is unique
	if !isBookIDUnique(newBook.ID) {
		ctx.WriteString("Invalid add id")
		return
	}

	// Insert the new book into MongoDB
	collection := mongoClient.Database(dbName).Collection("books")
	_, err = collection.InsertOne(context.Background(), newBook)
	if err != nil {
		ctx.Error("Error inserting data into MongoDB", fasthttp.StatusInternalServerError)
		return
	}

	// Respond with the added book
	responseJSON, err := json.Marshal(newBook)
	if err != nil {
		ctx.Error("Failed to marshal JSON", fasthttp.StatusInternalServerError)
		return
	}

	// Set response content type
	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(responseJSON)
}

// Patch
func BookPatch(ctx *fasthttp.RequestCtx) {
	// Extract the book ID from the URL path parameters
	idStr := ctx.UserValue("id").(string)

	// Parse the request body
	var updatedBook Book
	err := json.Unmarshal(ctx.PostBody(), &updatedBook)
	if err != nil {
		ctx.Error("Invalid JSON", fasthttp.StatusBadRequest)
		return
	}

	// Update the book in MongoDB
	collection := mongoClient.Database(dbName).Collection("books")
	filter := bson.M{"id": idStr}

	updateFields := bson.M{}
	if updatedBook.Title != "" {
		updateFields["title"] = updatedBook.Title
	}
	if updatedBook.Author != "" {
		updateFields["author"] = updatedBook.Author
	}

	update := bson.M{"$set": updateFields}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	log.Println("filter", filter)
	log.Println("update", update)
	if err != nil {
		ctx.Error("Error updating data in MongoDB", fasthttp.StatusInternalServerError)
		return
	}

	// Retrieve the updated book from MongoDB to include the id field
	var updatedBookWithID Book
	err = collection.FindOne(context.Background(), filter).Decode(&updatedBookWithID)
	if err != nil {
		ctx.Error("Error retrieving updated book", fasthttp.StatusInternalServerError)
		return
	}

	// Respond with the updated book
	responseJSON, err := json.Marshal(updatedBookWithID)
	if err != nil {
		ctx.Error("Failed to marshal JSON", fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	ctx.Write(responseJSON)
}

// Delete ById
func BookDeleteById(ctx *fasthttp.RequestCtx) {
	idStr, ok := ctx.UserValue("id").(string)
	if !ok {
		ctx.Error("Invalid ID", fasthttp.StatusBadRequest)
		return
	}

	// Delete the book from MongoDB
	collection := mongoClient.Database(dbName).Collection("books")
	filter := bson.M{"id": idStr}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		ctx.Error("Error deleting data from MongoDB", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.WriteString("Delete successful")
}

// delete all
func BookDeleteAll(ctx *fasthttp.RequestCtx) {
	// Delete all documents from the MongoDB collection
	collection := mongoClient.Database(dbName).Collection("books")
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		ctx.Error("Error deleting data from MongoDB", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.WriteString("All documents deleted successfully")
}

func main() {
	//run database

	connectToMongo()

	defer mongoClient.Disconnect(context.Background())
	//-------call function------//
	r := router.New()

	//------EndPoint-----------//
	r.GET("/books", BookGetAll)
	r.GET("/books/{id}", BookGetByID)
	r.POST("/books", BookPost)
	r.PATCH("/books/{id}", BookPatch)
	r.DELETE("/books/{id}", BookDeleteById)
	r.DELETE("/books", BookDeleteAll)

	fasthttp.ListenAndServe(":3000", r.Handler)

}

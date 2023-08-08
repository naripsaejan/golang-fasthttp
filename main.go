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

var books = []Book{
	{ID: "1", Title: "Harry Potter", Author: "J. K. Rowling"},
	{ID: "2", Title: "The Lord of the Rings", Author: "J. R. R. Tolkien"},
	{ID: "3", Title: "The Wizard of Oz", Author: "L. Frank Baum"},
}

//-------------------- function -----------------------//

// check duplicate ID
func isBookIDUnique(id string) bool {
	for _, book := range books {
		if book.ID == id {
			return false
		}
	}
	return true
}

// Get
func BookGet(ctx *fasthttp.RequestCtx) {
	b, err := json.Marshal(books)
	if err != nil {
		log.Fatal(err)
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(b)
}

// Post
func BookPost(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("application/json")

	var newBook Book
	err := json.Unmarshal(ctx.Request.Body(), &newBook)
	if err != nil {
		ctx.Error("Invalid JSON", fasthttp.StatusBadRequest)
		return
	}

	if !isBookIDUnique(newBook.ID) {
		ctx.WriteString("Invalid add id")
		return
	}

	books = append(books, newBook)
	responseJSON, _ := json.Marshal(newBook)
	ctx.Write(responseJSON)
}

// Patch
func BookPatch(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("application/json")

	idStr := ctx.UserValue("id").(string)

	var updatedBook Book

	b := ctx.Request.Body()
	err := json.Unmarshal(b, &updatedBook)

	if err != nil {
		ctx.Error("Invalid JSON", fasthttp.StatusBadRequest)
		return
	}

	// Find the book to update by ID
	index := -1
	for i := range books {
		if books[i].ID == idStr {
			index = i
			break
		}
	}

	if index == -1 {
		ctx.Error("Book not found", fasthttp.StatusNotFound)
		return
	}

	// Update the book's fields
	if updatedBook.Title != "" {
		books[index].Title = updatedBook.Title
	}
	if updatedBook.Author != "" {
		books[index].Author = updatedBook.Author
	}

	responseJSON, err := json.Marshal(books[index])
	if err != nil {
		ctx.Error("Failed to marshal JSON", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Write(responseJSON)
}

// delete
func BookDelete(ctx *fasthttp.RequestCtx) {
	idStr, ok := ctx.UserValue("id").(string)
	if !ok {
		ctx.Error("Invalid ID", fasthttp.StatusBadRequest)
		return
	}

	index := -1
	for i, book := range books {
		if book.ID == idStr {
			index = i
			break
		}
	}

	if index == -1 {
		ctx.Error("No id", fasthttp.StatusNotFound)
		return
	}

	books = append(books[:index], books[index+1:]...)
	ctx.WriteString("Delete successful")
}

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

// --------test endoint mongodb-----------------------//
func requestHandler(ctx *fasthttp.RequestCtx) {
	// Handle your API requests here
	fmt.Fprintf(ctx, "Hello, this is your API!")

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

func handleCreateBook(ctx *fasthttp.RequestCtx) {
	// Parse request body
	var book Book
	if err := json.Unmarshal(ctx.PostBody(), &book); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString("Error parsing request body")
		return
	}

	// Insert book into MongoDB
	collection := mongoClient.Database(dbName).Collection("books")
	_, err := collection.InsertOne(context.Background(), book)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString("Error inserting data into MongoDB")
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.WriteString("Book created successfully")
}

func main() {
	//run database

	connectToMongo()

	defer mongoClient.Disconnect(context.Background())
	//-------call function------//
	r := router.New()

	//test get in db
	r.GET("/bookDB", requestHandler)
	r.POST("/bookDB", handleCreateBook)
	//get data mock
	r.GET("/books", BookGet)
	r.POST("/books", BookPost)
	r.PATCH("/books/{id}", BookPatch)
	r.DELETE("/books/{id}", BookDelete)

	fasthttp.ListenAndServe(":3000", r.Handler)

}

package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
)

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
	collection := MongoClient.Database(dbName).Collection("books")
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

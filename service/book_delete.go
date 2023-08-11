package service

import (
	"context"
	"examp/hello-fast-http/utils"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
)

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
	collection := utils.MongoClient.Database(utils.DbName).Collection("books")
	filter := bson.M{"id": idStr}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}

	ctx.Write([]byte(successMessage))
}


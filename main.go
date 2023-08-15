package main

import (
	//lib

	//utils
	"context"
	"examp/hello-fast-http/service"
	"examp/hello-fast-http/utils"

	// BookGetAll

	//lib

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

//-------------------- function -----------------------//

func main() {
	r := router.New()
		utils.ConnectToMongo()
	defer utils.MongoClient.Disconnect(context.Background())

	r.GET("/books/read", service.BookGetAll)
	r.GET("/books/readConsumer", service.BookGetAllConsumer)
	r.GET("/books/readConsumerGroup", service.BookGetAllConsumerGroup)
	r.GET("/books/read/{id}", service.BookGetByID)
	r.POST("/books/create", service.BookPost)
	r.PATCH("/books/update/{id}", service.BookPatch)
	r.DELETE("/books/delete/{id}", service.BookDeleteById)

	fasthttp.ListenAndServe(":3000", r.Handler)
}
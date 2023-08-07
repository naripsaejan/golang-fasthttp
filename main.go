package main

import (
	"encoding/json"
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{ID: "1", Title: "Harry Potter", Author: "J. K. Rowling"},
	{ID: "2", Title: "The Lord of the Rings", Author: "J. R. R. Tolkien"},
	{ID: "3", Title: "The Wizard of Oz", Author: "L. Frank Baum"},
}

//check duplicate ID
func isBookIDUnique(id string) bool {
	for _, book := range books {
		if book.ID == id {
			return false
		}
	}
	return true
}

//Get
func BookGet(ctx *fasthttp.RequestCtx) {
    b,err := json.Marshal(books)
    if err != nil {
        log.Fatal(err)
    }
    ctx.Response.Header.SetContentType("application/json")
    ctx.Write(b)
}

//Post
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

//Patch
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

//delete
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

func main() {
    r := router.New()

    r.GET("/books",BookGet)
    r.POST("/books",BookPost)
    r.PATCH("/books/{id}",BookPatch)
	r.DELETE("/books/{id}", BookDelete)

    fasthttp.ListenAndServe(":3000", r.Handler)
	
}


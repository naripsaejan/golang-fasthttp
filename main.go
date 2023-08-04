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

func BookGet(ctx *fasthttp.RequestCtx) {
    b,err := json.Marshal(books)
    if err != nil {
        log.Fatal(err)
    }
    ctx.Response.Header.SetContentType("application/json")
    ctx.Write(b)
}

func BookPost(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("application/json")

	var newBook Book
	err := json.Unmarshal(ctx.Request.Body(), &newBook)
	if err != nil {
		ctx.Error("Invalid JSON", fasthttp.StatusBadRequest)
		return
	}

	if !isBookIDUnique(newBook.ID) {
		ctx.WriteString("เพิ่มไม่ได้ id นี้มีอยู่แล้ว")
		return
	}

	books = append(books, newBook)
	responseJSON, _ := json.Marshal(newBook)
	ctx.Write(responseJSON)
}

//TODO:ปรับ
// func BookPatch(ctx *fasthttp.RequestCtx) {
//     ctx.Response.Header.SetContentType("application/json")

//     idStr, ok := ctx.UserValue("id").(string)
//     if !ok {
//         ctx.Error("Invalid ID", fasthttp.StatusBadRequest)
//         return
//     }

//     _, err := strconv.Atoi(idStr)
//     if err != nil {
//         ctx.Error("Invalid ID", fasthttp.StatusBadRequest)
//         return
//     }

//     var updatedBook Book
//     b := ctx.Request.Body()
//     err = json.Unmarshal(b, &updatedBook)
//     if err != nil {
//         ctx.Error("Invalid JSON", fasthttp.StatusBadRequest)
//         return
//     }

//     // Find the book to update by ID
//     index := -1
//     for i := range books {
//         if books[i].ID == idStr {
//             index = i
//             break
//         }
//     }

//     if index == -1 {
//         ctx.Error("Book not found", fasthttp.StatusNotFound)
//         return
//     }

//     // Preserve the existing ID
//     updatedBook.ID = books[index].ID

//     // Validate and update the book
//     if updatedBook.Title != "" && updatedBook.Author != "" {
//         books[index] = updatedBook
//         responseJSON, err := json.Marshal(updatedBook)
//         if err != nil {
//             ctx.Error("Failed to marshal JSON", fasthttp.StatusInternalServerError)
//             return
//         }

//         ctx.Write(responseJSON)
//     } else {
//         ctx.Error("Title and Author cannot be blank", fasthttp.StatusBadRequest)
//     }
// }




// func BookPatch(ctx *fasthttp.RequestCtx) {
//     ctx.Response.Header.SetContentType("application/json")

//     var book Book
//     id :=ctx.UserValue("id")
//     b := ctx.Request.Body()

//     json.Unmarshal(b,&book)

//     for i, a := range books {
//         if a.ID == id {
//             books = append(books[:i], book)
//             break
//         }
//     }
//     ctx.Write(b)
// }

//TODO:ปรับ
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
        ctx.Error("ไม่มี id ที่ต้องการลบ", fasthttp.StatusNotFound)
        return
    }

    books = append(books[:index], books[index+1:]...)
    ctx.WriteString("Delete successful")
}

func main() {
    r := router.New()

    r.GET("/books",BookGet)
    r.POST("/books",BookPost)
    // r.PUT("/books", BookPut)
    // r.PATCH("/books/{id}",BookPatch)
	r.DELETE("/books/{id}", BookDelete)

    fasthttp.ListenAndServe(":3000", r.Handler)
	
}


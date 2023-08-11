package service

// Book defines the structure of a book.
type Book struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

// Constants for error messages and status codes
const (
	invalidRequest  = `{"status":"40000","error":"Invalid request"}`
	internalError   = `{"status":"50000","error":"Internal server error"}`
	notFoundError   = `{"status":"40004","error":"Book not found"}`
	successMessage = `{"status":"20000","message":"Operation successful"}`
)

// book.go

package service

import "go.mongodb.org/mongo-driver/mongo"

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

// MongoDB and Kafka configuration constants
const (
	dbName     = "test"
	kafkaName  = "rip-test"
)

var (
	MongoClient *mongo.Client
	kafkaBrokers = []string{"10.138.41.195:9092", "10.138.41.196:9092"}
)
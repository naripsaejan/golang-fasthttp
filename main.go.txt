package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB connection URI
	mongoURI := "mongodb://test:password@10.138.41.195:27017,10.138.41.196:27017,10.138.41.197:27017/?authSource=test&replicaSet=nmgw"
	// mongoURI := "mongodb://localhost:27017/"

	// Set up a MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB to ensure connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	// Close the MongoDB connection when the application exits
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Define an HTTP handler
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		// You can perform MongoDB operations here using the 'client' instance
		// For example, fetching data from the "testdb" and "testcollection"
		// collection and returning it as a JSON response.
		// You'll need to import the appropriate packages and handle errors properly.
	})

	// Start the HTTP server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

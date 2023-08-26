package service

import (
	"context"
	"encoding/json"
	"examp/hello-fast-http/utils"
	"fmt"
	"strings"

	"log"

	"github.com/valyala/fasthttp"

	"go.mongodb.org/mongo-driver/bson"
)

// // data check support
var (
	// support4G = []string{"LTE1800", "LTE2100", "LTE2600", "LTE700", "LTE900"}
	support4G = []string{"GSM 900MHz", "GSM 1800MHz", "LTE2600", "LTE700", "GSM 850MHz"}
	support5G = []string{"GSM 900MHz", "GSM 1800MHz", "LTE2600", "LTE700", "GSM 850MHz"}
	// support5G = []string{"NR 2600", "NR 700"}
)

// BookGetAllHandler handles retrieving all books from the MongoDB collection.
func BookGetAll(ctx *fasthttp.RequestCtx) {
	// Set response content type
	ctx.SetContentType("application/json")

	// Retrieve data from MongoDB
	collection := utils.MongoClient.Database(utils.DbName).Collection("books")
	filter := bson.D{} // You can add filtering criteria here

	var results []Book
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var book Book
		if err := cur.Decode(&book); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte(internalError))
			return
		}
		results = append(results, book)
	}

	// Respond with data
	jsonBytes, _ := json.Marshal(results)
	ctx.Write(jsonBytes)

}

// GetById
func BookGetByID(ctx *fasthttp.RequestCtx) {
	// Set response content type
	ctx.SetContentType("application/json")

	idStr := ctx.UserValue("id").(string)
	log.Println("bookID", idStr)
	// Retrieve book by ID from MongoDB
	collection := utils.MongoClient.Database(utils.DbName).Collection("books")

	// filter := bson.M{"_id": idStr}
	filter := bson.M{"id": idStr}

	var book Book
	err := collection.FindOne(context.Background(), filter).Decode(&book)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write([]byte(notFoundError))
		return
	}

	// Respond with the retrieved book
	jsonBytes, _ := json.Marshal(book)
	ctx.Write(jsonBytes)
}

// ----------------get consumer -------------------//
// ConsumeKafkaMessages consumes Kafka messages.
func ConsumeKafkaMessages() {
	err := utils.ConsumeMessagesFromKafka(utils.KafkaBrokers, utils.KafkaTopics)
	if err != nil {
		log.Println("Error consuming Kafka messages:", err)
	}
}

func BookGetAllConsumer(ctx *fasthttp.RequestCtx) {
	// Set response content type
	ctx.SetContentType("application/json")

	// Retrieve data from MongoDB
	collection := utils.MongoClient.Database(utils.DbName).Collection("books")
	filter := bson.D{} // You can add filtering criteria here

	var results []Book
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var book Book
		if err := cur.Decode(&book); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte(internalError))
			return
		}
		results = append(results, book)
	}

	// Respond with data
	jsonBytes, _ := json.Marshal(results)
	ctx.Write(jsonBytes)

	// Consume Kafka messages in the background
	go func() {
		log.Println("get=> ConsumeMessagesFromKafka")

		err := utils.ConsumeMessagesFromKafka(utils.KafkaBrokers, utils.KafkaTopics)
		if err != nil {
			log.Println("Error consuming Kafka messages:", err)
		}
	}()

}

func BookGetAllConsumerGroup(ctx *fasthttp.RequestCtx) {
	// Set response content type
	ctx.SetContentType("application/json")

	// Retrieve data from MongoDB
	collection := utils.MongoClient.Database(utils.DbName).Collection("books")
	filter := bson.D{} // You can add filtering criteria here

	var results []Book
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var book Book
		if err := cur.Decode(&book); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte(internalError))
			return
		}
		results = append(results, book)
	}

	// Respond with data
	jsonBytes, _ := json.Marshal(results)
	ctx.Write(jsonBytes)

	// Consume Kafka messages in the background using a consumer group
	go func() {
		log.Println("get=> ConsumeMessagesWithConsumerGroup")

		err := utils.ConsumeMessagesWithConsumerGroup(utils.KafkaBrokers, utils.KafkaTopics, "rip-group1")
		if err != nil {
			log.Println("Error consuming Kafka messages:", err)
		}
	}()
}

//logic compatibility

// BookGetAllHandler handles retrieving all books from the MongoDB collection.
func BookGetSupport(ctx *fasthttp.RequestCtx) {
	gridRow := []Grid{
		{
			Name:       false,
			MatchFound: false,
			IsDevice4G: false,
			IsDevice5G: false,
		},
	}

	// Set response content type
	ctx.SetContentType("application/json")

	// Retrieve data from MongoDB
	collection := utils.MongoClient.Database(utils.DbName).Collection("ais")

	// Define the filter criteria
	//---------data one---------\\
	// var dataIn = "XDA Argon"

	// filter := bson.M{"name": dataIn}
	// filter := bson.M{"resourceSpecCharacteristic.resourceSpecCharacteristicValue.value": "HTC Panda"}
	//---------data two---------\\

	// filter := bson.M{"resourceSpecCharacteristic.resourceSpecCharacteristicValue.value": "HTC Panda"}
	var dataIn = "HMD Nova Pro"
	filter := bson.M{"resourceSpecCharacteristic.resourceSpecCharacteristicValue.value": dataIn}
	fmt.Println("filter", filter)
	fmt.Println("idStr", dataIn)

	var books []IoTDeviceSpecification

	// Find documents based on the filter
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(internalError))
		return
	}
	defer cur.Close(context.Background())

	// Iterate through the result cursor and decode each document
	for cur.Next(context.Background()) {
		var book IoTDeviceSpecification
		if err := cur.Decode(&book); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte(internalError))
			return
		}
		books = append(books, book)
	}

	// Respond with the retrieved book(s)
	fmt.Println("books", books)
	//check support
	fmt.Println("--------before--------------")

	// ceda.MapToStruct(resultMasterGrid, &gridData)
	for _, book := range books {
		//check by name
		if strings.Contains(book.Name, dataIn) {
			fmt.Println(" Row ---> Name")

			gridRow[0].Name = true

			//NOTE: check support 4G
			for _, data4G := range support4G {
				fmt.Println("check support support4G ------", support4G)
				if strings.Contains(book.ResourceSpecCharacteristic[3].ResourceSpecCharacteristicValue[0].Value, data4G) {
					gridRow[0].IsDevice4G = true
				}
			}

			//NOTE: check support 5G
			for _, data5G := range support5G {
				if strings.Contains(book.ResourceSpecCharacteristic[3].ResourceSpecCharacteristicValue[0].Value, data5G) {
					gridRow[0].IsDevice5G = true

				}
			}
			break

		}

		//check by value
		if strings.Contains(book.ResourceSpecCharacteristic[0].ResourceSpecCharacteristicValue[0].Value, dataIn) {
			fmt.Println(" Row ---> ResourceSpecCharacteristic")

			// if strings.Contains(book.ResourceSpecCharacteristic[0].ResourceSpecCharacteristicValue[0].Value, dataIn) {
			gridRow[0].Value = true

			//NOTE: check support 4G
			for _, data4G := range support4G {
				fmt.Println("check support support4G ------", support4G)
				if strings.Contains(book.ResourceSpecCharacteristic[3].ResourceSpecCharacteristicValue[0].Value, data4G) {
					gridRow[0].IsDevice4G = true
				}
			}

			//NOTE: check support 5G
			for _, data5G := range support5G {
				if strings.Contains(book.ResourceSpecCharacteristic[3].ResourceSpecCharacteristicValue[0].Value, data5G) {
					gridRow[0].IsDevice5G = true

				}
			}
			break
		}

		// fmt.Println("for---> %+v", gridRow)
	}
	fmt.Printf("\n for---> %+v", gridRow)

	// if strings.Contains(dataIn, books ) {
	// 	// 	fmt.Println("resultMasterGrid ------", gridData.Name)
	// 	// 	gridRow[0].MatchFound = true
	// 	// 	for _, char := range gridData.ResourceSpecCharacteristic {

	// 	// 		//NOTE: check support 4G
	// 	// 		for _, data4G := range support4G {
	// 	// 			fmt.Println("check support support4G ------", support4G)
	// 	// 			if strings.Contains(char.ResourceSpecCharacteristicValue[0].Value, data4G) {
	// 	// 				gridRow[0].IsDevice4G = true
	// 	// 			}
	// 	// 		}

	// 	// 		//NOTE: check support 5G
	// 	// 		for _, data5G := range support5G {
	// 	// 			if strings.Contains(char.ResourceSpecCharacteristicValue[0].Value, data5G) {
	// 	// 				gridRow[0].IsDevice5G = true

	// 	// 			}
	// 	// 		}
	// 	// 	}

	// }

	jsonBytes, _ := json.Marshal(books)
	ctx.Write(jsonBytes)
}

package utils

import "go.mongodb.org/mongo-driver/mongo"

const (
	DbName   = "test"
	MongoURI = "mongodb://test:password@10.138.41.195:27017,10.138.41.196:27017,10.138.41.197:27017/?authSource=test&replicaSet=nmgw"
)

var (
	MongoClient  *mongo.Client
	KafkaTopics  = "rip-test"
	KafkaBrokers = []string{"10.138.41.195:9092", "10.138.41.196:9092"}
)
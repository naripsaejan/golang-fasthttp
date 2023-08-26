package utils

import "go.mongodb.org/mongo-driver/mongo"

// KafkaMessageHandler คือตัวจัดการข้อความของ Kafka
type KafkaMessageHandler struct{}

const (
	DbName   = "test" //คอเลคชั่น
	// MongoURI = "mongodb://test:password@10.138.41.195:27017,10.138.41.196:27017,10.138.41.197:27017/?authSource=test&replicaSet=nmgw"
	MongoURI = "mongodb+srv://admin:admin@cluster0.tiw82mg.mongodb.net/?retryWrites=true&w=majority"

)

var (
	MongoClient  *mongo.Client
	KafkaTopics  = "rip-test"
	GroupID = "rip-group"
	KafkaBrokers = []string{"10.138.41.195:9092", "10.138.41.196:9092"}
)
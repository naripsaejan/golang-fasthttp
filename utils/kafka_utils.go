package utils

import (
	"log"

	"github.com/IBM/sarama" // Import the correct package
)

// InitializeKafkaProducer initializes and returns a Kafka producer instance.
func InitializeKafkaProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Println("Failed to create Kafka producer:", err)
		return nil, err
	}

	return producer, nil
}

// SendMessageToKafka sends a message to the specified Kafka topic.
func SendMessageToKafka(producer sarama.SyncProducer, topic string, key string, value []byte) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	_, _, err := producer.SendMessage(message)
	if err != nil {
		log.Println("Failed to send Kafka message:", err)
		return err
	}

	return nil
}

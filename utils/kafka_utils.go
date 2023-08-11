package utils

import (
	"log"

	"github.com/IBM/sarama" // Import the correct package
)

// InitializeKafkaProducer initializes and returns a Kafka producer instance.
func InitializeKafkaProducer(brokers []string) (sarama.SyncProducer, error) {
			log.Println(("----------------InitializeKafkaProducer----------------"))

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
			log.Println(("----------------SendMessageToKafka----------------"))

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


// ConsumeMessagesFromKafka consumes messages from the specified Kafka topic.
func ConsumeMessagesFromKafka(brokers []string, topics []string) error {
			log.Println(("----------------ConsumeMessagesFromKafka----------------"))
	// log.Println("check auto topics",topics)
		consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return err
	}
	defer consumer.Close()

	// Create a partition consumer for each topic and partition
	for _, topic := range topics {
		partitions, err := consumer.Partitions(topic)
		if err != nil {
			return err
		}

		for _, partition := range partitions {
			partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				return err
			}

			go func() {
				for msg := range partitionConsumer.Messages() {
					log.Printf("Received Kafka message: Topic - %s, Partition - %d, Offset - %d, Key - %s, Value - %s\n",
						msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
					// Process the Kafka message here if needed
				}
			}()
		}
	}

	select {} // Keep the goroutine running
}

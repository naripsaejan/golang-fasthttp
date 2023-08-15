package utils

import (
	"context"
	"log"
	"sync"

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

//-------------------------consumer-------------------------------//

// ConsumeMessagesFromKafka consumes messages from the specified Kafka topic.
func ConsumeMessagesFromKafka(brokers []string, topic string) error {
	log.Println("---------------ConsumeMessagesFromKafka---------------")
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return err
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()

			partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Println("Error consuming partition:", err)
				return
			}

			for msg := range partitionConsumer.Messages() {
				log.Printf("Received Kafka message: Topic - %s, Partition - %d, Offset - %d, Key - %s, Value - %s\n",
					msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				// Process the Kafka message here if needed
			}
		}(partition)
	}

	wg.Wait()

	return nil
}

//-------------------------consumer group-------------------------------//

// ConsumeMessagesWithConsumerGroup consumes messages from the specified Kafka topic using a consumer group.
func ConsumeMessagesWithConsumerGroup(brokers []string, topic string, groupID string) error {
	log.Println("---------------ConsumeMessagesWithConsumerGroup---------------")
	config := sarama.NewConfig()
config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return err
	}
	defer consumerGroup.Close()

	ctx := context.Background()
	handler := &ConsumerGroupHandler{}

	for {
		err := consumerGroup.Consume(ctx, []string{topic}, handler)
		if err != nil {
			log.Println("Error from consumer group:", err)
		}
	}
}

// ConsumerGroupHandler implements the sarama.ConsumerGroupHandler interface.
type ConsumerGroupHandler struct{}

// Setup is called at the beginning of a new session.
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is called at the end of a session.
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim is called when a claim is available for processing.
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Received Kafka message: Topic - %s, Partition - %d, Offset - %d, Key - %s, Value - %s\n",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
		// Process the Kafka message here if needed

		session.MarkMessage(message, "")
	}

	return nil
}

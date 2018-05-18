package common

import (
	"log"
	"sync"
	"time"

	sarama "github.com/Shopify/sarama"
)

// KafkaProducer struct definition
type MessageProducerKafka struct {
	producer sarama.AsyncProducer
	wg       *sync.WaitGroup
}

func NewMessageProducerKafka(brokers []string) (*MessageProducerKafka, error) {
	// enable errors and notifications
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	wg := new(sync.WaitGroup)

	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		for err := range producer.Errors() {
			log.Printf("Failed to write message: %v", err)
		}
	}(wg)

	return &MessageProducerKafka{producer, wg}, nil
}

func (m *MessageProducerKafka) Close() error {
	defer m.wg.Wait()
	return m.producer.Close()
}

func (m *MessageProducerKafka) Publish(topic, partitionKey string, message []byte) error {
	msg := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(partitionKey),
		Value: sarama.ByteEncoder(message),
	}

	m.producer.Input() <- &msg
	return nil
}

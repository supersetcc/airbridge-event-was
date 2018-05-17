package common

type MessageProducer interface {
	Publish(topic, pk string, payload []byte) error
	Close() error
}

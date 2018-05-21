package common

type MessageProducerMock struct {
	IsClosed                  bool
	LastPublishedTopic        string
	LastPublishedPartitionKey string
	LastPublishedPayload      []byte
}

func (p *MessageProducerMock) Publish(topic, pk string, payload []byte) error {
	p.LastPublishedTopic = topic
	p.LastPublishedPartitionKey = pk
	p.LastPublishedPayload = payload
	return nil
}

func (p *MessageProducerMock) Close() error {
	p.IsClosed = true
	return nil
}

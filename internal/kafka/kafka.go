package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/model"
)

type (
	Producer struct {
		Topic    string
		Producer sarama.AsyncProducer
	}
)

func New(brokers []string, topic string) *Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 10

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatal("e: ", err)
	}

	return &Producer{
		Topic:    topic,
		Producer: producer,
	}
}

func (c *Producer) ProduceSamples(samples model.Samples) {
	for _, sample := range samples {
		fingerprint := EncodeFingerprint(sample.Metric.FastFingerprint())
		encodedSample := EncodeSample(sample)
		c.Producer.Input() <- &sarama.ProducerMessage{
			Topic: c.Topic,
			Key:   sarama.ByteEncoder(fingerprint),
			Value: sarama.ByteEncoder(encodedSample),
		}
	}
}

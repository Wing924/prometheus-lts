package kafka

import (
	"log"
	"testing"

	"github.com/prometheus/common/model"
)

func TestProducer_ProduceSamples(t *testing.T) {
	samples := model.Samples{
		&model.Sample{model.Metric{"__name__": "foo"}, 1, 0},
		&model.Sample{model.Metric{"__name__": "bar"}, 2, 0},
	}
	producer := NewProducer([]string{"localhost:9092"}, "topic")
	producer.ProduceSamples(samples)
	producer.AsyncClose()
	for i := 0; i < len(samples); i++ {
		select {
		case msg := <-producer.Successes():
			log.Println("success: ", msg)
		case err := <-producer.Errors():
			log.Println("error: ", err)
		}
	}
}

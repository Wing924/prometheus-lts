package writer

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Wing924/prometheus-lts/internal/kafka"

	"github.com/prometheus/common/model"

	"github.com/Wing924/prometheus-lts/internal/prom"
	"github.com/prometheus/prometheus/prompb"
)

type Writer struct {
	Kafka   kafka.Producer
	metrics *metrics
}

func New(k kafka.Producer, r prometheus.Registerer) *Writer {
	w := &Writer{
		Kafka:   k,
		metrics: newMetrics(r),
	}
	go func() {
		for {
			select {
			case <-w.Kafka.Producer.Successes():
				w.metrics.sentSamples.Add(1)
			case err := <-w.Kafka.Producer.Errors():
				log.Printf("e: %v", err.Err)
				w.metrics.failedSamples.Add(1)
			}
		}
	}()
	return w
}

func (w *Writer) RegisterHandler(mux *http.ServeMux) {
	mux.HandleFunc("/receive", func(resp http.ResponseWriter, req *http.Request) {
		begin := time.Now()
		w.receiveHandler(resp, req)
		duration := time.Since(begin).Seconds()
		w.metrics.sentBatchDuration.Observe(duration)
	})
}

func (w *Writer) receiveHandler(resp http.ResponseWriter, req *http.Request) {
	body, err := readBody(req)
	if err != nil {
		log.Print("e: Prom write handler unable to read bytes from request body")
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	wr, err := prom.DecodeWriteRequest(body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	samples := protoToSamples(wr)
	w.metrics.receivedSamples.Add(float64(len(samples)))

	w.Kafka.ProduceSamples(samples)
}

func readBody(r *http.Request) ([]byte, error) {
	var bs []byte
	if r.ContentLength > 0 {
		bs = make([]byte, 0, r.ContentLength)
	}
	buf := bytes.NewBuffer(bs)
	_, err := buf.ReadFrom(r.Body)
	return buf.Bytes(), err
}

func protoToSamples(req prompb.WriteRequest) model.Samples {
	var samples model.Samples
	for _, ts := range req.Timeseries {
		metric := make(model.Metric, len(ts.Labels))
		for _, l := range ts.Labels {
			metric[model.LabelName(l.Name)] = model.LabelValue(l.Value)
		}

		for _, s := range ts.Samples {
			samples = append(samples, &model.Sample{
				Metric:    metric,
				Value:     model.SampleValue(s.Value),
				Timestamp: model.Time(s.Timestamp),
			})
		}
	}
	return samples
}

//
//func rehashSamples(samples model.Samples, partition int64) []model.Samples {
//	partitions := make([]model.Samples, partition)
//
//	for _, sample := range samples {
//		pos := hashMetrics(sample) % partition
//		partitions[pos] = append(partitions[pos], sample)
//	}
//
//	return partitions
//}
//
//func hashMetrics(sample *model.Sample) int64 {
//	return int64(sample.Metric.FastFingerprint())
//}

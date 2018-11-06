package writer

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	receivedSamples   prometheus.Counter
	sentSamples       prometheus.Counter
	failedSamples     prometheus.Counter
	sentBatchDuration prometheus.Histogram
}

const namespace = "prom_lts_writer"

func newMetrics(r prometheus.Registerer) *metrics {
	var m metrics
	m.receivedSamples = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "received_samples_total",
			Help:      "Total number of received samples.",
		},
	)
	m.sentSamples = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "sent_samples_total",
			Help:      "Total number of processed samples sent to remote storage.",
		},
	)
	m.failedSamples = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "failed_samples_total",
			Help:      "Total number of processed samples which failed on send to remote storage.",
		},
	)
	m.sentBatchDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "sent_batch_duration_seconds",
			Help:      "Duration of sample batch send calls to the remote storage.",
			Buckets:   prometheus.DefBuckets,
		},
	)
	if r != nil {
		r.MustRegister(
			m.receivedSamples,
			m.sentSamples,
			m.failedSamples,
			m.sentBatchDuration,
		)
	}
	return &m
}

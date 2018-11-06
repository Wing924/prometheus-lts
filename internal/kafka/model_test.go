package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prometheus/common/model"
)

var tests = []struct {
	name   string
	sample model.Sample
}{
	{`{"metric":{},"value":[0,"123"]}`,
		model.Sample{Metric: model.Metric{}, Value: 123, Timestamp: 0}},
	{`{"metric":{"__name__":"foo"},"value":[0,"123"]}`,
		model.Sample{Metric: model.Metric{"__name__": "foo"}, Value: 123, Timestamp: 0}},
	{`{"metric":{"__name__":"foo","key1":"value1"},"value":[0,"123"]}`,
		model.Sample{Metric: model.Metric{"__name__": "foo", "key1": "value1"}, Value: 123, Timestamp: 0}},
	{`{"metric":{"__name__":"foo","key1":"a b","key2":"a\"b"},"value":[0,"123"]}`,
		model.Sample{Metric: model.Metric{"__name__": "foo", "key1": "a b", "key2": `a"b`}, Value: 123, Timestamp: 0}},
	{`{"metric":{"key1":"a","key2":"b"},"value":[0,"123"]}`,
		model.Sample{Metric: model.Metric{"key1": "a", "key2": `b`}, Value: 123, Timestamp: 0}},
}

func TestEncodeSample(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.name, string(EncodeSample(&test.sample)))
		})
	}
}

func TestDecodeSample(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded := EncodeSample(&test.sample)
			decoded, err := DecodeSample(encoded)
			assert.NoError(t, err)
			assert.Equal(t, test.sample, *decoded)
		})
	}
}

func TestDecodeFingerprint(t *testing.T) {
	tests := []struct {
		name        string
		fingerprint model.Fingerprint
	}{
		{"0x0", 0},
		{"0x123456", 0x123456},
		{"0x123456789", 0x123456789},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded := EncodeFingerprint(test.fingerprint)
			decoded := DecodeFingerprint(encoded)
			assert.Equal(t, test.fingerprint, decoded)
		})
	}
}

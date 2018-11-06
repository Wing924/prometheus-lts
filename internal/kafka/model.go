package kafka

import (
	"encoding/binary"
	"encoding/json"

	"github.com/prometheus/common/model"
)

type (
	Sample = model.Sample
)

func EncodeSample(s *model.Sample) []byte {
	b, _ := json.Marshal(s)
	return b
}

func DecodeSample(b []byte) (*model.Sample, error) {
	sample := &model.Sample{}
	err := json.Unmarshal(b, sample)
	return sample, err
}

func EncodeFingerprint(fingerprint model.Fingerprint) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(fingerprint))
	return buf
}

func DecodeFingerprint(b []byte) model.Fingerprint {
	return model.Fingerprint(binary.LittleEndian.Uint64(b))
}

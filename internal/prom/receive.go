package prom

import (
	"log"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

func DecodeWriteRequest(in []byte) (prompb.WriteRequest, error) {
	var req prompb.WriteRequest
	reqBuf, err := snappy.Decode(nil, in)
	if err != nil {
		log.Print("e: ", err.Error())
		return req, err
	}

	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		log.Print("e: ", err.Error())
		return req, err
	}

	return req, nil
}

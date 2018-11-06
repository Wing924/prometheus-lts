package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

var (
	addr = flag.String("listen-address", ":18086", "The address to listen on for HTTP requests.")
)

func init() {
	prometheus.MustRegister(version.NewCollector("prom_lts_writer"))
}

func main() {
	if os.Getenv("DEBUG") != "" {
		runtime.SetBlockProfileRate(20)
		runtime.SetMutexProfileFraction(20)
	}

	//writer.RegisterHandler()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal("e: ", http.ListenAndServe(*addr, nil))
}

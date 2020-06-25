package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	hDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "www_method_monitor_duration",
		Help: "www_method_monitor help",
	}, []string{"code", "method", "name"})
	// hCount := prometheus.NewCounterVec(prometheus.CounterOpts{
	// 	Name: "www_method_monitor_count",
	// 	Help: "www_method_monitor help",
	// }, []string{"code", "method", "name"})

	reg := prometheus.NewRegistry()
	reg.MustRegister(hDuration)
	// reg.MustRegister(hCount)

	prom := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	hello := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(201)
		w.Write([]byte("i'm your hello world web"))
	})

	m := http.NewServeMux()
	m.HandleFunc("/metrics", basic(
		os.Getenv("WWW_AUTH_BASIC__USERNAME"),
		os.Getenv("WWW_AUTH_BASIC__PASSWORD"),
		prom,
	))

	const HOME = "/"
	m.HandleFunc(HOME, func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		wrapper := responseWriter{ResponseWriter: w}

		hello(&wrapper, req)
		end := time.Since(start)
		hDuration.WithLabelValues(
			strconv.Itoa(wrapper.status), req.Method, req.URL.EscapedPath()).
			Observe(end.Seconds())
		// hCount.WithLabelValues(
		// 	strconv.Itoa(200), req.Method, req.URL.EscapedPath()).
		// 	Inc()

	})
	log.Fatalln(http.ListenAndServe(":8080", m))
}

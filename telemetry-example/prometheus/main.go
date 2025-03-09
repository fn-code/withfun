package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fn-code/withfun/telemetry-example/prometheus/handler"
	"github.com/fn-code/withfun/telemetry-example/prometheus/metric"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	mux := mux.NewRouter()
	metricsBuilder := metric.NewMetricsBuilder("metapp", "metapp_", "1.0", "1q2w3e4r321", time.Now().String())

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		res := metric.NewLoggingResponseWriter(w)
		metricsBuilder.Instrument(res, handler.UserHandler, "/user")(res, r)
	})

	mux.HandleFunc("/sub", func(w http.ResponseWriter, r *http.Request) {
		res := metric.NewLoggingResponseWriter(w)
		metricsBuilder.Instrument(res, handler.SubHandler, "/sub")(res, r)
	})

	hn := http.HandlerFunc(promhttp.Handler().ServeHTTP)

	// http.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/metrics", handler.BasicAuth(hn))

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

}

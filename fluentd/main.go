package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

type SimpleHandler struct {
	logger *zap.Logger
}

func (s *SimpleHandler) helloHandler(w http.ResponseWriter, r *http.Request) {
	with := r.URL.Query().Get("with")
	if with == "err" {
		s.logger.Error(fmt.Sprintf("helloHandler: %s\n", "error accessing data"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	s.logger.Info(fmt.Sprintf("helloHandler: %s", hostname))
	_, err = fmt.Fprintf(w, "Hello from host %s", hostname)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

}

func main() {

	port := flag.Int("port", 8080, "-port and then server port")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer logger.Sync() // flushes buffer, if any

	smplHandler := &SimpleHandler{logger: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("/", smplHandler.helloHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
		Handler: mux,
	}

	logger.Info(fmt.Sprintf("Server running on port: %d\n", *port))
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

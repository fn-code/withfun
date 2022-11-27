package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	log.Printf("helloHandler: %s\n", hostname)
	_, err = fmt.Fprintf(w, "Hello from host %s", hostname)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

}

func main() {

	port := flag.Int("port", 8080, "-port and then server port")

	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
		Handler: mux,
	}

	log.Printf("Server running on port: %d\n", *port)
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

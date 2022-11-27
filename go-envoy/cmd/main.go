package main

import (
	"fmt"
	"github.com/fn-code/go-envoy-example/handler"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	args := os.Args
	if len(args) < 3 {
		log.Fatal("invalid arguments")
	}

	p, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	h := &handler.Server{
		Name: args[2],
		Port: p,
	}

	mux.Handle("/", h)

	log.Printf("server running on port:%d", p)
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", p),
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

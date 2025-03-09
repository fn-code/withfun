package main

import (
	"github.com/fn-code/swagger-example/data"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	h := mux.NewRouter()

	h.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	h.HandleFunc("/v1/api/data", data.GetData).Methods("GET")


	// documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	h.Handle("/docs", sh)

	// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "/swagger.yaml", Path: "docs1"}
	sh1 := middleware.Redoc(opts1, nil)
	h.Handle("/docs1", sh1)

	if err := http.ListenAndServe(":8080", h); err != nil {
		panic(err)
	}
}
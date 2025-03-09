package main

import (
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/fn-code/telemetry-example/open_telemetry/userservice"
)

func main() {
	restful.Add(userservice.New())

	log.Fatal(http.ListenAndServe(":8080", nil))
}

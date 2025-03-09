package main

import (
	"fmt"
	"go.elastic.co/apm/module/apmhttp/v2"
	"log"
	"net/http"
	"net/url"
	"os"

	"go.elastic.co/apm/transport"
	"go.elastic.co/apm/v2"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "welcome")
}

func New(serviceName string, serviceURL string, apmEnvironment string) (*apm.Tracer, error) {
	apm.DefaultTracer().Close()

	u, _ := url.Parse(serviceURL)
	httpTransport, err := setTransport(u)
	if err != nil {
		return nil, err
	}

	tracer, err := apm.NewTracerOptions(apm.TracerOptions{
		ServiceName:        serviceName,
		ServiceVersion:     "1",
		ServiceEnvironment: apmEnvironment,
		Transport:          httpTransport,
	})
	if err != nil {
		return nil, err
	}

	return tracer, nil

}

func setTransport(url2 *url.URL) (*transport.HTTPTransport, error) {
	httpTransport, err := transport.NewHTTPTransport()
	if err != nil {
		return nil, err
	}

	httpTransport.SetServerURL(url2)
	return httpTransport, nil
}

func main() {
	wt, err := New("elastic-testing", "http://localhost:8200", "staging")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", index)

	opts := apmhttp.WithTracer(wt)

	if err := http.ListenAndServe(":8080", apmhttp.Wrap(mux, opts)); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}

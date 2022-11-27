package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Song struct {
	ID     string `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Artist string `json:"artist,omitempty"`
}

type Tracer struct {
	*sdktrace.TracerProvider
}

func main() {
	mux := http.NewServeMux()

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("tikadesapi-song"),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String("abcdef12345"),
		attribute.String("environment", "demo"),
	)

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		log.Println(err)
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resources),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
	)

	tr := &Tracer{provider}
	mux.HandleFunc("/songs/", tr.songHandler)

	if err := http.ListenAndServe(":9002", mux); err != nil {
		tr.Shutdown(context.Background())
		log.Fatal(err)
	}

}

func (tr *Tracer) songHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	log.Println(id)
	tc := propagation.TraceContext{}
	ctx := tc.Extract(r.Context(), HTTPCarier(r.Header))

	ctx, span := tr.Tracer("service/songs/songHandler").Start(ctx, "songHandler")
	defer span.End()

	res := tr.getSongDB(ctx, id)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
		return
	}

}

func (tr *Tracer) getSongDB(ctx context.Context, id string) []Song {
	tc := tr.Tracer("service/songs/getSongDB")

	_, span := tc.Start(ctx, "getSongDB")
	defer span.End()

	var songs = map[string][]Song{
		"001": []Song{
			Song{ID: "001", Title: "Easy On Me", Artist: "Adele"},
			Song{ID: "002", Title: "Jodie", Artist: "Ben Flat"},
		},
		"002": []Song{
			Song{ID: "001", Title: "Easy On Me", Artist: "Adele"},
			Song{ID: "002", Title: "Jodie", Artist: "Ben Flat"},
		},
		"003": []Song{
			Song{ID: "001", Title: "Easy On Me", Artist: "Adele"},
		},
	}

	val, ok := songs[id]
	if !ok {
		return nil
	}

	return val
}

type HTTPCarier http.Header

// Get returns the value associated with the passed key.
func (c HTTPCarier) Get(key string) string {
	h := http.Header(c)
	return h.Get(key)
}

// DO NOT CHANGE: any modification will not be backwards compatible and
// must never be done outside of a new major release.

// Set stores the key-value pair.
func (c HTTPCarier) Set(key string, value string) {
	h := http.Header(c)
	h.Set(key, value)
}

// DO NOT CHANGE: any modification will not be backwards compatible and
// must never be done outside of a new major release.

// Keys lists the keys stored in this carrier.
func (c HTTPCarier) Keys() []string {
	keys := make([]string, 0)
	for k, _ := range c {
		keys = append(keys, k)
	}
	return keys
}

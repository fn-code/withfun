package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type User struct {
	ID       string     `json:"id,omitempty"`
	Username string     `json:"username,omitempty"`
	Playlist []Playlist `json:"playlist,omitempty"`
}

type Playlist struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Song []Song `json:"song,omitempty"`
}

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
		semconv.ServiceNameKey.String("tikadesapi-user"),
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

	mux.HandleFunc("/user", tr.userHandler)
	mux.HandleFunc("/sub", tr.subHandler)

	if err := http.ListenAndServe(":9000", mux); err != nil {
		tr.Shutdown(context.Background())
		log.Fatal(err)
	}
}

func (tr *Tracer) subHandler(w http.ResponseWriter, r *http.Request) {

	ctx, span := tr.Tracer("service.subscribe").Start(context.Background(), "subHandler")
	defer span.End()

	res := tr.getSubDB(ctx)

	span.RecordError(fmt.Errorf("error skali"))

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	err := enc.Encode(res)
	if err != nil {
		log.Println(err)
		return
	}

}

func (tr *Tracer) getSubDB(ctx context.Context) []*User {
	ts := tr.Tracer("service.user.getSubDB")

	_, span := ts.Start(ctx, "getSubDB")
	defer span.End()

	time.Sleep(5 * time.Second)
	return []*User{
		{ID: "001", Username: "ludinnento"},
	}
}

func (tr *Tracer) userHandler(w http.ResponseWriter, r *http.Request) {

	jenis := r.URL.Query().Get("j")

	ctx, span := tr.Tracer("service.user").Start(context.Background(), "userHandler")
	defer span.End()

	res := tr.getUserDB(ctx)

	tc := propagation.TraceContext{}
	for _, val := range res {
		urls := ""

		if jenis == "" || jenis == "all" {
			urls = "http://localhost:9001/playlist/" + val.ID
		} else if jenis == "fav" {
			urls = "http://localhost:9001/playlist/favorite/" + val.ID
		}

		req, err := http.NewRequest("GET", urls, nil)
		if err != nil {
			panic(err)
		}

		tc.Inject(ctx, HTTPCarier(req.Header))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		listPlaylist := []Playlist{}
		if err := json.NewDecoder(resp.Body).Decode(&listPlaylist); err != nil {
			log.Println(err)
		}

		if val.Playlist == nil {
			val.Playlist = make([]Playlist, 0)
		}

		val.Playlist = listPlaylist
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	err := enc.Encode(res)
	if err != nil {
		log.Println(err)
		return
	}
}

func (tr *Tracer) getUserDB(ctx context.Context) []*User {
	ts := tr.Tracer("service.user.get_user_db")

	_, span := ts.Start(ctx, "getUserDB")
	defer span.End()

	return []*User{
		{ID: "001", Username: "ludinnento"},
	}
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

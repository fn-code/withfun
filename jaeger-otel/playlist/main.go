package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Song struct {
	ID     string `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Artist string `json:"artist,omitempty"`
}

type Playlist struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Song []Song `json:"song,omitempty"`
}

type Tracer struct {
	*sdktrace.TracerProvider
}

func main() {
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("tikadesapi-playlist"),
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

	mux := http.NewServeMux()
	mux.HandleFunc("/playlist/", tr.playlishHandler)
	mux.HandleFunc("/playlist/favorite/", tr.playlishFavoriteHandler)

	if err := http.ListenAndServe(":9001", mux); err != nil {
		tr.Shutdown(context.Background())
		log.Fatal(err)
	}

}

func (tr *Tracer) playlishHandler(w http.ResponseWriter, r *http.Request) {

	id := filepath.Base(r.URL.Path)

	log.Println("all : ", id)
	var span trace.Span

	tc := propagation.TraceContext{}
	ctx := tc.Extract(r.Context(), HTTPCarier(r.Header))

	ctx, span = tr.Tracer("service.playlist").Start(ctx, "playlistHandler", trace.WithAttributes(
		attribute.String("http.host", r.Host),
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.Path),
	))

	defer span.End()

	res := tr.getPlaylist(ctx, id)

	for _, val := range res {
		req, err := http.NewRequest("GET", "http://localhost:9002/songs/"+val.ID, nil)
		if err != nil {
			panic(err)
		}

		tc.Inject(ctx, HTTPCarier(req.Header))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		listSongs := []Song{}
		if err := json.NewDecoder(resp.Body).Decode(&listSongs); err != nil {
			log.Println(err)
		}

		if val.Song == nil {
			val.Song = make([]Song, 0)
		}

		val.Song = listSongs
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	err := enc.Encode(res)
	if err != nil {

		log.Println(err)
		return
	}

}

func (tr *Tracer) playlishFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	log.Println("fav : ", id)
	var span trace.Span

	tc := propagation.TraceContext{}
	ctx := tc.Extract(r.Context(), HTTPCarier(r.Header))

	ctx, span = tr.Tracer("service.playlist.favorite").Start(ctx, "playlishFavoriteHandler", trace.WithAttributes(
		attribute.String("http.host", r.Host),
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.Path),
	))

	defer span.End()

	res := tr.getPlaylistFavorite(ctx, id)

	for _, val := range res {
		req, err := http.NewRequest("GET", "http://localhost:9002/songs/"+val.ID, nil)
		if err != nil {
			panic(err)
		}

		tc.Inject(ctx, HTTPCarier(req.Header))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		listSongs := []Song{}
		if err := json.NewDecoder(resp.Body).Decode(&listSongs); err != nil {
			log.Println(err)
		}

		if val.Song == nil {
			val.Song = make([]Song, 0)
		}

		val.Song = listSongs
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	err := enc.Encode(res)
	if err != nil {
		log.Println(err)
		return
	}

}

func (tr *Tracer) getPlaylist(ctx context.Context, id string) []*Playlist {
	ts := tr.Tracer("service.playlist.get_playlist")

	_, span := ts.Start(ctx, "getPlaylist")
	defer span.End()

	time.Sleep(3 * time.Second)
	playlist := []*Playlist{
		{
			ID:   "001",
			Name: "POP",
		},
		{
			ID:   "002",
			Name: "ROCK",
		},
	}

	var playlists = map[string][]*Playlist{
		"001": playlist,
	}

	val, ok := playlists[id]
	if !ok {
		return nil
	}

	return val
}

func (tr *Tracer) getPlaylistFavorite(ctx context.Context, id string) []*Playlist {
	ts := tr.Tracer("service.playlistFavorite.get_playlist_favorite")

	_, span := ts.Start(ctx, "getPlaylistFavorite")
	defer span.End()

	time.Sleep(3 * time.Second)
	playlist := []*Playlist{
		{
			ID:   "001",
			Name: "POP",
		},
		{
			ID:   "002",
			Name: "ROCK",
		},
		{
			ID:   "003",
			Name: "METAL",
		},
	}

	var playlists = map[string][]*Playlist{
		"001": playlist,
	}

	val, ok := playlists[id]
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

package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
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

func SubHandler(w http.ResponseWriter, r *http.Request) {

	res := getSubDB(context.TODO())

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	err := enc.Encode(res)
	if err != nil {
		log.Println(err)
		return
	}

}

func getSubDB(ctx context.Context) []*User {

	time.Sleep(3 * time.Second)
	return []*User{
		{ID: "001", Username: "ludinnento"},
		{ID: "002", Username: "ludinnento 2"},
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {

	res := getUserDB(context.TODO())

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	err := enc.Encode(res)
	if err != nil {
		log.Println(err)
		return
	}
}

func getUserDB(ctx context.Context) []*User {
	return []*User{
		{ID: "001", Username: "ludinnento"},
		{ID: "002", Username: "ludinnento 2"},
	}
}

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the username and password from the request
		// Authorization header. If no Authentication header is present
		// or the header value is invalid, then the 'ok' return value
		// will be false.
		username, password, ok := r.BasicAuth()
		if ok {

			expectedUsername := "admin"
			expectedPassword := "admin"

			// If the username and password are correct, then call
			// the next handler in the chain. Make sure to return
			// afterwards, so that none of the code below is run.
			if username == expectedUsername && password == expectedPassword {
				next.ServeHTTP(w, r)
				return
			}
		}

		// If the Authentication header is not present, is invalid, or the
		// username or password is wrong, then set a WWW-Authenticate
		// header to inform the client that we expect them to use basic
		// authentication and send a 401 Unauthorized response.
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

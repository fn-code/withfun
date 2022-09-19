package handler

import (
	"fmt"
	"net/http"
)

type Server struct {
	Name string
	Port int
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from port %s %d \n", s.Name, s.Port)
}
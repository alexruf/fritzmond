package ui

import (
	"net/http"
)

func (s *server) RegisterRoutes() {
	http.HandleFunc("/", s.handleGetIndex())
}

package server

import "net/http"

func (s *Server) handleIndexView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))

}

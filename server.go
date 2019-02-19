package main

import (
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"net/http"
)

type server struct {
	db *pg.DB
	mux *mux.Router
}

func newServer(db *pg.DB, mux *mux.Router) *server {
	s:= server{db, mux}
	s.routes() //add routes
	return &s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
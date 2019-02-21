package main

import (
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	db := pg.Connect(&pg.Options{
		Database: "pgsimon",
		User:     "simon",
		Password: "simon",
	})

	mux := mux.NewRouter()
	server := newServer(db, mux)
	http.ListenAndServe(":8080", server)

}

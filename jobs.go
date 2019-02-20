package main

import (
	"encoding/json"
	"github.com/azbshiri/common/db"
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
)

type job struct {
	db.Model
	Name string `json:"name"`
}

func (s *server) getJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []*job
	err := s.db.Model(&jobs).Select()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(DatabaseError)
		return
	}

	ffjson.NewEncoder(w).Encode(jobs)
}

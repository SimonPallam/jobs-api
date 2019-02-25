package main

import (
	"encoding/json"
	"github.com/azbshiri/common/db"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
	"strconv"
)

type job struct {
	db.Model
	id         uint64
	Name       string `json:"name"`
	created_at string
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

func (s *server) getJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadParamError)
		return
	}

	_job := job{}
	err = s.db.Model(&_job).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(JobNotFoundError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(DatabaseError)
		return
	}

	ffjson.NewEncoder(w).Encode(_job)
}

func (s *server) createJobs(w http.ResponseWriter, r *http.Request) {
	var param struct {
		Name string `json:"name"`
	}
	err := ffjson.NewDecoder().DecodeReader(r.Body, &param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadParamError)
		return
	}

	job := job{Name: param.Name}
	err = s.db.Insert(&job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(DatabaseError)
		return
	}

	ffjson.NewEncoder(w).Encode(job)
}

func (s *server) deleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadParamError)
		return
	}

	_job := job{id: id}
	err = s.db.Delete(_job)
	if err != nil {
		if err == pg.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(JobNotFoundError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(DatabaseError)
		return
	}

	ffjson.NewEncoder(w).Encode(_job)
}

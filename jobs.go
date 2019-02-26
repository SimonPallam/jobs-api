package main

import (
	"encoding/json"
	"fmt"
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
	debug("the id of the Job to delete is : %v", _job.id)

	err = s.db.Model(&_job).Where("id = ?", &id).Select()
	debug("the Job to be deleted is %v", _job)
	err = s.db.Delete(&_job)

	if err != nil {
		debug("The error is %v", err)
		if err == pg.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(JobNotFoundError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(DatabaseError)
		return
	}
	debug("The job has been deleted %v", nil)
	w.WriteHeader(http.StatusOK)
}

func (s *server) updateJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadParamError)
		return
	}

	var param struct {
		Name string `json:"name"`
	}
	err2 := ffjson.NewDecoder().DecodeReader(r.Body, &param)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadParamError)
		return
	}
	debug("The name param is %v \n", param.Name)

	_job := job{}
	_job.id = id
	_job.Name = param.Name
	debug("The prepared job is %v \n", _job)

	_, err3 := s.db.Model(&_job).Set("name = ?", param.Name).Where("id = ?", id).Update()
	if err3 != nil {
		fmt.Println(err3)
		if err3 == pg.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(JobNotFoundError)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(DatabaseError)
		return
	}
	_job2 := job{}
	err = s.db.Model(&_job2).Where("id = ?", id).Select()
	debug("The updated job is %v \n", _job2)
	ffjson.NewEncoder(w).Encode(_job2)
}

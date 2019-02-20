package main

import (
	"github.com/azbshiri/common/test"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

var testServer *server
var badServer *server

func TestMain(m *testing.M) {
	testServer = newServer(
		pg.Connect(&pg.Options{
			User:     "simon",
			Password: "simon",
			Database: "pgsimon",
		}),
		mux.NewRouter(),
	)
	badServer = newServer(
		pg.Connect(&pg.Options{
			User:     "noone",
			Password: "simon",
			Database: "pgsimon",
		}),
		mux.NewRouter(),
	)

	//temporaryy table for test isolation
	testServer.db.CreateTable(&job{}, &orm.CreateTableOptions{
		Temp: true,
	})

	os.Exit(m.Run())
}

func TestListJobs_emptyResponse(t *testing.T) {
	var body []job
	res, err := test.DoRequest(testServer, "GET", JobPath, nil)

	ffjson.NewDecoder().DecodeReader(res.Body, &body)
	assert.NoError(t, err)
	assert.Len(t, body, 0)
	assert.Equal(t, res.Code, http.StatusOK)
}

func TestListJobs_NormalResponse(t *testing.T) {
	var body []job
	jobs, err := CreateJobListFactory(testServer.db, 10)
	assert.NoError(t, err)

	res, err := test.DoRequest(testServer, "GET", JobPath, nil)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.NoError(t, err)

	assert.Len(t, body, 10)
	assert.Equal(t, jobs, &body)
}

func TestListJobs_DatabaseError(t *testing.T) {
	var body Error
	res, err := test.DoRequest(badServer, "GET", JobPath, nil)
	assert.NoError(t, err)
	assert.Equal(t, DatabaseError, &body)
	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

package main

import (
	"github.com/azbshiri/common/test"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
	"os"
	"testing"
)

var testServer *server
var badServer *server

func TestMain(m *testing.M) {
	testServer = newServer(
		pg.Connect(&pg.Options{
			User:"simon",
			Password:"simon",
			Database:"pgsimon",
		}),
		mux.NewRouter(),
		)
	badServer = newServer(
		pg.Connect(&pg.Options{
			User:"noone",
			Password:"simon",
			Database:"pgsimon",
		}),
		mux.NewRouter(),
	)

	// Here we create a temporary table to store each test case
	// data and follow isolation which would be dropped after.
	testServer.db.CreateTable(&job{}, &orm.CreateTableOptions{
		Temp: true,
	})

	os.Exit(m.Run())
}

func TestListJobs(t *testing.T){
	var body []job
	res, err := test.DoRequest(testServer, "GET", "/jobs", nil)

	ffjson.NewDecoder().DecodeReader(res.Body, &body)
	assert.NoError(t, err)
	assert.Len(t, body, 0)
	assert.Equal(t, res.Code, http.StatusOK)
}
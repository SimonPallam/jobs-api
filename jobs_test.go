package main

import (
	"bytes"
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
			User:     "simonbad",
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

//-----Test List Jobs-----//
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

	ffjson.NewDecoder().DecodeReader(res.Body, &body)
	assert.Len(t, body, 10)
	assert.Equal(t, jobs, &body)
}

func TestListJobs_DatabaseError(t *testing.T) {
	var body Error
	res, err := test.DoRequest(badServer, "GET", JobPath, nil)

	ffjson.NewDecoder().DecodeReader(res.Body, &body)
	assert.NoError(t, err)
	assert.Equal(t, DatabaseError, &body)
	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

//-----Test Create Jobs-----//

func TestCreateJob(t *testing.T) {
	var body job
	byt, err := ffjson.Marshal(&job{Name: "Developper"})
	rdr := bytes.NewReader(byt)

	res, err := test.DoRequest(testServer, "POST", JobPath, rdr)

	ffjson.NewDecoder().DecodeReader(res.Body, &body)
	assert.NoError(t, err)
	assert.Equal(t, "Developper", body.Name)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestCreateJob_BadParamError(t *testing.T) {
	var body Error
	res, err := test.DoRequest(testServer, "POST", JobPath, bytes.NewReader([]byte{}))

	ffjson.NewDecoder().DecodeReader(res.Body, &body)
	assert.NoError(t, err)
	assert.Equal(t, BadParramError, &body)
	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCreateJob_DatabaseError(t *testing.T) {
	var body Error
	byt, err := ffjson.Marshal(&job{Name: "bad developper"})
	rdr := bytes.NewReader(byt)

	res, err := test.DoRequest(badServer, "POST", JobPath, rdr)

	ffjson.NewDecoder().DecodeReader(res.Body, &body)
	assert.NoError(t, err)
	assert.Equal(t, DatabaseError, &body)
	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

//-----Test Read Job-----//
//@todo create job and get id

//@todo read job

//-----Test Update Job-----//
//@todo create job and get id

//@todo read job

//@todo update job

//@todo compare
//-----Test Delete Job-----//
//@todo create job and get id

//@todo delete job

//@todo fail to retreive job

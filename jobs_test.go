package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	// Temporary table for test isolation
	testServer.db.CreateTable(&job{}, &orm.CreateTableOptions{
		Temp: true,
	})

	os.Exit(m.Run())
}

// Test List Jobs
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

// Test Create Job
func TestCreateJob(t *testing.T) {
	var body job

	byt, err := ffjson.Marshal(&job{Name: "Developer"})
	rdr := bytes.NewReader(byt)

	res, err := test.DoRequest(testServer, "POST", JobPath, rdr)
	ffjson.NewDecoder().DecodeReader(res.Body, &body)

	assert.NoError(t, err)
	assert.Equal(t, "Developer", body.Name)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestCreateJob_BadParamError(t *testing.T) {
	var body Error

	res, err := test.DoRequest(testServer, "POST", JobPath, bytes.NewReader([]byte{}))
	ffjson.NewDecoder().DecodeReader(res.Body, &body)

	assert.NoError(t, err)
	assert.Equal(t, BadParamError, &body)
	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCreateJob_DatabaseError(t *testing.T) {
	var body Error

	byt, err := ffjson.Marshal(&job{Name: "bad developer"})
	rdr := bytes.NewReader(byt)

	res, err := test.DoRequest(badServer, "POST", JobPath, rdr)
	ffjson.NewDecoder().DecodeReader(res.Body, &body)

	assert.NoError(t, err)
	assert.Equal(t, DatabaseError, &body)
	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

// Test Commons
func HelperCreateJob() (job, error) {
	var _job job

	byt, err := ffjson.Marshal(&job{Name: "Developer"})
	if err != nil {
		return _job, err
	}
	rdr := bytes.NewReader(byt)
	res, err := test.DoRequest(testServer, "POST", JobPath, rdr)
	ffjson.NewDecoder().DecodeReader(res.Body, &_job)
	return _job, err
}

func HelperClearTable() error {
	err := testServer.db.DropTable(&job{}, &orm.DropTableOptions{
		IfExists: true,
		Cascade:  true,
	})
	if err != nil {
		return err
	}
	err = testServer.db.CreateTable(&job{}, &orm.CreateTableOptions{
		Temp: true,
	})

	return err
}

// Test Read Job
func TestReadJob(t *testing.T) {
	_job, err := HelperCreateJob()

	var body job
	assert.NoError(t, err)

	res, err := test.DoRequest(testServer, "GET", JobPath+`/`+fmt.Sprint(_job.ID), nil)
	ffjson.NewDecoder().DecodeReader(res.Body, &body)

	debug("the body of the response is %v", res.Body)
	assert.NoError(t, err)
	assert.Equal(t, _job, body)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestReadJob_NotFound(t *testing.T) {
	var body Error

	err := HelperClearTable()
	assert.NoError(t, err)

	res, err := test.DoRequest(testServer, "GET", JobPath+"/1", nil)
	ffjson.NewDecoder().DecodeReader(res.Body, &body)

	assert.NoError(t, err)
	assert.Equal(t, JobNotFoundError, &body)
	assert.Equal(t, http.StatusNotFound, res.Code)
}

// Test Update Job
func TestUpdateJob(t *testing.T) {
	_job, err := HelperCreateJob()

	type Request struct {
		Name string `json:"name"`
	}

	request := Request{"changed_name"}
	byt, err := json.Marshal(&request)
	rdr := bytes.NewReader(byt)

	var body2 job
	res, err := test.DoRequest(testServer, "PATCH", JobPath+`/`+fmt.Sprint(_job.ID), rdr)
	ffjson.NewDecoder().DecodeReader(res.Body, &body2)

	debug("the updated job is %v", body2)
	assert.NoError(t, err)
	assert.Equal(t, "changed_name", body2.Name)
	assert.Equal(t, http.StatusOK, res.Code)
}

// Test Delete Job
func TestDeleteJob(t *testing.T) {
	job, err := HelperCreateJob()
	debug("An error occurred in Helper CreateJob %v", err)
	assert.NoError(t, err)

	res, _ := test.DoRequest(testServer, "DELETE", JobPath+"/"+fmt.Sprint(job.ID), nil)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestDeleteJob_Idempotent(t *testing.T) {
	job, err := HelperCreateJob()
	debug("An error occurred in Helper CreateJob %v", err)
	assert.NoError(t, err)

	test.DoRequest(testServer, "DELETE", JobPath+"/"+fmt.Sprint(job.ID), nil)
	res, _ := test.DoRequest(testServer, "DELETE", JobPath+"/"+fmt.Sprint(job.ID), nil)
	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestDeleteJob_BadDatabase(t *testing.T) {
	job, err := HelperCreateJob()
	assert.NoError(t, err)

	res, _ := test.DoRequest(badServer, "DELETE", JobPath+"/"+fmt.Sprint(job.ID), nil)
	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestDeleteJob_BadParam(t *testing.T) {
	res, _ := test.DoRequest(badServer, "DELETE", JobPath+"/bad_param", nil)
	assert.Equal(t, http.StatusBadRequest, res.Code)
}

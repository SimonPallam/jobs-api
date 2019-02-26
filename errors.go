package main

import "net/http"

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func Err(message string, status int) *Error {
	return &Error{message, status}
}

var DatabaseError = Err("Request failed : Database Error", http.StatusInternalServerError)
var BadParamError = Err("Request failed : Bad parameters were sent", http.StatusBadRequest)
var JobNotFoundError = Err("Request failed : No job was found", http.StatusNotFound)

//@todo create proper Response class or subclass
var JobDeletedResponse = Err("Job deleted", http.StatusOK)

package main

import "net/http"

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func Err(message string, status int) *Error {
	return &Error{message, status}
}

var DatabaseError = Err("Request failled : Database Error", http.StatusInternalServerError)
var BadParramError = Err("Request failled : Bad Parameters where sent", http.StatusInternalServerError)
var JobNotFoundError = Err("Request failled : No job was found", http.StatusNotFound)

//@todo create proper Response class or subclass
var JobDeletedResponse = Err("Job deleted", http.StatusOK)

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

package main

import (
	"github.com/go-pg/pg"
	"math/rand"
	"time"
)

func CreateJobListFactory(db *pg.DB, length int) (*[]job, error) {
	jobs := make([]job, length)
	for key, _ := range jobs {
		jobs[key].Name = RandStringRunes(8)
	}

	err := db.Insert(&jobs)
	if err != nil {
		return nil, err
	}

	return &jobs, nil
}

//---tools---//
func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

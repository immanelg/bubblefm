package main

import "math/rand"

type status struct {
	id    int
	isErr bool
	text  string
}

func newStatus(text string, err bool) status {
	return status{text: text, isErr: err, id: rand.Int()}
}

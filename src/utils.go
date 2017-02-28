package main

import (
	"math/rand"
	"time"
)

// random generates a random number between two range
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

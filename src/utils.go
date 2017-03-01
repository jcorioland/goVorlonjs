package main

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"time"
)

// random generates a random number between two range
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func computeSHA256Hash(toHash string, cypher string) string {
	hasher := sha256.New()
	hasher.Write([]byte(toHash))

	// get sha256 as base64 encoded string
	sha := base64.URLEncoding.EncodeToString(hasher.Sum([]byte(cypher)))

	return sha
}

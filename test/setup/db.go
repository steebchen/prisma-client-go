package setup

import (
	"math/rand"
	"time"
)

var source = rand.NewSource(time.Now().UnixNano())
var rng = rand.New(source)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func RandomString() string {
	return "go_client_db_" + randSeq(7)
}

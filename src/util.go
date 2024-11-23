package src

import (
	"time"

	"math/rand"
)

func RandomReplciateSeedId() string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomId := make([]byte, 40)

	for i := 0; i < 40; i++ {
		randomId[i] = charset[random.Intn(len(charset))]
	}

	return string(randomId)
}

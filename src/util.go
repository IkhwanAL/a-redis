package src

import (
	"errors"
	"strconv"
	"time"
	"unicode"

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

func IsStringCanBeEncodedAsUInteger(words string) (int64, error) {
	if words == "" {
		return 0, errors.New("value cannot be empty string")
	}

	firstRune := words[0]

	if firstRune == '0' && len(words) > 1 {
		return 0, errors.New("value should not start with 0")
	}

	if firstRune == '+' || firstRune == '-' {
		return 0, errors.New("value should not start with '+' or '-'")
	}

	for _, word := range words {
		if !unicode.IsDigit(word) {
			return 0, errors.New("one of the character is not digit")
		}
	}

	intLength, err := strconv.ParseInt(words, 10, 64)

	if err != nil {
		return 0, err
	}

	return intLength, nil
}

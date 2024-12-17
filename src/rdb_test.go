package src

import (
	"fmt"
	"testing"
)

func TestRDBLengthEncoding(t *testing.T) {
	testTable := []map[string]interface{}{
		{
			"test":   10,
			"expect": []byte{10},
		},
		{
			"test":   700,
			"expect": []byte{'\x42', '\xbc'},
		},
		{
			"test":   17000,
			"expect": []byte{'\x80', 00, 00, '\x42', '\x68'},
		},
	}

	for _, toTest := range testTable {
		response := sizeBitMask(toTest["test"].(int))
		// Data Type Problem

		actualValue := fmt.Sprintf("%x", response)
		expectValue := fmt.Sprintf("%x", toTest["expect"])

		if actualValue != expectValue {
			t.Errorf("%x is not %x", response, toTest["expect"])
		}
	}
}

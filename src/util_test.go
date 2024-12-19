package src

import (
	"testing"
)

func TestEncodeStringToInt64(t *testing.T) {
	testTable := []map[string]interface{}{
		{
			"test":   "10",
			"expect": int64(10),
		},
		{
			"test":   "100",
			"expect": int64(100),
		},
		{
			"test":   "+100",
			"expect": int64(0),
		},
		{
			"test":   "-100",
			"expect": int64(0),
		},
		{
			"test":   "007",
			"expect": int64(0),
		},
		{
			"test":   "0",
			"expect": int64(0),
		},
	}

	for _, toTest := range testTable {
		response, _ := IsStringCanBeEncodedAsUInteger(toTest["test"].(string))

		if response != toTest["expect"].(int64) {
			t.Errorf("%v is not %v", response, toTest["expect"])
		}
	}
}

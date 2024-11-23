package src

import (
	"testing"
)

func TestProtocolGeneratorRESP3(t *testing.T) {
	testTable := []map[string]interface{}{
		{
			"test":   "*1\r\n$4\r\nPING\r\n",
			"expect": "PING",
		},
		{
			"test":   "*3\r\n$3\r\nSET\r\n$4\r\nAABC\r\n$2\r\naa\r\n",
			"expect": []string{"SET", "AABC", "aa"},
		},
		{
			"test":   "*2\r\n$3\r\nGET\r\n$4\r\nAABC\r\n",
			"expect": []string{"GET", "AABC"},
		},
	}

	for _, toTest := range testTable {
		request := ParseReadRESP([]byte(toTest["test"].(string)))

		switch toTest["expect"].(type) {
		case string:
			if request[0] != toTest["expect"] {
				t.Errorf("%s is not %s", request[0], toTest["expect"])
			}
		case []string:
			for index, expectResult := range toTest["expect"].([]string) {
				if request[index] != expectResult {
					t.Errorf("%s is not %s", request[0], expectResult)
				}
			}
		}
	}
}

func TestProtocolGeneratorMultipleValueRESP3(t *testing.T) {
	testTable := []map[string]interface{}{
		{
			"test":   []string{"replicaId:123", "offset:0"},
			"expect": "$25\r\nreplicaId:123\r\noffset:0\r\n\r\n",
		},
	}

	for _, toTest := range testTable {
		response := ParseGenerateMultipleValue(toTest["test"].([]string)...)

		if response != toTest["expect"].(string) {
			t.Errorf("%v is not %v", response, toTest["expect"])
		}
	}
}

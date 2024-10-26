package src

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
)

var totalCommandLengthIndicator = byte('*')
var totalCommandCharLengthIndicator = byte('$')

var crlf = "\r\n"

func ParseReadRESP(requests []byte) []string {
	var bufRead bytes.Buffer

	bufRead.Write(requests)

	scanner := bufio.NewScanner(&bufRead)

	totalCommand := -1

	var messages []string

	for scanner.Scan() {
		buffers := scanner.Bytes()

		if len(buffers) == 0 {
			continue
		}

		if buffers[0] == totalCommandLengthIndicator {
			totalCommand = int(buffers[1] - '0')
			continue
		}

		if buffers[0] == totalCommandCharLengthIndicator {
			continue
		}

		if (len(messages)) != totalCommand {
			messages = append(messages, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return messages
}

func ParseGenerateRESP(value string) string {

	lengthOfValue := len(value)

	return fmt.Sprintf("$%d%s%s%s", lengthOfValue, crlf, value, crlf)
}

func ParseGenerateRESPError(message string) string {
	return fmt.Sprintf("$%s%s", message, crlf)
}

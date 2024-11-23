package src

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
)

var totalCommandLengthIndicator = byte('*')
var totalCommandCharLengthIndicator = byte('$')
var simpleStringIndicator = byte('+')

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

		if buffers[0] == simpleStringIndicator {
			messages = append(messages, string(buffers[1:]))
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

func ParseGenerateMultipleValue(multipleValue ...string) string {
	var response string

	lengthValue := 0

	for _, value := range multipleValue {
		tmp := fmt.Sprintf("%s%s", value, crlf)
		len := len(value)

		response += tmp
		lengthValue += len + 2
	}

	// Add One Crlf for Redis Cli Accept The Message
	return fmt.Sprintf("$%d%s%s%s", lengthValue, crlf, response, crlf)
}

func ParseGenerateArrayValueRESP(arrayValue ...string) string {
	reponses := fmt.Sprintf("*%d%s", len(arrayValue), crlf)

	for _, value := range arrayValue {
		tmp := fmt.Sprintf("$%d%s%s%s", len(value), crlf, value, crlf)

		reponses += tmp
	}

	return reponses
}

func ParseGenerateRESPError(message string) string {
	return fmt.Sprintf("$%s%s", message, crlf)
}

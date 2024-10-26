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

func ParseReadRESP(requests []byte) ([]int, []string) {
	var bufRead bytes.Buffer

	bufRead.Write(requests)

	scanner := bufio.NewScanner(&bufRead)

	totalCommand := -1

	var cmdLengthSlice []int

	var command []string

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
			cmdLengthSlice = append(cmdLengthSlice, int(buffers[1]-'0'))
			continue
		}

		if (len(command)) != totalCommand {
			command = append(command, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return cmdLengthSlice, command
}

func ParseGenerateRESP(value string) string {

	lengthOfValue := len(value)

	return fmt.Sprintf("$%d%s%s%s", lengthOfValue, crlf, value, crlf)
}

func ParseGenerateRESPError(message string) string {
	return fmt.Sprintf("$%s%s", message, crlf)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

func decodeBencode(bencodedString string, startIndex int) (interface{}, int, error) {
	strLen := len(bencodedString)
	firstByte := bencodedString[startIndex]

	if startIndex == strLen {
		return nil, startIndex, io.ErrUnexpectedEOF
	}

	switch {
	case firstByte == 'i':
		return decodeInt(bencodedString, startIndex);
	case unicode.IsDigit(rune(firstByte)):
		return decodeString(bencodedString, startIndex);
	default:
		return nil, startIndex, fmt.Errorf("unexpected value: %q", firstByte)
	}
}

func decodeString(bencodedString string, startIndex int) (string, int, error) {
	strLen := len(bencodedString)
	firstByte := bencodedString[startIndex]

	if unicode.IsDigit(rune(firstByte)) {
		var colonIndex int

		for i := startIndex; i < strLen; i++ {
			if bencodedString[i] == ':' {
				colonIndex = i
				break
			}
		}
		
		if (colonIndex == 0) {
			return "", startIndex, fmt.Errorf("bad string, didn't found semicolon")
		}

		lengthStr := bencodedString[startIndex:colonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", startIndex, err
		}

		if length > strLen {
			return "", startIndex, fmt.Errorf("bad string, length is greater than string length")
		}
		str := string([]rune(bencodedString)[colonIndex+1:colonIndex+1+length])
		return str, colonIndex + 1 + length, nil
	} else {
		return "", startIndex, fmt.Errorf("bad string")
	}
}

func decodeInt(bencodedString string, startIndex int) (int, int, error) {
	strLen := len(bencodedString)
	firstByte := bencodedString[startIndex]

	if firstByte == 'i' {
		var endIndex int

		for i := startIndex+1; i < strLen; i++ {
			if bencodedString[i] == 'e' {
				endIndex = i
				break
			}
		}

		numStr := bencodedString[startIndex+1:endIndex]

		num, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, startIndex, err
		}

		return num, endIndex+1, nil
	} else {
		return 0, startIndex, fmt.Errorf("bad int")
	}
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]
		
		decoded,_, err := decodeBencode(bencodedValue, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}

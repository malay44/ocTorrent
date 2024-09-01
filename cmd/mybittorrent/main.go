package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

func decodeBencode(bencodedString string) (interface{}, error) {
	strLen := len(bencodedString)
	firstByte := bencodedString[0]
	lastByte := bencodedString[strLen-1]


	if unicode.IsDigit(rune(firstByte)) {
		var colonIndex int

		for i := 0; i < strLen; i++ {
			if bencodedString[i] == ':' {
				colonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:colonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}

		return bencodedString[colonIndex+1 : colonIndex+1+length], nil
	} else if firstByte == 'i' && lastByte == 'e' {
		num, err := strconv.Atoi(bencodedString[1 : strLen-1])
		if err != nil {
			return "", err
		}

		return num, nil
	} else {
		return "", fmt.Errorf("only strings are supported at the moment")
	}
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]
		
		decoded, err := decodeBencode(bencodedValue)
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

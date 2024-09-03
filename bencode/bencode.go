package bencode

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func Decode(bencodedString string, startIndex int) (interface{}, int, error) {
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
	case firstByte == 'l':
		return decodeList(bencodedString, startIndex);
	case firstByte == 'd':
		return decodeDict(bencodedString, startIndex);
	default:
		fmt.Println(bencodedString);
		// explanation of the following line:
		// strconv.Itoa(startIndex+1) - converts the startIndex to a string
		fmt.Println(fmt.Sprintf("%"+strconv.Itoa(startIndex+1)+"s", "^"))

		return nil, startIndex, fmt.Errorf("unexpected value: %q at index %d", firstByte, startIndex)
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
		str := bencodedString[colonIndex+1:colonIndex+1+length]
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

func decodeList(bencodedString string, startIndex int) ([]interface{}, int, error) {
	strLen := len(bencodedString)
	firstByte := bencodedString[startIndex]

	if firstByte == 'l' {
		var endIndex int
		list := []interface{}{}

		for i := startIndex+1; i < strLen; {
			if bencodedString[i] == 'e' {
				endIndex = i
				break
			}
			decoded, nextIndex, err := Decode(bencodedString, i)
			if err != nil {
				return nil, startIndex, err
			}
			list = append(list, decoded)
			i = nextIndex
		}

		return list, endIndex+1, nil
	} else {
		return nil, startIndex, fmt.Errorf("bad list")
	}
}

func decodeDict(bencodedString string, startIndex int) (map[string]interface{}, int, error) {
	strLen := len(bencodedString)
	firstByte := bencodedString[startIndex]

	if firstByte == 'd' {
		var endIndex int
		dict := map[string]interface{}{}

		for i := startIndex+1; i < strLen; {
			if bencodedString[i] == 'e' {
				endIndex = i
				break
			}

			key, nextIndex, err := decodeString(bencodedString, i)
			if err != nil {
				return nil, startIndex, err
			}
			i = nextIndex

			decoded, nextIndex, err := Decode(bencodedString, i)
			if err != nil {
				return nil, startIndex, err
			}
			dict[key] = decoded
			i = nextIndex
		}

		return dict, endIndex+1, nil
	} else {
		return nil, startIndex, fmt.Errorf("bad dict")
	}
}

func Encode(data interface{}) (string, error) {
	// "type switch" to determine the type of data
	switch v := data.(type) {
	case int:
		return encodeInt(v), nil
	case string:
		return encodeString(v), nil
	case []interface{}:
		return encodeList(v)
	case map[string]interface{}:
		return encodeDict(v)
	default:
		return "", fmt.Errorf("unsupported type: %T", v)
	}
}

func encodeInt(num int) string {
	return fmt.Sprintf("i%de", num)
}

func encodeString(str string) string {
	return fmt.Sprintf("%d:%s", len(str), str)
}

func encodeList(list []interface{}) (string, error) {
	encodedList := "l"

	for _, v := range list {
		encoded, err := Encode(v)
		if err != nil {
			return "", err
		}
		encodedList += encoded
	}

	return encodedList + "e", nil
}

func encodeDict(dict map[string]interface{}) (string, error) {
	var encodedDictSb strings.Builder
	encodedDictSb.WriteString("d")

	// Sort the keys in alphabetical order
	keys := make([]string, 0, len(dict))
	for k := range dict {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := dict[k]

		encodedKey := encodeString(k)
		encodedDictSb.WriteString(encodedKey)

		encodedValue, err := Encode(v)
		if err != nil {
			return "", err
		}
		encodedDictSb.WriteString(encodedValue)
	}

	encodedDictSb.WriteString("e")
	return encodedDictSb.String(), nil
}
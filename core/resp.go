package core

import (
	"errors"
	"fmt"
)

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}
	return nil, 0, nil

}

func readLength(data []byte) (int, int) {
	pos, length := 0, 0
	for ; data[pos] != '\r'; pos++ {
		if data[pos] >= '0' && data[pos] <= '9' {
			length = length*10 + int(data[pos]-'0')
		}
	}
	return length, pos + 2
}

func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for ; data[pos] != '\r'; pos++ {
	}
	return string(data[1:pos]), pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readInt64(data []byte) (int64, int, error) {
	pos := 1
	var value int64 = 0
	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}
	return value, pos + 2, nil
}

func readBulkString(data []byte) (string, int, error) {
	pos := 1
	len, delta := readLength(data[pos:])
	pos += delta
	return string(data[pos:(pos + len)]), pos + len + 2, nil
}

func readArray(data []byte) (interface{}, int, error) {
	pos := 1

	count, delta := readLength(data[pos:])
	pos += delta
	var elements []interface{} = make([]interface{}, count)
	for i := range elements {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elements[i] = elem
		pos += delta
	}
	return elements, pos, nil

}

func DecodeArrayString(data []byte) ([]string, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}
	ts := value.([]interface{})
	tokens := make([]string, len(ts))
	for i := range tokens {
		tokens[i] = ts[i].(string)
	}
	return tokens, nil
}

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	value, _, err := DecodeOne(data)
	return value, err
}

func Encode(value interface{}, isSimple bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", v))
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
	}
	return []byte{}
}

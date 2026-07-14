package core

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"time"
)

var RESP_NIL []byte = []byte("$-1\r\n")
var RESP_OK []byte = []byte("+OK\r\n")
var RESP_EXPIRED []byte = []byte(":-2\r\n")
var RESP_DOES_NOT_EXIST []byte = []byte(":-1\r\n")

func evalPING(args []string) []byte {
	var b []byte
	if len(args) >= 2 {
		return Encode(errors.New("ERR wrong number of argumnets for 'ping' command"), false)
	}
	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	return b
}

func evalSet(args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SET' command"), false)
	}
	var key, value string
	var exDurationMs int64 = -1

	key = args[0]
	value = args[1]

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			i++
			if i == len(args) {
				return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
			}
			exDurationSec, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
			}
			exDurationMs = exDurationSec * 1000
		default:
			return Encode(errors.New("(error) ERR syntax error"), false)
		}

	}
	Put(key, NewObj(value, exDurationMs))
	return RESP_OK
}

func evalGet(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command"), false)
	}

	var key string = args[0]
	obj := Get(key)

	if obj == nil {
		return RESP_NIL
	}
	if obj.ExpiresAt != -1 && obj.ExpiresAt < time.Now().UnixMilli() {
		return RESP_NIL
	}
	return Encode(obj.Value, false)
}

func evalTTL(args []string) []byte {
	if len(args) != 1 {
		Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command"), false)
	}
	var key string = args[0]
	obj := Get(key)

	if obj == nil {
		return RESP_EXPIRED
	}

	if obj.ExpiresAt == -1 {
		return RESP_DOES_NOT_EXIST
	}

	durationMs := obj.ExpiresAt - time.Now().UnixMilli()
	if durationMs < 0 {
		return RESP_EXPIRED
	}

	return Encode(int64(durationMs/1000), false)
}

func evalDEL(args []string) []byte {
	if len(args) < 1 {
		Encode(errors.New("(error) ERR wrong number of arguments for 'DEL' command"), false)
	}
	return Encode(Del(args), false)
}

func evalExpire(args []string) []byte {
	if len(args) <= 1 {
		Encode(errors.New("(error) ERR wrong number of arguments for 'EXPIRE' command"), false)
	}
	var key string = args[0]

	exDurationSec, err := strconv.ParseInt(args[1], 10, 64)

	if err != nil {
		return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
	}

	obj := Get(key)
	if obj == nil {
		return []byte(":0\r\n")
	}

	obj.ExpiresAt = time.Now().UnixMilli() + exDurationSec*1000
	return []byte(":1\r\n")
}

func EvalAndRespond(cmds []*RedisCmd, c io.ReadWriter) {
	var responses []byte
	buf := bytes.NewBuffer(responses)
	for _, cmd := range cmds {
		switch cmd.Cmd {
		case "PING":
			buf.Write(evalPING(cmd.Args))
		case "SET":
			buf.Write(evalSet(cmd.Args))
		case "GET":
			buf.Write(evalGet(cmd.Args))
		case "TTL":
			buf.Write(evalTTL(cmd.Args))
		case "DEL":
			buf.Write(evalDEL(cmd.Args))
		case "EXPIRE":
			buf.Write(evalExpire(cmd.Args))
		default:
			buf.Write(evalPING(cmd.Args))
		}
	}
	c.Write(buf.Bytes())
}

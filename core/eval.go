package core

import (
	"errors"
	"io"
)

func EvalPING(args []string, c io.ReadWriter) error {
	var b []byte
	if len(args) >= 2 {
		return errors.New("ERR wrong number of argumnets for 'ping' command")
	}
	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	_, err := c.Write(b)
	return err
}
func EvalAndRespond(cmd *RedisCmd, c io.ReadWriter) error {
	switch cmd.Cmd {
	case "PING":
		return EvalPING(cmd.Args, c)
	default:
		return EvalPING(cmd.Args, c)
	}
}

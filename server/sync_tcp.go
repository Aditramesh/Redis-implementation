package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/Aditramesh/Redis-implementation/core"
)

func readCommand(c io.ReadWriter) (*core.RedisCmd, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return nil, err
	}
	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}
	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, c io.ReadWriter) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmd *core.RedisCmd, c io.ReadWriter) error {
	err := core.EvalAndRespond(cmd, c)
	if err != nil {
		respondError(err, c)
	}
	return nil
}

func RunTCPSyncServer() {
	var conn_clients int = 0

	lsnr, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}

	for {
		c, err := lsnr.Accept()
		if err != nil {
			panic(err)
		}

		conn_clients++
		fmt.Println("client connected with addr", c.RemoteAddr(), "conn_clients", conn_clients)

		for {
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				conn_clients--
				fmt.Println("client disconnected", c.RemoteAddr(), "conn_clients", conn_clients)
				if err == io.EOF {
					break
				}
				log.Println("err:", err)
			}
			log.Println("command", cmd)
			respond(cmd, c)
		}
	}
}

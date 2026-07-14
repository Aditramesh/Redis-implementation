package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/Aditramesh/Redis-implementation/core"
)

func readCommands(c io.ReadWriter) ([]*core.RedisCmd, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return nil, err
	}
	commands, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func respondError(err error, c io.ReadWriter) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmds []*core.RedisCmd, c io.ReadWriter) error {
	core.EvalAndRespond(cmds, c)
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
			cmds, err := readCommands(c)
			if err != nil {
				c.Close()
				conn_clients--
				fmt.Println("client disconnected", c.RemoteAddr(), "conn_clients", conn_clients)
				if err == io.EOF {
					break
				}
				log.Println("err:", err)
			}
			log.Println("command", cmds)
			respond(cmds, c)
		}
	}
}

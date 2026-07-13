//go:build linux

package server

import (
	"log"
	"net"
	"syscall"
	"time"

	"github.com/Aditramesh/Redis-implementation/core"
)

var con_clients = 0
var cronFrequency time.Duration = 1 * time.Second
var lastCronExecTime time.Time = time.Now()

func RunTCPASyncServer() error {
	log.Println("starting an asynchronous TCP server on localhost port 9999")

	max_clients := 20000

	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_clients)
	serverFd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFd)
	if err = syscall.SetNonblock(serverFd, true); err != nil {
		return err
	}
	ip4 := net.ParseIP("127.0.0.1")
	if err = syscall.Bind(serverFd, &syscall.SockaddrInet4{
		Port: 9999,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	if err = syscall.Listen(serverFd, max_clients); err != nil {
		return err
	}

	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(epollFD)

	var socketServerEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFd),
	}

	if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFd, &socketServerEvent); err != nil {
		return err
	}

	for {
		if time.Now().After(lastCronExecTime.Add(cronFrequency)) {
			core.DeleteExpiredKeys()
			lastCronExecTime = time.Now()
		}
		nevents, e := syscall.EpollWait(epollFD, events[:], -1)
		if e != nil {
			continue
		}
		for i := 0; i < nevents; i++ {
			if int(events[i].Fd) == serverFd {
				fd, _, err := syscall.Accept(serverFd)
				if err != nil {
					log.Println("err", err)
					continue
				}
				con_clients++
				syscall.SetNonblock(serverFd, true)

				var socketClientEvent syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(fd),
				}

				if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClientEvent); err != nil {
					log.Fatal(err)
				}
			} else {
				comm := core.FDComm{Fd: int(events[i].Fd)}
				cmd, err := readCommand(comm)
				if err != nil {
					syscall.Close(int(events[i].Fd))
					con_clients -= 1
					continue
				}
				respond(cmd, comm)
			}

		}
	}

}

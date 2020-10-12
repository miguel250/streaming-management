package util

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
)

type EchoServer struct {
	sync.RWMutex
	shutdown       chan bool
	serverResponse *bytes.Buffer
	listener       net.Listener
	t              *testing.T
	addr           string
}

func (e *EchoServer) Start() {
	go func() {
		for {
			select {
			case _, ok := <-e.shutdown:
				if !ok {
					e.listener.Close()
					return
				}
			default:
				conn, err := e.listener.Accept()

				if err != nil {
					e.t.Errorf("Failed to start listening for connections with %s", err)
					close(e.shutdown)
					continue
				}

				go func(c net.Conn) {
					e.Lock()
					defer e.Unlock()
					if e.serverResponse == nil {
						io.Copy(c, c)
						c.Close()
						return
					}

					fmt.Println(e.serverResponse.String())

					io.Copy(c, e.serverResponse)
				}(conn)
			}
		}
	}()
}

func (e *EchoServer) SetResponse(msg *bytes.Buffer) {
	e.Lock()
	defer e.Unlock()

	msg.WriteByte('\n')
	e.serverResponse = msg
}

func (e *EchoServer) Shutdown() {
	close(e.shutdown)
}

func MockTwitchChatServer(t *testing.T) *EchoServer {
	l, err := net.Listen("tcp", ":0")

	if err != nil {
		t.Fatalf("Failed to listen for connections with %s", err)
	}

	return &EchoServer{
		shutdown: make(chan bool, 1),
		listener: l,
		t:        t,
		addr:     l.Addr().String(),
	}
}

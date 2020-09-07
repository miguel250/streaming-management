package irc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"sync"
	"testing"
)

func TestSimpleClient(t *testing.T) {
	ts := testIRCEchoServer(t)
	ts.Start()
	defer ts.Shutdown()

	conn, err := net.Dial("tcp", ts.addr)

	if err != nil {
		t.Fatalf("Failed to create dialer with %s", err)
	}

	want := "Hello World!"
	fmt.Fprintf(conn, "%s\n", want)

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	line, err := tp.ReadLine()

	if err != nil {
		t.Fatalf("Failed to read message with %s", err)
	}

	if line != want {
		t.Errorf("Lines didn't match want: %s, got: %s", want, line)
	}
}

type echoServer struct {
	sync.RWMutex
	shutdown       chan bool
	serverResponse *bytes.Buffer
	listener       net.Listener
	t              *testing.T
	addr           string
}

func (e *echoServer) Start() {
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
					e.RLock()
					defer e.RUnlock()
					if e.serverResponse == nil {
						io.Copy(c, c)
						c.Close()
						return
					}

					io.Copy(c, e.serverResponse)
				}(conn)
			}
		}
	}()
}

func (e *echoServer) setResponse(msg *bytes.Buffer) {
	e.Lock()
	defer e.Unlock()

	msg.WriteByte('\n')
	e.serverResponse = msg
}

func (e *echoServer) Shutdown() {
	close(e.shutdown)
}

func testIRCEchoServer(t *testing.T) *echoServer {
	l, err := net.Listen("tcp", ":0")

	if err != nil {
		t.Fatalf("Failed to listen for connections with %s", err)
	}

	return &echoServer{
		shutdown: make(chan bool, 1),
		listener: l,
		t:        t,
		addr:     l.Addr().String(),
	}
}

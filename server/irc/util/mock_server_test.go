package util

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"testing"
)

func TestSimpleServer(t *testing.T) {
	ts := MockTwitchChatServer(t)
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

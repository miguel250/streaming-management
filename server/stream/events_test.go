package stream

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEvent(t *testing.T) {
	event := New()
	event.Start()

	server := httptest.NewServer(event)
	defer server.Close()

	requestDone := make(chan bool)

	go func() {
		res, err := http.Get(server.URL)
		if err != nil {
			t.Error(err)
		}

		want := "event: new_follower\ndata: Hi\n\n"
		buf := make([]byte, len(want))

		n, err := res.Body.Read(buf)
		defer res.Body.Close()

		if n != len(want) {
			t.Errorf("Read len: %d, err: %v Want len: %d, Body: %s", n, err, len(want), string(buf))
			requestDone <- true
			return
		}

		if string(buf) != want {
			t.Errorf("Got %s, want: %s", string(buf), want)
		}

		requestDone <- true
	}()

	<-event.ClientConnected
	message := Message{
		Type: NewFollower,
		Text: "Hi",
	}
	event.Message <- message
	<-requestDone
	event.Close()
}

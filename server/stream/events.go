package stream

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

type EventType int

const (
	NewFollower EventType = iota
	NewSubscriber
	NewChatMessage
)

type Event struct {
	sync.RWMutex
	Message          chan Message
	clients          map[chan Message]struct{}
	ClientConnect    chan chan Message
	ClientConnected  chan bool
	ClientDisconnect chan chan Message
	shutdown         chan os.Signal
	isRunning        bool
	Once             sync.Once
}

type Message struct {
	Type EventType
	Text string
}

var EventTypeToString = map[EventType]string{
	NewFollower:    "new_follower",
	NewSubscriber:  "new_subscriber",
	NewChatMessage: "new_chat_message",
}

func (e *Event) Start() error {
	if e.isRunning {
		return errors.New("Server is already running")
	}

	e.Lock()
	defer e.Unlock()
	e.isRunning = true
	e.Once = sync.Once{}

	go func() {
		for {
			select {
			case message, ok := <-e.Message:
				if ok {
					for client := range e.clients {
						client <- message
					}
				}
			case <-e.ClientConnected:
			case client, ok := <-e.ClientConnect:
				if ok {
					e.Lock()
					e.clients[client] = struct{}{}
					e.Unlock()
					e.ClientConnected <- true
				}
			case client, ok := <-e.ClientDisconnect:
				if ok {
					e.Lock()
					delete(e.clients, client)
					close(client)
					e.Unlock()
				}
			case <-e.shutdown:
				e.Close()
				return
			}
		}
	}()
	return nil
}

func (e *Event) Close() {
	e.Once.Do(e.stop)
}

func (e *Event) stop() {
	e.Lock()
	defer e.Unlock()
	if e.isRunning {
		e.isRunning = false
		close(e.ClientConnect)
		close(e.ClientConnected)
		close(e.ClientDisconnect)
		close(e.Message)

		for client := range e.clients {
			close(client)
		}
	}
}

func (e *Event) Send(event EventType, msg string) {
	message := Message{
		Type: event,
		Text: msg,
	}
	e.Message <- message
}

func (e *Event) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	client := make(chan Message)

	rw.Header().Add("Cache-Control", "no-cache")
	rw.Header().Add("Content-Type", "text/event-stream")

	e.ClientConnect <- client
	for {
		select {
		case message, ok := <-client:

			if !ok {
				return
			}

			fmt.Fprintf(rw, "event: %s\n", EventTypeToString[message.Type])
			fmt.Fprintf(rw, "data: %s\n", message.Text)
			fmt.Fprint(rw, "\n")
			rw.(http.Flusher).Flush()
		case <-req.Context().Done():
			e.RLock()
			defer e.RUnlock()
			if e.isRunning {
				e.ClientDisconnect <- client
			}
			return
		}
	}
}

func New() *Event {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	return &Event{
		Message:          make(chan Message),
		clients:          make(map[chan Message]struct{}),
		ClientConnect:    make(chan chan Message),
		ClientConnected:  make(chan bool, 100),
		ClientDisconnect: make(chan chan Message),
		shutdown:         shutdown,
	}
}

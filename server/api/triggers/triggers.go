package triggers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/miguel250/streaming-setup/server/config"
	"github.com/miguel250/streaming-setup/server/stream"
)

type Triggers struct {
	conf  *config.Config
	event *stream.Event
}

func (t *Triggers) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	pathPieces := strings.Split(req.URL.Path, "/")

	if len(pathPieces) != 4 {
		http.Error(rw, "missing action", http.StatusBadRequest)
		return
	}

	action := pathPieces[3]

	switch action {
	case "new-follower":
		t.event.Send(stream.NewFollower, t.conf.Twitch.IRC.Channel)
		fmt.Fprintln(rw, "Action Triggered")
	case "new-subscriber":
		t.event.Send(stream.NewSubscriber, t.conf.Twitch.IRC.Channel)
		fmt.Fprintln(rw, "Action Triggered")
	default:
		http.Error(rw, "unknown action", http.StatusNotFound)
	}
}

func New(event *stream.Event, conf *config.Config) *Triggers {
	return &Triggers{
		conf:  conf,
		event: event,
	}
}

package triggers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/miguel250/streaming-setup/server/config"
	"github.com/miguel250/streaming-setup/server/irc"
	"github.com/miguel250/streaming-setup/server/stream"
	"github.com/miguel250/streaming-setup/server/twitch"
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
	case "chat-message":
		msg := irc.Message{
			Message: "This is a test",
			Badges: []*twitch.Badge{
				{
					Title:   "3-Month Subscriber",
					Image1X: "https://static-cdn.jtvnw.net/badges/v1/a2b9b912-4d2a-4103-b741-8b1ebe42fdcc/1",
					Image2X: "https://static-cdn.jtvnw.net/badges/v1/a2b9b912-4d2a-4103-b741-8b1ebe42fdcc/2",
					Image4X: "https://static-cdn.jtvnw.net/badges/v1/a2b9b912-4d2a-4103-b741-8b1ebe42fdcc/3",
				},
			},
			DisplayName:  "MiguelCodeTV",
			ProfileImage: "https://static-cdn.jtvnw.net/jtv_user_pictures/7345d8af-adc7-41b0-a342-57e829941608-profile_image-300x300.png",
			Channel:      "miguelcodetv",
		}

		b, err := json.Marshal(msg)
		if err != nil {
			log.Println("error:", err)
		}

		t.event.Send(stream.NewChatMessage, string(b))
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

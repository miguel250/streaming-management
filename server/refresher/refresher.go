package refresher

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/miguel250/streaming-setup/server/cache"
	"github.com/miguel250/streaming-setup/server/config"
	"github.com/miguel250/streaming-setup/server/stream"
	"github.com/miguel250/streaming-setup/server/twitch"
)

const (
	FollowerRefreshTime = 15
)

type Worker struct {
	conf      *config.Config
	cache     *cache.Cache
	client    *twitch.API
	event     *stream.Event
	isRunning bool
}

func (w *Worker) Refresher() {
	newFollower, err := w.currentFollower()

	if err != nil {
		log.Printf("Failed to get current follower with %v\n", err)
	}

	if newFollower != nil && w.isRunning {
		fmt.Printf("New follower!! %s\n", newFollower.DisplayName)
		message := stream.Message{
			Type: stream.NewFollower,
			Text: newFollower.DisplayName,
		}
		w.event.Message <- message
	}
	time.AfterFunc(FollowerRefreshTime*time.Second, func() {
		w.Refresher()
	})
	w.isRunning = true
}

func (w *Worker) currentFollower() (*twitch.User, error) {

	currentFollowers, err := w.client.Channel.Followers(w.conf.Twitch.ChannelID, 1)

	if err != nil {
		return nil, fmt.Errorf("Failed to get current followers")
	}
	if len(currentFollowers.Follows) == 0 {
		return nil, nil
	}

	currentFollower := currentFollowers.Follows[0].User
	currentFollowerID := currentFollower.ID

	oldFollowerID, _ := w.cache.Get(cache.LastFollowerIDKey)

	if oldFollowerID != currentFollowerID {
		w.cache.Set(cache.LastFollowerNameKey, currentFollower.DisplayName)
		w.cache.Set(cache.LastFollowerIDKey, currentFollowerID)
		w.cache.Set(cache.TotalFollowerKey, strconv.Itoa(currentFollowers.Total))
		return &currentFollower, nil
	}
	return nil, nil
}

func New(conf *config.Config, c *cache.Cache, client *twitch.API, event *stream.Event) *Worker {
	return &Worker{
		conf:   conf,
		cache:  c,
		client: client,
		event:  event,
	}
}

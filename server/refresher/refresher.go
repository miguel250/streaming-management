package refresher

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
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
	once      sync.Once
	started   sync.Once
}

func (w *Worker) Refresher() {
	if w.apiAuthSuccess() {
		w.started.Do(func() {
			log.Println("Authentication completed. Refresh worker started")
		})

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
	}

	time.AfterFunc(FollowerRefreshTime*time.Second, func() {
		w.Refresher()
	})
	w.isRunning = true
}

func (w *Worker) apiAuthSuccess() bool {
	if _, err := w.cache.Get(cache.UserAccessCode); err != nil {
		w.once.Do(func() {
			commands := []string{}
			switch runtime.GOOS {
			case "darwin":
				commands = append(commands, "/usr/bin/open")
			case "windows":
				commands = append(commands, "cmd", "/c", "start")
			}
			url := w.client.AuthURL()
			cmd := exec.Command(commands[0], append(commands[1:], url)...)
			cmd.Start()
			log.Printf("Please go to %s to authenticate your twitch account", url)
		})
		return false
	}

	return true
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

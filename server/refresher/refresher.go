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
	refreshTime := 1 * time.Second

	if w.apiAuthSuccess() {
		refreshTime = FollowerRefreshTime * time.Second
		w.started.Do(func() {
			log.Println("Authentication completed. Refresh worker started")
		})

		newFollower, err := w.currentFollower()
		if err != nil {
			log.Printf("Failed to get current follower with %v\n", err)
		}

		if newFollower != nil && w.isRunning {
			log.Printf("New follower!! %s\n", newFollower.DisplayName)
			w.event.Send(stream.NewFollower, newFollower.DisplayName)
		}

		newSubscriber, err := w.currentSubscriber()
		if err != nil {
			log.Printf("failed to get new subscriber with %s", err)
		}

		if newSubscriber != nil && w.isRunning {
			log.Printf("New subscriber!! %s\n", newSubscriber.DisplayName)
			w.event.Send(stream.NewSubscriber, newSubscriber.DisplayName)
		}
		w.isRunning = true
	}

	time.AfterFunc(refreshTime, func() {
		w.Refresher()
	})
}

func (w *Worker) apiAuthSuccess() bool {
	if _, err := w.cache.Get(cache.UserAccessCode); err != nil {
		w.once.Do(func() {
			commands := []string{}
			url := w.client.AuthURL()
			log.Printf("Please go to %s to authenticate your twitch account", url)

			switch runtime.GOOS {
			case "darwin":
				commands = append(commands, "/usr/bin/open")
			case "windows":
				commands = append(commands, "cmd", "/c", "start")
			default:
				log.Println("On a browser")
				return
			}
			cmd := exec.Command(commands[0], append(commands[1:], url)...)
			err := cmd.Start()
			if err != nil {
				log.Printf("failed to run command with %s", err)
			}
			log.Printf("Please go to %s to authenticate your twitch account", url)
		})
		return false
	}

	return true
}

func (w *Worker) currentFollower() (*twitch.User, error) {

	currentFollowers, err := w.client.Channel.Followers(w.conf.Twitch.ChannelID, 1)

	if err != nil {
		return nil, fmt.Errorf("failed to get current followers")
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

func (w *Worker) currentSubscriber() (*twitch.User, error) {

	currentSubscribers, err := w.client.Channel.Subscribers(w.conf.Twitch.ChannelID, 1)

	if err != nil {
		return nil, fmt.Errorf("failed to get current followers")
	}
	if len(currentSubscribers.Subscriptions) == 0 {
		return nil, nil
	}

	currentSubscriber := currentSubscribers.Subscriptions[0]
	currentID := currentSubscriber.User.ID

	oldSubscriberID, _ := w.cache.Get(cache.LastSubscribeIDKey)

	if oldSubscriberID != currentID {
		w.cache.Set(cache.LastSubscribeIDKey, currentID)
		w.cache.Set(cache.LastSubscribeNameKey, currentSubscriber.User.DisplayName)
		w.cache.Set(cache.TotalSubscribersKey, strconv.Itoa(currentSubscribers.Total))
		return &currentSubscriber.User, nil
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

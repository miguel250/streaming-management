package goals

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/miguel250/streaming-setup/server/cache"
	"github.com/miguel250/streaming-setup/server/config"
)

type Goals struct {
	conf  *config.Config
	cache *cache.Cache
}

func (api *Goals) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	cacheFollowerCount, _ := api.cache.Get(cache.TotalFollowerKey)
	totalFollowerCount, _ := strconv.Atoi(cacheFollowerCount)
	followerName, _ := api.cache.Get(cache.LastFollowerNameKey)
	disableFollowerGoal := false

	if api.conf.Twitch.FollowerGoalTotal == 0 {
		disableFollowerGoal = true
	}

	cacheTotalSubscribers := "0"
	totalSubscribers, _ := strconv.Atoi(cacheTotalSubscribers)
	subscriberName := ""
	disableSubscriberGoal := false

	if api.conf.Twitch.SubscriberGoalTotal == 0 {
		disableSubscriberGoal = true
	}

	response := struct {
		DisableSubscriberGoal bool   `json:"disable_subscriber_follower_goal"`
		DisableFollowerGoal   bool   `json:"disable_follower_goal"`
		FollowerName          string `json:"follower_name"`
		SubscriberName        string `json:"subscriber_name"`
		FollowerGoal          int    `json:"follower_goal"`
		FollowerTotal         int    `json:"follower_total"`
		SubscriberGoal        int    `json:"subscribe_goal"`
		SubscriberTotal       int    `json:"subscriber_total"`
	}{
		disableSubscriberGoal,
		disableFollowerGoal,
		followerName,
		subscriberName,
		api.conf.Twitch.FollowerGoalTotal,
		totalFollowerCount,
		api.conf.Twitch.SubscriberGoalTotal,
		totalSubscribers,
	}

	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Printf("failed to encode json with %s", err)
		http.Error(rw, "Server error", http.StatusInternalServerError)
	}
}

func New(conf *config.Config, cache *cache.Cache) *Goals {
	return &Goals{
		conf,
		cache,
	}
}

package goals

import (
	"encoding/json"
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
	cacheCount, _ := api.cache.Get(cache.TotalFollowerKey)
	totalCount, _ := strconv.Atoi(cacheCount)
	followerName, _ := api.cache.Get(cache.LastFollowerNameKey)
	DisableFollowerGoal := false

	if api.conf.Twitch.FollowerGoalTotal == 0 {
		DisableFollowerGoal = true
	}

	response := struct {
		DisableFollowerGoal bool   `json:"disable_follower_goal"`
		FollowerName        string `json:"follower_name"`
		Goal                int    `json:"goal"`
		Total               int    `json:"total"`
	}{

		DisableFollowerGoal,
		followerName,
		api.conf.Twitch.FollowerGoalTotal,
		totalCount,
	}

	json.NewEncoder(rw).Encode(response)
}

func New(conf *config.Config, cache *cache.Cache) *Goals {
	return &Goals{
		conf,
		cache,
	}
}

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

	response := struct {
		FollowerName string `json:"follower_name"`
		Goal         int    `json:"goal"`
		Total        int    `json:"total"`
	}{
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

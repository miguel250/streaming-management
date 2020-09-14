package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/miguel250/streaming-setup/server/cache"
	"github.com/miguel250/streaming-setup/server/config"
	"github.com/miguel250/streaming-setup/server/twitch"
)

type Auth struct {
	conf      *config.Config
	cache     *cache.Cache
	twitchAPI *twitch.API
}

func (api *Auth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")

	if code == "" {
		http.Error(rw, "missing code", http.StatusBadRequest)
		return
	}

	resp, err := api.twitchAPI.AuthUser(code)

	if err != nil {
		log.Printf("Failed to get access token with %s", err)
		http.Error(rw, "invalid response from twitch API", http.StatusInternalServerError)
		return
	}

	api.cache.SetAccessToken(resp.AccessToken, resp.RefreshToken, resp.ExpiresIn)
	fmt.Fprintf(rw, `<html lang="en"><head><meta charset="utf-8"></head><body>This window can be close.</body></html>`)
}

func New(conf *config.Config, twitchAPI *twitch.API, cache *cache.Cache) *Auth {
	return &Auth{
		conf,
		cache,
		twitchAPI,
	}
}

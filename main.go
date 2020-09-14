package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/miguel250/kuma/http/server"
	"github.com/miguel250/streaming-setup/server/api/auth"
	"github.com/miguel250/streaming-setup/server/api/goals"
	"github.com/miguel250/streaming-setup/server/cache"
	"github.com/miguel250/streaming-setup/server/config"
	"github.com/miguel250/streaming-setup/server/irc"
	"github.com/miguel250/streaming-setup/server/refresher"
	"github.com/miguel250/streaming-setup/server/stream"
	"github.com/miguel250/streaming-setup/server/twitch"
	"github.com/miguel250/streaming-setup/server/twitchemotes"
)

func main() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("obs-assets/overlays"))
	mux.Handle("/overlays/", http.StripPrefix("/overlays", fs))

	event := stream.New()
	mux.Handle("/events", event)
	event.Start()
	defer event.Close()

	conf, err := config.New("streaming_config.json")

	if err != nil {
		log.Fatalf("Failed to load configuration file with %s", err)
	}

	srvConf := &server.Config{
		Addr: "localhost",
		Port: 8080,
	}

	srv := server.New(srvConf, mux)
	err = srv.Start()
	if err != nil {
		log.Fatalf("Failed to start server with %s", err)
	}

	twitchConf := &twitch.Config{
		AuthURL:     conf.Twitch.AuthURL,
		TwitchURL:   conf.Twitch.APIURL,
		RedirectURL: fmt.Sprintf("%s/api/auth", srv.Addr),
		BadgeURL:    conf.Twitch.BadgesURL,
		ClientID:    conf.Twitch.ClientID,
		Secret:      conf.Twitch.Secret,
	}

	c := cache.New()
	apiClient, err := twitch.New(twitchConf, c)

	fmt.Println(apiClient.AuthURL())

	if err != nil {
		log.Fatalf("Failed to create Twitch API client with %s", err)
	}

	globalBadges, err := apiClient.GetGlobalBadges()

	if err != nil {
		log.Fatalf("Failed to load badges for channel %s", err)
	}

	channelBadges, err := apiClient.Channel.GetBadges(conf.Twitch.ChannelID)

	if err != nil {
		log.Fatalf("Failed to load badges for channel %s", err)
	}

	for key, val := range channelBadges {
		globalBadges[key] = val
	}

	worker := refresher.New(conf, c, apiClient, event)
	worker.Refresher()

	mux.Handle("/api/goals", goals.New(conf, c))
	mux.Handle("/api/auth", auth.New(conf, apiClient, c))
	emotesAPI, err := twitchemotes.New(conf.Twitch.Emote.URL)

	if err != nil {
		log.Fatalf("Failed to create instance of emote API with %s", err)
	}

	ircConf := conf.Twitch.IRC
	ircConf.Badges = globalBadges
	ircConf.TwitchAPI = apiClient
	ircConf.TwitchEmotes = emotesAPI
	chatClient, err := irc.New(&ircConf)

	if err != nil {
		log.Fatalf("Failed to connect to twitch server: %s", err)
	}

	log.Println("Chat auth")
	chatClient.Start()
	chatClient.Auth()

	go func() {
		for {
			msg := <-chatClient.OnMessage

			if msg.DisplayName == conf.Twitch.IRC.Name {
				continue
			}

			b, err := json.Marshal(msg)
			if err != nil {
				log.Println("error:", err)
			}

			e := stream.Message{
				Type: stream.NewChatMessage,
				Text: string(b),
			}

			event.Message <- e
		}
	}()

	srv.StartAndWait()
}

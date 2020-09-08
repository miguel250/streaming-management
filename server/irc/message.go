package irc

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/miguel250/streaming-setup/server/irc/parser"
	"github.com/miguel250/streaming-setup/server/twitch"
)

func (c *Client) handleEmotes(parse *parser.Message) {
	emoteTag := parse.Tags["emotes"]

	if emoteTag == "" {
		return
	}

	emotes := strings.Split(emoteTag, "/")
	chatEmotes := make([]string, 0)
	cacheMiss := make([]string, 0)

	for _, val := range emotes {
		emoteParts := strings.Split(val, ":")
		emoteID := emoteParts[0]
		chatEmotes = append(chatEmotes, emoteParts[0])

		c.RLock()
		if _, ok := c.emotesCache[emoteID]; !ok {
			cacheMiss = append(cacheMiss, emoteID)
		}
		c.RUnlock()
	}

	if len(cacheMiss) > 0 {
		data, err := c.twitchEmotes.Emotes.GetByID(cacheMiss)

		if err != nil {
			fmt.Printf("Failed to get emote with %s\n", err)
			return
		}

		for _, val := range data {
			emoteURL := fmt.Sprintf("https://static-cdn.jtvnw.net/emoticons/v1/%d/2.0", val.ID)
			e := &emote{
				id:       strconv.Itoa(val.ID),
				code:     val.Code,
				imageURL: emoteURL,
			}
			c.Lock()
			c.emotesCache[e.id] = e
			c.Unlock()
		}
	}

	for _, val := range chatEmotes {
		if emote, ok := c.emotesCache[val]; ok {
			imgElem := fmt.Sprintf("<img src='%s'>", emote.imageURL)
			parse.Message = strings.ReplaceAll(parse.Message, emote.code, imgElem)
		}
	}
}

func (c *Client) handleBadges(parse *parser.Message) ([]*twitch.Badge, error) {
	badgeTags, ok := parse.Tags["badges"]

	if !ok || badgeTags == "" {
		return nil, nil
	}

	badges := strings.Split(badgeTags, ",")
	userBadges := make([]*twitch.Badge, 0, len(badges))

	for _, badge := range badges {
		badgeSlice := strings.Split(badge, "/")

		if len(badgeSlice) < 2 {
			log.Printf("IRC: unable to parse badge %s", badge)
			continue
		}

		name := badgeSlice[0]
		version := badgeSlice[1]

		if badgeVersions, ok := c.badges[name]; ok {
			if badge, ok := badgeVersions.Versions[version]; ok {
				userBadges = append(userBadges, badge)
			}
		}
	}
	return userBadges, nil
}

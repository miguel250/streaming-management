package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	channelPath       = "/kraken/channels"
	userPath          = "/kraken/users"
	channelFollows    = "/follows"
	globalBadgesPath  = "/v1/badges/global/display"
	channelBadgesPath = "/v1/badges/channels"
)

type API struct {
	clientID string
	url      *url.URL
	badgeURL *url.URL
	Channel  *Channel
}

type Channel struct {
	api *API
}

type User struct {
	ID          string `json:"_id"`
	DisplayName string `json:"display_name"`
	Logo        string `json:"logo"`
}

func (api *API) GetUser(id string) (*User, error) {
	path := fmt.Sprintf("%s/%s", userPath, id)

	resp, err := api.makeGetRequest(api.url, path, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to make request to twitch with %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to parse body for followers with %w", err)
	}

	user := &User{}
	err = json.Unmarshal(body, user)

	if err != nil {
		return nil, fmt.Errorf("failed to parse json for followers with %w", err)
	}

	return user, nil
}

type TwitchChannelResponse struct {
	Total   int        `json:"_total"`
	Follows []UserInfo `json:"follows"`
}

type UserInfo struct {
	User User `json:"user"`
}

func (c *Channel) Followers(channelID string, limit int) (*TwitchChannelResponse, error) {
	path := fmt.Sprintf("%s/%s%s", channelPath, channelID, channelFollows)
	queryParam := map[string]string{
		"limit": strconv.Itoa(limit),
	}

	resp, err := c.api.makeGetRequest(c.api.url, path, queryParam)

	if err != nil {
		return nil, fmt.Errorf("failed to make request to twitch with %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to parse body for followers with %w", err)
	}

	responseData := &TwitchChannelResponse{}
	err = json.Unmarshal(body, responseData)

	if err != nil {
		return nil, fmt.Errorf("failed to parse json for followers with %w", err)
	}

	return responseData, nil
}

type BadgesResponse struct {
	BadgeSet map[string]*BadgeVersion `json:"badge_sets"`
}

type BadgeVersion struct {
	Versions map[string]*Badge `json:"versions"`
}

type Badge struct {
	Title   string `json:"title"`
	Image1X string `json:"image_url_1x"`
	Image2X string `json:"image_url_2x"`
	Image4X string `json:"image_url_4x"`
}

func (c *Channel) GetBadges(channelID string) (map[string]*BadgeVersion, error) {
	path := fmt.Sprintf("%s/%s/%s", channelBadgesPath, channelID, "display")
	return getBadges(c.api, path)
}

func (api *API) GetGlobalBadges() (map[string]*BadgeVersion, error) {
	return getBadges(api, globalBadgesPath)
}

func getBadges(api *API, path string) (map[string]*BadgeVersion, error) {
	badgesResponse := &BadgesResponse{}

	resp, err := api.makeGetRequest(api.badgeURL, path, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to get twitch badges with %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to parse body for badges with %w", err)
	}

	err = json.Unmarshal(body, &badgesResponse)

	if err != nil {
		return nil, fmt.Errorf("failed to parse json for badges with %w", err)
	}
	return badgesResponse.BadgeSet, nil
}

func (c *API) makeGetRequest(requestURL *url.URL, path string, query map[string]string) (*http.Response, error) {
	u := *requestURL
	q := u.Query()

	for key, val := range query {
		q.Set(key, val)
	}

	u.RawQuery = q.Encode()
	u.Path = path

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request with %w", err)
	}

	req.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Add("Client-ID", c.clientID)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to make request with %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to make with status code %d", resp.StatusCode)
	}
	return resp, nil
}

func New(conf *Config) (*API, error) {
	if conf == nil {
		return nil, ErrNilConf
	}

	if err := conf.validate(); err != nil {
		return nil, err
	}

	u, err := url.Parse(conf.TwitchURL)

	if err != nil {
		return nil, fmt.Errorf("twitch: Invalid URL %w", err)
	}

	badgeURL, err := url.Parse(conf.BadgeURL)

	if err != nil {
		return nil, fmt.Errorf("twitch: Invalid URL %w", err)
	}

	api := &API{
		url:      u,
		clientID: conf.ClientID,
		badgeURL: badgeURL,
	}

	api.Channel = &Channel{api: api}

	return api, nil
}

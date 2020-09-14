package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/miguel250/streaming-setup/server/cache"
)

const (
	channelPath       = "/kraken/channels"
	userPath          = "/kraken/users"
	channelFollows    = "/follows"
	globalBadgesPath  = "/v1/badges/global/display"
	channelBadgesPath = "/v1/badges/channels"
	authPath          = "/oauth2/token"
)

type API struct {
	client      *http.Client
	authClient  *http.Client
	clientID    string
	secret      string
	url         *url.URL
	authURL     *url.URL
	badgeURL    *url.URL
	redirectURL *url.URL
	Channel     *Channel
}

type Channel struct {
	api *API
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type User struct {
	ID          string `json:"_id"`
	DisplayName string `json:"display_name"`
	Logo        string `json:"logo"`
}

func (api *API) AuthUser(code string) (*AuthResponse, error) {
	return api.auth(code, "authorization_code")
}

func (api *API) AuthTokenRefresh(refreshToken string) (*AuthResponse, error) {
	return api.auth(refreshToken, "refresh_token")
}

func (api *API) auth(code string, grantType string) (*AuthResponse, error) {
	queryParams := map[string]string{
		"client_id":     api.clientID,
		"client_secret": api.secret,
		"grant_type":    grantType,
		"redirect_uri":  api.redirectURL.String(),
	}

	switch grantType {
	case "authorization_code":
		queryParams["code"] = code
	case "refresh_token":
		queryParams["refresh_token"] = code
	default:
		return nil, fmt.Errorf("invalid grant type %s", grantType)
	}

	req := &request{
		method:      "POST",
		url:         api.authURL,
		path:        authPath,
		queryParams: queryParams,
	}

	resp, err := api.handleRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token with %w", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body for followers with %w", err)
	}

	authResp := &AuthResponse{}
	err = json.Unmarshal(body, authResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json with %w", err)
	}

	return authResp, nil
}

func (api *API) AuthURL() string {
	u := *api.authURL
	q := u.Query()

	q.Set("response_type", "code")
	q.Set("client_id", api.clientID)
	q.Set("redirect_uri", api.redirectURL.String())
	q.Set("scope", strings.Join([]string{"channel_subscriptions", "channel_read"}, " "))

	u.RawQuery = q.Encode()
	u.Path = "/oauth2/authorize"
	return u.String()
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

type SubscribersResponse struct {
	Total         int `json:"_total"`
	Subscriptions []*UserInfo
}

func (c *Channel) Subscribers(channelID string, limit int) (*SubscribersResponse, error) {
	path := fmt.Sprintf("%s/%s/subscriptions", channelPath, channelID)
	req := &request{
		client: c.api.authClient,
		method: "GET",
		url:    c.api.url,
		path:   path,
		queryParams: map[string]string{
			"limit":     strconv.Itoa(limit),
			"direction": "desc",
		},
	}

	resp, err := c.api.handleRequest(req)
	if err != nil {
		return nil, fmt.Errorf("twitch failed to make request with %w", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body for subscriptions with %w", err)
	}

	subResp := &SubscribersResponse{}
	err = json.Unmarshal(body, subResp)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response for subscriptions with %w", err)
	}

	return subResp, nil
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
	return c.handleRequest(&request{
		method:      "GET",
		path:        path,
		url:         requestURL,
		queryParams: query,
	})
}

func (c *API) makePostRequest(requestURL *url.URL, path string, query map[string]string, body io.Reader) (*http.Response, error) {
	return c.handleRequest(&request{
		method:      "POST",
		path:        path,
		url:         requestURL,
		body:        body,
		queryParams: query,
	})
}

type request struct {
	client      *http.Client
	method      string
	path        string
	url         *url.URL
	queryParams map[string]string
	headers     map[string]string
	body        io.Reader
}

func (c *API) handleRequest(req *request) (*http.Response, error) {
	u := *req.url
	q := u.Query()

	for key, val := range req.queryParams {
		q.Set(key, val)
	}

	u.RawQuery = q.Encode()
	u.Path = req.path

	httpReq, err := http.NewRequest(req.method, u.String(), req.body)

	if err != nil {
		return nil, fmt.Errorf("failed to create request with %w", err)
	}

	httpReq.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
	httpReq.Header.Add("Client-ID", c.clientID)

	for key, val := range req.headers {
		httpReq.Header.Add(key, val)
	}

	client := req.client
	if client == nil {
		client = c.client
	}

	resp, err := client.Do(httpReq)

	if err != nil {
		return nil, fmt.Errorf("failed to make request with %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to make with status code %d - %s", resp.StatusCode, string(b))
	}
	return resp, nil
}

func New(conf *Config, c *cache.Cache) (*API, error) {
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

	authURL, err := url.Parse(conf.AuthURL)

	if err != nil {
		return nil, fmt.Errorf("twitch: Invalid URL %w", err)
	}

	redirectURL, err := url.Parse(conf.RedirectURL)

	if err != nil {
		return nil, fmt.Errorf("twitch: Invalid URL %w", err)
	}

	api := &API{
		url:         u,
		secret:      conf.Secret,
		clientID:    conf.ClientID,
		badgeURL:    badgeURL,
		authURL:     authURL,
		redirectURL: redirectURL,
		client:      &http.Client{},
	}

	api.Channel = &Channel{api: api}
	tr := newTransport(api, c)
	api.authClient = &http.Client{Transport: tr}
	return api, nil
}

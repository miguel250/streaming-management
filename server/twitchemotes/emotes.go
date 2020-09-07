package twitchemotes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	initialPath = "/api/v4"
	emotes      = "emotes"
)

type API struct {
	url    *url.URL
	Emotes *Emotes
}

type Emotes struct {
	api *API
}

type Response struct {
	Code string `json:"code"`
	ID   int    `json:"id"`
}

func (e *Emotes) GetByID(emoteIDs []string) ([]*Response, error) {
	u := *e.api.url
	q := u.Query()
	ids := strings.Join(emoteIDs, ",")
	q.Set("id", ids)
	u.RawQuery = q.Encode()

	u.Path = fmt.Sprintf("%s/%s", initialPath, emotes)

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request with %w", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to get followers with %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get followers with status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to parse body for followers with %w", err)
	}

	responseData := make([]*Response, 0, len(emoteIDs))
	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return nil, fmt.Errorf("failed to parse json for followers with %w", err)
	}

	return responseData, nil
}

func New(urlStr string) (*API, error) {
	u, err := url.Parse(urlStr)

	if err != nil {
		return nil, fmt.Errorf("twitch emotes: Invalid URL %w", err)
	}

	api := &API{
		url: u,
	}

	api.Emotes = &Emotes{api: api}

	return api, nil
}

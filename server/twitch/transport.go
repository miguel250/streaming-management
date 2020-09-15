package twitch

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/miguel250/streaming-setup/server/cache"
)

type transport struct {
	sync.Mutex
	baseURL *url.URL
	api     *API
	cache   *cache.Cache
	tr      http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.token()

	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", token))
	return t.tr.RoundTrip(req)
}

func (t *transport) token() (string, error) {
	t.Lock()
	defer t.Unlock()
	expiresAtStr, err := t.cache.Get(cache.UserAccessExpiresAt)
	if err != nil {
		return "", err
	}

	expiresAt, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", expiresAtStr)
	if err != nil {
		return "", err
	}

	token, err := t.cache.Get(cache.UserAccessCode)
	if err != nil {
		return "", err
	}

	if expiresAt.Add(-time.Minute).Before(time.Now()) {
		refreshToken, err := t.cache.Get(cache.UserRefreshCode)
		if err != nil {
			return "", err
		}
		resp, err := t.api.AuthTokenRefresh(refreshToken)

		if err != nil {
			return "", fmt.Errorf("failed to refresh token with %w", err)
		}
		token = resp.AccessToken
		t.cache.SetAccessToken(token, resp.RefreshToken, resp.ExpiresIn)
	}

	return token, nil
}

func newTransport(api *API, c *cache.Cache) *transport {
	return &transport{
		api:   api,
		tr:    http.DefaultTransport,
		cache: c,
	}
}

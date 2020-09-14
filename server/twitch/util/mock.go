package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miguel250/streaming-setup/server/cache"
	"github.com/miguel250/streaming-setup/server/twitch"
)

func TestCreateClient(t *testing.T, responsePath string, endpoint string, channeID string) (*twitch.API, *httptest.Server) {
	return TestCreateClientQueryParams(t, responsePath, endpoint, channeID, nil, nil, 3600)
}

func TestCreateClientAuth(t *testing.T, responsePath string, endpoint string, channeID string, headers map[string]string, expiresAt int64) (*twitch.API, *httptest.Server) {
	return TestCreateClientQueryParams(t, responsePath, endpoint, channeID, nil, headers, expiresAt)
}

//TODO: move to using t.Cleanup(). Don't return httptest server.
func TestCreateClientQueryParams(t *testing.T, responsePath string, endpoint string, channeID string, queryParams map[string]string, headers map[string]string, expiresAt int64) (*twitch.API, *httptest.Server) {
	clientID := "test_client_id"

	ts := TestServerQueryParam(clientID, endpoint, responsePath, t, queryParams, headers)

	conf := &twitch.Config{
		TwitchURL:   ts.URL,
		ClientID:    clientID,
		BadgeURL:    ts.URL,
		AuthURL:     ts.URL,
		RedirectURL: "http://localhost/api/auth",
		Secret:      "nyo51xcdrerl8z9m56w9w6wg",
	}

	c := cache.New()
	c.SetAccessToken(
		"test_access_token",
		"test_refres_token",
		expiresAt,
	)

	api, err := twitch.New(conf, c)

	if err != nil {
		t.Fatalf("Failed to create API struct %v", err)
	}
	return api, ts
}

func TestServer(clientID string, endpoint string, responsePath string, t *testing.T) *httptest.Server {
	return TestServerQueryParam(clientID, endpoint, responsePath, t, nil, nil)
}

func TestServerQueryParam(clientID string, endpoint string, responsePath string, t *testing.T, wantQueryParams map[string]string, headers map[string]string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, handler(t, responsePath, clientID, wantQueryParams, headers))

	if endpoint != "/oauth2/token" {
		mux.HandleFunc("/oauth2/token", handler(t, "refresh_response", clientID, nil, nil))
	}

	return httptest.NewServer(mux)
}

func handler(t *testing.T, responsePath string, clientID string, wantQueryParams, headers map[string]string) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", responsePath))

		if err != nil {
			t.Errorf("Failed to get response with %s", err)
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		}

		wantHanders := map[string]string{
			"Accept":    "application/vnd.twitchtv.v5+json",
			"Client-ID": clientID,
		}

		for key, val := range wantHanders {
			if gotHeaderValue := req.Header.Get(key); gotHeaderValue != val {
				http.Error(rw, "", http.StatusBadRequest)
				return
			}
		}

		for key, val := range headers {
			if gotHeaderValue := req.Header.Get(key); gotHeaderValue != val {
				t.Errorf("header %s don't match value got: %s, want: %s", key, gotHeaderValue, val)
				http.Error(rw, "", http.StatusBadRequest)
				return
			}
		}

		for key, val := range wantQueryParams {
			if gotValue := req.URL.Query().Get(key); val != gotValue {
				t.Errorf("Query param %s doesn't match got: %s, want: %s", key, gotValue, val)
			}
		}

		fmt.Fprintln(rw, string(body))
	}
}

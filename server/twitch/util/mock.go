package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miguel250/streaming-setup/server/twitch"
)

//TODO: move to using t.Cleanup(). Don't return httptest server.
func TestCreateClient(t *testing.T, responsePath string, endpoint string, channeID string) (*twitch.API, *httptest.Server) {
	clientID := "test_client_id"

	ts := TestServer(clientID, endpoint, responsePath, t)

	conf := &twitch.Config{
		TwitchURL: ts.URL,
		ClientID:  clientID,
		BadgeURL:  ts.URL,
	}

	api, err := twitch.New(conf)

	if err != nil {
		t.Fatalf("Failed to create API struct %v", err)
	}
	return api, ts
}

func TestServer(clientID string, endpoint string, responsePath string, t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, func(rw http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", responsePath))

		if err != nil {
			t.Errorf("Failed to get response with %s", err)
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		}

		gotAcceptHeader := req.Header.Get("Accept")
		wantAcceptHeader := "application/vnd.twitchtv.v5+json"

		if gotAcceptHeader != "application/vnd.twitchtv.v5+json" {
			t.Errorf("Missing accept header got: %s, want: %s", gotAcceptHeader, wantAcceptHeader)
		}

		gotClientIDHeader := req.Header.Get("Client-ID")

		if gotClientIDHeader != clientID {
			http.Error(rw, "", http.StatusBadRequest)
			return
		}

		fmt.Fprintln(rw, string(body))
	})

	return httptest.NewServer(mux)
}

package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miguel250/streaming-setup/server/twitchemotes"
)

func TestCreateClient(t *testing.T, responsePath string, expectedQueryParams []string) *twitchemotes.API {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/emotes", func(rw http.ResponseWriter, req *http.Request) {

		queryParams := req.URL.Query()
		emoteIDs, ok := queryParams["id"]

		if !ok {
			t.Fatalf("Missing id query param")
		}

		for _, expectedID := range expectedQueryParams {
			found := false

			for _, gotID := range emoteIDs {
				if gotID == expectedID {
					found = true
				}
			}

			if !found {
				http.Error(rw, "", http.StatusNotFound)
				return
			}
		}

		body, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", responsePath))

		if err != nil {
			t.Errorf("Failed to get response with %s", err)
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		}

		fmt.Fprintln(rw, string(body))
	})

	ts := httptest.NewServer(mux)

	emoteClient, err := twitchemotes.New(ts.URL)

	if err != nil {
		t.Fatalf("failed to create twitch emotes client with: %s", err)
	}

	t.Cleanup(func() {
		ts.Close()
	})

	return emoteClient
}

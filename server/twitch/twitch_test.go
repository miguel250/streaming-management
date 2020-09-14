package twitch_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/miguel250/streaming-setup/server/cache"
	"github.com/miguel250/streaming-setup/server/twitch"
	"github.com/miguel250/streaming-setup/server/twitch/util"
)

var update = flag.Bool("update", false, "update golden files")

func TestInvalidURL(t *testing.T) {
	_, err := twitch.New(&twitch.Config{
		TwitchURL: "::////super_invalid_url",
		BadgeURL:  "::////super_invalid_url",
		ClientID:  "client-id",
	}, cache.New())

	if err == nil {
		t.Error("we are expecting an url parsing error here")
	}

	_, err = twitch.New(&twitch.Config{
		TwitchURL: "invalid_url",
		BadgeURL:  "::////super_invalid_url",
		ClientID:  "client-id",
	}, cache.New())

	if err == nil {
		t.Error("we are expecting an url parsing error here")
	}
}

func TestGetUser(t *testing.T) {
	channeID := "558843277"
	testEndpoint := "/kraken/users/1"
	api, ts := util.TestCreateClient(t, "user_response", testEndpoint, channeID)
	defer ts.Close()

	user, err := api.GetUser("1")

	if err != nil {
		t.Fatalf("failed to run GetUser with %s", err)
	}
	wantFilename := "testdata/user_expected.json"
	testCompareGoldenFiles(wantFilename, user, t)
}

func TestTwitchChannel(t *testing.T) {
	channeID := "558843277"
	testEndpoint := fmt.Sprintf("/kraken/channels/%s/follows", channeID)

	api, ts := util.TestCreateClient(t, "follower_response", testEndpoint, channeID)
	defer ts.Close()

	data, err := api.Channel.Followers(channeID, 1)

	if err != nil {
		t.Fatalf("Failed to %v", err)
	}

	if len(data.Follows) != 1 {
		t.Fatalf("Followers is not equal to 1, got: %d", len(data.Follows))
	}

	follower := data.Follows[0]
	wantFollowerID := "565688138"

	if follower.User.ID != wantFollowerID {
		t.Errorf("Follower ID don't match, want: %s got: %s", wantFollowerID, follower.User.ID)
	}

	wantDisplayName := "angelicahill95"

	if follower.User.DisplayName != wantDisplayName {
		t.Errorf("Follower DisplayName don't match, want: %s got: %s", wantDisplayName, follower.User.DisplayName)
	}
}

func TestSubscriptionChannel(t *testing.T) {
	channeID := "558843277"
	testEndpoint := fmt.Sprintf("/kraken/channels/%s/subscriptions", channeID)

	for _, test := range []struct {
		name        string
		wantHeaders map[string]string
		expiredAt   int64
	}{
		{
			"valid token",
			map[string]string{
				"Authorization": "OAuth test_access_token",
			},
			3600,
		},
		{
			"refresh token",
			map[string]string{
				"Authorization": "OAuth asdfasdf",
			},
			-900,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			api, ts := util.TestCreateClientAuth(t, "subscriptions_response", testEndpoint, channeID, test.wantHeaders, test.expiredAt)
			defer ts.Close()

			data, err := api.Channel.Subscribers(channeID, 1)

			if err != nil {
				t.Fatalf("Failed to %v", err)
			}

			if len(data.Subscriptions) != 1 {
				t.Fatalf("Followers is not equal to 1, got: %d", len(data.Subscriptions))
			}

			subscriber := data.Subscriptions[0]
			wantFollowerID := "202015635"

			if subscriber.User.ID != wantFollowerID {
				t.Errorf("Follower ID don't match, want: %s got: %s", wantFollowerID, subscriber.User.ID)
			}

			wantDisplayName := "williamconnelly"

			if subscriber.User.DisplayName != wantDisplayName {
				t.Errorf("Follower DisplayName don't match, want: %s got: %s", wantDisplayName, subscriber.User.DisplayName)
			}
		})
	}
}

func TestChannelGetBadges(t *testing.T) {
	channeID := "558843277"
	testEndpoint := fmt.Sprintf("/v1/badges/channels/%s/display", channeID)
	api, ts := util.TestCreateClient(t, "channel_badges_response", testEndpoint, channeID)
	defer ts.Close()

	badges, err := api.Channel.GetBadges(channeID)

	if err != nil {
		t.Fatalf("failed to get channel badges %s", err)
	}

	wantFilename := "testdata/channel_badges_expected.json"
	testCompareGoldenFiles(wantFilename, badges, t)
}

func TestGetGlobalBadges(t *testing.T) {
	channeID := "558843277"
	testEndpoint := "/v1/badges/global/display"
	api, ts := util.TestCreateClient(t, "global_badges_response", testEndpoint, channeID)
	defer ts.Close()

	badges, err := api.GetGlobalBadges()

	if err != nil {
		t.Fatalf("failed to get channel badges %s", err)
	}

	wantFilename := "testdata/global_badges_expected.json"
	testCompareGoldenFiles(wantFilename, badges, t)
}

func TestBadRequest(t *testing.T) {
	clientID := "test_client_id"
	endpoint := "/kraken/users/1"
	ts := util.TestServer(clientID, endpoint, "user_response", t)

	conf := &twitch.Config{
		TwitchURL:   ts.URL,
		Secret:      "not_really_secret",
		ClientID:    "test",
		BadgeURL:    ts.URL,
		AuthURL:     ts.URL,
		RedirectURL: ts.URL,
	}

	api, err := twitch.New(conf, cache.New())

	if err != nil {
		t.Fatalf("Failed to create API struct %v", err)
	}

	_, err = api.GetUser("1")

	if err == nil {
		t.Error("expected an error but didn't get one")
	}
}

func TestAuthUser(t *testing.T) {
	testEndpoint := "/oauth2/token"
	code := "394a8bc98028f39660e53025de824134fb46313"
	queryParams := map[string]string{
		"client_id":     "test_client_id",
		"client_secret": "nyo51xcdrerl8z9m56w9w6wg",
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  "http://localhost/api/auth",
	}
	api, ts := util.TestCreateClientQueryParams(t, "auth_response", testEndpoint, "", queryParams, nil, 100)
	defer ts.Close()

	resp, err := api.AuthUser(code)

	if err != nil {
		t.Fatalf("Failed to get user access token with: %s", err)
	}

	if resp == nil {
		t.Fatal("Response shouldn't be nil")
	}

	testCompareGoldenFiles("testdata/auth_token.json", resp, t)
}

func TestAuthTokenRefresh(t *testing.T) {
	testEndpoint := "/oauth2/token"
	code := "eyJfaWQmNzMtNGCJ9%6VFV5LNrZFUj8oU231/3Aj"
	queryParams := map[string]string{
		"client_id":     "test_client_id",
		"client_secret": "nyo51xcdrerl8z9m56w9w6wg",
		"refresh_token": code,
		"grant_type":    "refresh_token",
		"redirect_uri":  "http://localhost/api/auth",
	}
	api, ts := util.TestCreateClientQueryParams(t, "refresh_response", testEndpoint, "", queryParams, nil, 100)
	defer ts.Close()

	resp, err := api.AuthTokenRefresh(code)

	if err != nil {
		t.Fatalf("Failed to get user access token with: %s", err)
	}

	if resp == nil {
		t.Fatal("Response shouldn't be nil")
	}

	testCompareGoldenFiles("testdata/refresh_token.json", resp, t)
}

func testCompareGoldenFiles(wantFilename string, v interface{}, t *testing.T) {
	got, err := json.Marshal(v)

	if err != nil {
		t.Fatalf("Failed to marshal to json with %s", err)
	}

	if *update {
		err := ioutil.WriteFile(wantFilename, got, 0644)
		if err != nil {
			t.Fatalf("Failed to save golden file with %s", err)
		}
	}

	want, err := ioutil.ReadFile(wantFilename)

	if err != nil {
		t.Fatalf("Failed to open golden file with %s", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf("User doesn't match got (%s), want (%s)", string(got), string(want))
	}
}

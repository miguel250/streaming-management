package twitch_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/miguel250/streaming-setup/server/twitch"
	"github.com/miguel250/streaming-setup/server/twitch/util"
)

var update = flag.Bool("update", false, "update golden files")

func TestInvalidURL(t *testing.T) {
	_, err := twitch.New(&twitch.Config{
		TwitchURL: "::////super_invalid_url",
		BadgeURL:  "::////super_invalid_url",
		ClientID:  "client-id",
	})

	if err == nil {
		t.Error("we are expecting an url parsing error here")
	}

	_, err = twitch.New(&twitch.Config{
		TwitchURL: "invalid_url",
		BadgeURL:  "::////super_invalid_url",
		ClientID:  "client-id",
	})

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
		TwitchURL: ts.URL,
		ClientID:  "test",
		BadgeURL:  ts.URL,
	}

	api, err := twitch.New(conf)

	if err != nil {
		t.Fatalf("Failed to create API struct %v", err)
	}

	_, err = api.GetUser("1")

	if err == nil {
		t.Error("expected an error but didn't get one")
	}
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

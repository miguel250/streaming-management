package twitchemotes_test

import (
	"testing"

	"github.com/miguel250/streaming-setup/server/twitchemotes/util"
)

type emote struct {
	id   int
	code string
}

func TestEmotesAPI(t *testing.T) {
	for _, test := range []struct {
		name,
		fileName string
		expectedIDs []string
		queryIDs    []string
		wantErr     bool
		wantEmotes  []emote
	}{
		{
			"invalid json",
			"invalid_response",
			[]string{"303365132"},
			[]string{"303365132"},
			true,
			[]emote{},
		},
		{
			"not found",
			"emote_response",
			[]string{"303365132"},
			[]string{"303365"},
			true,
			[]emote{},
		},
		{
			"get emote by id",
			"emote_response",
			[]string{"303365132"},
			[]string{"303365132"},
			false,
			[]emote{{
				303365132,
				"miguel156Hero",
			}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			client := util.TestCreateClient(t, test.fileName, test.expectedIDs)
			emotes, err := client.Emotes.GetByID(test.queryIDs)

			if test.wantErr && err == nil {
				t.Fatalf("expected an error but got none")
			}

			if !test.wantErr {
				if err != nil {
					t.Fatalf("failed to get emotes with %s", err)
				}

				if len(emotes) != len(test.wantEmotes) {
					t.Fatalf("Expect emotes len to be equal got: %d", len(emotes))
				}

				for i, val := range test.wantEmotes {
					wantCode := val.code

					if wantCode != emotes[i].Code {
						t.Errorf("Emote code don't match want: %s, got: %s", wantCode, emotes[0].Code)
					}

					wantID := val.id
					if wantID != emotes[i].ID {
						t.Errorf("Emote ID don't match want: %d, got: %d", wantID, emotes[0].ID)
					}
				}
			}
		})
	}
}

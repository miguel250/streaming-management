package cache

import "testing"

func TestGetSet(t *testing.T) {
	wantKey := "Hello"
	wantValue := "World"

	c := New()

	c.Set(wantKey, wantValue)

	value, err := c.Get(wantKey)

	if err != nil {
		t.Fatalf("Failed to get key %s with err %s", wantKey, err)
	}

	if value != wantValue {
		t.Fatalf("Value didn't match want: %s, got: %s", wantValue, value)
	}
}

func TestParallel(t *testing.T) {
	wantKey := "Hello"
	wantValue := "World"

	wantKey2 := "HI"
	wantValue2 := "twitch"

	c := New()

	done := make(chan bool)

	go func() {
		c.Set(wantKey2, wantValue2)
		value, err := c.Get(wantKey2)

		if err != nil {
			t.Errorf("Failed to get key %s with err %s", wantKey2, err)
		}

		if value != wantValue2 {
			t.Errorf("Value didn't match want: %s, got: %s", wantValue2, value)
		}
		done <- true
	}()

	c.Set(wantKey, wantValue)

	value, err := c.Get(wantKey)

	if err != nil {
		t.Fatalf("Failed to get key %s with err %s", wantKey, err)
	}

	if value != wantValue {
		t.Fatalf("Value didn't match want: %s, got: %s", wantValue, value)
	}
	<-done
}

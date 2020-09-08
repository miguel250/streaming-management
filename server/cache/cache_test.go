package cache

import (
	"fmt"
	"sync"
	"testing"
)

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
	c := New()

	start := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			<-start
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)

			c.Set(key, value)
			actualValue, err := c.Get(key)

			if err != nil {
				t.Errorf("Failed to get key %s with err %s", key, err)
			}

			if actualValue != value {
				t.Errorf("Value didn't match want: %s, got: %s", value, actualValue)

			}
			wg.Done()
		}()
	}

	close(start)
	wg.Wait()
}

package watcher_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/escaletech/casbin-redis-watcher/watcher"
	"github.com/stretchr/testify/assert"
)

func TestWatcher(t *testing.T) {
	t.Run("updates other watchers", func(t *testing.T) {
		mr, err := miniredis.Run()
		assert.NoError(t, err)
		t.Cleanup(mr.Close)

		// set up blue watcher, which will be the source of the update
		blueWatcher, err := watcher.New(watcher.Options{
			RedisURL: "redis://" + mr.Addr(),
			LocalID:  "blue",
		})
		assert.NoError(t, err)
		t.Cleanup(blueWatcher.Close)

		// set up red watcher, which should be notified of the update
		redWatcher, err := watcher.New(watcher.Options{
			RedisURL: "redis://" + mr.Addr(),
			LocalID:  "red",
		})
		assert.NoError(t, err)
		t.Cleanup(redWatcher.Close)

		updateChan := make(chan string, 1)
		redWatcher.SetUpdateCallback(func(payload string) {
			updateChan <- payload
		})

		// update blue watcher
		err = blueWatcher.Update()
		assert.NoError(t, err)

		select {
		case msg := <-updateChan:
			assert.Equal(t, "blue", msg)
		case <-time.After(500 * time.Millisecond):
			assert.FailNow(t, "didn't detect update after 500ms")
		}
	})

	t.Run("does not update itself", func(t *testing.T) {
		mr, err := miniredis.Run()
		assert.NoError(t, err)
		t.Cleanup(mr.Close)

		// set up blue watcher, which will be the source of the update
		blueWatcher, err := watcher.New(watcher.Options{
			RedisURL: "redis://" + mr.Addr(),
			LocalID:  "blue",
		})
		assert.NoError(t, err)
		t.Cleanup(blueWatcher.Close)

		blueWatcher.SetUpdateCallback(func(payload string) {
			assert.FailNow(t, "blue watcher should not have been updated")
		})

		// update blue watcher
		err = blueWatcher.Update()
		assert.NoError(t, err)
	})

	t.Run("does not update after closed", func(t *testing.T) {
		mr, err := miniredis.Run()
		assert.NoError(t, err)
		t.Cleanup(mr.Close)

		// set up blue watcher, which will be the source of the update
		blueWatcher, err := watcher.New(watcher.Options{
			RedisURL: "redis://" + mr.Addr(),
			LocalID:  "blue",
		})
		assert.NoError(t, err)

		blueWatcher.Close()

		<-time.After(20 * time.Millisecond)
		assert.EqualError(t, blueWatcher.Update(), "redis: client is closed")
	})

	t.Run("returns error for", func(t *testing.T) {
		t.Run("invalid Redis address", func(t *testing.T) {
			_, err := watcher.New(watcher.Options{
				RedisURL: "not a redis URL",
			})
			assert.EqualError(t, err, "invalid redis URL scheme: ")
		})

		t.Run("Redis address of inexistent server", func(t *testing.T) {
			_, err := watcher.New(watcher.Options{
				RedisURL: "redis://not-a-redis-server:6379",
			})
			assert.EqualError(t, err, "dial tcp: lookup not-a-redis-server: no such host")
		})
	})
}

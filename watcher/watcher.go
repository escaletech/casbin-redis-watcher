package watcher

import (
	"context"
	"sync"

	"github.com/casbin/casbin/v2/persist"
	"github.com/go-redis/redis/v8"
)

var _ persist.Watcher = &Watcher{}

// New creates a Watcher
func New(opt Options) (*Watcher, error) {
	opt = opt.validate()

	w := &Watcher{
		opt: opt,
	}

	if err := w.init(); err != nil {
		return nil, err
	}

	return w, nil
}

// Watcher implements Casbin's persist.Watcher using Redis as a backend
type Watcher struct {
	publisher  redis.UniversalClient
	subscriber redis.UniversalClient
	callback   func(string)
	opt        Options
	once       sync.Once
	closeChan  chan struct{}
}

// SetUpdateCallback sets the callback function that the watcher will call
// when the policy in DB has been changed by other instances.
// A classic callback is Enforcer.LoadPolicy().
func (w *Watcher) SetUpdateCallback(callback func(string)) error {
	w.callback = callback
	return nil
}

// Update calls the update callback of other instances to synchronize their policy.
// It is usually called after changing the policy in DB, like Enforcer.SavePolicy(),
// Enforcer.AddPolicy(), Enforcer.RemovePolicy(), etc.
func (w *Watcher) Update() error {
	ctx := context.Background()
	if err := w.publisher.Publish(ctx, w.opt.Channel, w.opt.LocalID).Err(); err != nil {
		return err
	}

	return nil
}

// Close stops and releases the watcher, the callback function will not be called any more.
func (w *Watcher) Close() {
	w.once.Do(func() {
		close(w.closeChan)
	})
}

func (w *Watcher) init() error {
	ctx := context.Background()

	var err error
	w.publisher, err = w.opt.NewClient()
	if err != nil {
		return err
	}

	w.subscriber, err = w.opt.NewClient()
	if err != nil {
		return err
	}

	pubsub := w.subscriber.Subscribe(ctx, w.opt.Channel)

	if _, err := pubsub.Receive(ctx); err != nil {
		return err
	}

	w.closeChan = make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-pubsub.Channel():
				w.messageReceived(msg.Payload)

			case <-w.closeChan:
				pubsub.Close()
				w.subscriber.Close()
				w.publisher.Close()
				return
			}
		}
	}()

	return nil
}

func (w *Watcher) messageReceived(publisherID string) {
	if publisherID == w.opt.LocalID {
		// ignore messages from itself
		return
	}

	if w.callback != nil {
		w.callback(publisherID)
	}
}

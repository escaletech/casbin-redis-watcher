package watcher

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// Options are used to configure the Watcher
type Options struct {
	// RedisURL is used to establish a Redis connection (default: "redis://127.0.0.1:6379")
	RedisURL string

	// Cluster tells whether to use a Redis Cluster client or not (default: false)
	Cluster bool

	// Channel is the Redis channel for pub/sub (default: "/casbin")
	Channel string

	// LocalID is the identifier used for avoiding unnecessary updates (default is an auto-generated UUID)
	LocalID string

	// NewClient is the function used to create a new redis.UniversalClient (optional)
	NewClient func() (redis.UniversalClient, error)
}

func (opt Options) validate() Options {
	if opt.RedisURL == "" {
		opt.RedisURL = "redis://127.0.0.1:6379"
	}

	if opt.Channel == "" {
		opt.Channel = "/casbin"
	}

	if opt.LocalID == "" {
		opt.LocalID = uuid.New().String()
	}

	if opt.NewClient == nil {
		opt.NewClient = newRedisClient(opt.RedisURL, opt.Cluster)
	}

	return opt
}

package watcher

import "github.com/go-redis/redis/v8"

func newRedisClient(redisURL string, cluster bool) func() (redis.UniversalClient, error) {
	return func() (redis.UniversalClient, error) {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			return nil, err
		}

		if cluster {
			return redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    []string{opt.Addr},
				Username: opt.Username,
				Password: opt.Password,
			}), nil
		}

		return redis.NewClient(opt), nil
	}
}

package db

import (
	"fmt"
	"time"

	"github.com/mediocregopher/radix.v2/redis"
)

// Redis can communicate with a Redis instance.
type Redis struct {
	*redis.Client
}

// NewRedis returns an instance of db.Redis that communicates with the Redis instance running on the specified host.
func NewRedis(host, proto string, timeout time.Duration) (*Redis, error) {
	client, err := redis.DialTimeout(proto, host, timeout)
	if err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}

// Exist returns true if the given URL exists in the Redis instance.
// Otherwise, it returns false.
func (r *Redis) Exist(url string) (bool, error) {
	v, err := r.Cmd("GET", url).Str()
	if err != nil {
		fmt.Printf("%v\n", err)
		return false, err
	}

	fmt.Printf("%v\n", v)
	return false, nil
}

// Close terminates the connection to the Redis instance.
func (r *Redis) Close() error {
	return r.Client.Close()
}

package db

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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
	arg := strings.TrimSuffix(url, "\n")
	_, err := r.Cmd("GET", arg).Str()
	if err != nil {
		if err == redis.ErrRespNil {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Close terminates the connection to the Redis instance.
func (r *Redis) Close() error {
	return r.Client.Close()
}

// Load inserts data from the reader into the Redis instance.
func (r *Redis) Load(rd io.Reader) error {
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		arg := strings.TrimSuffix(scanner.Text(), "\n")
		r.PipeAppend("SET", arg, "")
	}

	var (
		cmdErr error
		allErr []string
	)
	for cmdErr != redis.ErrPipelineEmpty {
		cmdErr = r.PipeResp().Err
		if cmdErr != nil && cmdErr != redis.ErrPipelineEmpty {
			allErr = append(allErr, cmdErr.Error())
		}
	}

	if err := scanner.Err(); err != nil {
		allErr = append(allErr, err.Error())
	}

	if len(allErr) > 0 {
		return fmt.Errorf("%s", strings.Join(allErr, "\n"))
	}

	return nil
}

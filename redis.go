// Author:
// - Dong Fei <dongfei@hotpu.cn>
package redis

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/garyburd/redigo/redis"
)

var typeOfJsonRawMessage reflect.Type = reflect.TypeOf(json.RawMessage{})

const (
	EXPIRE_FLAG       int64 = -2
	NEVER_EXPIRE_FLAG int64 = -1

	HISTORY_EMP_LENGTH = 30
	HISTORY_SUB_LENGTH = 30
)

var ErrNotExist = errors.New("Not exists in redis")

type Builder struct {
	Server      string        // Redis server address
	Password    string        // Redis Server password
	MaxIdle     int64         // Maximum number of idle connections in the pool.
	MaxActive   int64         // Maximum number of connections allocated by the pool at a given time.When zero, there is no limit on the number of connections in the pool.
	IdleTimeout time.Duration // Close connections after remaining idle for this duration. If the value is zero, then idle connections are not closed. Applications should set the timeout to a value less than the server's timeout.
}

// cache object
type Cache struct {
	Builder *Builder
	Pool    *redis.Pool
}

func (builder *Builder) Build() (*Cache, error) {
	pool := &redis.Pool{
		MaxIdle:     int(builder.MaxIdle),
		MaxActive:   int(builder.MaxActive),
		IdleTimeout: builder.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", builder.Server)
			if err != nil {
				return nil, err
			}

			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	//test if we can connect to redis
	if _, err := pool.Dial(); err != nil {
		return nil, err
	}
	return &Cache{
		Builder: builder,
		Pool:    pool,
	}, nil
}

func (cache *Cache) Close() error {
	return cache.Pool.Close()
}

func (cache *Cache) GetConn() redis.Conn {
	conn := cache.Pool.Get()
	return conn
}

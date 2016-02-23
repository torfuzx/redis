package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// -- --------------------------------------------------------------------------
// -- Strings
// -- --------------------------------------------------------------------------

func (cache *Cache) DoSET(conn redis.Conn, key string, val interface{}, expire time.Duration) error {
	_, err := redis.String(conn.Do("SET", MyArgs{}.Add(key).Add(val)...))
	if err != nil {
		return err
	}

	// set expire when necessary(greater than 0)
	// NOTE: The expire is set for the whole hash, not for a single hash key.
	if err := cache.DoEXPIRE(conn, key, expire); err != nil {
		return err
	}

	return nil
}

func (cache *Cache) DoGET(conn redis.Conn, key string) (interface{}, error) {
	reply, err := conn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

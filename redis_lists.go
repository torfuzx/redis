package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// -- --------------------------------------------------------------------------
// -- Lists
// -- --------------------------------------------------------------------------

func (cache *Cache) DoRPUSH(conn redis.Conn, key string, value interface{}, expire time.Duration) error {
	args := MyArgs{}
	args = args.Add(key).Add(value)
	_, err := conn.Do("RPUSH", args...)

	if err != nil {
		return err
	}

	return cache.DoEXPIRE(conn, key, expire)
}

func (cache *Cache) DoLRANGE(conn redis.Conn, key string, start, stop int64) []interface{} {
	args := MyArgs{}
	args = args.Add(key).Add(start).Add(stop)

	vals, err := redis.Values((conn.Do("LRANGE", args...)))
	if err != nil {
		vals = []interface{}{}
	}
	return vals
}

func (cache *Cache) DoLTRIM(conn redis.Conn, key string, start, stop int64) error {
	args := MyArgs{}
	args = args.Add(key).Add(start).Add(stop)

	_, err := conn.Do("LTRIM", args...)
	return err
}

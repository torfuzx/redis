package redis

import (
	"github.com/garyburd/redigo/redis"
)

// -- --------------------------------------------------------------------------
// -- Pub/Sub
// -- --------------------------------------------------------------------------

func (cache *Cache) DoPublish(channel string, msg string) (int64, error) {
	conn := cache.GetConn()
	defer conn.Close()

	args := MyArgs{}.Add(channel).Add(msg)
	return redis.Int64(conn.Do("PUBLISH", args...))
}

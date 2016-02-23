package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"

	"hotpu.cn/xkefu/common/log"
)

// -- --------------------------------------------------------------------------
// -- Sorted Sets
// -- --------------------------------------------------------------------------

// The `score` and `member` must be given in pairs.
// Score is type int64, and member is type string.
// Both single member and multiple members are supported.
type ScoreMemberPair struct {
	Score  int64
	Member string
}

// Add a single member to a sorted set.
func (cache *Cache) DoZADD(conn redis.Conn, key string, expire time.Duration, pair *ScoreMemberPair) error {
	args := MyArgs{}.Add(key).Add(pair.Score).Add(pair.Member)
	reply, err := redis.Int64(conn.Do("ZADD", args...))
	if err != nil {
		return err
	}

	// set expire when necessary(greater than 0)
	if err := cache.DoEXPIRE(conn, key, expire); err != nil {
		return err
	}

	log.Debug("common.store.redis", "Cache.DoZADD", "%d members are added to the sorted sets: %s.", reply, key)

	return nil
}

// Add multiple members to a sorted list.
func (cache *Cache) DoMultiZADD(conn redis.Conn, key string, expire time.Duration, pairs []*ScoreMemberPair) error {
	args := MyArgs{}.Add(key)
	for _, pair := range pairs {
		args.Add(pair.Score).Add(pair.Member)
	}

	reply, err := redis.Int64(conn.Do("ZADD", args...))
	if err != nil {
		return err
	}

	// set expire when necessary(greater than 0)
	if err := cache.DoEXPIRE(conn, key, expire); err != nil {
		return err
	}

	log.Debug("common.store.redis", "Cache.DoMultiZADD", "%d members are added to the sorted sets: %s.", reply, key)

	return nil
}

func (cache *Cache) DoZRANGE(conn redis.Conn, key string, start int64, stop int64) ([]interface{}, error) {
	args := MyArgs{}.Add(key).Add(start).Add(stop)
	reply, err := redis.Values(conn.Do("ZRANGE", args...))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (cache *Cache) DoZREVRANGE(conn redis.Conn, key string, start, stop int64) ([]interface{}, error) {
	args := MyArgs{}.Add(key).Add(start).Add(stop)
	reply, err := redis.Values(conn.Do("ZREVRANGE", args...))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (cache *Cache) DoZREM(conn redis.Conn, key, member string) (int64, error) {
	args := MyArgs{}.Add(key).Add(member)
	num, err := redis.Int64(conn.Do("ZREM", args...))
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (cahce *Cache) DoZRANK(conn redis.Conn, key, member string) (int64, error) {
	args := MyArgs{}.Add(key).Add(member)
	return redis.Int64(conn.Do("ZRANK", args...))
}

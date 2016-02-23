package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"

	"hotpu.cn/xkefu/common/log"
)

// -- --------------------------------------------------------------------------
// -- Hashes
// -- --------------------------------------------------------------------------

func (cache *Cache) DoHMSET(conn redis.Conn, key string, data interface{}, expire time.Duration) error {
	ret, err := redis.String(conn.Do("HMSET", MyArgs{}.Add(key).AddFlat(data)...))
	if err != nil {
		return err
	}
	if ret != "OK" {
		log.Warn("common.store.redis", "Cache.DoHMSET", "HMSET should return OK on success man!!! (ret: %s)", ret)
	}

	// set expire when necessary(greater than 0)
	if err := cache.DoEXPIRE(conn, key, expire); err != nil {
		return err
	}

	return nil
}

func (cache *Cache) DoHSET(conn redis.Conn, key string, fieldKey string, fieldVal string, expire time.Duration) error {
	ret, err := redis.Int64(conn.Do("HSET", key, fieldKey, fieldVal))
	if err != nil {
		return err
	}

	if ret == 0 {
		log.Info("common.store.redis", "Cache.DoHSET", "The field %q already exists in the hash %q and the value was updated.", fieldKey, key)
	} else if ret == 1 {
		log.Debug("common.store.redis", "Cache.DoHSET", "The field %q is a new field in the hash %q and value was set.", fieldKey, key)
	}

	// set expire when necessary(greater than 0)
	// NOTE: The expire is set for the whole hash, not for a single hash key.
	if err := cache.DoEXPIRE(conn, key, expire); err != nil {
		return err
	}

	return nil
}

func (cache *Cache) DoHGETALL(conn redis.Conn, key string) ([]interface{}, error) {
	reply, err := conn.Do("HGETALL", key)
	if err != nil {
		return nil, err
	}
	return redis.Values(reply, err)
}

func (cache *Cache) DoHGET(conn redis.Conn, key, fieldKey string) (interface{}, error) {
	reply, err := conn.Do("HGET", key, fieldKey)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (cache *Cache) DoHDEL(conn redis.Conn, key, fieldKey string) (int64, error) {
	num, err := redis.Int64(conn.Do("HDEL", key, fieldKey))
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (cache *Cache) DoHINCRBY(conn redis.Conn, key, fieldKey string, increment int64) (int64, error) {
	val, err := redis.Int64(conn.Do("HINCRBY", key, fieldKey, increment))
	if err != nil {
		return 0, err
	}
	log.Info("common.store.redis", "Cache.DoHINCRBY", "The reply is: %d", val)
	return val, nil
}

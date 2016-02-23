package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"

	"hotpu.cn/xkefu/common/log"
)

// -- --------------------------------------------------------------------------
// -- Keys
// -- --------------------------------------------------------------------------

func (cache *Cache) DoEXPIRE(conn redis.Conn, key string, expire time.Duration) error {
	if expire > time.Duration(0) {
		ret, err := redis.Int64(conn.Do("EXPIRE", key, expire.Seconds()))
		if err != nil {
			return err
		}
		if ret == 1 {
			return nil
		} else {
			return fmt.Errorf("DoEXPIRE error: The key %q doesn't exist or the timeout could not be set")
		}
	}

	return nil
}

func (cache *Cache) DoDEL(conn redis.Conn, keys ...string) error {
	reply, err := redis.Int64(conn.Do("DEL", MyArgs{}.AddFlat(keys)...))
	if err != nil {
		return err
	}
	log.Debug("common.store.redis", "Cache.DoDEL", "%d keys were affected.", reply)

	return nil
}

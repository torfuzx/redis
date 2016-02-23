package redis

import (
	"math"

	"github.com/garyburd/redigo/redis"
)

// -- --------------------------------------------------------------------------
// -- Bitmaps
// -- --------------------------------------------------------------------------

// Return the original bit value stored at offset
func (cache *Cache) DoSETBIT(conn redis.Conn, key string, offset int64, bit int64) (int64, error) {
	args := MyArgs{}.Add(key).Add(offset).Add(bit)
	return redis.Int64(conn.Do("SETBIT", args...))
}

// Returns the bit value at offset in the string value stored at key
func (cache *Cache) DoGETBIT(conn redis.Conn, key string, offset int64) (int64, error) {
	args := MyArgs{}.Add(key).Add(offset)
	return redis.Int64(conn.Do("GETBIT", args...))
}

// Count the number of set bits (population counting) in a string.
func (cache *Cache) DoBITCOUNT(conn redis.Conn, key string, start, end int64) (int64, error) {
	start = int64(math.Floor(float64(start/8))) + 1
	end = int64(math.Floor(float64(end/8))) + 1

	args := MyArgs{}.Add(key).Add(start).Add(end)
	return redis.Int64(conn.Do("BITCOUNT", args...))
}

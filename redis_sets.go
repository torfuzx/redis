package redis

import "github.com/garyburd/redigo/redis"

// -- --------------------------------------------------------------------------
// -- Sets
// -- --------------------------------------------------------------------------

// Get all memebers in a set.
func (cache *Cache) DoSMEMBERS(conn redis.Conn, key string) ([]interface{}, error) {
	reply, err := redis.Values(conn.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// Returns the set cardinality (number of elements) of the set stored at key.
func (cache *Cache) DoSCARD(conn redis.Conn, key string) (int64, error) {
	reply, err := redis.Int64(conn.Do("SCARD", key))
	if err != nil {
		return 0, err
	}
	return reply, nil
}

// Returns if member is a member of the set stored at key.
func (cache *Cache) DoSISMEMBER(conn redis.Conn, key string, member string) (bool, error) {
	reply, err := redis.Bool(conn.Do("SISMEMBER", key, member))
	if err != nil {
		return false, err
	}
	return reply, nil
}

package redis

import "github.com/garyburd/redigo/redis"

// -- --------------------------------------------------------------------------
// -- Server
// -- --------------------------------------------------------------------------

func (cache *Cache) DoCLIENTSETNAME(conn redis.Conn, name string) (bool, error) {
	reply, err := redis.String(conn.Do("CLIENT", "SETNAME", name))
	if err != nil {
		return false, err
	}

	if reply == "OK" {
		return true, nil
	} else {
		return false, nil
	}
}

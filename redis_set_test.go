package redis

import (
	"fmt"
	. "gopkg.in/check.v1"
	"sync"
	//"testing"
	"time"
)

// To run this suite, execute:
// > go test -gocheck.f=RedisSetSuite
type RedisSetSuite struct {
	Cache *Cache
}

type Wrapper struct {
	sync.WaitGroup
}

var _ = Suite(&RedisSetSuite{})
var wg = new(sync.WaitGroup)

func (wg *Wrapper) Wrap(cb func()) {
	wg.Add(1)
	go func() {
		cb()
		wg.Done()
	}()
}

// -- --------------------------------------------------------------------------
// -- Fixtures
// -- --------------------------------------------------------------------------

func (s *RedisSetSuite) SetUpSuite(c *C) {
	cache, err := GetTestCache()
	c.Assert(err, IsNil)
	s.Cache = cache
}

// simulating consumer's redis lock
func (s *RedisSetSuite) Test(c *C) {
	key := "foo"
	now := fmt.Sprintf("%d", time.Now().Unix())
	isset, err := s.Cache.DoSETNX(key, now)
	c.Assert(err, IsNil)
	c.Assert(isset, Equals, int64(1))

	result, err := s.Cache.DoGETSET(key, fmt.Sprintf("%d", time.Now().Unix()))
	c.Assert(err, IsNil)
	c.Assert(now, Equals, result)

	entcode := "yto"
	subid := int64(10001)

	lockkey, err := GetKeySubScriberLock(entcode, subid)
	c.Assert(err, IsNil)
	c.Assert(lockkey, Equals, "yto:subscriber_lock:10001")

	_, err = s.Cache.DoDelete(key)
	c.Assert(err, IsNil)
}

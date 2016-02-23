//author:dongfei
package redis

import (
	"testing"

	. "gopkg.in/check.v1"
)

// The redis package leverages the gocheck to do the testing.
func Test(t *testing.T) { TestingT(t) }

func GetTestCache() (*Cache, error) {
	builer := &Builder{
		Server:      "127.0.0.1:6379",
		Password:    "",
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: 0,
	}
	return builer.Build()
}

//type Addr struct {
//	Where string
//	Code  string
//}
//
//func TestRedis(t *testing.T) {
//	builder := &Builder{
//		Server:      "localhost:6379",
//		MaxIdle:     3,
//		IdleTimeout: 240 * time.Second,
//	}
//	cache, err := builder.Build()
//	if err != nil {
//		t.Fatal("can not connect to redis")
//	}
//	// var ID int64 = 1
//	// enterprise_cache := CreateEnterpriseCache(ID, cache)
//
//	// var emps int64 = 60
//	// var name string = "hotpu"
//	// var email string = "xxxx@hotpu.cn"
//	// var timeout int64 = 10
//	// addr := &Addr{
//	// 	Where: "shanghai of china",
//	// 	Code:  "21321312",
//	// }
//
//	// if err := enterprise_cache.Set("emps", emps); err != nil { //set int
//	// 	t.Fatal(err.Error())
//	// }
//	// enterprise_cache.Set("name", name) //set string
//	// enterprise_cache.Set("email", email)
//	// enterprise_cache.Set("addr", addr) //set object
//	// enterprise_cache.SetWithTimeout("timeout", timeout, 100)
//
//	// e, err := enterprise_cache.GetInt("emps")  //get int
//	// n, _ := enterprise_cache.GetString("name") //get string
//	// em, _ := enterprise_cache.GetString("email")
//	// if err != nil {
//	// 	t.Fatal(err)
//	// }
//	// t.Logf("found <%d,%s,%s>", e, n, em)
//	// dest := &Addr{}
//	// err_get_object := enterprise_cache.GetObject("addr", dest) //get object
//	// if err_get_object != nil {
//	// 	t.Fatal(err_get_object)
//	// }
//	// t.Logf("given key addr, found<%s,%s>", dest.Where, dest.Code)
//	// if err := enterprise_cache.Expire("addr", 10); err != nil {
//	// 	t.Fatal(err)
//	// }
//	// ttl, err := enterprise_cache.TTL("addr")
//	// tt, _ := enterprise_cache.TTL("timeout")
//	// if err != nil {
//	// 	t.Fatal(err)
//	// }
//	// t.Logf("addr:time to live %d ====== timeout: %d", ttl, tt)
//	// enterprise_cache.Close() //put back connection to pool,don't forget!!!!
//	if err := cache.SetWithTimeout("name", "dongfei", 10); err != nil {
//		t.Fatal("SetWithTimeout error")
//	}
//	name, errGet := cache.GetString("name")
//	if err != nil {
//		t.Fatal("GetString error :", errGet.Error())
//	}
//	t.Log("name : ", name)
//
//	sysfolderKey := "sys:folder"
//	if err := cache.Rpush(sysfolderKey, 1); err != nil {
//		t.Fatalf(err.Error())
//	}
//
//	if folders, err := cache.List(sysfolderKey); err != nil {
//		t.Fatalf(err.Error())
//	} else {
//		t.Logf("%+v", folders)
//	}
//
//}
//
//func TestGetSource(t *testing.T) {
//	// set up
//	builder := &Builder{
//		Server:      "localhost:6379",
//		MaxIdle:     3,
//		IdleTimeout: 240 * time.Second,
//	}
//	cache, err := builder.Build()
//	if err != nil {
//		t.Fatal("can not connect to redis")
//	}
//
//	type ApiCustomResourceSetData struct {
//		Key      string          `json:"key"        mapstructure:"key"`
//		Value    json.RawMessage `json:"value"      mapstructure:"value"`
//		ExpireIn int64           `json:"expire_in"  mapstructure:"expire_in"`
//	}
//
//	key := "foo:bar:baz"
//	jsonStr := `
//	{
//	    "action": "customresource.set",
//	    "data": {
//	        "key": "foo",
//	        "value": {
//	            "name": "张三",
//	            "age": 23,
//	            "favourite_food": [
//	                "apple",
//	                "banana"
//	            ],
//	            "company_address": {
//	                "city": "Shanghai",
//	                "street": "WestRenmingRd."
//	            }
//	        },
//	        "expire_in": 3600
//	    }
//	}`
//	val := json.RawMessage(jsonStr)
//
//	err = cache.SetRawMessage(key, val, 3600)
//	assert.Nil(t, err)
//
//	// test
//	raw, err := cache.GetRawMessage(key)
//	assert.Nil(t, err)
//	assert.NotNil(t, raw)
//	b, err := raw.MarshalJSON()
//	assert.Nil(t, err)
//	assert.Equal(t, val, string(b))
//
//}

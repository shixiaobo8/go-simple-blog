package lib

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
)

var Cache cache.Cache

func init() {
	var err error
	redisHost := beego.AppConfig.String("RedisHost")
	redisPort, _ := beego.AppConfig.Int("RedisPort")
	redisAuth := beego.AppConfig.String("RedisAuth")

	redisCollectionName := "blog_cache"
	config := map[string]string{
		"key":      redisCollectionName,
		"conn":     fmt.Sprintf("%s:%d", redisHost, redisPort),
		"dbNum":    "0",
		"password": redisAuth,
	}

	configJson, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	Cache, err = cache.NewCache("redis", string(configJson))
	if err != nil {
		panic(err)
	}
}

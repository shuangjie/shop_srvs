package initialize

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"srvs/order_srv/global"
)

// InitRedisCli : 初始化 redis
func InitRedisCli() {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		DB:   9, // db9 -> order_srv
	})

	global.RedisCli = rdb
}

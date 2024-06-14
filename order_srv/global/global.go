package global

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"srvs/order_srv/config"
	"srvs/order_srv/proto"
)

var (
	DB           *gorm.DB
	RedisCli     *redis.Client
	RS           *redsync.Redsync
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig

	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)

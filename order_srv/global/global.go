package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"srvs/order_srv/config"
	"srvs/order_srv/proto"
)

var (
	DB           *gorm.DB
	RS           *redsync.Redsync
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig

	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)

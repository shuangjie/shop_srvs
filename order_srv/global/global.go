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

//func init() {
//	dsn := "root:123456@tcp(127.0.0.1:3306)/shop_order_srv?charset=utf8mb4&parseTime=True&loc=Local"
//
//	newLogger := logger.New(
//		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
//		logger.Config{
//			SlowThreshold:             time.Second, // Slow SQL threshold
//			LogLevel:                  logger.Info, // Log level
//			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
//			ParameterizedQueries:      false,       // Don't include params in the SQL log
//			Colorful:                  true,        // Disable color
//		},
//	)
//	var err error
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//		Logger: newLogger,
//	})
//	if err != nil {
//		panic("failed to connect database")
//	}
//}

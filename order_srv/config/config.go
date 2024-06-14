package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	DbName   string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type SrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type ServerConfig struct {
	Name       string       `mapstructure:"name" json:"name"`
	Host       string       `mapstructure:"host" json:"host"`
	Tags       []string     `mapstructure:"tags" json:"tags"`
	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	RedisInfo  RedisConfig  `mapstructure:"redis" json:"redis"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`

	// 商品微服务
	GoodsSrvInfo SrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	// 库存微服务
	InventorySrvInfo SrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
}

type NacosConfig struct {
	Host        string `mapstructure:"host"`
	Port        uint64 `mapstructure:"port"`
	NamespaceId string `mapstructure:"namespace_id"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Group       string `mapstructure:"group"`
	DataId      string `mapstructure:"data_id"`
}

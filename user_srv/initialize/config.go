package initialize

import (
	"fmt"

	"github.com/spf13/viper"

	"srvs/user_srv/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	dev := GetEnvInfo("SHOP_ENV_DEV")
	configFilePrefix := "user_srv/config-"
	configFileName := fmt.Sprintf("%spro.yaml", configFilePrefix)
	if dev {
		configFileName = fmt.Sprintf("%sdev.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}

}

package utils

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"qbot/pkg/logger"
)

type globalConfig struct {
	Mysql struct {
		DBAddress      string `yaml:"dbAddress"`
		DBUserName     string `yaml:"dbUserName"`
		DBPassword     string `yaml:"dbPassword"`
		DBName         string `yaml:"dbName"`
		DBMaxOpenConns int    `yaml:"dbMaxOpenConns"`
		DBMaxIdleConns int    `yaml:"dbMaxIdleConns"`
		DBMaxLifeTime  int    `yaml:"dbMaxLifeTime"`
	}

	Server struct {
		Ip   string `yaml:"ip"`
		Port string `yaml:"port"`
	}

	ThirdParty struct {
		GaoDe struct {
			Key          string            `yaml:"key"`
			LocationCode map[string]string `yaml:"locationCode"`
		} `yaml:"gaoDe"`
	}

	QqBot struct {
		QqId            string            `yaml:"qqId"`
		WeatherLocation map[string]string `yaml:"weatherLocation"`
	}
}

var GlobalConf = new(globalConfig)

func init() {

	filePath := "config/config.yaml"
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("load config file failed,err:%s", err))
	}

	if err := viper.Unmarshal(GlobalConf); err != nil {
		panic(fmt.Errorf("unmarshal config file failed,err:%s", err))
	}

	viper.WatchConfig()

	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("config changed")
		if err := viper.Unmarshal(GlobalConf); err != nil {
			panic(fmt.Errorf("unmarshal config file failed,err:%s", err))
		}
	})
	logger.Log.Infof("配置文件初始化成功...")
	logger.Log.Infof("当前机器人QQ号：%s", GlobalConf.QqBot.QqId)
	var supportCity []string
	for k, _ := range GlobalConf.QqBot.WeatherLocation {
		supportCity = append(supportCity, k)
	}
	logger.Log.Infof("天气查询支持城市：%v", supportCity)

}

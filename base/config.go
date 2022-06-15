package base

import (
	"log"

	"github.com/spf13/viper"
)

var ConfigFile = "./configs/config.toml"

var Config = &Configuration{}

type Configuration struct {
	MainConfig
	MysqlConf
	KafkaConf
}

type MainConfig struct {
	Listen             string
	Log                string
	Level              string
	LogRotateMegaBytes int
	LogRetainDays      int
}

func initConfig() {
	viper.SetConfigFile(ConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("failed to read config file -> ", err.Error())
	}
	Config.MainConfig.Listen = viper.GetString("main.listen")
	Config.MainConfig.Log = viper.GetString("main.log")
	Config.MainConfig.Level = viper.GetString("main.level")

	Config.MysqlConf.Host = viper.GetString("mysql.host")
	Config.MysqlConf.Port = viper.GetInt("mysql.port")
	Config.MysqlConf.Database = viper.GetString("mysql.database")
	Config.MysqlConf.User = viper.GetString("mysql.user")
	Config.MysqlConf.Password = viper.GetString("mysql.password")

	Config.KafkaConf.Brokers = viper.GetStringSlice("kafka.brokers")
}

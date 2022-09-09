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
	NotificationConf
}

type MainConfig struct {
	Listen             string
	Log                string
	Level              string
	Domain             string
	Easy               string
	LogRotateMegaBytes int
	LogRetainDays      int
	StaticDir          string
}

func initConfig() {
	viper.SetConfigFile(ConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("failed to read config file -> ", err.Error())
	}
	Config.MainConfig.Listen = viper.GetString("main.listen")
	Config.MainConfig.Log = viper.GetString("main.log")
	Config.MainConfig.Level = viper.GetString("main.level")
	Config.MainConfig.Domain = viper.GetString("main.domain")
	Config.MainConfig.Easy = viper.GetString("main.easy")
	Config.MainConfig.LogRotateMegaBytes = viper.GetInt("main.log_rotate_mega_bytes")
	Config.MainConfig.LogRetainDays = viper.GetInt("log_retain_days")

	Config.MysqlConf.Host = viper.GetString("mysql.host")
	Config.MysqlConf.Port = viper.GetInt("mysql.port")
	Config.MysqlConf.Database = viper.GetString("mysql.database")
	Config.MysqlConf.User = viper.GetString("mysql.user")
	Config.MysqlConf.Password = viper.GetString("mysql.password")

	Config.KafkaConf.Brokers = viper.GetStringSlice("kafka.brokers")

	Config.NotificationConf.DefaultDingTalkAppKey = viper.GetString("notification.default_dingtalk_appkey")
	Config.NotificationConf.DefaultDingTalkAppSecret = viper.GetString("notification.default_dingtalk_appsecret")
	Config.NotificationConf.DefaultDingtalkChatID = viper.GetString("notification.default_dingtalk_chatid")
}

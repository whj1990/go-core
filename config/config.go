package config

import (
	"flag"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/spf13/viper"
)

var apolloClient agollo.Client
var configPath string

func Init() {
	path := flag.String("c", "conf", "config path, eg: -c conf")
	flag.Parse()
	configPath = *path

	viper.AddConfigPath(*path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	c := &config.AppConfig{
		AppID:          viper.GetString("appid"),
		Cluster:        viper.GetString("cluster"),
		IP:             viper.GetString("url"),
		NamespaceName:  viper.GetString("namespace"),
		IsBackupConfig: false,
		Secret:         viper.GetString("secret"),
	}
	//agollo.SetLogger(&log.DefaultLogger{})
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		panic(err)
	}
	apolloClient = client
}

func GetString(key string, defaultValue string) string {
	return apolloClient.GetStringValue(key, defaultValue)
}

func GetInt(key string, defaultValue int) int {
	return apolloClient.GetIntValue(key, defaultValue)
}

func GetBool(key string, defaultValue bool) bool {
	return apolloClient.GetBoolValue(key, defaultValue)
}

func GetConfigPath() string {
	return configPath
}

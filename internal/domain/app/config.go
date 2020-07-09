package app

import (
	"github.com/spf13/viper"
)

type DBType string

const (
	DbTypePostgres DBType = "postgres"
)

type PostgresMaster struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Db       string `mapstructure:"db"`
}

type PostgresSlave struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type config struct {
	DBType   DBType `mapstructure:"db_type"`
	Postgres struct {
		Master PostgresMaster  `mapstructure:"master"`
		Slaves []PostgresSlave `mapstructure:"slaves"`
	} `mapstructure:"postgres"`
}

type appConf struct {
	App config `mapstructure:"app"`
}

func readConfig() (conf config, err error) {
	v := viper.New()

	// because there is no file extension in a stream of bytes,
	// supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	v.SetConfigType("yaml")
	v.SetConfigFile("config.yaml")

	err = v.ReadInConfig()
	if err != nil {
		return
	}

	var c = &appConf{}
	err = v.Unmarshal(&c)
	if err != nil {
		return
	}

	conf = c.App
	return
}

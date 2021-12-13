package config

import (
	"log"

	"gopkg.in/ini.v1"
)

type ServerConfig struct {
	Port string
}

type DbConfig struct {
	DbUser string
	DbName string
	DbPass string
	DbAddr string
}

type ESConfig struct {
	ServerConfig ServerConfig
	DbConfig     DbConfig
}

func LoadConfig() *ESConfig {

	var config ESConfig
	cfg, err := ini.Load("config.ini")

	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	config.ServerConfig.Port = cfg.Section("server").Key("port").String()

	//load database info
	config.DbConfig.DbName = cfg.Section("database").Key("db_name").String()
	config.DbConfig.DbPass = cfg.Section("database").Key("db_password").String()
	config.DbConfig.DbUser = cfg.Section("database").Key("db_username").String()
	config.DbConfig.DbAddr = cfg.Section("database").Key("db_address").String()

	return &config
}

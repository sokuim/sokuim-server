package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	LoggerConfig     LoggerConfig
	MqProducerConfig MqProducerConfig
}

type LoggerConfig struct {
	Path         string
	Folder       string
	Stdout       bool
	RotateSize   uint
	RotateLimit  uint
	RotateExpire uint
}

type MqProducerConfig struct {
	Brokers []string
}

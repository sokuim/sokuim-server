package svc

import (
	"sokuim/sokuim-server/app/gateway/internal/config"
	"sokuim/sokuim-server/pkg/logger"
	"sokuim/sokuim-server/pkg/mq"
)

type ServiceContext struct {
	Config     config.Config
	Logger     *logger.Logger
	MqProducer *mq.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		Logger:     newLogger(c.LoggerConfig),
		MqProducer: newMqProducer(c.MqProducerConfig.Brokers, c.LoggerConfig),
	}
}

func newMqProducer(brokers []string, c config.LoggerConfig) *mq.Producer {
	p, err := mq.NewProducer(brokers, covertLogConf(c))
	if err != nil {
		panic(err)
	}
	return p
}

func newLogger(c config.LoggerConfig) *logger.Logger {
	return logger.NewLogger(covertLogConf(c))
}

func covertLogConf(c config.LoggerConfig) logger.Config {
	return logger.Config{
		Path:         c.Path,
		Folder:       c.Folder,
		Stdout:       c.Stdout,
		RotateSize:   c.RotateSize,
		RotateLimit:  c.RotateLimit,
		RotateExpire: c.RotateExpire,
	}
}

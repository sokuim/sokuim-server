package svc

import (
	"sokuim/sokuim-server/app/gateway/internal/config"
	"sokuim/sokuim-server/pkg/logger"
)

type ServiceContext struct {
	Config config.Config
	Logger *logger.Logger
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Logger: NewLogger(c.LoggerConfig),
	}
}

func NewLogger(c config.LoggerConfig) *logger.Logger {
	logConf := logger.Config{
		Path:         c.Path,
		Folder:       c.Folder,
		Stdout:       c.Stdout,
		RotateSize:   c.RotateSize,
		RotateLimit:  c.RotateLimit,
		RotateExpire: c.RotateExpire,
	}
	return logger.NewLogger(logConf)
}

package svc

import (
	"sokuim/sokuim-server/app/comet/internal/config"
	"sokuim/sokuim-server/app/core/client/socket"
	"sokuim/sokuim-server/pkg/logger"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	SocketRPC socket.Socket
	Log       *logger.Logger
}

func NewServiceContext(c config.Config) *ServiceContext {
	coreClient := zrpc.MustNewClient(c.CoreRpcConfig)
	return &ServiceContext{
		Config:    c,
		SocketRPC: socket.NewSocket(coreClient),
		Log:       newLogger(c),
	}
}

func newLogger(conf config.Config) *logger.Logger {
	logConf := logger.Config{
		Path:         conf.LoggerConfig.Path,
		Folder:       conf.LoggerConfig.Folder,
		Stdout:       conf.LoggerConfig.Stdout,
		RotateSize:   conf.LoggerConfig.RotateSize,
		RotateLimit:  conf.LoggerConfig.RotateLimit,
		RotateExpire: conf.LoggerConfig.RotateExpire,
	}
	log := logger.NewLogger(logConf)
	return log
}

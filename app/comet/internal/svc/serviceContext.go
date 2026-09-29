package svc

import (
	"sokuim/sokuim-server/app/comet/internal/config"
	"sokuim/sokuim-server/app/core/client/socket"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	SocketRPC socket.Socket
}

func NewServiceContext(c config.Config) *ServiceContext {
	coreClient := zrpc.MustNewClient(c.CoreRpcConfig)
	return &ServiceContext{
		Config:    c,
		SocketRPC: socket.NewSocket(coreClient),
	}
}

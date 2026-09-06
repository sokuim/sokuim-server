package main

import (
	"flag"
	"fmt"

	"sokuim/sokuim-server/app/core/internal/config"
	socketServer "sokuim/sokuim-server/app/core/internal/server/socket"
	"sokuim/sokuim-server/app/core/internal/svc"
	"sokuim/sokuim-server/app/core/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/core.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterSocketServer(grpcServer, socketServer.NewSocketServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting core rpc server at %s...\n", c.ListenOn)
	s.Start()
}

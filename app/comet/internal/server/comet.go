// Package server
// File comet.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-03 00:29:19
// Modified 2026-09-03 00:29:19

package server

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"sokuim/sokuim-server/app/comet/internal/config"
	cometServer "sokuim/sokuim-server/app/comet/internal/server/comet"
	"sokuim/sokuim-server/app/comet/internal/svc"
	"sokuim/sokuim-server/app/comet/pb"
)

type Comet struct {
	conf      config.Config
	svcCtx    *svc.ServiceContext
	rpcServer *zrpc.RpcServer
}

func NewComet(conf config.Config) *Comet {
	svcCtx := svc.NewServiceContext(conf)
	return &Comet{
		conf:   conf,
		svcCtx: svcCtx,
	}
}

func (c *Comet) Start() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	cs := NewServer(c.conf, c.svcCtx.SocketRPC)
	if err := c.initTCP(cs); err != nil {
		panic(err)
	}
	if err := c.initWebSocket(cs); err != nil {
		panic(err)
	}
	s := zrpc.MustNewServer(c.conf.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterCometServer(grpcServer, cometServer.NewCometServer(c.svcCtx))
		if c.conf.Mode == service.DevMode || c.conf.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	c.rpcServer = s
	fmt.Printf("Starting comet.rpc server at %s...\n", c.conf.ListenOn)
	s.Start()
}

func (c *Comet) initTCP(server *Server) error {
	return InitTCP(server, c.conf.TcpConfig.Bind, runtime.NumCPU())
}

func (c *Comet) initWebSocket(server *Server) error {
	return InitWebsocket(server, c.conf.WebsocketConfig.Bind, runtime.NumCPU())
}

func (c *Comet) Stop() {
	fmt.Println("Comet stopping...")
	if c.rpcServer == nil {
		return
	}
	c.rpcServer.Stop()
}

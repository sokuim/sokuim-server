// Package server
// File server_tcp.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 22:29:29
// Modified 2026-09-29 22:29:29

package server

import (
	"context"
	"net"
	"sokuim/sokuim-server/pkg/logger"
)

func InitWebsocket(s *Server, addrs []string, accept int) (err error) {
	var (
		bind     string
		listener *net.TCPListener
		addr     *net.TCPAddr
	)

	log := logger.NewLogger(logger.Config{
		Path:         s.conf.LoggerConfig.Path,
		Folder:       "websocket",
		Stdout:       s.conf.LoggerConfig.Stdout,
		RotateSize:   s.conf.LoggerConfig.RotateSize,
		RotateLimit:  s.conf.LoggerConfig.RotateLimit,
		RotateExpire: s.conf.LoggerConfig.RotateExpire,
	})

	ctx := context.Background()

	for _, bind = range addrs {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			log.Errorf(ctx, "net.ResolveTCPAddr(tcp, %s) error(%v)", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			log.Errorf(ctx, "net.ListenTCP(tcp, %s) error(%v)", bind, err)
			return
		}
		log.Infof(ctx, "start ws listen: %s", bind)
	}
	log.Info(ctx, listener)
	return
}

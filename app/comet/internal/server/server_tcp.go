// Package server
// File server_tcp.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 22:38:59
// Modified 2026-09-29 22:38:59

package server

import (
	"context"
	"net"
	"sokuim/sokuim-server/pkg/logger"
)

func InitTCP(s *Server, addrs []string, accept int) (err error) {
	var (
		bind     string
		listener *net.TCPListener
		addr     *net.TCPAddr
	)

	log := logger.NewLogger(logger.Config{
		Path:         s.conf.LoggerConfig.Path,
		Folder:       "tcp",
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
		for i := 0; i < accept; i++ {
			go acceptTCP(s, listener, log)
		}
	}
	log.Info(ctx, listener)
	return
}

func acceptTCP(s *Server, lis *net.TCPListener, log *logger.Logger) {
	var (
		conn *net.TCPConn
		err  error
		r    int
	)
	for {
		if conn, err = lis.AcceptTCP(); err != nil {
			return
		}
		if err = conn.SetKeepAlive(s.conf.TcpConfig.Keepalive); err != nil {
			return
		}
		if err = conn.SetReadBuffer(s.conf.TcpConfig.RcvBuf); err != nil {
			return
		}
		if err = conn.SetWriteBuffer(s.conf.TcpConfig.SndBuf); err != nil {
			return
		}
		go serveTCP(s, conn, r, log)
		if r++; r == maxInt {
			r = 0
		}
	}
}

func serveTCP(s *Server, conn *net.TCPConn, r int, log *logger.Logger) {
	log.Debug(context.Background(), "serveTCP", conn.RemoteAddr())
}

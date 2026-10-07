// Package config
// File config.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-03 00:24:50
// Modified 2026-09-03 00:24:50

package config

import (
	xtime "time"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	CoreRpcConfig   zrpc.RpcClientConf
	LoggerConfig    LoggerConfig
	WebsocketConfig WebsocketConfig
	TcpConfig       TcpConfig
	BucketConfig    BucketConfig
	ProtocolConfig  ProtocolConfig
}

type LoggerConfig struct {
	Path         string
	Folder       string
	Stdout       bool
	RotateSize   uint
	RotateLimit  uint
	RotateExpire uint
}

type TcpConfig struct {
	Bind         []string
	SndBuf       int
	RcvBuf       int
	Keepalive    bool
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
}

type WebsocketConfig struct {
	Bind    []string
	TlsOpen bool
	TlsBind []string
}

type BucketConfig struct {
	Size          int
	Channel       int
	Room          int
	RoutineAmount uint64
	RoutineSize   int
}

type ProtocolConfig struct {
	Timer            int
	TimerSize        int
	SrvProto         int
	CliProto         int
	HandshakeTimeout xtime.Duration
}

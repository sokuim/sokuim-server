// Package server
// File server.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-28 23:33:52
// Modified 2026-09-28 23:33:52

package server

import (
	"sokuim/sokuim-server/app/comet/internal/config"
	"time"
)

const (
	minServerHeartbeat = time.Minute * 10
	maxServerHeartbeat = time.Minute * 30
	maxInt             = 1<<31 - 1
)

type Server struct {
	conf config.Config
}

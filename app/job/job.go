// Package main
// File job.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-08-24 23:41:19
// Modified 2026-08-24 23:41:19

package main

import (
	"flag"
	"sokuim/sokuim-server/app/job/internal/config"
	"sokuim/sokuim-server/app/job/internal/server"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/job.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	group := service.NewServiceGroup()
	defer group.Stop()

	port := c.Port
	health := server.NewHealth(port)
	group.Add(health)

	group.Start()
}

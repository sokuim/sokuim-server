// Package server
// File comet.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-03 00:29:19
// Modified 2026-09-03 00:29:19

package server

import "log"

type Comet struct{}

func NewComet() *Comet {
	return &Comet{}
}

func (c *Comet) Start() {
	log.Printf("Comet starting...")
}

func (c *Comet) Stop() {
	log.Printf("Comet stopping...")
}

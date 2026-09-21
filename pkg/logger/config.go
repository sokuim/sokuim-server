// Package logger
// File config.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-20 23:38:52
// Modified 2026-09-20 23:38:52

package logger

type Config struct {
	Path         string
	Folder       string
	Stdout       bool
	RotateSize   uint // split size(MB)
	RotateLimit  uint // old log saved count
	RotateExpire uint // auto delete time(day)
}

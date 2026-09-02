// Package params
// File params.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-03 00:16:49
// Modified 2026-09-03 00:16:49

package params

type HealthData struct {
	Status string `json:"status"`
}

type HealthResp struct {
	Status  int        `json:"status"`
	Success bool       `json:"success"`
	Msg     string     `json:"msg"`
	Data    HealthData `json:"data"`
}

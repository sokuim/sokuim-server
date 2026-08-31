// Package result
// File response_bean.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-01 01:27:00
// Modified 2026-09-01 01:27:00

package result

import "strings"

type ResponseBean struct {
	Status  uint32      `json:"status"`
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func Success(msg string, data interface{}) *ResponseBean {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		msg = "success"
	}
	return &ResponseBean{
		Status:  uint32(0),
		Success: true,
		Msg:     msg,
		Data:    data,
	}
}

func Error(code uint32, msg string) *ResponseBean {
	return &ResponseBean{
		Status:  code,
		Success: false,
		Msg:     msg,
		Data:    nil,
	}
}

// Package xerr
// File msg.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-01 01:30:19
// Modified 2026-09-01 01:30:19

package xerr

var message map[int32]string

func MapErrMsg(errCode int32) string {
	if msg, ok := message[errCode]; ok {
		return msg
	}
	return "Server Error"
}

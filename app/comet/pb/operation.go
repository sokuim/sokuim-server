// Package pb
// File operation.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 00:24:05
// Modified 2026-09-29 00:24:05

package pb

const (
	OpHandshake      = int32(0)
	OpHandshakeResp  = int32(1)
	OpHeartbeat      = int32(2)
	OpHeartbeatResp  = int32(3)
	OpSendMsg        = int32(4)
	OpSendMsgResp    = int32(5)
	OpDisconnectResp = int32(6)
	OpAuth           = int32(7)
	OpAuthResp       = int32(8)
	OpRaw            = int32(9)
	OpProtoReady     = int32(10)
	OpProtoFinish    = int32(11)
	OpChangeRoom     = int32(12)
	OpChangeRoomResp = int32(13)
	OpSub            = int32(14)
	OpSubResp        = int32(15)
	OpUnsub          = int32(16)
	OpUnsubResp      = int32(17)
)

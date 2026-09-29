// Package server
// File errors.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 00:01:51
// Modified 2026-09-29 00:01:51

package server

import "sokuim/sokuim-server/pkg/xerr"

var (
	ErrHandshake            = xerr.NewErrCodeMsg(xerr.ServerErr, "handshake failed")
	ErrOperation            = xerr.NewErrCodeMsg(xerr.ServerErr, "request operation not valid")
	ErrRingEmpty            = xerr.NewErrCodeMsg(xerr.ServerErr, "ring buffer empty")
	ErrRingFull             = xerr.NewErrCodeMsg(xerr.ServerErr, "ring buffer full")
	ErrTimerFull            = xerr.NewErrCodeMsg(xerr.ServerErr, "timer full")
	ErrTimerEmpty           = xerr.NewErrCodeMsg(xerr.ServerErr, "timer empty")
	ErrTimerNoItem          = xerr.NewErrCodeMsg(xerr.ServerErr, "timer item not exist")
	ErrPushMsgArg           = xerr.NewErrCodeMsg(xerr.ServerErr, "rpc push msg arg error")
	ErrPushMsgsArg          = xerr.NewErrCodeMsg(xerr.ServerErr, "rpc push msgs arg error")
	ErrMPushMsgArg          = xerr.NewErrCodeMsg(xerr.ServerErr, "rpc multi push msg arg error")
	ErrMPushMsgsArg         = xerr.NewErrCodeMsg(xerr.ServerErr, "rpc multi push msgs arg error")
	ErrSignalFullMsgDropped = xerr.NewErrCodeMsg(xerr.ServerErr, "signal channel full, msg dropped")
	ErrBroadCastArg         = xerr.NewErrCodeMsg(xerr.ServerErr, "rpc broadcast arg error")
	ErrBroadCastRoomArg     = xerr.NewErrCodeMsg(xerr.ServerErr, "rpc broadcast  room arg error")
	ErrRoomDropped          = xerr.NewErrCodeMsg(xerr.ServerErr, "room dropped")
	ErrLogic                = xerr.NewErrCodeMsg(xerr.ServerErr, "logic rpc is not available")
	ErrRingFullMsgDropped   = xerr.NewErrCodeMsg(xerr.ServerErr, "signal channel full, msg dropped")
)

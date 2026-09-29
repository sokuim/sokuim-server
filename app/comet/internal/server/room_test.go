// Package server
// File room_test.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 10:39:09
// Modified 2026-09-29 10:39:09

package server

import (
	"sokuim/sokuim-server/app/comet/pb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomPutDel(t *testing.T) {
	r := NewRoom("room1")
	ch := NewChannel(1, 1)

	assert.Nil(t, r.Put(ch))
	assert.Equal(t, int32(1), r.Online)

	assert.True(t, r.Del(ch))
	assert.Equal(t, int32(0), r.Online)
}

func TestRoomPush(t *testing.T) {
	r := NewRoom("room1")
	ch1 := NewChannel(1, 1)
	ch2 := NewChannel(1, 1)
	assert.Nil(t, r.Put(ch1))
	assert.Nil(t, r.Put(ch2))

	r.Push(&pb.CometMsgProto{Op: 7})

	// 每个 channel 的 signal 都能收到该消息
	assert.Equal(t, int32(7), ch1.Ready().Op)
	assert.Equal(t, int32(7), ch2.Ready().Op)
}

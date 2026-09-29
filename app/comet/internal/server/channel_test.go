// Package server
// File channel_test.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 10:14:53
// Modified 2026-09-29 10:14:53

package server

import (
	"sokuim/sokuim-server/app/comet/pb"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestChannelPushDrop
func TestChannelPushDrop(t *testing.T) {
	ch := NewChannel(1, 1)

	// success. not dropped
	err := ch.Push(&pb.CometMsgProto{Op: 1})
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), ch.DropCount())

	// dropped when cache full
	err = ch.Push(&pb.CometMsgProto{Op: 2})
	assert.Equal(t, ErrSignalFullMsgDropped, err)
	assert.Equal(t, uint64(1), ch.DropCount())
}

// TestChannelWatch
func TestChannelWatch(t *testing.T) {
	ch := NewChannel(1, 1)
	ch.Watch(1, 2, 3)

	assert.True(t, ch.NeedPush(1))
	assert.True(t, ch.NeedPush(3))
	assert.False(t, ch.NeedPush(9))

	ch.UnWatch(1)
	assert.False(t, ch.NeedPush(1))
	assert.True(t, ch.NeedPush(2))
}

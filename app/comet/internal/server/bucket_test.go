// Package server
// File bucket_test.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-10-07 22:08:10
// Modified 2026-10-07 22:08:10

package server

import (
	"sokuim/sokuim-server/app/comet/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBucketPutDel 验证 bucket 对 channel 的存取、计数与删除。
func TestBucketPutDel(t *testing.T) {
	b := NewBucket(config.BucketConfig{Channel: 10, Room: 10, RoutineAmount: 1, RoutineSize: 10})

	ch := NewChannel(1, 1)
	ch.Key = "key1"
	ch.IP = "127.0.0.1"

	assert.Nil(t, b.Put("", ch))
	assert.Equal(t, 1, b.ChannelCount())
	assert.Equal(t, ch, b.Channel("key1"))

	b.Del(ch)
	assert.Equal(t, 0, b.ChannelCount())
	assert.Nil(t, b.Channel("key1"))
}

// TestBucketPutRoom 验证 bucket 将 channel 放入指定房间。
func TestBucketPutRoom(t *testing.T) {
	b := NewBucket(config.BucketConfig{Channel: 10, Room: 10, RoutineAmount: 1, RoutineSize: 10})

	ch := NewChannel(1, 1)
	ch.Key = "key1"
	ch.IP = "127.0.0.1"

	assert.Nil(t, b.Put("room1", ch))
	assert.NotNil(t, ch.Room)
	assert.Equal(t, "room1", ch.Room.ID)
	assert.Equal(t, 1, b.RoomCount())
}

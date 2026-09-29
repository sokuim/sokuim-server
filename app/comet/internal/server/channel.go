// Package server
// File channel.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-28 23:50:50
// Modified 2026-09-28 23:50:50

package server

import (
	"sokuim/sokuim-server/app/comet/pb"
	"sokuim/sokuim-server/pkg/bufio"
	"sync"
	"sync/atomic"
)

// Channel used by message pusher send msg to write goroutine.
// Expression a user.
type Channel struct {
	signal    chan *pb.CometMsgProto
	watchOps  map[int32]struct{}
	lock      sync.RWMutex
	Room      *Room
	CliProto  Ring
	Writer    bufio.Writer
	Reader    bufio.Reader
	Next      *Channel
	Prev      *Channel
	Mid       string
	Key       string
	IP        string
	CType     string
	dropCount uint64 // dropped msg count with channel full.
}

// NewChannel new a channel
func NewChannel(cli, svr int) *Channel {
	c := new(Channel)
	c.CliProto.Init(cli)
	c.signal = make(chan *pb.CometMsgProto, svr)
	c.watchOps = make(map[int32]struct{})
	return c
}

// Watch Set this channel subscribe op collection
func (c *Channel) Watch(accepts ...int32) {
	c.lock.Lock()
	for _, op := range accepts {
		c.watchOps[op] = struct{}{}
	}
	c.lock.Unlock()
}

func (c *Channel) UnWatch(accepts ...int32) {
	c.lock.Lock()
	for _, op := range accepts {
		delete(c.watchOps, op)
	}
	c.lock.Unlock()
}

// NeedPush verify if in watch
func (c *Channel) NeedPush(op int32) bool {
	c.lock.RLock()
	if _, ok := c.watchOps[op]; ok {
		c.lock.RUnlock()
		return true
	}
	c.lock.RUnlock()
	return false
}

func (c *Channel) Push(p *pb.CometMsgProto) (err error) {
	select {
	case c.signal <- p:
	default:
		atomic.AddUint64(&c.dropCount, 1)
		err = ErrRingFullMsgDropped
	}
	return
}

// DropCount dropped msg count
func (c *Channel) DropCount() uint64 {
	return atomic.LoadUint64(&c.dropCount)
}

func (c *Channel) Ready() *pb.CometMsgProto {
	return <-c.signal
}

func (c *Channel) Signal() {
	c.signal <- pb.ProtoReady
}

func (c *Channel) Close() {
	c.signal <- pb.ProtoFinish
}

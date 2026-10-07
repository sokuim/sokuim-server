// Package server
// File round.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-10-07 22:11:58
// Modified 2026-10-07 22:11:58

package server

import (
	"sokuim/sokuim-server/app/comet/internal/config"
	"sokuim/sokuim-server/pkg/bytes"
	"sokuim/sokuim-server/pkg/time"
)

// RoundOptions round options.
type RoundOptions struct {
	Timer        int
	TimerSize    int
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
}

// Round used for connection round-robin get a reader/writer/timer for split big lock.
type Round struct {
	readers []bytes.Pool
	writers []bytes.Pool
	timers  []time.Timer
	options RoundOptions
}

// NewRound new a round struct.
func NewRound(c config.Config) (r *Round) {
	var i int
	r = &Round{
		options: RoundOptions{
			Reader:       c.TcpConfig.Reader,
			ReadBuf:      c.TcpConfig.ReadBuf,
			ReadBufSize:  c.TcpConfig.ReadBufSize,
			Writer:       c.TcpConfig.Writer,
			WriteBuf:     c.TcpConfig.WriteBuf,
			WriteBufSize: c.TcpConfig.WriteBufSize,
			Timer:        c.ProtocolConfig.Timer,
			TimerSize:    c.ProtocolConfig.TimerSize,
		}}
	// reader
	r.readers = make([]bytes.Pool, r.options.Reader)
	for i = 0; i < r.options.Reader; i++ {
		r.readers[i].Init(r.options.ReadBuf, r.options.ReadBufSize)
	}
	// writer
	r.writers = make([]bytes.Pool, r.options.Writer)
	for i = 0; i < r.options.Writer; i++ {
		r.writers[i].Init(r.options.WriteBuf, r.options.WriteBufSize)
	}
	// timer
	r.timers = make([]time.Timer, r.options.Timer)
	for i = 0; i < r.options.Timer; i++ {
		r.timers[i].Init(r.options.TimerSize)
	}
	return
}

// Timer get a timer.
func (r *Round) Timer(rn int) *time.Timer {
	return &(r.timers[rn%r.options.Timer])
}

// Reader get a reader memory buffer.
func (r *Round) Reader(rn int) *bytes.Pool {
	return &(r.readers[rn%r.options.Reader])
}

// Writer get a writer memory buffer pool.
func (r *Round) Writer(rn int) *bytes.Pool {
	return &(r.writers[rn%r.options.Writer])
}

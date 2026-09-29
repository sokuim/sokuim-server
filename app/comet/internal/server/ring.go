// Package server
// File ring.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-28 23:57:54
// Modified 2026-09-28 23:57:54

package server

import "sokuim/sokuim-server/app/comet/pb"

// Ring cache channel's data and push to client.
type Ring struct {
	rp   uint64
	num  uint64
	mask uint64
	wp   uint64
	data []pb.CometMsgProto
}

func NewRing(num int) *Ring {
	r := new(Ring)
	r.init(uint64(num))
	return r
}

func (r *Ring) Init(num int) {
	r.init(uint64(num))
}

func (r *Ring) init(num uint64) {
	// 2^N
	if num&(num-1) != 0 {
		for num&(num-1) != 0 {
			num &= num - 1
		}
		num <<= 1
	}
	r.data = make([]pb.CometMsgProto, num)
	r.num = num
	r.mask = r.num - 1
}

func (r *Ring) Get() (proto *pb.CometMsgProto, err error) {
	if r.rp == r.wp {
		return nil, ErrRingEmpty
	}
	proto = &r.data[r.rp&r.mask]
	return
}

func (r *Ring) GetAdv() {
	r.rp++
}

func (r *Ring) Set() (proto *pb.CometMsgProto, err error) {
	if r.wp-r.rp >= r.num {
		return nil, ErrRingFull
	}
	proto = &r.data[r.wp&r.mask]
	return
}

func (r *Ring) SetAdv() {
	r.wp++
}

func (r *Ring) Reset() {
	r.rp = 0
	r.wp = 0
}

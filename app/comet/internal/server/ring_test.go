// Package server
// File ring_test.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 00:04:12
// Modified 2026-09-29 00:04:12

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRing validate ring buffer write&read
func TestRing(t *testing.T) {
	r := NewRing(4)

	_, err := r.Get()
	assert.Equal(t, ErrRingEmpty, err)

	for i := 0; i < 4; i++ {
		p, e := r.Set()
		assert.Nil(t, e)
		p.Op = int32(i)
		r.SetAdv()
	}

	_, err = r.Set()
	assert.Equal(t, ErrRingFull, err)

	for i := 0; i < 4; i++ {
		p, e := r.Get()
		assert.Nil(t, e)
		assert.Equal(t, int32(i), p.Op)
		r.GetAdv()
	}

	r.Reset()
	_, err = r.Get()
	assert.Equal(t, ErrRingEmpty, err)
}

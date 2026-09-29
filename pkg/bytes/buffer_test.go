// Package bytes
// File buffer_test.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 00:25:30
// Modified 2026-09-29 00:25:30

package bytes

import "testing"

func TestBuffer(t *testing.T) {
	p := NewPool(2, 10)
	b := p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}
	b = p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}
	b = p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}
}

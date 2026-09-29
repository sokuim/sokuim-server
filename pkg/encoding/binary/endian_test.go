// Package binary
// File endian_test.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 00:30:28
// Modified 2026-09-29 00:30:28

package binary

import "testing"

func TestInt8(t *testing.T) {
	b := make([]byte, 1)
	BigEndian.PutInt8(b, 100)
	i := BigEndian.Int8(b)
	if i != 100 {
		t.FailNow()
	}
}

func TestInt16(t *testing.T) {
	b := make([]byte, 2)
	BigEndian.PutInt16(b, 100)
	i := BigEndian.Int16(b)
	if i != 100 {
		t.FailNow()
	}
}

func TestInt32(t *testing.T) {
	b := make([]byte, 4)
	BigEndian.PutInt32(b, 100)
	i := BigEndian.Int32(b)
	if i != 100 {
		t.FailNow()
	}
}

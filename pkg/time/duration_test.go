// Package time
// File duration_test.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-28 23:38:44
// Modified 2026-09-28 23:38:44

package time

import (
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	var (
		input  = []byte("10s")
		output = time.Second * 10
		d      Duration
	)
	if err := d.UnmarshalText(input); err != nil {
		t.FailNow()
	}
	if int64(output) != int64(d) {
		t.FailNow()
	}
}

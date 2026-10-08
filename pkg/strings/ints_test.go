// Package strings
// File ints_test.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-10-08 22:39:10
// Modified 2026-10-08 22:39:10

package strings

import (
	"reflect"
	"testing"
)

func TestInt32(t *testing.T) {
	i := []int32{1, 2, 3}
	s := JoinInt32s(i, ",")
	ii, _ := SplitInt32s(s, ",")
	if !reflect.DeepEqual(i, ii) {
		t.FailNow()
	}
}

func TestInt64(t *testing.T) {
	i := []int64{1, 2, 3}
	s := JoinInt64s(i, ",")
	ii, _ := SplitInt64s(s, ",")
	if !reflect.DeepEqual(i, ii) {
		t.FailNow()
	}
}

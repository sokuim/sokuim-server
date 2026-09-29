// Package time
// File timer_test.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-28 23:40:22
// Modified 2026-09-28 23:40:22

package time

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := NewTimer(100)
	tds := make([]*TimerData, 100)
	for i := 0; i < 100; i++ {
		tds[i] = timer.Add(time.Duration(i)*time.Second+5*time.Minute, nil)
	}
	printTimer(timer)
	for i := 0; i < 100; i++ {
		fmt.Printf("td: %s, %s, %d\n", tds[i].Key, tds[i].ExpireString(), tds[i].index)
		timer.Del(tds[i])
	}
	printTimer(timer)
	for i := 0; i < 100; i++ {
		tds[i] = timer.Add(time.Duration(i)*time.Second+5*time.Minute, nil)
	}
	printTimer(timer)
	for i := 0; i < 100; i++ {
		timer.Del(tds[i])
	}
	printTimer(timer)
	timer.Add(time.Second, nil)
	time.Sleep(time.Second * 2)
	if len(timer.timers) != 0 {
		t.FailNow()
	}
}

func printTimer(timer *Timer) {
	fmt.Printf("----------timers: %d ----------\n", len(timer.timers))
	for i := 0; i < len(timer.timers); i++ {
		fmt.Printf("timer: %s, %s, index: %d\n", timer.timers[i].Key, timer.timers[i].ExpireString(), timer.timers[i].index)
	}
	fmt.Printf("--------------------\n")
}

// Package main
// File core.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-08-17 00:41:00
// Modified 2026-08-17 00:41:00

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "sokuim.core: "+time.Now().String())
	})
	log.Println("comet start at :6071")
	err := http.ListenAndServe(":6071", nil)
	if err != nil {
		panic(err)
	}
}

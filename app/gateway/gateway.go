// Package main
// File gateway.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-08-17 00:42:59
// Modified 2026-08-17 00:42:59

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "sokuim.gateway: "+time.Now().String())
	})
	log.Println("comet start at :6072")
	err := http.ListenAndServe(":6072", nil)
	if err != nil {
		panic(err)
	}
}

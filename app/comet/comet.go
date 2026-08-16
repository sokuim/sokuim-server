// Package main
// File comet.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-08-17 00:31:03
// Modified 2026-08-17 00:31:03

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "sokuim.comet: "+time.Now().String())
	})
	log.Println("comet start at :6070")
	err := http.ListenAndServe(":6070", nil)
	if err != nil {
		panic(err)
	}
}

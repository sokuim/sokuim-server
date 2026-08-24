// Package main
// File job.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-08-24 23:41:19
// Modified 2026-08-24 23:41:19

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Job: "+time.Now().String())
	})
	log.Println("Listening on :7073")
	err := http.ListenAndServe(":7073", nil)
	if err != nil {
		panic(err)
	}
}

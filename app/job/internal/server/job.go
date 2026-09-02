// Package server
// File job.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-02 21:58:21
// Modified 2026-09-02 21:58:21

package server

import "log"

type Job struct {
}

func NewJob() *Job {
	return &Job{}
}

func (j *Job) Start() {
	log.Println("start job")
}

func (j *Job) Stop() {
	log.Println("stop job")
}

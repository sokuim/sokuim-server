// Package server
// File health.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-02 10:30:07
// Modified 2026-09-02 10:30:07

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type healthResp struct {
	Status  int        `json:"status"`
	Success bool       `json:"success"`
	Msg     string     `json:"msg"`
	Data    healthData `json:"data"`
}

type healthData struct {
	Status string `json:"status"`
}

type Health struct {
	port   int
	server *http.Server
	once   sync.Once
}

func NewHealth(port int) *Health {
	return &Health{
		port: port,
	}
}

func (h *Health) Start() {
	hr := healthResp{
		Status:  0,
		Success: true,
		Msg:     "success",
		Data: healthData{
			Status: "ok",
		},
	}
	hrb, _ := json.Marshal(hr)

	mux := http.NewServeMux()
	mux.HandleFunc("/common/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(hrb); err != nil {
			log.Printf("health response write error: %v", err)
		}
	})
	h.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", h.port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	log.Printf("job health check on port %d", h.port)

	h.once.Do(func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Health server failed: %v", err)
		}
	})
}

func (h *Health) Stop() {
	if h.server != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Printf("Health check server shutting down on port %d", h.port)
	if err := h.server.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown health server: %v", err)
	}
}

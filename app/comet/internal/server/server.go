// Package server
// File server.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-28 23:33:52
// Modified 2026-09-28 23:33:52

package server

import (
	"context"
	"math/rand"
	"os"
	"sokuim/sokuim-server/app/comet/internal/config"
	"sokuim/sokuim-server/app/core/client/socket"
	"sokuim/sokuim-server/pkg/logger"
	"time"

	"github.com/cespare/xxhash/v2"
)

const (
	minServerHeartbeat = time.Minute * 10
	maxServerHeartbeat = time.Minute * 30
	maxInt             = 1<<31 - 1
)

type Server struct {
	conf        config.Config
	socketRPC   socket.Socket
	buckets     []*Bucket
	bucketIndex uint32
	serverID    string
	round       *Round
	done        chan struct{} // 用于通知后台 goroutine（如 onlineproc）在进程退出时优雅停止。
	log         *logger.Logger
}

func NewServer(conf config.Config, socketRPC socket.Socket, log *logger.Logger) *Server {
	s := &Server{
		conf:      conf,
		round:     NewRound(conf),
		socketRPC: socketRPC,
		serverID:  getHostName(),
		done:      make(chan struct{}),
		log:       log,
	}
	// init bucket
	s.buckets = make([]*Bucket, conf.BucketConfig.Size)
	s.bucketIndex = uint32(conf.BucketConfig.Size)
	for i := 0; i < conf.BucketConfig.Size; i++ {
		s.buckets[i] = NewBucket(conf.BucketConfig)
	}
	go s.onlineproc()
	return s
}

func (s *Server) ServerID() string {
	return s.serverID
}

func (s *Server) Buckets() []*Bucket {
	return s.buckets
}

func (s *Server) Bucket(subKey string) *Bucket {
	idx := uint32(xxhash.Sum64String(subKey)) % s.bucketIndex
	return s.buckets[idx]
}

func (s *Server) RandServerHearbeat() time.Duration {
	return (minServerHeartbeat + time.Duration(rand.Int63n(int64(maxServerHeartbeat-minServerHeartbeat))))
}

func (s *Server) Close() (err error) {
	close(s.done)
	for _, b := range s.buckets {
		b.Close()
	}
	return
}

func (s *Server) onlineproc() {
	for {
		var (
			allRoomsCount map[string]int32
			err           error
		)
		roomCount := make(map[string]int32)
		for _, bucket := range s.buckets {
			for roomID, count := range bucket.RoomsCount() {
				roomCount[roomID] += count
			}
		}
		if allRoomsCount, err = s.RenewOnline(context.Background(), s.serverID, roomCount); err != nil {
			select {
			case <-s.done:
				return
			case <-time.After(time.Second):
			}
			continue
		}
		for _, bucket := range s.buckets {
			bucket.UpRoomsCount(allRoomsCount)
		}
		select {
		case <-s.done:
			return
		case <-time.After(time.Second * 10):
		}
	}
}

func getHostName() string {
	name, _ := os.Hostname()
	return name
}

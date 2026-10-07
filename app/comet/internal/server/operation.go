// Package server
// File operation.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-10-07 22:37:59
// Modified 2026-10-07 22:37:59

package server

import (
	"context"
	"sokuim/sokuim-server/app/core/client/socket"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

func (s *Server) RenewOnline(ctx context.Context, serverID string, roomCount map[string]int32) (allRoom map[string]int32, err error) {
	reply, err := s.socketRPC.RenewOnline(ctx, &socket.SocketOnlineReq{
		Server:    serverID,
		RoomCount: roomCount,
	}, grpc.UseCompressor(gzip.Name))
	if err != nil {
		return
	}
	return reply.AllRoomCount, nil
}

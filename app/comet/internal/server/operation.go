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
	"sokuim/sokuim-server/app/comet/pb"
	"sokuim/sokuim-server/app/core/client/socket"
	"sokuim/sokuim-server/pkg/strings"

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

func (s *Server) Receive(ctx context.Context, mid string, p *pb.CometMsgProto) (err error) {
	_, err = s.socketRPC.Receive(ctx, &socket.SocketReceiveReq{
		Mid: mid,
		Proto: &socket.SocketMsgProto{
			Op:   p.Op,
			Body: p.Body,
		},
	})
	return
}

func (s *Server) Operate(ctx context.Context, p *pb.CometMsgProto, ch *Channel, b *Bucket) (err error) {
	switch p.Op {
	case pb.OpChangeRoom:
		err = b.ChangeRoom(string(p.Body), ch)
		p.Op = pb.OpChangeRoomResp
	case pb.OpSub:
		if ops, err := strings.SplitInt32s(string(p.Body), ","); err == nil {
			ch.Watch(ops...)
		}
		p.Op = pb.OpSubResp
	case pb.OpUnsub:
		if ops, err := strings.SplitInt32s(string(p.Body), ","); err == nil {
			ch.UnWatch(ops...)
		}
		p.Op = pb.OpUnsubResp
	default:
		err = s.Receive(ctx, ch.Mid, p)
		p.Body = nil
	}
	return
}

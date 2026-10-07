// Package server
// File server_tcp.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 22:38:59
// Modified 2026-09-29 22:38:59

package server

import (
	"context"
	"net"
	"sokuim/sokuim-server/app/comet/pb"
	"sokuim/sokuim-server/pkg/bufio"
	"sokuim/sokuim-server/pkg/bytes"
	"sokuim/sokuim-server/pkg/logger"
	xtime "sokuim/sokuim-server/pkg/time"
	"time"

	"github.com/google/uuid"
)

func InitTCP(s *Server, addrs []string, accept int, log *logger.Logger) (err error) {
	var (
		bind     string
		listener *net.TCPListener
		addr     *net.TCPAddr
	)

	ctx := context.Background()

	for _, bind = range addrs {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			log.Errorf(ctx, "net.ResolveTCPAddr(tcp, %s) error(%v)", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			log.Errorf(ctx, "net.ListenTCP(tcp, %s) error(%v)", bind, err)
			return
		}
		for i := 0; i < accept; i++ {
			go acceptTCP(s, listener, log)
		}
	}
	log.Info(ctx, listener)
	return
}

func acceptTCP(s *Server, lis *net.TCPListener, log *logger.Logger) {
	var (
		conn *net.TCPConn
		err  error
		r    int
	)
	for {
		if conn, err = lis.AcceptTCP(); err != nil {
			return
		}
		if err = conn.SetKeepAlive(s.conf.TcpConfig.Keepalive); err != nil {
			return
		}
		if err = conn.SetReadBuffer(s.conf.TcpConfig.RcvBuf); err != nil {
			return
		}
		if err = conn.SetWriteBuffer(s.conf.TcpConfig.SndBuf); err != nil {
			return
		}
		go serveTCP(s, conn, r, log)
		if r++; r == maxInt {
			r = 0
		}
	}
}

func serveTCP(s *Server, conn *net.TCPConn, r int, log *logger.Logger) {
	var (
		tr    = s.round.Timer(r)
		rp    = s.round.Reader(r)
		wp    = s.round.Writer(r)
		lAddr = conn.LocalAddr().String()
		rAddr = conn.RemoteAddr().String()
	)
	log.Infof(context.Background(), "start tcp server %v with %v", lAddr, rAddr)
	s.ServeTCP(conn, rp, wp, tr)
}

func (s *Server) ServeTCP(conn *net.TCPConn, rp, wp *bytes.Pool, tr *xtime.Timer) {
	var (
		err     error
		rid     string
		accepts []int32
		hb      time.Duration
		p       *pb.CometMsgProto
		b       *Bucket
		trd     *xtime.TimerData
		lastHb  = time.Now()
		rb      = rp.Get()
		wb      = wp.Get()
		ch      = NewChannel(s.conf.ProtocolConfig.CliProto, s.conf.ProtocolConfig.SrvProto)
		rr      = &ch.Reader
		wr      = &ch.Writer
	)
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	ch.Writer.ResetBuffer(conn, wb.Bytes())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// handshake
	step := 0
	trd = tr.Add(time.Duration(s.conf.ProtocolConfig.HandshakeTimeout), func() {
		conn.Close()
		s.log.Infof(ctx, "key: %s remote: %s, step: %d tcp handshake timeout", ch.Key, conn.RemoteAddr().String(), step)
	})
	ch.IP, _, _ = net.SplitHostPort(conn.LocalAddr().String())

	// must not setadv, only used in auth
	step = 1
	if p, err = ch.CliProto.Set(); err == nil {
		if ch.Mid, ch.Key, rid, accepts, hb, err = s.authTCP(ctx, rr, wr, p); err == nil {
			ch.Watch(accepts...)
			b = s.Bucket(ch.Key)
			err = b.Put(rid, ch)
			s.log.Infof(ctx, "tcp connected key: %s, mid: %s, proto: %+v", ch.Key, ch.Mid, ch)
		}
	}
	s.log.Info(ctx, "aaa: ", ch.IP)
	s.log.Info(ctx, err, rid, accepts, hb, p, b, trd, lastHb, rr, wr, step)
}

func (s *Server) authTCP(ctx context.Context, rr *bufio.Reader, wr *bufio.Writer, p *pb.CometMsgProto) (mid, key, rid string, accepts []int32, hb time.Duration, err error) {
	mid = uuid.New().String()
	return
}

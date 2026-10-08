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
	"io"
	"net"
	"sokuim/sokuim-server/app/comet/pb"
	"sokuim/sokuim-server/pkg/bufio"
	"sokuim/sokuim-server/pkg/bytes"
	"sokuim/sokuim-server/pkg/logger"
	xtime "sokuim/sokuim-server/pkg/time"
	"strings"
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
		}
	}
	step = 2
	if err != nil {
		conn.Close()
		rp.Put(rb)
		wp.Put(wb)
		tr.Del(trd)
		s.log.Infof(ctx, "tcp key: %s handshake failed error: %v", ch.Key, err)
	}
	trd.Key = ch.Key
	tr.Set(trd, hb)
	step = 3
	// hanshake ok start dispatch goroutine
	go s.dispatchTCP(ctx, conn, wr, wp, wb, ch)
	serverHeartbeat := s.RandServerHearbeat()
	for {
		if p, err = ch.CliProto.Set(); err != nil {
			break
		}
		if err = p.ReadTCP(rr); err != nil {
			return
		}
		if p.Op == pb.OpHeartbeat {
			tr.Set(trd, hb)
			p.Op = pb.OpHeartbeatResp
			p.Body = nil
			if now := time.Now(); now.Sub(lastHb) > serverHeartbeat {
				if err1 := s.Heartbeat(ctx, ch.Mid, ch.Key); err1 == nil {
					lastHb = now
				}
			}
			s.log.Infof(ctx, "tcp heartbeat receive key:%s, mid:%d", ch.Key, ch.Mid)
			step++
		} else {
			if err = s.Operate(ctx, p, ch, b); err != nil {
				break
			}
		}
		ch.CliProto.SetAdv()
		ch.Signal()
	}
	if err != nil && err != io.EOF && !strings.Contains(err.Error(), "closed") {
		s.log.Errorf(ctx, "key: %s server tcp failed error(%v)", ch.Key, err)
	}
	b.Del(ch)
	tr.Del(trd)
	rp.Put(rb)
	conn.Close()
	ch.Close()
	if err = s.Disconnect(ctx, ch.Mid, ch.Key); err != nil {
		s.log.Errorf(ctx, "key: %s mid: %d operator do disconnect error(%v)", ch.Key, ch.Mid, err)
	}
	s.log.Info(ctx, "tcp disconnected key: %s mid: %d", ch.Key, ch.Mid)
}

func (s *Server) dispatchTCP(ctx context.Context, conn *net.TCPConn, wr *bufio.Writer, wp *bytes.Pool, wb *bytes.Buffer, ch *Channel) {
	var (
		err    error
		finish bool
		online int32
	)
	s.log.Infof(ctx, "key %s dispatch tcp goroutine", ch.Key)
	for {
		var p = ch.Ready()
		s.log.Infof(ctx, "key %s dispatch msg: %v", ch.Key, p)
		switch p {
		case pb.ProtoFinish:
			s.log.Infof(ctx, "key: %s wakeup exit dispatch goroutine", ch.Key)
			finish = true
			goto field
		case pb.ProtoReady:
			// fetch message from svrbox(client send)
			for {
				if p, err = ch.CliProto.Get(); err != nil {
					break
				}
				if p.Op == pb.OpHeartbeatResp {
					if ch.Room != nil {
						online = ch.Room.OnlineNum()
					}
					if err = p.WriteTCPHeader(wr, online); err != nil {
						goto field
					}
				} else {
					if err = p.WriteTCP(wr); err != nil {
						goto field
					}
				}
				p.Body = nil
				ch.CliProto.GetAdv()
			}
		default:
			if err = p.WriteTCP(wr); err != nil {
				goto field
			}
			s.log.Infof(ctx, "tcp sent a message key:%s mid:%d proto:%+v", ch.Key, ch.Mid, p)
		}
	}
field:
	if err != nil {
		s.log.Errorf(ctx, "key: %s dispatch tcp error(%v)", ch.Key, err)
	}
	conn.Close()
	wp.Put(wb)
	for !finish {
		finish = (ch.Ready() == pb.ProtoFinish)
	}
	s.log.Infof(ctx, "key: %s dispatch goroutine exit", ch.Key)
}

func (s *Server) authTCP(ctx context.Context, rr *bufio.Reader, wr *bufio.Writer, p *pb.CometMsgProto) (mid, key, rid string, accepts []int32, hb time.Duration, err error) {
	mid = uuid.New().String()
	key = uuid.New().String()
	hb = time.Duration(5 * time.Second)
	return
}

// Package pb
// File protocol.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 00:18:07
// Modified 2026-09-29 00:18:07

package pb

import (
	"sokuim/sokuim-server/pkg/bytes"
	"sokuim/sokuim-server/pkg/encoding/binary"
	"sokuim/sokuim-server/pkg/xerr"
)

const MaxBodySize = int32(1 << 12) // 4096

const (
	_packSize      = 4                                   // cap of data pack
	_headerSize    = 2                                   // cap of header
	_opSize        = 4                                   // cap of op
	_heartSize     = 4                                   // cap of heartbeat
	_rawHeaderSize = _packSize + _headerSize + _opSize   // cap of header
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize) // max cap of data pack
)

const (
	_packOffset   = 0                           // offset of data pack
	_headerOffset = _packOffset + _packSize     // offset of header
	_opOffset     = _headerOffset + _headerSize // offset of op
	_heartOffset  = _opOffset + _opSize         // offset of heartbeat
)

var (
	ErrorProtoPackLen = xerr.NewErrCodeMsg(xerr.RequestErr, "Default server codec pack length error")
	ErrProtoHeaderLen = xerr.NewErrCodeMsg(xerr.RequestErr, "Default server codec header length error")
)

var (
	ProtoReady  = &CometMsgProto{Op: OpProtoReady}
	ProtoFinish = &CometMsgProto{Op: OpProtoFinish}
)

func (p *CometMsgProto) WriteTo(b *bytes.Writer) {
	var (
		packLen = _rawHeaderSize + int32(len(p.Body))
		buf     = b.Peek(_rawHeaderSize)
	)
	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	if p.Body != nil {
		b.Write(p.Body)
	}
}

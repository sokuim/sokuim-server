// Package server
// File bucket.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-10-07 21:05:12
// Modified 2026-10-07 21:05:12

// Bucket is a channel holder.
// 每个comet拥有若干bucket(session manager), 用于记录当前comet服务于哪些room和channel
// 作用: Channel和Room的管理者
package server

import (
	"sokuim/sokuim-server/app/comet/internal/config"
	"sokuim/sokuim-server/app/comet/pb"
	"sync"
	"sync/atomic"
)

type Bucket struct {
	c           config.BucketConfig
	rwLock      sync.RWMutex                     // protect the channels for chs
	chs         map[string]*Channel              // 每个Channel相当于一个长链接用户
	rooms       map[string]*Room                 // bucket room channels
	routines    []chan *pb.CometBroadcastRoomReq // 信号传输管道, 用于广播信息给bucket下的所有room
	routinesNum uint64                           // 用于广播信息的chan的数量
	ipCnts      map[string]int32                 // 同一IP链接数量(e.g: 相同IP下多客户端链接的数量)
}

// NewBucket new a bucket struct. store the key with im channel.
func NewBucket(c config.BucketConfig) (b *Bucket) {
	b = new(Bucket)
	b.chs = make(map[string]*Channel, c.Channel)
	b.ipCnts = make(map[string]int32)
	b.c = c
	b.rooms = make(map[string]*Room, c.Room)
	b.routines = make([]chan *pb.CometBroadcastRoomReq, c.RoutineAmount)
	for i := uint64(0); i < c.RoutineAmount; i++ {
		c := make(chan *pb.CometBroadcastRoomReq, c.RoutineSize)
		b.routines[i] = c
		go b.roomproc(c)
	}
	return
}

// ChannelCount channel count in the bucket
func (b *Bucket) ChannelCount() int {
	return len(b.chs)
}

// Close close all channels in the bucket.
// 关闭该 bucket 下的所有连接：向每个 channel 的 signal 发送 ProtoFinish，
// 通知对应的读写 goroutine 优雅退出（配合 Server.Close 做平滑下线）。
func (b *Bucket) Close() {
	b.rwLock.RLock()
	for _, ch := range b.chs {
		ch.Close()
	}
	b.rwLock.RUnlock()
}

// RoomCount room count in the bucket
// 获取该bucket下的房间数量
func (b *Bucket) RoomCount() int {
	return len(b.rooms)
}

// RoomsCount get all room id where online number > 0.
// 获取该bucket下, channel不为0的所有room信息(ID号以及该ID下的在线数量)
func (b *Bucket) RoomsCount() (res map[string]int32) {
	var (
		roomID string
		room   *Room
	)
	b.rwLock.RLock()
	res = make(map[string]int32)
	for roomID, room = range b.rooms {
		if room.Online > 0 {
			res[roomID] = room.Online
		}
	}
	b.rwLock.RUnlock()
	return
}

// ChangeRoom change ro room
// nrid: new room id 新房间ID号
// 将channel替换到新的room中, 如果新房间ID为空, 则channel的房间为空
func (b *Bucket) ChangeRoom(nrid string, ch *Channel) (err error) {
	var (
		nroom *Room // new room: 新房间
		ok    bool
		oroom = ch.Room // old room: 旧房间
	)
	// change to no room
	if nrid == "" {
		if oroom != nil && oroom.Del(ch) {
			b.DelRoom(oroom)
		}
		ch.Room = nil
		return
	}
	b.rwLock.Lock()
	// 新房间如果不存在则创建, 并添加到bucket中
	if nroom, ok = b.rooms[nrid]; !ok {
		nroom = NewRoom(nrid)
		b.rooms[nrid] = nroom
	}
	b.rwLock.Unlock()
	// 删除旧房间
	if oroom != nil && oroom.Del(ch) {
		b.DelRoom(oroom)
	}
	// channel 添加到新房间中
	if err = nroom.Put(ch); err != nil {
		return
	}
	ch.Room = nroom
	return
}

// Put put a channel according with sub key.
// 将channel添加到bucket所属的room中
func (b *Bucket) Put(rid string, ch *Channel) (err error) {
	var (
		room *Room
		ok   bool
	)
	b.rwLock.Lock()
	// close old channel
	if dch := b.chs[ch.Key]; dch != nil {
		dch.Close()
	}
	b.chs[ch.Key] = ch
	if rid != "" {
		if room, ok = b.rooms[rid]; !ok {
			room = NewRoom(rid)
			b.rooms[rid] = room
		}
		ch.Room = room
	}
	b.ipCnts[ch.IP]++
	b.rwLock.Unlock()
	if room != nil {
		err = room.Put(ch)
	}
	return
}

// Del delete the channel by sub key.
// 删除当前bucket的channel, 同时删除channel对应的room
func (b *Bucket) Del(dch *Channel) {
	room := dch.Room
	b.rwLock.Lock()
	if ch, ok := b.chs[dch.Key]; ok {
		if ch == dch {
			delete(b.chs, ch.Key)
		}
		// ip counter
		if b.ipCnts[ch.IP] > 1 {
			b.ipCnts[ch.IP]--
		} else {
			delete(b.ipCnts, ch.IP)
		}
	}
	b.rwLock.Unlock()
	if room != nil && room.Del(dch) {
		// if empty room, must delete from bucket
		b.DelRoom(room)
	}
}

// Channel get a channel by sub key.
// 根据key获取当前bucket的channel
func (b *Bucket) Channel(key string) (ch *Channel) {
	b.rwLock.RLock()
	ch = b.chs[key]
	b.rwLock.RUnlock()
	return
}

// Broadcast push msgs to all channels in the bucket.
// 广播消息给当前bucket下的所有channel
func (b *Bucket) Broadcast(p *pb.CometMsgProto, op int32) {
	var ch *Channel
	b.rwLock.RLock()
	for _, ch = range b.chs {
		if !ch.NeedPush(op) {
			continue
		}
		_ = ch.Push(p)
	}
	b.rwLock.RUnlock()
}

// Room get a room by roomid.
// 根据roomID获取当前bucket下的room
func (b *Bucket) Room(rid string) (room *Room) {
	b.rwLock.RLock()
	room = b.rooms[rid]
	b.rwLock.RUnlock()
	return
}

// DelRoom delete a room by roomid.
// 根据roomID删除当前当前bucket的room
func (b *Bucket) DelRoom(room *Room) {
	b.rwLock.Lock()
	delete(b.rooms, room.ID)
	b.rwLock.Unlock()
	room.Close()
}

// BroadcastRoom broadcast a message to specified room
// 根据roomID广播消息给当前bucket下的room的所有
// routines: []chan *pb.CometBroadcastRoomReq
// 有一个协程在专门监听routines数据: roomproc
// 此处写入数据后, roomproc方法会执行
func (b *Bucket) BroadcastRoom(arg *pb.CometBroadcastRoomReq) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.c.RoutineAmount
	b.routines[num] <- arg
}

// Rooms get all room id where online number > 0.
// 获取当前bucket中的online > 0 的所有room的ID
func (b *Bucket) Rooms() (res map[string]struct{}) {
	var (
		roomID string
		room   *Room
	)
	res = make(map[string]struct{})
	b.rwLock.RLock()
	for roomID, room = range b.rooms {
		if room.Online > 0 {
			res[roomID] = struct{}{}
		}
	}
	b.rwLock.RUnlock()
	return
}

// IPCount get ip count.
// 获取bucket下存有的所有客户端IP
func (b *Bucket) IPCount() (res map[string]struct{}) {
	var (
		ip string
	)
	b.rwLock.RLock()
	res = make(map[string]struct{}, len(b.ipCnts))
	for ip = range b.ipCnts {
		res[ip] = struct{}{}
	}
	b.rwLock.RUnlock()
	return
}

// UpRoomsCount update all room count
// 更新bucket下的room的在线数量
func (b *Bucket) UpRoomsCount(roomCountMap map[string]int32) {
	var (
		roomID string
		room   *Room
	)
	b.rwLock.RLock()
	for roomID, room = range b.rooms {
		room.AllOnline = roomCountMap[roomID]
	}
	b.rwLock.RUnlock()
}

// roomproc
// 协程监听房间推送
func (b *Bucket) roomproc(c chan *pb.CometBroadcastRoomReq) {
	for {
		arg := <-c
		if room := b.Room(arg.RoomID); room != nil {
			room.Push(arg.Proto)
		}
	}
}

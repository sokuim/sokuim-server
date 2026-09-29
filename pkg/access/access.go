// Package access
// File access.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 23:37:51
// Modified 2026-09-29 23:37:51

// 名单数据通过 Store 接口从外部存储（如 Redis SET）读取，并缓存到本地内存，
// 供连接鉴权等高频场景做 O(1) 查询；支持管理员动态增删，增删后立即刷新缓存、近乎实时生效，
// 同时由后台 goroutine 定期刷新，保证多实例最终一致。
//
// 使用示例：
//
//	store := myRedisStore{} // 实现 Store 接口
//	m := access.New(store, access.Options{WhiteEnable: true, BlackEnable: true})
//	if err := m.Check(mid); err != nil { /* 拒绝 */ }
package access

import (
	"context"
	"errors"
	"sync"
	"time"

	log "github.com/golang/glog"
)

// ErrDenied 账号未通过黑白名单准入校验
var ErrDenied = errors.New("access denied")

// 黑/白名单在存储中的默认 key（可被 Options 覆盖）。
const (
	DefaultWhiteKey = "whitelist_mids"
	DefaultBlackKey = "blacklist_mids"
)

// Store 名单存储接口：不同服务可用各自的后端实现（如 Redis SET）。
// key 用于区分名单（黑名单 / 白名单），具体取值由 Options 决定。
type Store interface {
	// Mids 全量加载某个名单的所有成员。
	Mids(ctx context.Context, key string) ([]string, error)
	// Add 向某个名单添加成员。
	Add(ctx context.Context, key string, mid string) error
	// Remove 从某个名单移除成员。
	Remove(ctx context.Context, key string, mid string) error
}

// Options 黑白名单配置。
type Options struct {
	WhiteEnable bool          // 是否启用白名单（启用后仅名单内 mid 放行）
	BlackEnable bool          // 是否启用黑名单（启用后名单内 mid 一律拒绝）
	Refresh     time.Duration // 本地缓存刷新周期
	WhiteKey    string        // 白名单存储 key（空则用 DefaultWhiteKey）
	BlackKey    string        // 黑名单存储 key（空则用 DefaultBlackKey）
}

// Result 账号访问状态：黑白名单综合判断后的最终结论与原因。
type Result struct {
	Mid         string // 用户 mid(UUID)
	Blacklisted bool   // 是否命中黑名单
	Whitelisted bool   // 是否命中白名单
	Allowed     bool   // 是否允许接入（最终结论）
	Reason      string // 拒绝/放行原因
}

// list 名单本地缓存（黑/白名单通用）。
type list struct {
	mu   sync.RWMutex
	mids map[string]struct{}
}

func newList() *list {
	return &list{mids: make(map[string]struct{})}
}

func (l *list) contains(mid string) bool {
	l.mu.RLock()
	_, ok := l.mids[mid]
	l.mu.RUnlock()
	return ok
}

// Manager 黑白名单管理器。非并发安全的配置在 New 时确定，此后只读。
type Manager struct {
	opts  Options
	store Store
	white *list // nil 表示未启用白名单
	black *list // nil 表示未启用黑名单
}

// New 创建管理器：按配置建立黑/白名单缓存，初始加载一次并启动定期刷新。
// 若黑白名单均未启用，则不启动后台刷新（返回一个 Check 恒放行的空管理器）。
func New(store Store, opts Options) *Manager {
	if opts.WhiteKey == "" {
		opts.WhiteKey = DefaultWhiteKey
	}
	if opts.BlackKey == "" {
		opts.BlackKey = DefaultBlackKey
	}
	if opts.Refresh <= 0 {
		opts.Refresh = 10 * time.Second
	}
	m := &Manager{opts: opts, store: store}
	if opts.WhiteEnable {
		m.white = newList()
	}
	if opts.BlackEnable {
		m.black = newList()
	}
	if !opts.WhiteEnable && !opts.BlackEnable {
		return m
	}
	_ = m.load(context.Background())
	go m.refresh()
	return m
}

// load 从 Store 全量加载黑/白名单到本地缓存。
// 任一名单加载失败即返回错误（保留旧缓存，避免误清空已有名单）。
func (m *Manager) load(ctx context.Context) error {
	if m.white != nil {
		mids, err := m.store.Mids(ctx, m.opts.WhiteKey)
		if err != nil {
			return err
		}
		m.white.mu.Lock()
		m.white.mids = toSet(mids)
		m.white.mu.Unlock()
	}
	if m.black != nil {
		mids, err := m.store.Mids(ctx, m.opts.BlackKey)
		if err != nil {
			return err
		}
		m.black.mu.Lock()
		m.black.mids = toSet(mids)
		m.black.mu.Unlock()
	}
	return nil
}

func toSet(mids []string) map[string]struct{} {
	m := make(map[string]struct{}, len(mids))
	for _, mid := range mids {
		m[mid] = struct{}{}
	}
	return m
}

// refresh 后台定期刷新缓存。
func (m *Manager) refresh() {
	t := time.NewTicker(m.opts.Refresh)
	defer t.Stop()
	for range t.C {
		if err := m.load(context.Background()); err != nil {
			log.Errorf("access load error(%v)", err)
		}
	}
}

// Check 连接准入校验：黑名单优先于白名单，未通过返回 ErrDenied。
// 规则：
//  1. 命中黑名单 → 一律拒绝（即使同时命中白名单）；
//  2. 启用白名单且未命中 → 拒绝（白名单语义为「名单内才放行」）；
//  3. 其余情况 → 放行。
func (m *Manager) Check(mid string) error {
	r := m.Account(mid)
	if !r.Allowed {
		log.Warningf("access denied mid:%d reason:%s", mid, r.Reason)
		return ErrDenied
	}
	return nil
}

// Account 获取账号访问信息（黑白名单综合结论），不修改任何状态。
func (m *Manager) Account(mid string) *Result {
	r := &Result{Mid: mid, Allowed: true}
	// 1. 黑名单优先：命中即拒绝，无论白名单如何
	if m.black != nil && m.black.contains(mid) {
		r.Blacklisted = true
		r.Allowed = false
		r.Reason = "mid in blacklist"
		return r
	}
	// 2. 白名单：仅在启用时（m.white != nil）生效
	if m.white != nil {
		if m.white.contains(mid) {
			r.Whitelisted = true
			r.Reason = "mid in whitelist"
		} else {
			r.Allowed = false
			r.Reason = "mid not in whitelist"
		}
	}
	return r
}

// AddWhite 将 mid 加入白名单：写存储后立即刷新缓存，近乎实时生效。
func (m *Manager) AddWhite(ctx context.Context, mid string) error {
	if err := m.store.Add(ctx, m.opts.WhiteKey, mid); err != nil {
		return err
	}
	return m.load(ctx)
}

// RemoveWhite 将 mid 移出白名单：写存储后立即刷新缓存，近乎实时生效。
func (m *Manager) RemoveWhite(ctx context.Context, mid string) error {
	if err := m.store.Remove(ctx, m.opts.WhiteKey, mid); err != nil {
		return err
	}
	return m.load(ctx)
}

// AddBlack 将 mid 加入黑名单：写存储后立即刷新缓存，近乎实时生效。
func (m *Manager) AddBlack(ctx context.Context, mid string) error {
	if err := m.store.Add(ctx, m.opts.BlackKey, mid); err != nil {
		return err
	}
	return m.load(ctx)
}

// RemoveBlack 将 mid 移出黑名单：写存储后立即刷新缓存，近乎实时生效。
func (m *Manager) RemoveBlack(ctx context.Context, mid string) error {
	if err := m.store.Remove(ctx, m.opts.BlackKey, mid); err != nil {
		return err
	}
	return m.load(ctx)
}

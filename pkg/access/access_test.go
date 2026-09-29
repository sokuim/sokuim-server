// Package access
// File access_test.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-29 23:45:04
// Modified 2026-09-29 23:45:04

package access

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// memStore 内存实现的 Store，用于测试（不依赖 Redis）。
type memStore struct {
	mu   sync.Mutex
	data map[string]map[string]struct{}
}

func newMemStore() *memStore {
	return &memStore{data: make(map[string]map[string]struct{})}
}

func (s *memStore) Mids(ctx context.Context, key string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.data[key]
	mids := make([]string, 0, len(m))
	for mid := range m {
		mids = append(mids, mid)
	}
	return mids, nil
}

func (s *memStore) Add(ctx context.Context, key string, mid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data[key] == nil {
		s.data[key] = make(map[string]struct{})
	}
	s.data[key][mid] = struct{}{}
	return nil
}

func (s *memStore) Remove(ctx context.Context, key string, mid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data[key], mid)
	return nil
}

// TestBlacklistPriority 验证黑名单优先于白名单的准入规则。
func TestBlacklistPriority(t *testing.T) {
	m := New(newMemStore(), Options{WhiteEnable: true, BlackEnable: true, Refresh: time.Hour})
	ctx := context.Background()

	_ = m.AddWhite(ctx, "1")
	_ = m.AddBlack(ctx, "2")
	_ = m.AddWhite(ctx, "2") // 2 同时命中黑白名单

	assert.Nil(t, m.Check("1"))              // 白名单内 → 放行
	assert.Equal(t, ErrDenied, m.Check("2")) // 同时命中 → 黑名单优先，拒绝
	assert.Equal(t, ErrDenied, m.Check("3")) // 不在白名单 → 白名单启用，拒绝
}

// TestWhitelistDisabled 验证未启用白名单时仅黑名单生效。
func TestWhitelistDisabled(t *testing.T) {
	m := New(newMemStore(), Options{BlackEnable: true, Refresh: time.Hour})
	ctx := context.Background()

	assert.Nil(t, m.Check("999")) // 未启用白名单，默认放行

	_ = m.AddBlack(ctx, "999")
	assert.Equal(t, ErrDenied, m.Check("999")) // 命中黑名单 → 拒绝
}

// TestAccount 验证账号访问信息查询。
func TestAccount(t *testing.T) {
	m := New(newMemStore(), Options{WhiteEnable: true, Refresh: time.Hour})
	ctx := context.Background()

	_ = m.AddWhite(ctx, "1")

	r := m.Account("1")
	assert.True(t, r.Whitelisted)
	assert.True(t, r.Allowed)

	r2 := m.Account("2")
	assert.False(t, r2.Whitelisted)
	assert.False(t, r2.Allowed)
	assert.Equal(t, "mid not in whitelist", r2.Reason)
}

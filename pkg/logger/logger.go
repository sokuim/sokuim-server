// Package logger
// File logger.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-21 21:32:39
// Modified 2026-09-21 21:32:39

package logger

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

type Logger struct {
	conf Config
	log  *glog.Logger
}

func NewLogger(conf Config) *Logger {
	log, err := newLogger(conf)
	if err != nil {
		panic(err)
	}
	return &Logger{
		conf: conf,
		log:  log,
	}
}

func (l *Logger) Info(ctx context.Context, args ...interface{}) {
	l.log.Info(ctx, args)
}

func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.log.Infof(ctx, format, args)
}

func (l *Logger) Error(ctx context.Context, args ...interface{}) {
	l.log.Error(ctx, args)
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.log.Errorf(ctx, format, args)
}

func (l *Logger) Warning(ctx context.Context, args ...interface{}) {
	l.log.Warning(ctx, args)
}

func (l *Logger) Warningf(ctx context.Context, format string, args ...interface{}) {
	l.log.Warningf(ctx, format, args)
}

func (l *Logger) Notice(ctx context.Context, args ...interface{}) {
	l.log.Notice(ctx, args)
}

func (l *Logger) Noticef(ctx context.Context, format string, args ...interface{}) {
	l.log.Noticef(ctx, format, args)
}

func (l *Logger) Debug(ctx context.Context, args ...interface{}) {
	l.log.Debug(ctx, args)
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.log.Debugf(ctx, format, args)
}

func newLogger(conf Config) (log *glog.Logger, err error) {
	log = glog.New()
	folder := strings.TrimSpace(conf.Folder)
	if folder == "" {
		err = errors.New("log folder is empty")
		return
	}
	path := strings.TrimSpace(conf.Path)
	if path == "" {
		err = errors.New("log path is empty")
		return
	}

	logPath := filepath.Join(path, folder)
	fmt.Println("logPath:", logPath)
	err = log.SetConfigWithMap(g.Map{
		"path":               logPath,
		"stdout":             conf.Stdout,
		"StStatus":           0,
		"rotateSize":         fmt.Sprintf("%dM", conf.RotateSize),
		"rotateBackupLimit":  conf.RotateLimit,
		"RotateBackupExpire": fmt.Sprintf("%dd", conf.RotateExpire),
	})
	return
}

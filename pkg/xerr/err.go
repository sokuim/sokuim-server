// Package xerr
// File err.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-01 01:32:01
// Modified 2026-09-01 01:32:01

package xerr

import "fmt"

type CodeError struct {
	errCode int32
	errMsg  string
}

func (e *CodeError) GetErrCode() int32 {
	return e.errCode
}

func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("code:%d msg:%s", e.errCode, e.errMsg)
}

func NewErrCodeMsg(code int32, msg string) *CodeError {
	return &CodeError{errCode: code, errMsg: msg}
}

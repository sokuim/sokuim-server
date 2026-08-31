// Package result
// File http_result.go
// Copyright 2026 sokuim.com - All Rights Reserved
// Link https://www.sokuim.com
// Author stiffer.chen <stiffer@sokuim.com>
// Created 2026-09-01 01:34:08
// Modified 2026-09-01 01:34:08

package result

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"

	"net/http"
	"sokuim/sokuim-server/pkg/xerr"
)

func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, msg string, err error) {
	if err == nil {
		res := Success(msg, resp)
		httpx.WriteJson(w, http.StatusOK, res)
		return
	}

	errCode := uint32(http.StatusInternalServerError)
	errMsg := "Server Error"

	causeErr := errors.Cause(err)
	if e, ok := causeErr.(*xerr.CodeError); ok {
		errCode = uint32(e.GetErrCode())
		errMsg = e.GetErrMsg()
	} else {
		if rpcStatus, rpcOk := status.FromError(causeErr); rpcOk {
			gprCode := uint32(rpcStatus.Code())
			errCode = gprCode
			errMsg = rpcStatus.Message()
		}
	}

	logx.WithContext(r.Context()).Errorf("[api-err]: %x+v", err)
	httpx.WriteJson(w, http.StatusOK, Error(errCode, errMsg))
}

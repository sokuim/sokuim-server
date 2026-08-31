package common

import (
	"context"

	"sokuim/sokuim-server/app/gateway/internal/svc"
	"sokuim/sokuim-server/app/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthzLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Healthz check
func NewHealthzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthzLogic {
	return &HealthzLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthzLogic) Healthz(req *types.CommonHealthReq) (resp *types.CommonHealthResp, msg string, err error) {
	resp = &types.CommonHealthResp{
		Status: "ok",
	}
	return
}

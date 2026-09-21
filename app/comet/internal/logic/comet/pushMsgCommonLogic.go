package cometlogic

import (
	"context"

	"sokuim/sokuim-server/app/comet/internal/svc"
	"sokuim/sokuim-server/app/comet/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushMsgCommonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushMsgCommonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushMsgCommonLogic {
	return &PushMsgCommonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushMsgCommonLogic) PushMsgCommon(in *pb.CometPushCommonReq) (*pb.CometPushCommonResp, error) {
	return &pb.CometPushCommonResp{}, nil
}

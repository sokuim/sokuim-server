package socketlogic

import (
	"context"

	"sokuim/sokuim-server/app/core/internal/svc"
	"sokuim/sokuim-server/app/core/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HeartbeatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHeartbeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HeartbeatLogic {
	return &HeartbeatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HeartbeatLogic) Heartbeat(in *pb.SocketHeartbeatReq) (*pb.SocketHeartbeatResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SocketHeartbeatResp{}, nil
}

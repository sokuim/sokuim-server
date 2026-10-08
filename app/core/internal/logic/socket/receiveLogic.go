package socketlogic

import (
	"context"

	"sokuim/sokuim-server/app/core/internal/svc"
	"sokuim/sokuim-server/app/core/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveLogic {
	return &ReceiveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReceiveLogic) Receive(in *pb.SocketReceiveReq) (*pb.SocketReceiveResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SocketReceiveResp{}, nil
}

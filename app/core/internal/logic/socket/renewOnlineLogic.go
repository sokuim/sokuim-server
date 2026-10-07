package socketlogic

import (
	"context"

	"sokuim/sokuim-server/app/core/internal/svc"
	"sokuim/sokuim-server/app/core/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RenewOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRenewOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RenewOnlineLogic {
	return &RenewOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RenewOnlineLogic) RenewOnline(in *pb.SocketOnlineReq) (*pb.SocketOnlineResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SocketOnlineResp{}, nil
}

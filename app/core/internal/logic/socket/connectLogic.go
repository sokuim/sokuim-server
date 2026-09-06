package socketlogic

import (
	"context"

	"sokuim/sokuim-server/app/core/internal/svc"
	"sokuim/sokuim-server/app/core/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConnectLogic {
	return &ConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Connect valid socket connect
func (l *ConnectLogic) Connect(in *pb.SocketConnectReq) (*pb.SocketConnectResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SocketConnectResp{}, nil
}

package common

import (
	"net/http"
	"sokuim/sokuim-server/pkg/result"

	"sokuim/sokuim-server/app/gateway/internal/logic/common"
	"sokuim/sokuim-server/app/gateway/internal/svc"
	"sokuim/sokuim-server/app/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Healthz check
func HealthzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CommonHealthReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := common.NewHealthzLogic(r.Context(), svcCtx)
		resp, msg, err := l.Healthz(&req)
		result.HttpResult(r, w, resp, msg, err)
	}
}

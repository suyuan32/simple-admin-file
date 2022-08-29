package file

import (
	"net/http"

	"github.com/suyuan32/simple-admin-file/api/internal/logic/file"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// swagger:route POST /file/status file changePublicStatus
// Change file public status | 改变文件公开状态
// Responses:
//   200: SimpleMsg
//   401: SimpleMsg
//   500: SimpleMsg

func ChangePublicStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangeStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := file.NewChangePublicStatusLogic(r.Context(), svcCtx)
		resp, err := l.ChangePublicStatus(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

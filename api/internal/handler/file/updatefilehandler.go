package file

import (
	"net/http"

	"github.com/suyuan32/simple-admin-file/api/internal/logic/file"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// swagger:route POST /file/list file updateFile
// Update file information | 更新文件信息
// Responses:
//   200: SimpleMsg
//   401: SimpleMsg
//   500: SimpleMsg

func UpdateFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateFileReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := file.NewUpdateFileLogic(r.Context(), svcCtx)
		resp, err := l.UpdateFile(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

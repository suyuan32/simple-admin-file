package file

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/suyuan32/simple-admin-file/internal/logic/file"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
)

// swagger:route post /file/delete_by_url file DeleteFileByUrl
//
// Delete file information | 删除文件信息
//
// Delete file information | 删除文件信息
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: FileDeleteReq
//
// Responses:
//  200: BaseMsgResp

func DeleteFileByUrlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileDeleteReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewDeleteFileByUrlLogic(r.Context(), svcCtx)
		resp, err := l.DeleteFileByUrl(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

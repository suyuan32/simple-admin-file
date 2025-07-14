package cloudfile

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/suyuan32/simple-admin-file/internal/logic/cloudfile"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
)

// swagger:route post /cloud_file/delete_by_url cloudfile DeleteCloudFileByUrl
//
// Delete cloud file information | 删除云文件信息
//
// Delete cloud file information | 删除云文件信息
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: CloudFileDeleteReq
//
// Responses:
//  200: BaseMsgResp

func DeleteCloudFileByUrlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CloudFileDeleteReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cloudfile.NewDeleteCloudFileByUrlLogic(r.Context(), svcCtx)
		resp, err := l.DeleteCloudFileByUrl(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

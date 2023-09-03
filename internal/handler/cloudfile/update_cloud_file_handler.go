package cloudfile

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/suyuan32/simple-admin-file/internal/logic/cloudfile"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
)

// swagger:route post /cloud_file/update cloudfile UpdateCloudFile
//
// Update cloud file information | 更新云文件
//
// Update cloud file information | 更新云文件
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: CloudFileInfo
//
// Responses:
//  200: BaseMsgResp

func UpdateCloudFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CloudFileInfo
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cloudfile.NewUpdateCloudFileLogic(r.Context(), svcCtx)
		resp, err := l.UpdateCloudFile(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

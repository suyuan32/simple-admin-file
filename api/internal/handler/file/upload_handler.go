package file

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/suyuan32/simple-admin-file/api/internal/logic/file"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
)

// swagger:route post /upload file Upload
//
// Upload file | 上传文件
//
// Upload file | 上传文件
//
// Responses:
//  200: UploadResp
//  401: SimpleMsg
//  500: SimpleMsg

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewUploadLogic(r, svcCtx)
		resp, err := l.Upload()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

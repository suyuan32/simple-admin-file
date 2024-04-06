package file

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/suyuan32/simple-admin-file/internal/logic/file"
	"github.com/suyuan32/simple-admin-file/internal/svc"
)

// swagger:route post /upload file Upload
//
// Upload file | 上传文件
//
// Upload file | 上传文件
//
// Responses:
//  200: UploadResp

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	fmt.Println("=========enter upload file handler")
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewUploadLogic(r, svcCtx)
		resp, err := l.Upload()
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

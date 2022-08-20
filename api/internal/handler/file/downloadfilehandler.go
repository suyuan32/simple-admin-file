package file

import (
	"github.com/suyuan32/simple-admin-file/api/internal/logic/file"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io/ioutil"
	"net/http"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
)

func DownloadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := file.NewDownloadFileLogic(r.Context(), svcCtx)
		filePath, err := l.DownloadFile(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		} else {
			body, err := ioutil.ReadFile(filePath)
			if err != nil {
				httpx.Error(w, httpx.NewApiError(http.StatusInternalServerError, err.Error()))
				return
			}
			w.Write(body)
			return
		}
	}
}

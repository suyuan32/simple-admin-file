package file

import (
	"net/http"
	"os"

	"github.com/suyuan32/simple-admin-file/api/internal/logic/file"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// swagger:route GET /file/download/{id} file downloadFile
// Download file | 下载文件
// Parameters:
//  + name: id
//    require: true
//    in: path
// Responses:
//   200: SimpleMsg
//   401: SimpleMsg
//   500: SimpleMsg

func DownloadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IDPathReq
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
			body, err := os.ReadFile(filePath)
			if err != nil {
				httpx.Error(w, httpx.NewApiError(http.StatusInternalServerError, err.Error()))
				return
			}
			w.Header().Set("Accept-Encoding", "identity;q=1, *;q=0")
			w.Write(body)
			return
		}
	}
}

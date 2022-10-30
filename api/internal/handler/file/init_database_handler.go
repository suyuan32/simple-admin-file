package file

import (
	"net/http"

	"github.com/suyuan32/simple-admin-file/api/internal/logic/file"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// swagger:route get /init/database file InitDatabase
//
// Initialize database | 初始化数据库
//
// Initialize database | 初始化数据库
//
// Responses:
//  200: SimpleMsg
//  401: SimpleMsg
//  500: SimpleMsg

func InitDatabaseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewInitDatabaseLogic(r.Context(), svcCtx)
		resp, err := l.InitDatabase()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

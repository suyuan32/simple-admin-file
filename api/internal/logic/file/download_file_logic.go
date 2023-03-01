package file

import (
	"context"
	"net/http"
	"path"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
	"github.com/suyuan32/simple-admin-file/pkg/utils/dberrorhandler"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lang   string
}

func NewDownloadFileLogic(r *http.Request, svcCtx *svc.ServiceContext) *DownloadFileLogic {
	return &DownloadFileLogic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		svcCtx: svcCtx,
		lang:   r.Header.Get("Accept-Language"),
	}
}

func (l *DownloadFileLogic) DownloadFile(req *types.IDPathReq) (filePath string, err error) {
	file, err := l.svcCtx.DB.File.Get(l.ctx, req.Id)

	if err != nil {
		return "", dberrorhandler.DefaultEntError(err, req)
	}

	if file.Status == 1 {
		logx.Infow("public download", logx.Field("fileName", file.Name),
			logx.Field("userId", l.ctx.Value("userId").(string)),
			logx.Field("filePath", file.Path))
		return path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, file.Path), nil
	} else {
		logx.Infow("private download", logx.Field("fileName", file.Name),
			logx.Field("userId", l.ctx.Value("userId").(string)),
			logx.Field("filePath", file.Path))
		return path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, file.Path), nil
	}
}

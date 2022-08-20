package file

import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type DownloadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadFileLogic {
	return &DownloadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadFileLogic) DownloadFile(req *types.DownloadReq) (filePath string, err error) {
	var target model.FileInfo
	check := l.svcCtx.DB.Where("id = ?", req.Id).First(&target)
	if check.Error != nil {
		return "", httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if check.RowsAffected == 0 {
		return "", httpx.NewApiErrorWithoutMsg(http.StatusNotFound)
	}
	// judge the user is admin or owner
	// 只有管理员和拥有者能下载文件
	if l.ctx.Value("roleId").(json.Number).String() != "1" && l.ctx.Value("userId").(string) != target.UserUUID {
		return "", httpx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
	}

	if target.Status {
		return path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, target.Path), nil
	} else {
		return path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, target.Path), nil
	}
}

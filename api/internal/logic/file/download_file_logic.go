package file

import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/suyuan32/simple-message/core/log"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
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

func (l *DownloadFileLogic) DownloadFile(req *types.IDPathReq) (filePath string, err error) {
	var target model.FileInfo
	check := l.svcCtx.DB.Where("id = ?", req.ID).First(&target)
	if check.Error != nil {
		return "", errorx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if check.RowsAffected == 0 {
		return "", errorx.NewApiErrorWithoutMsg(http.StatusNotFound)
	}

	// only admin and owner can do it
	roleId := l.ctx.Value("roleId").(json.Number).String()
	userId := l.ctx.Value("userId").(string)
	if roleId != "1" && userId != target.UserUUID {
		logx.Errorw(log.OperationNotAllow, logx.Field("RoleId", roleId),
			logx.Field("UserId", userId))
		return "", errorx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
	}

	if target.Status {
		logx.Infow("Public download", logx.Field("FileName", target.Name), logx.Field("UserId", userId),
			logx.Field("FilePath", target.Path))
		return path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, target.Path), nil
	} else {
		logx.Infow("Private download", logx.Field("FileName", target.Name), logx.Field("UserId", userId),
			logx.Field("FilePath", target.Path))
		return path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, target.Path), nil
	}
}

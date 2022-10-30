package file

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/suyuan32/simple-message/core/log"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
)

type DeleteFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFileLogic) DeleteFile(req *types.IDReq) (resp *types.SimpleMsg, err error) {
	var target model.FileInfo
	result := l.svcCtx.DB.Where("id = ?", req.ID).First(&target)
	if result.Error != nil {
		logx.Errorw(log.DatabaseError, logx.Field("detail", result.Error.Error()))
		return nil, errorx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if result.RowsAffected == 0 {
		logx.Errorw("file cannot be found", logx.Field("fileId", req.ID))
		return nil, errorx.NewApiErrorWithoutMsg(http.StatusNotFound)
	}

	// only admin and owner can do it
	roleId := l.ctx.Value("roleId").(json.Number).String()
	userId := l.ctx.Value("userId").(string)
	if roleId != "1" && userId != target.UserUUID {
		logx.Errorw(log.OperationNotAllow, logx.Field("roleId", roleId),
			logx.Field("userId", userId))
		return nil, errorx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
	}

	if target.Status {
		err = os.RemoveAll(l.svcCtx.Config.UploadConf.PublicStorePath + target.Path)
		if err != nil {
			logx.Errorw("fail to remove the file", logx.Field("path",
				l.svcCtx.Config.UploadConf.PublicStorePath+target.Path), logx.Field("detail", err.Error()))
			return nil, errorx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	} else {
		err = os.RemoveAll(l.svcCtx.Config.UploadConf.PrivateStorePath + target.Path)
		if err != nil {
			logx.Errorw("fail to remove the file", logx.Field("path",
				l.svcCtx.Config.UploadConf.PrivateStorePath+target.Path), logx.Field("detail", err.Error()))
			return nil, errorx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	}

	result = l.svcCtx.DB.Delete(&target)
	if result.Error != nil {
		logx.Errorw(log.DatabaseError, logx.Field("detail", result.Error.Error()))
		return nil, errorx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}

	logx.Infow("delete file successfully", logx.Field("detail", target))
	return &types.SimpleMsg{Msg: errorx.DeleteSuccess}, nil
}

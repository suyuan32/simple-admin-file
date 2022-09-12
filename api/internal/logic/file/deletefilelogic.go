package file

import (
	"context"
	"encoding/json"
	"github.com/suyuan32/simple-admin-file/api/internal/util/logmessage"
	"net/http"
	"os"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
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
		logx.Errorw(logmessage.DatabaseError, logx.Field("Detail", result.Error.Error()))
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if result.RowsAffected == 0 {
		logx.Errorw("File does not find", logx.Field("FileId", req.ID))
		return nil, httpx.NewApiErrorWithoutMsg(http.StatusNotFound)
	}

	// only admin and owner can do it
	roleId := l.ctx.Value("roleId").(json.Number).String()
	userId := l.ctx.Value("userId").(string)
	if roleId != "1" && userId != target.UserUUID {
		logx.Errorw(logmessage.OperationNotAllow, logx.Field("RoleId", roleId),
			logx.Field("UserId", userId))
		return nil, httpx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
	}

	if target.Status {
		err = os.RemoveAll(l.svcCtx.Config.UploadConf.PublicStorePath + target.Path)
		if err != nil {
			logx.Errorw("Fail to remove the file", logx.Field("Path",
				l.svcCtx.Config.UploadConf.PublicStorePath+target.Path), logx.Field("Detail", err.Error()))
			return nil, httpx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	} else {
		err = os.RemoveAll(l.svcCtx.Config.UploadConf.PrivateStorePath + target.Path)
		if err != nil {
			logx.Errorw("Fail to remove the file", logx.Field("Path",
				l.svcCtx.Config.UploadConf.PrivateStorePath+target.Path), logx.Field("Detail", err.Error()))
			return nil, httpx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	}

	result = l.svcCtx.DB.Delete(&target)
	if result.Error != nil {
		logx.Errorw(logmessage.DatabaseError, logx.Field("Detail", result.Error.Error()))
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}

	logx.Infow("Delete file successfully", logx.Field("Detail", target))
	return &types.SimpleMsg{Msg: errorx.DeleteSuccess}, nil
}

package file

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path"

	"github.com/suyuan32/simple-message/core/log"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePublicStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePublicStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePublicStatusLogic {
	return &ChangePublicStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePublicStatusLogic) ChangePublicStatus(req *types.ChangeStatusReq) (resp *types.SimpleMsg, err error) {
	var origin model.FileInfo
	result := l.svcCtx.DB.Where("id = ?", req.ID).First(&origin)
	if result.Error != nil {
		logx.Errorw(log.DatabaseError, logx.Field("detail", result.Error.Error()))
		return nil, errorx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if result.RowsAffected == 0 {
		logx.Errorw("file dose not find", logx.Field("fileId", req.ID))
		return nil, errorx.NewApiErrorWithoutMsg(http.StatusNotFound)
	}

	// only admin and owner can do it
	roleId := l.ctx.Value("roleId").(json.Number).String()
	userId := l.ctx.Value("userId").(string)
	if roleId != "1" && userId != origin.UserUUID {
		logx.Errorw(log.OperationNotAllow, logx.Field("roleId", roleId),
			logx.Field("userId", userId))
		return nil, errorx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
	}

	if req.Status {
		err = os.Rename(path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, origin.Path),
			path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, origin.Path))
		if err != nil {
			logx.Errorw("fail to change the path of file", logx.Field("detail", err.Error()))
			return nil, errorx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	} else {
		err = os.Rename(path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, origin.Path),
			path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, origin.Path))
		if err != nil {
			logx.Errorw("fail to change the path of file", logx.Field("detail", err.Error()))
			return nil, errorx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	}
	result = l.svcCtx.DB.Model(&model.FileInfo{}).Where("id = ?", req.ID).Update("status", req.Status)
	if result.Error != nil {
		logx.Errorw(log.DatabaseError, logx.Field("detail", result.Error.Error()))
		return nil, errorx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}

	if result.RowsAffected == 0 {
		logx.Errorw("update file status fail", logx.Field("detail", req))
		return &types.SimpleMsg{Msg: errorx.UpdateFailed}, nil
	}
	return &types.SimpleMsg{Msg: errorx.UpdateSuccess}, nil
}

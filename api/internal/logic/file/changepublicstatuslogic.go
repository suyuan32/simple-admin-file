package file

import (
	"context"
	"encoding/json"
	"github.com/suyuan32/simple-admin-file/api/internal/util/logmessage"
	"net/http"
	"os"
	"path"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
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
		logx.Errorw(logmessage.DatabaseError, logx.Field("Detail", result.Error.Error()))
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if result.RowsAffected == 0 {
		logx.Errorw("File dose not find", logx.Field("FileId", req.ID))
		return nil, httpx.NewApiErrorWithoutMsg(http.StatusNotFound)
	}

	// only admin and owner can do it
	roleId := l.ctx.Value("roleId").(json.Number).String()
	userId := l.ctx.Value("userId").(string)
	if roleId != "1" && userId != origin.UserUUID {
		logx.Errorw(logmessage.OperationNotAllow, logx.Field("RoleId", roleId),
			logx.Field("UserId", userId))
		return nil, httpx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
	}

	if req.Status {
		err = os.Rename(path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, origin.Path),
			path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, origin.Path))
		if err != nil {
			logx.Errorw("Fail to change the path of file", logx.Field("Detail", err.Error()))
			return nil, httpx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	} else {
		err = os.Rename(path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, origin.Path),
			path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, origin.Path))
		if err != nil {
			logx.Errorw("Fail to change the path of file", logx.Field("Detail", err.Error()))
			return nil, httpx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	}
	result = l.svcCtx.DB.Model(&model.FileInfo{}).Where("id = ?", req.ID).Update("status", req.Status)
	if result.Error != nil {
		logx.Errorw(logmessage.DatabaseError, logx.Field("Detail", result.Error.Error()))
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}

	if result.RowsAffected == 0 {
		logx.Errorw("Update file status fail", logx.Field("Detail", req))
		return &types.SimpleMsg{Msg: errorx.UpdateFailed}, nil
	}
	return &types.SimpleMsg{Msg: errorx.UpdateSuccess}, nil
}

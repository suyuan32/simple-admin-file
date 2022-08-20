package file

import (
	"context"
	"encoding/json"
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
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return nil, httpx.NewApiErrorWithoutMsg(http.StatusNotFound)
	}

	// only admin and owner can do it
	if l.ctx.Value("roleId").(json.Number).String() != "1" && l.ctx.Value("userId").(string) != origin.UserUUID {
		return nil, httpx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
	}

	if req.Status {
		err = os.Rename(path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, origin.Path),
			path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, origin.Path))
		if err != nil {
			return nil, httpx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	} else {
		err = os.Rename(path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, origin.Path),
			path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, origin.Path))
		if err != nil {
			return nil, httpx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	}
	result = l.svcCtx.DB.Model(&model.FileInfo{}).Where("id = ?", req.ID).Update("status", req.Status)
	if result.Error != nil {
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}

	if result.RowsAffected == 0 {
		return &types.SimpleMsg{Msg: errorx.UpdateFailed}, nil
	}
	return &types.SimpleMsg{Msg: errorx.UpdateSuccess}, nil
}

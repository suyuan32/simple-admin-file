package file

import (
	"context"
	"encoding/json"
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

func (l *DeleteFileLogic) DeleteFile(req *types.IdReq) (resp *types.SimpleMsg, err error) {
	var target model.FileInfo
	result := l.svcCtx.DB.Where("id = ?", req.ID).First(&target)
	if result.Error != nil {
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return nil, httpx.NewApiErrorWithoutMsg(http.StatusNotFound)
	}

	// judge the user is admin or owner
	// 只有管理员和拥有者能删除信息
	if l.ctx.Value("roleId").(json.Number).String() != "1" && l.ctx.Value("userId").(string) != target.UserUUID {
		return nil, httpx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
	}

	if target.Status {
		err = os.Remove(l.svcCtx.Config.UploadConf.PublicStorePath + target.Path)
		if err != nil {
			return nil, httpx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	} else {
		err = os.Remove(l.svcCtx.Config.UploadConf.PrivateStorePath + target.Path)
		if err != nil {
			return nil, httpx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	}

	result = l.svcCtx.DB.Delete(&target)
	if result.Error != nil {
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	return &types.SimpleMsg{Msg: errorx.DeleteSuccess}, nil
}

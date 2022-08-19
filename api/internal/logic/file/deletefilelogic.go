package file

import (
	"context"
	"github.com/suyuan32/simple-admin-core/api/common/errorx"
	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"os"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	if target.Status {
		err = os.Remove(l.svcCtx.Config.UploadConf.PublicStorePath + target.Path)
		if err != nil {
			return nil, httpx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	}

	result = l.svcCtx.DB.Delete(target)
	if result.Error != nil {
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	return &types.SimpleMsg{Msg: errorx.DeleteSuccess}, nil
}

package file

import (
	"context"
	"github.com/suyuan32/simple-admin-core/api/common/errorx"
	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileLogic {
	return &UpdateFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFileLogic) UpdateFile(req *types.UpdateFileReq) (resp *types.SimpleMsg, err error) {
	result := l.svcCtx.DB.Model(&model.FileInfo{}).Where("id = ?", req.ID).Update("name", req.Name)
	if result.Error != nil {
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return &types.SimpleMsg{Msg: errorx.UpdateSuccess}, nil
	}
	return
}

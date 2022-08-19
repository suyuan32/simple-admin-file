package file

import (
	"context"
	"github.com/suyuan32/simple-admin-core/api/common/errorx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

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
	result := l.svcCtx.DB.Where("id = ?", req.ID).Update("status", req.Status)
	if result.Error != nil {
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}

	if result.RowsAffected == 0 {
		return &types.SimpleMsg{Msg: errorx.UpdateFailed}, nil
	}
	return &types.SimpleMsg{Msg: errorx.UpdateSuccess}, nil
}

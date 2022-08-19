package file

import (
	"context"
	"github.com/suyuan32/simple-admin-core/api/common/errorx"
	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"net/http"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitDatabaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInitDatabaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitDatabaseLogic {
	return &InitDatabaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitDatabaseLogic) InitDatabase() (resp *types.SimpleMsg, err error) {
	var file model.FileInfo
	check := l.svcCtx.DB.First(&file)
	if check.RowsAffected != 0 {
		return &types.SimpleMsg{Msg: errorx.AlreadyInit}, nil
	}

	err = l.svcCtx.DB.AutoMigrate(&model.FileInfo{})
	if err != nil {
		return nil, errorx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	return &types.SimpleMsg{Msg: errorx.Success}, nil
}

package file

import (
	"context"
	"net/http"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
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
		logx.Errorw("Initialize database error", logx.Field("Detail", err.Error()))
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	logx.Infow("Initialize database successfully")
	return &types.SimpleMsg{Msg: errorx.Success}, nil
}

package file

import (
	"context"
	"net/http"

	"github.com/suyuan32/simple-message/core/log"
	"github.com/suyuan32/simple-models/model/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/errorx"
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
		logx.Errorw("initialize database error", logx.Field("detail", err.Error()))
		return nil, errorx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}

	err = l.insertApiData()
	if err != nil {
		logx.Errorw(log.DatabaseError, logx.Field("detail", err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}

	logx.Infow("initialize database successfully")
	return &types.SimpleMsg{Msg: errorx.Success}, nil
}

func (l *InitDatabaseLogic) insertApiData() error {
	apis := []core.Api{
		// user
		{
			Path:        "/upload",
			Description: "apiDesc.uploadFile",
			ApiGroup:    "file",
			Method:      "POST",
		},
		{
			Path:        "/file/list",
			Description: "apiDesc.fileList",
			ApiGroup:    "file",
			Method:      "POST",
		},
		{
			Path:        "/file",
			Description: "apiDesc.updateFileInfo",
			ApiGroup:    "file",
			Method:      "POST",
		},
		{
			Path:        "/file/status",
			Description: "apiDesc.setPublicStatus",
			ApiGroup:    "file",
			Method:      "POST",
		},
		{
			Path:        "/file",
			Description: "apiDesc.deleteFile",
			ApiGroup:    "file",
			Method:      "DELETE",
		},
		{
			Path:        "/file/download",
			Description: "apiDesc.downloadFile",
			ApiGroup:    "file",
			Method:      "GET",
		},
	}

	result := l.svcCtx.DB.CreateInBatches(apis, 100)
	if result.Error != nil {
		logx.Errorw(log.DatabaseError, logx.Field("detail", result.Error.Error()))
		return status.Error(codes.Internal, result.Error.Error())
	} else {
		return nil
	}
}

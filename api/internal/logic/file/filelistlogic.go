package file

import (
	"context"
	"github.com/suyuan32/simple-admin-core/api/common/errorx"
	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"time"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileListLogic {
	return &FileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileListLogic) FileList(req *types.FileListReq) (resp *types.FileListResp, err error) {
	db := l.svcCtx.DB.Model(&model.FileInfo{})

	if req.FileName != "" {
		db = db.Where("name like ?", req.FileName)
	}

	if req.FileType != "" {
		db = db.Where("file_type = ?", req.FileType)
	}

	if req.BeginDate != 0 && req.EndDate != 0 {
		db = db.Where("create_at between ? and ?", time.UnixMilli(req.BeginDate), time.UnixMilli(req.EndDate))
	}

	var fileInfos []model.FileInfo
	result := db.Limit(int(req.PageSize)).Offset(int((req.Page - 1) * req.PageSize)).
		Order("create_at desc").Find(&fileInfos)

	if result.Error != nil {
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}

	resp = &types.FileListResp{}
	resp.Total = uint64(result.RowsAffected)

	for _, v := range fileInfos {
		resp.Data = append(resp.Data, types.FileInfo{
			ID:       int64(v.ID),
			UUID:     v.UUID,
			UserUUID: v.UserUUID,
			Name:     v.Name,
			FileType: v.FileType,
			Size:     v.Size,
			Path:     v.Path,
		})
	}

	return resp, nil
}

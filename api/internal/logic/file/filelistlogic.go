package file

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
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
	// only admin can view the list
	if l.ctx.Value("roleId").(json.Number).String() != "1" {
		return nil, httpx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
	}

	db := l.svcCtx.DB.Model(&model.FileInfo{})

	if req.FileName != "" {
		db = db.Where("name like ?", req.FileName)
	}

	if req.FileType != "" {
		db = db.Where("file_type = ?", req.FileType)
	}

	if req.Period[0] != "" && req.Period[1] != "" {
		begin, err := time.Parse("2006-01-02 15:04:05", req.Period[0])
		if err != nil {
			return nil, httpx.NewApiErrorWithoutMsg(http.StatusBadRequest)
		}
		end, err := time.Parse("2006-01-02 15:04:05", req.Period[1])
		if err != nil {
			return nil, httpx.NewApiErrorWithoutMsg(http.StatusBadRequest)
		}

		db = db.Where("created_at between ? and ?", begin, end)
	}

	var fileInfos []model.FileInfo
	result := db.Limit(int(req.PageSize)).Offset(int((req.Page - 1) * req.PageSize)).
		Order("created_at desc").Find(&fileInfos)

	if result.Error != nil {
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}

	resp = &types.FileListResp{}
	resp.Total = uint64(result.RowsAffected)

	for _, v := range fileInfos {
		resp.Data = append(resp.Data, types.FileInfo{
			ID:        int64(v.ID),
			UUID:      v.UUID,
			UserUUID:  v.UserUUID,
			Name:      v.Name,
			FileType:  v.FileType,
			Size:      v.Size,
			Path:      v.Path,
			Status:    v.Status,
			CreatedAt: v.CreatedAt.UnixMilli(),
		})
	}

	return resp, nil
}

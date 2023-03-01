package file

import (
	"context"
	"net/http"
	"time"

	"github.com/suyuan32/simple-admin-core/pkg/enum"
	"github.com/suyuan32/simple-admin-core/pkg/i18n"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
	"github.com/suyuan32/simple-admin-file/pkg/ent/file"
	"github.com/suyuan32/simple-admin-file/pkg/ent/predicate"
	"github.com/suyuan32/simple-admin-file/pkg/utils/dberrorhandler"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lang   string
}

func NewFileListLogic(r *http.Request, svcCtx *svc.ServiceContext) *FileListLogic {
	return &FileListLogic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		svcCtx: svcCtx,
		lang:   r.Header.Get("Accept-Language"),
	}
}

func (l *FileListLogic) FileList(req *types.FileListReq) (resp *types.FileListResp, err error) {
	var predicates []predicate.File

	if req.FileType != 0 {
		predicates = append(predicates, file.FileTypeEQ(req.FileType))
	}

	if req.FileName != "" {
		predicates = append(predicates, file.NameContains(req.FileName))
	}

	if req.Period != nil {
		begin, err := time.Parse("2006-01-02 15:04:05", req.Period[0])
		if err != nil {
			return nil, errorx.NewCodeError(enum.InvalidArgument, i18n.Failed)
		}
		end, err := time.Parse("2006-01-02 15:04:05", req.Period[1])
		if err != nil {
			return nil, errorx.NewCodeError(enum.InvalidArgument, i18n.Failed)
		}
		predicates = append(predicates, file.CreatedAtGT(begin), file.CreatedAtLT(end))
	}

	files, err := l.svcCtx.DB.File.Query().Where(predicates...).Page(l.ctx, req.Page, req.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(err, req)
	}

	resp = &types.FileListResp{}
	resp.Msg = l.svcCtx.Trans.Trans(l.lang, i18n.Success)
	resp.Data.Total = files.PageDetails.Total

	for _, v := range files.List {
		resp.Data.Data = append(resp.Data.Data, types.FileInfo{
			BaseInfo: types.BaseInfo{
				Id:        v.ID,
				CreatedAt: v.CreatedAt.UnixMilli(),
				UpdatedAt: v.UpdatedAt.UnixMilli(),
			},
			UUID:     v.UUID,
			UserUUID: v.UserUUID,
			Name:     v.Name,
			FileType: v.FileType,
			Size:     v.Size,
			Path:     v.Path,
			Status:   v.Status,
		})
	}

	return resp, nil
}

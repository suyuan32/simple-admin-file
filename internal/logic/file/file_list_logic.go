package file

import (
	"context"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/suyuan32/simple-admin-file/ent"
	"github.com/suyuan32/simple-admin-file/ent/tag"
	"time"

	"github.com/suyuan32/simple-admin-common/enum/errorcode"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/suyuan32/simple-admin-file/internal/utils/dberrorhandler"

	"github.com/suyuan32/simple-admin-file/ent/file"
	"github.com/suyuan32/simple-admin-file/ent/predicate"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lang   string
}

func NewFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileListLogic {
	return &FileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileListLogic) FileList(req *types.FileListReq) (resp *types.FileListResp, err error) {
	var predicates []predicate.File

	if req.FileType != nil && *req.FileType != 0 {
		predicates = append(predicates, file.FileTypeEQ(*req.FileType))
	}

	if req.FileName != nil {
		predicates = append(predicates, file.NameContains(*req.FileName))
	}

	if req.TagIds != nil {
		predicates = append(predicates, file.HasTagsWith(tag.IDIn(req.TagIds...)))
	}

	if req.Status != nil {
		predicates = append(predicates, file.StatusEQ(*req.Status))
	}

	if req.Period != nil {
		begin, err := time.Parse("2006-01-02 15:04:05", req.Period[0])
		if err != nil {
			return nil, errorx.NewCodeError(errorcode.InvalidArgument, i18n.Failed)
		}
		end, err := time.Parse("2006-01-02 15:04:05", req.Period[1])
		if err != nil {
			return nil, errorx.NewCodeError(errorcode.InvalidArgument, i18n.Failed)
		}
		predicates = append(predicates, file.CreatedAtGT(begin), file.CreatedAtLT(end))
	}

	files, err := l.svcCtx.DB.File.Query().WithTags().Where(predicates...).Page(l.ctx, req.Page, req.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	resp = &types.FileListResp{}
	resp.Msg = l.svcCtx.Trans.Trans(l.ctx, i18n.Success)
	resp.Data.Total = files.PageDetails.Total

	for _, v := range files.List {
		resp.Data.Data = append(resp.Data.Data, types.FileInfo{
			BaseUUIDInfo: types.BaseUUIDInfo{
				Id:        pointy.GetPointer(v.ID.String()),
				CreatedAt: pointy.GetPointer(v.CreatedAt.UnixMilli()),
				UpdatedAt: pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			},
			UserUUID:   &v.UserUUID,
			Name:       &v.Name,
			FileType:   &v.FileType,
			Size:       &v.Size,
			Path:       &v.Path,
			Status:     &v.Status,
			TagIds:     l.getTagIds(v.Edges.Tags),
			PublicPath: pointy.GetPointer(l.svcCtx.Config.UploadConf.ServerURL + v.Path),
		})
	}

	return resp, nil
}

func (l *FileListLogic) getTagIds(tags []*ent.Tag) (result []uint64) {
	for _, v := range tags {
		result = append(result, v.ID)
	}
	return result
}

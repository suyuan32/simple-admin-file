package file

import (
	"context"
	"time"

	"github.com/suyuan32/simple-admin-common/utils/pointy"

	"github.com/suyuan32/simple-admin-file/ent"
	"github.com/suyuan32/simple-admin-file/ent/filetag"

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
	currentUserId := l.ctx.Value("userId").(string)

	// Privacy filter: show public files (status=1) OR private files (status=2) owned by current user
	// 隐私过滤：显示公开文件 (status=1) 或 当前用户自己的私有文件 (status=2)
	privacyPredicate := file.Or(
		file.StatusEQ(1), // public files
		file.And(
			file.StatusEQ(2),          // private files
			file.UserIDEQ(currentUserId), // owned by current user
		),
	)
	predicates = append(predicates, privacyPredicate)

	if req.FileType != nil && *req.FileType != 0 {
		predicates = append(predicates, file.FileTypeEQ(*req.FileType))
	}

	if req.FileName != nil {
		predicates = append(predicates, file.NameContains(*req.FileName))
	}

	if req.FileTagIds != nil {
		predicates = append(predicates, file.HasTagsWith(filetag.IDIn(req.FileTagIds...)))
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
		// Only return public path for public files (status=1)
		// Private files (status=2) should be downloaded via API
		var publicPath *string
		if v.Status == 1 {
			publicPath = pointy.GetPointer(l.svcCtx.Config.UploadConf.ServerURL + v.Path)
		} else {
			publicPath = pointy.GetPointer("")
		}

		resp.Data.Data = append(resp.Data.Data, types.FileInfo{
			BaseUUIDInfo: types.BaseUUIDInfo{
				Id:        pointy.GetPointer(v.ID.String()),
				CreatedAt: pointy.GetPointer(v.CreatedAt.UnixMilli()),
				UpdatedAt: pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			},
			UserUUID:   &v.UserID,
			Name:       &v.Name,
			FileType:   &v.FileType,
			Size:       &v.Size,
			Path:       &v.Path,
			Status:     &v.Status,
			FileTagIds: l.getFileTagIds(v.Edges.Tags),
			PublicPath: publicPath,
		})
	}

	return resp, nil
}

func (l *FileListLogic) getFileTagIds(tags []*ent.FileTag) (result []uint64) {
	for _, v := range tags {
		result = append(result, v.ID)
	}
	return result
}

package file

import (
	"context"
	"net/url"

	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/suyuan32/simple-admin-file/ent/file"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
	"github.com/suyuan32/simple-admin-file/internal/utils/dberrorhandler"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileByUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFileByUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileByUrlLogic {
	return &DeleteFileByUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFileByUrlLogic) DeleteFileByUrl(req *types.FileDeleteReq) (resp *types.BaseMsgResp, err error) {
	urlData, err := url.Parse(req.Url)
	if err != nil {
		return nil, errorx.NewApiBadRequestError("failed to parse url")
	}

	data, err := l.svcCtx.DB.File.Query().Where(file.PathEQ(urlData.Path)).Only(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	logic := NewDeleteFileLogic(l.ctx, l.svcCtx)
	return logic.DeleteFile(&types.UUIDsReq{Ids: []string{data.ID.String()}})
}

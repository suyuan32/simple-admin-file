package cloudfile

import (
	"context"

	"github.com/suyuan32/simple-admin-file/ent/cloudfile"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
	"github.com/suyuan32/simple-admin-file/internal/utils/dberrorhandler"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCloudFileByUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCloudFileByUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCloudFileByUrlLogic {
	return &DeleteCloudFileByUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCloudFileByUrlLogic) DeleteCloudFileByUrl(req *types.CloudFileDeleteReq) (resp *types.BaseMsgResp, err error) {
	data, err := l.svcCtx.DB.CloudFile.Query().Where(cloudfile.URLEQ(req.Url)).Only(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	logic := NewDeleteCloudFileLogic(l.ctx, l.svcCtx)
	return logic.DeleteCloudFile(&types.UUIDsReq{Ids: []string{data.ID.String()}})
}

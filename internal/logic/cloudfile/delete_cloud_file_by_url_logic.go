package cloudfile

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
	"github.com/suyuan32/simple-admin-file/internal/utils/filex"
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
	fileId, err := filex.ConvertUrlStringToFileUUID(req.Url)
	if err != nil {
		return nil, err
	}

	logic := NewDeleteCloudFileLogic(l.ctx, l.svcCtx)
	return logic.DeleteCloudFile(&types.UUIDsReq{Ids: []string{fileId}})
}

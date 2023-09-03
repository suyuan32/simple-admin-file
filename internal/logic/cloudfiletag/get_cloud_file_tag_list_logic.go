package cloudfiletag

import (
	"context"

	"github.com/suyuan32/simple-admin-file/ent/cloudfiletag"
	"github.com/suyuan32/simple-admin-file/ent/predicate"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
	"github.com/suyuan32/simple-admin-file/internal/utils/dberrorhandler"

	"github.com/suyuan32/simple-admin-common/i18n"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetCloudFileTagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCloudFileTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudFileTagListLogic {
	return &GetCloudFileTagListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCloudFileTagListLogic) GetCloudFileTagList(req *types.CloudFileTagListReq) (*types.CloudFileTagListResp, error) {
	var predicates []predicate.CloudFileTag
	if req.Name != nil {
		predicates = append(predicates, cloudfiletag.NameContains(*req.Name))
	}
	if req.Remark != nil {
		predicates = append(predicates, cloudfiletag.RemarkContains(*req.Remark))
	}
	data, err := l.svcCtx.DB.CloudFileTag.Query().Where(predicates...).Page(l.ctx, req.Page, req.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	resp := &types.CloudFileTagListResp{}
	resp.Msg = l.svcCtx.Trans.Trans(l.ctx, i18n.Success)
	resp.Data.Total = data.PageDetails.Total

	for _, v := range data.List {
		resp.Data.Data = append(resp.Data.Data,
			types.CloudFileTagInfo{
				BaseIDInfo: types.BaseIDInfo{
					Id:        &v.ID,
					CreatedAt: pointy.GetPointer(v.CreatedAt.Unix()),
					UpdatedAt: pointy.GetPointer(v.UpdatedAt.Unix()),
				},
				Status: &v.Status,
				Name:   &v.Name,
				Remark: &v.Remark,
			})
	}

	return resp, nil
}

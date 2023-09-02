package filetag

import (
	"context"

	"github.com/suyuan32/simple-admin-file/ent/filetag"
	"github.com/suyuan32/simple-admin-file/ent/predicate"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
	"github.com/suyuan32/simple-admin-file/internal/utils/dberrorhandler"

	"github.com/suyuan32/simple-admin-common/i18n"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileTagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileTagListLogic {
	return &GetFileTagListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFileTagListLogic) GetFileTagList(req *types.FileTagListReq) (*types.FileTagListResp, error) {
	var predicates []predicate.FileTag
	if req.Name != nil {
		predicates = append(predicates, filetag.NameContains(*req.Name))
	}
	if req.Remark != nil {
		predicates = append(predicates, filetag.RemarkContains(*req.Remark))
	}
	data, err := l.svcCtx.DB.FileTag.Query().Where(predicates...).Page(l.ctx, req.Page, req.PageSize)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	resp := &types.FileTagListResp{}
	resp.Msg = l.svcCtx.Trans.Trans(l.ctx, i18n.Success)
	resp.Data.Total = data.PageDetails.Total

	for _, v := range data.List {
		resp.Data.Data = append(resp.Data.Data,
			types.FileTagInfo{
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

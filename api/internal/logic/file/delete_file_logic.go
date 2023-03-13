package file

import (
	"context"
	"os"

	"github.com/suyuan32/simple-admin-common/i18n"

	"github.com/suyuan32/simple-admin-file/api/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-file/api/internal/utils/entx"

	"github.com/suyuan32/simple-admin-file/api/ent"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lang   string
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFileLogic) DeleteFile(req *types.IDReq) (resp *types.BaseMsgResp, err error) {
	err = entx.WithTx(l.ctx, l.svcCtx.DB, func(tx *ent.Tx) error {
		file, err := tx.File.Get(l.ctx, req.Id)

		if err != nil {
			return err
		}

		err = tx.File.DeleteOneID(req.Id).Exec(l.ctx)

		if err != nil {
			return err
		}

		if file.Status == 1 {
			err = os.RemoveAll(l.svcCtx.Config.UploadConf.PublicStorePath + file.Path)
			if err != nil {
				return err
			}
		} else {
			err = os.RemoveAll(l.svcCtx.Config.UploadConf.PrivateStorePath + file.Path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	return &types.BaseMsgResp{
		Code: 0,
		Msg:  l.svcCtx.Trans.Trans(l.ctx, i18n.DeleteSuccess),
	}, nil
}

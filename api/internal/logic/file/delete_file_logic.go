package file

import (
	"context"
	"net/http"
	"os"

	"github.com/suyuan32/simple-admin-core/pkg/i18n"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
	"github.com/suyuan32/simple-admin-file/pkg/ent"
	"github.com/suyuan32/simple-admin-file/pkg/utils"
	"github.com/suyuan32/simple-admin-file/pkg/utils/dberrorhandler"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lang   string
}

func NewDeleteFileLogic(r *http.Request, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		svcCtx: svcCtx,
		lang:   r.Header.Get("Accept-Language"),
	}
}

func (l *DeleteFileLogic) DeleteFile(req *types.IDReq) (resp *types.BaseMsgResp, err error) {
	err = utils.WithTx(l.ctx, l.svcCtx.DB, func(tx *ent.Tx) error {
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
		return nil, dberrorhandler.DefaultEntError(err, req)
	}

	return &types.BaseMsgResp{
		Code: 0,
		Msg:  l.svcCtx.Trans.Trans(l.lang, i18n.DeleteSuccess),
	}, nil
}

package file

import (
	"context"
	"net/http"
	"os"
	"path"

	"github.com/suyuan32/simple-admin-core/pkg/i18n"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
	"github.com/suyuan32/simple-admin-file/pkg/ent"
	"github.com/suyuan32/simple-admin-file/pkg/utils"
	"github.com/suyuan32/simple-admin-file/pkg/utils/dberrorhandler"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePublicStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lang   string
}

func NewChangePublicStatusLogic(r *http.Request, svcCtx *svc.ServiceContext) *ChangePublicStatusLogic {
	return &ChangePublicStatusLogic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		svcCtx: svcCtx,
		lang:   r.Header.Get("Accept-Language"),
	}
}

func (l *ChangePublicStatusLogic) ChangePublicStatus(req *types.StatusCodeReq) (resp *types.BaseMsgResp, err error) {
	err = utils.WithTx(l.ctx, l.svcCtx.DB, func(tx *ent.Tx) error {
		file, err := tx.File.UpdateOneID(req.Id).SetStatus(uint8(req.Status)).Save(l.ctx)

		if err != nil {
			return err
		}

		if req.Status == 1 {
			err = os.Rename(path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, file.Path),
				path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, file.Path))
			if err != nil {
				logx.Errorw("fail to change the path of file", logx.Field("detail", err.Error()))
				return err
			}
		} else {
			err = os.Rename(path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, file.Path),
				path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, file.Path))
			if err != nil {
				logx.Errorw("fail to change the path of file", logx.Field("detail", err.Error()))
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
		Msg:  l.svcCtx.Trans.Trans(l.lang, i18n.Success),
	}, nil
}

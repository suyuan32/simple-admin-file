package file

import (
	"context"
	"net/http"
	"os"

	"github.com/suyuan32/simple-admin-core/pkg/i18n"
	"github.com/suyuan32/simple-admin-core/pkg/msg/logmsg"
	"github.com/suyuan32/simple-admin-core/pkg/statuserr"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
	"github.com/suyuan32/simple-admin-file/pkg/ent"
	"github.com/suyuan32/simple-admin-file/pkg/utils"

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
			switch {
			case ent.IsNotFound(err):
				logx.Errorw(err.Error(), logx.Field("detail", req))
				return err
			default:
				logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
				return err
			}
		}

		err = tx.File.DeleteOneID(req.Id).Exec(l.ctx)

		if err != nil {
			switch {
			case ent.IsNotFound(err):
				logx.Errorw(err.Error(), logx.Field("detail", req))
				return err
			default:
				logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
				return err
			}
		}

		if file.Status == 1 {
			err = os.RemoveAll(l.svcCtx.Config.UploadConf.PublicStorePath + file.Path)
			if err != nil {
				logx.Errorw("fail to remove the file", logx.Field("path",
					l.svcCtx.Config.UploadConf.PublicStorePath+file.Path), logx.Field("detail", err.Error()))
				return err
			}
		} else {
			err = os.RemoveAll(l.svcCtx.Config.UploadConf.PrivateStorePath + file.Path)
			if err != nil {
				logx.Errorw("fail to remove the file", logx.Field("path",
					l.svcCtx.Config.UploadConf.PrivateStorePath+file.Path), logx.Field("detail", err.Error()))
				return err
			}
		}

		return nil
	})

	if err != nil {
		logx.Errorf("failed to delete file, error : %s", err.Error())
		return nil, statuserr.NewInternalError(i18n.DatabaseError)
	}

	return &types.BaseMsgResp{
		Code: 0,
		Msg:  l.svcCtx.Trans.Trans(l.lang, i18n.DeleteSuccess),
	}, nil
}

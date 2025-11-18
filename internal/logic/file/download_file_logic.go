package file

import (
	"context"
	"path"

	"github.com/suyuan32/simple-admin-common/enum/errorcode"
	"github.com/suyuan32/simple-admin-common/utils/uuidx"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/suyuan32/simple-admin-file/internal/utils/dberrorhandler"

	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadFileLogic {
	return &DownloadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadFileLogic) DownloadFile(req *types.UUIDPathReq) (filePath string, err error) {
	file, err := l.svcCtx.DB.File.Get(l.ctx, uuidx.ParseUUIDString(req.Id))

	if err != nil {
		return "", dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	currentUserId := l.ctx.Value("userId").(string)

	if file.Status == 1 {
		// Public file: anyone can download
		logx.Infow("public download", logx.Field("fileName", file.Name),
			logx.Field("userId", currentUserId),
			logx.Field("filePath", file.Path))
		return path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, file.Path), nil
	} else {
		// Private file: only owner can download
		if file.UserID != currentUserId {
			logx.Errorw("unauthorized download attempt", logx.Field("fileName", file.Name),
				logx.Field("fileOwner", file.UserID),
				logx.Field("requestUser", currentUserId))
			return "", errorx.NewCodeError(errorcode.PermissionDenied, "file.permissionDenied")
		}

		logx.Infow("private download", logx.Field("fileName", file.Name),
			logx.Field("userId", currentUserId),
			logx.Field("filePath", file.Path))
		return path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath, file.Path), nil
	}
}

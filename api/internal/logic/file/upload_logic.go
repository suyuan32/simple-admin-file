package file

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/suyuan32/knife/core/date/format"
	filex2 "github.com/suyuan32/knife/core/io/filex"
	"github.com/suyuan32/simple-admin-common/enum/errorcode"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-common/utils/uuidx"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
	"github.com/suyuan32/simple-admin-file/api/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-file/api/internal/utils/filex"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadLogic(r *http.Request, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *UploadLogic) Upload() (resp *types.UploadResp, err error) {
	err = l.r.ParseMultipartForm(l.svcCtx.Config.UploadConf.MaxVideoSize)
	if err != nil {
		logx.Error("fail to parse the multipart form")
		return nil, errorx.NewCodeError(errorcode.InvalidArgument,
			l.svcCtx.Trans.Trans(l.ctx, "file.parseFormFailed"))
	}

	file, handler, err := l.r.FormFile("file")
	if err != nil {
		logx.Error("the value of file cannot be found")
		return nil, errorx.NewCodeError(errorcode.InvalidArgument,
			l.svcCtx.Trans.Trans(l.ctx, "file.parseFormFailed"))
	}
	defer file.Close()

	// judge if the suffix is legal
	// 校验后缀是否合法
	dotIndex := strings.LastIndex(handler.Filename, ".")
	// if there is no suffix, reject it
	// 拒绝无后缀文件
	if dotIndex == -1 {
		logx.Errorw("reject the file which does not have suffix")
		return nil, errorx.NewCodeError(errorcode.InvalidArgument,
			l.svcCtx.Trans.Trans(l.ctx, "file.wrongTypeError"))
	}

	fileName, fileSuffix := handler.Filename[:dotIndex], handler.Filename[dotIndex+1:]
	fileUUID := uuidx.NewUUID().String()
	storeFileName := fileUUID + "." + fileSuffix
	timeString := time.Now().Format(format.DashYearToDay)
	userId := l.ctx.Value("userId").(string)

	// judge if the file size is over max size
	// 判断文件大小是否超过设定值
	fileType := strings.Split(handler.Header.Get("Content-Type"), "/")[0]
	if fileType != "image" && fileType != "video" && fileType != "audio" {
		fileType = "other"
	}
	err = filex.CheckOverSize(l.ctx, l.svcCtx, fileType, handler.Size)
	if err != nil {
		logx.Errorw("the file is over size", logx.Field("type", fileType),
			logx.Field("userId", userId), logx.Field("size", handler.Size),
			logx.Field("fileName", handler.Filename))
		return nil, err
	}

	// generate path
	// 生成路径

	//judge if the file directory exists, if not create it. Both private and public need
	//to be created in order to move files when status changed
	//判断文件夹是否已创建, 同时创建好私人和公开文件夹防止文件状态改变时无法移动

	publicStoreDir := path.Join(l.svcCtx.Config.UploadConf.PublicStorePath,
		l.svcCtx.Config.Name, fileType, timeString)
	privateStoreDir := path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath,
		l.svcCtx.Config.Name, fileType, timeString)

	err = filex2.MkdirIfNotExist(publicStoreDir, filex2.SuperPerm)
	if err != nil {
		logx.Errorw("failed to create directory for storing public files", logx.Field("path", publicStoreDir))
		return nil, errorx.NewCodeError(errorcode.Internal,
			l.svcCtx.Trans.Trans(l.ctx, i18n.Failed))
	}

	err = filex2.MkdirIfNotExist(privateStoreDir, filex2.SuperPerm)
	if err != nil {
		logx.Errorw("failed to create directory for storing private files", logx.Field("path", privateStoreDir))
		return nil, errorx.NewCodeError(errorcode.Internal,
			l.svcCtx.Trans.Trans(l.ctx, i18n.Failed))
	}

	// default is public
	// 默认是公开的
	targetFile, err := os.Create(path.Join(publicStoreDir, storeFileName))
	_, err = io.Copy(targetFile, file)
	if err != nil {
		logx.Errorw("fail to create file", logx.Field("path", path.Join(publicStoreDir, storeFileName)))
		return nil, errorx.NewCodeError(errorcode.Internal,
			l.svcCtx.Trans.Trans(l.ctx, i18n.Failed))
	}

	// store in database
	// 提交数据库
	relativePath := fmt.Sprintf("/%s/%s/%s/%s", l.svcCtx.Config.Name,
		fileType, timeString, storeFileName)

	err = l.svcCtx.DB.File.Create().
		SetUUID(fileUUID).
		SetName(fileName).
		SetFileType(filex.ConvertFileTypeToUint8(fileType)).
		SetPath(relativePath).
		SetUserUUID(userId).
		SetMd5(l.r.MultipartForm.Value["md5"][0]).
		SetStatus(1).
		SetSize(uint64(handler.Size)).
		Exec(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, "upload failed")
	}

	return &types.UploadResp{
		BaseDataInfo: types.BaseDataInfo{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.Success)},
		Data:         types.UploadInfo{Name: handler.Filename, Url: relativePath},
	}, nil
}

package file

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/suyuan32/simple-message/core/log"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
	"github.com/suyuan32/simple-admin-file/api/internal/util/msg"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UploadLogic struct {
	logx.Logger
	r      *http.Request
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(r *http.Request, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		r:      r,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload() (resp *types.UploadResp, err error) {
	err = l.r.ParseMultipartForm(l.svcCtx.Config.UploadConf.MaxVideoSize)
	if err != nil {
		logx.Error("fail to parse the multipart form")
		return nil, errorx.NewApiError(http.StatusBadRequest, "sys.api.apiRequestFailed")
	}
	file, handler, err := l.r.FormFile("file")
	if err != nil {
		logx.Error("the value of file cannot be found")
		return nil, errorx.NewApiError(http.StatusBadRequest, "sys.api.apiRequestFailed")
	}
	defer file.Close()

	// judge if the suffix is legal
	// 校验后缀是否合法
	nameData := strings.Split(handler.Filename, ".")
	// if there is no suffix, reject it
	// 拒绝无后缀文件
	if len(nameData) < 2 {
		logx.Errorw("reject the file which does not have suffix")
		return nil, errorx.NewApiError(http.StatusBadRequest, msg.WrongTypeError)
	}

	fileName, fileSuffix := nameData[0], nameData[1]
	fileUUID := uuid.NewString()
	storeFileName := fileUUID + "." + fileSuffix
	newTime := time.Now()
	timeString := fmt.Sprintf("%d-%d-%d", newTime.Year(), newTime.Month(), newTime.Day())
	userId := l.ctx.Value("userId").(string)

	// judge if the file size is over max size
	// 判断文件大小是否超过设定值
	fileType := strings.Split(handler.Header.Get("Content-Type"), "/")[0]
	if fileType == "image" && handler.Size > l.svcCtx.Config.UploadConf.MaxImageSize {
		logx.Errorw("the file is over size", logx.Field("type", "image"),
			logx.Field("userId", userId), logx.Field("size", handler.Size),
			logx.Field("fileName", handler.Filename))
		return nil, errorx.NewApiError(http.StatusBadRequest, msg.OverSizeError)
	} else if fileType == "video" && handler.Size > l.svcCtx.Config.UploadConf.MaxVideoSize {
		logx.Errorw("the file is over size", logx.Field("type", "video"),
			logx.Field("userId", userId), logx.Field("size", handler.Size),
			logx.Field("fileName", handler.Filename))
		return nil, errorx.NewApiError(http.StatusBadRequest, msg.OverSizeError)
	} else if fileType == "audio" && handler.Size > l.svcCtx.Config.UploadConf.MaxAudioSize {
		logx.Errorw("the file is over size", logx.Field("type", "audio"),
			logx.Field("userId", userId), logx.Field("size", handler.Size),
			logx.Field("fileName", handler.Filename))
		return nil, errorx.NewApiError(http.StatusBadRequest, msg.OverSizeError)
	} else if fileType != "image" && fileType != "video" && fileType != "audio" &&
		handler.Size > l.svcCtx.Config.UploadConf.MaxOtherSize {
		logx.Errorw("the file is over size", logx.Field("type", "other"),
			logx.Field("userId", userId), logx.Field("size", handler.Size),
			logx.Field("fileName", handler.Filename))
		return nil, errorx.NewApiError(http.StatusBadRequest, msg.OverSizeError)
	}
	if fileType != "image" && fileType != "video" && fileType != "audio" {
		fileType = "other"
	}

	// generate path
	// 生成路径

	//judge if the file directory exists, if not create it. Both private and public need
	//to be created in order to move files when status changed
	//判断文件夹是否已创建, 同时创建好私人和公开文件夹防止文件状态改变时无法移动

	_, err = os.Stat(path.Join(l.svcCtx.Config.UploadConf.PublicStorePath,
		l.svcCtx.Config.Name, fileType, timeString))
	if os.IsNotExist(err) {
		mask := syscall.Umask(0)
		defer syscall.Umask(mask)

		err = os.MkdirAll(path.Join(l.svcCtx.Config.UploadConf.PublicStorePath,
			l.svcCtx.Config.Name, fileType, timeString), 0777)
		if err != nil {
			logx.Errorw("fail to make directory", logx.Field("path", path.Join(l.svcCtx.Config.UploadConf.PublicStorePath,
				l.svcCtx.Config.Name, fileType, timeString)))
			return nil, errorx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	}

	_, err = os.Stat(path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath,
		l.svcCtx.Config.Name, fileType, timeString))
	if os.IsNotExist(err) {
		mask2 := syscall.Umask(0)
		defer syscall.Umask(mask2)

		err = os.MkdirAll(path.Join(l.svcCtx.Config.UploadConf.PrivateStorePath,
			l.svcCtx.Config.Name, fileType, timeString), 0777)
		if err != nil {
			logx.Errorw("fail to create directory", logx.Field("Path", path.Join(l.svcCtx.Config.UploadConf.PublicStorePath,
				l.svcCtx.Config.Name, fileType, timeString)))
			return nil, errorx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
		}
	}

	// default is public
	// 默认是公开的
	targetFile, err := os.Create(path.Join(l.svcCtx.Config.UploadConf.PublicStorePath, l.svcCtx.Config.Name,
		fileType, timeString, storeFileName))
	if err != nil {
		logx.Errorw("fail to create directory", logx.Field("path", path.Join(l.svcCtx.Config.UploadConf.PublicStorePath,
			l.svcCtx.Config.Name, fileType, timeString, storeFileName)))
		return nil, errorx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
	}
	_, err = io.Copy(targetFile, file)
	if err != nil {
		logx.Errorw("fail to create file", logx.Field("path", path.Join(l.svcCtx.Config.UploadConf.PublicStorePath,
			l.svcCtx.Config.Name, fileType, timeString, storeFileName)))
		return nil, errorx.NewApiErrorWithoutMsg(http.StatusInternalServerError)
	}

	// store in database
	// 提交数据库
	relativePath := fmt.Sprintf("/%s/%s/%s/%s", l.svcCtx.Config.Name,
		fileType, timeString, storeFileName)
	var fileInfo model.FileInfo
	fileInfo = model.FileInfo{
		Model:    gorm.Model{},
		UUID:     fileUUID,
		Name:     fileName,
		FileType: fileType,
		Size:     handler.Size,
		Path:     relativePath,
		UserUUID: userId,
		Md5:      l.r.MultipartForm.Value["md5"][0],
		Status:   true,
	}
	result := l.svcCtx.DB.Create(&fileInfo)

	if result.Error != nil {
		logx.Errorw(log.DatabaseError, logx.Field("detail", result.Error.Error()))
		return nil, errorx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}

	logx.Infow("create file successfully", logx.Field("detail", fileInfo))
	return &types.UploadResp{
		Msg:  "ok",
		Name: handler.Filename,
		Path: relativePath,
	}, nil
}

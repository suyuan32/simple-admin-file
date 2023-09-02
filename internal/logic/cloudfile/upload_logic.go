package cloudfile

import (
	"context"
	"errors"
	"fmt"
	"github.com/suyuan32/simple-admin-common/enum/errorcode"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/suyuan32/simple-admin-common/utils/uuidx"
	"github.com/suyuan32/simple-admin-file/internal/enum/cloudprovider"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
	"github.com/suyuan32/simple-admin-file/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-file/internal/utils/filex"
	"github.com/zeromicro/go-zero/core/errorx"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *UploadLogic) Upload() (resp *types.CloudFileInfoResp, err error) {
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
	fileUUID := uuidx.NewUUID()
	storeFileName := fileUUID.String() + "." + fileSuffix
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

	var provider string
	if l.r.MultipartForm.Value["provider"] != nil {
		provider = l.r.MultipartForm.Value["provider"][0]
	} else {
		return nil, errors.New("provider is not set")
	}

	url, err := l.UploadToProvider(file, storeFileName, provider)
	if err != nil {
		return nil, err
	}

	// store to database
	data, err := l.svcCtx.DB.CloudFile.Create().
		SetName(fileName).
		SetFileType(filex.ConvertFileTypeToUint8(fileType)).
		SetStorageProvidersID(l.svcCtx.CloudUploader.ProviderIdData[provider]).
		SetURL(url).
		SetSize(uint64(handler.Size)).
		SetUserID(userId).
		Save(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, nil)
	}

	return &types.CloudFileInfoResp{
		BaseDataInfo: types.BaseDataInfo{
			Code: 0,
			Msg:  i18n.Success,
			Data: "",
		},
		Data: types.CloudFileInfo{
			BaseUUIDInfo: types.BaseUUIDInfo{
				Id:        pointy.GetPointer(data.ID.String()),
				CreatedAt: pointy.GetPointer(data.CreatedAt.Unix()),
			},
			State:    pointy.GetPointer(data.State),
			Name:     pointy.GetPointer(data.Name),
			Url:      pointy.GetPointer(data.URL),
			Size:     pointy.GetPointer(data.Size),
			FileType: pointy.GetPointer(data.FileType),
			UserId:   pointy.GetPointer(data.UserID),
		},
	}, nil
}

func (l *UploadLogic) UploadToProvider(file multipart.File, fileName, provider string) (url string, err error) {
	if strings.Contains(provider, cloudprovider.Tencent) {
		if client, ok := l.svcCtx.CloudUploader.TencentCOS[provider]; ok {
			resp, err := client.Object.Put(l.ctx, fileName, file, nil)
			if err != nil {
				return url, err
			}

			fmt.Println(resp)
		}
	}

	return url, nil
}

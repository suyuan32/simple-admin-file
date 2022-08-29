package file

import (
	"context"
	"encoding/json"
	"github.com/suyuan32/simple-admin-file/api/internal/util/logmessage"
	"net/http"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type UpdateFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileLogic {
	return &UpdateFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFileLogic) UpdateFile(req *types.UpdateFileReq) (resp *types.SimpleMsg, err error) {
	var target model.FileInfo
	check := l.svcCtx.DB.Where("id = ?", req.ID).First(&target)
	if check.Error != nil {
		logx.Errorw(logmessage.DatabaseError, logx.Field("Detail", check.Error.Error()))
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if check.RowsAffected == 0 {
		logx.Errorw("File does not found", logx.Field("FileId", req.ID))
		return nil, httpx.NewApiErrorWithoutMsg(http.StatusNotFound)
	}

	// only admin and owner can do it
	roleId := l.ctx.Value("roleId").(json.Number).String()
	userId := l.ctx.Value("userId").(string)
	if roleId != "1" && userId != target.UserUUID {
		logx.Errorw(logmessage.OperationNotAllow, logx.Field("RoleId", roleId),
			logx.Field("UserId", userId))
		return nil, httpx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
	}

	// update data
	result := l.svcCtx.DB.Model(&model.FileInfo{}).Where("id = ?", req.ID).Update("name", req.Name)
	if result.Error != nil {
		logx.Errorw(logmessage.DatabaseError, logx.Field("Detail", result.Error.Error()))
		return nil, httpx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if result.RowsAffected == 0 {
		logx.Errorw("Fail to update the file", logx.Field("Detail", req))
		return &types.SimpleMsg{Msg: errorx.UpdateFailed}, nil
	}
	return &types.SimpleMsg{Msg: errorx.UpdateSuccess}, nil
}

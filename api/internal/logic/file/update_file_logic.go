package file

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/suyuan32/simple-message/core/log"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/simple-admin-file/api/internal/model"
	"github.com/suyuan32/simple-admin-file/api/internal/svc"
	"github.com/suyuan32/simple-admin-file/api/internal/types"
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
		logx.Errorw(log.DatabaseError, logx.Field("detail", check.Error.Error()))
		return nil, errorx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if check.RowsAffected == 0 {
		logx.Errorw("file does not found", logx.Field("FileId", req.ID))
		return nil, errorx.NewApiErrorWithoutMsg(http.StatusNotFound)
	}

	// only admin and owner can do it
	roleId := l.ctx.Value("roleId").(json.Number).String()
	userId := l.ctx.Value("userId").(string)
	if roleId != "1" && userId != target.UserUUID {
		logx.Errorw(log.OperationNotAllow, logx.Field("RoleId", roleId),
			logx.Field("userId", userId))
		return nil, errorx.NewApiErrorWithoutMsg(http.StatusUnauthorized)
	}

	// update data
	result := l.svcCtx.DB.Model(&model.FileInfo{}).Where("id = ?", req.ID).Update("name", req.Name)
	if result.Error != nil {
		logx.Errorw(log.DatabaseError, logx.Field("detail", result.Error.Error()))
		return nil, errorx.NewApiError(http.StatusInternalServerError, errorx.DatabaseError)
	}
	if result.RowsAffected == 0 {
		logx.Errorw("fail to update the file", logx.Field("detail", req))
		return &types.SimpleMsg{Msg: errorx.UpdateFailed}, nil
	}
	return &types.SimpleMsg{Msg: errorx.UpdateSuccess}, nil
}

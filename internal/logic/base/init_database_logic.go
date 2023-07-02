package base

import (
	"context"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/suyuan32/simple-admin-common/enum/errorcode"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-common/msg/logmsg"
	"github.com/suyuan32/simple-admin-core/rpc/types/core"
	"github.com/suyuan32/simple-admin-file/internal/svc"
	"github.com/suyuan32/simple-admin-file/internal/types"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitDatabaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lang   string
}

func NewInitDatabaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitDatabaseLogic {
	return &InitDatabaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitDatabaseLogic) InitDatabase() (resp *types.BaseMsgResp, err error) {
	if l.svcCtx.Config.CoreRpc.Enabled {
		err = l.initApi()
		if err != nil {
			if status.Code(err) == codes.InvalidArgument {
				return nil, errorx.NewCodeError(errorcode.InvalidArgument,
					l.svcCtx.Trans.Trans(l.ctx, "init.alreadyInit"))
			}
			return nil, err
		}
	}

	if err := l.svcCtx.DB.Schema.Create(l.ctx, schema.WithForeignKeys(false)); err != nil {
		logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
		return nil, errorx.NewCodeError(errorcode.Internal, err.Error())
	}

	return &types.BaseMsgResp{
		Code: 0,
		Msg:  l.svcCtx.Trans.Trans(l.ctx, i18n.Success),
	}, nil
}

func (l *InitDatabaseLogic) initApi() error {
	// create API in core service
	_, err := l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/upload"),
		Description: pointy.GetPointer("apiDesc.uploadFile"),
		ApiGroup:    pointy.GetPointer("file"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/file/list"),
		Description: pointy.GetPointer("apiDesc.fileList"),
		ApiGroup:    pointy.GetPointer("file"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/file/update"),
		Description: pointy.GetPointer("apiDesc.updateFileInfo"),
		ApiGroup:    pointy.GetPointer("file"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/file/status"),
		Description: pointy.GetPointer("apiDesc.setPublicStatus"),
		ApiGroup:    pointy.GetPointer("file"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/file/delete"),
		Description: pointy.GetPointer("apiDesc.deleteFile"),
		ApiGroup:    pointy.GetPointer("file"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/file/download/:id"),
		Description: pointy.GetPointer("apiDesc.downloadFile"),
		ApiGroup:    pointy.GetPointer("file"),
		Method:      pointy.GetPointer("GET"),
	})

	if err != nil {
		return err
	}

	return nil
}

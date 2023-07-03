package base

import (
	"context"
	"github.com/suyuan32/simple-admin-common/enum/common"
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

		err = l.initMenu()
		if err != nil {
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

	// Tag

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/tag/create"),
		Description: pointy.GetPointer("apiDesc.createTag"),
		ApiGroup:    pointy.GetPointer("tag"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/tag/update"),
		Description: pointy.GetPointer("apiDesc.updateTag"),
		ApiGroup:    pointy.GetPointer("tag"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/tag/delete"),
		Description: pointy.GetPointer("apiDesc.deleteTag"),
		ApiGroup:    pointy.GetPointer("tag"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/tag/list"),
		Description: pointy.GetPointer("apiDesc.getTagList"),
		ApiGroup:    pointy.GetPointer("tag"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/tag"),
		Description: pointy.GetPointer("apiDesc.getTagById"),
		ApiGroup:    pointy.GetPointer("tag"),
		Method:      pointy.GetPointer("Post"),
	})

	if err != nil {
		return err
	}

	return nil
}

func (l *InitDatabaseLogic) initMenu() error {
	menuData, err := l.svcCtx.CoreRpc.CreateMenu(l.ctx, &core.MenuInfo{
		Level:     pointy.GetPointer(uint32(1)),
		ParentId:  pointy.GetPointer(common.DefaultParentId),
		Path:      pointy.GetPointer("/fms_dir"),
		Name:      pointy.GetPointer("FileManagementDirectory"),
		Component: pointy.GetPointer("LAYOUT"),
		Sort:      pointy.GetPointer(uint32(3)),
		Disabled:  pointy.GetPointer(false),
		Meta: &core.Meta{
			Title:              pointy.GetPointer("route.fileManagement"),
			Icon:               pointy.GetPointer("ant-design:folder-open-outlined"),
			HideMenu:           pointy.GetPointer(false),
			HideBreadcrumb:     pointy.GetPointer(false),
			IgnoreKeepAlive:    pointy.GetPointer(false),
			HideTab:            pointy.GetPointer(false),
			CarryParam:         pointy.GetPointer(false),
			HideChildrenInMenu: pointy.GetPointer(false),
			Affix:              pointy.GetPointer(false),
		},
		MenuType: pointy.GetPointer(uint32(1)),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateMenu(l.ctx, &core.MenuInfo{
		Level:     pointy.GetPointer(uint32(2)),
		ParentId:  pointy.GetPointer(menuData.Id),
		Path:      pointy.GetPointer("/fms/file"),
		Name:      pointy.GetPointer("FileManagement"),
		Component: pointy.GetPointer("/fms/file/index"),
		Sort:      pointy.GetPointer(uint32(1)),
		Disabled:  pointy.GetPointer(false),
		Meta: &core.Meta{
			Title:              pointy.GetPointer("route.fileManagement"),
			Icon:               pointy.GetPointer("ant-design:folder-open-outlined"),
			HideMenu:           pointy.GetPointer(false),
			HideBreadcrumb:     pointy.GetPointer(false),
			IgnoreKeepAlive:    pointy.GetPointer(false),
			HideTab:            pointy.GetPointer(false),
			CarryParam:         pointy.GetPointer(false),
			HideChildrenInMenu: pointy.GetPointer(false),
			Affix:              pointy.GetPointer(false),
		},
		MenuType: pointy.GetPointer(uint32(1)),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateMenu(l.ctx, &core.MenuInfo{
		Level:     pointy.GetPointer(uint32(2)),
		ParentId:  pointy.GetPointer(menuData.Id),
		Path:      pointy.GetPointer("/fms/tag"),
		Name:      pointy.GetPointer("FileTagManagement"),
		Component: pointy.GetPointer("/fms/tag/index"),
		Sort:      pointy.GetPointer(uint32(2)),
		Disabled:  pointy.GetPointer(false),
		Meta: &core.Meta{
			Title:              pointy.GetPointer("route.fileTagManagement"),
			Icon:               pointy.GetPointer("ant-design:book-outlined"),
			HideMenu:           pointy.GetPointer(false),
			HideBreadcrumb:     pointy.GetPointer(false),
			IgnoreKeepAlive:    pointy.GetPointer(false),
			HideTab:            pointy.GetPointer(false),
			CarryParam:         pointy.GetPointer(false),
			HideChildrenInMenu: pointy.GetPointer(false),
			Affix:              pointy.GetPointer(false),
		},
		MenuType: pointy.GetPointer(uint32(1)),
	})

	if err != nil {
		return err
	}

	return err
}

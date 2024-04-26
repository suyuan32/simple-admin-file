package base

import (
	"context"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/suyuan32/simple-admin-common/msg/logmsg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/suyuan32/simple-admin-common/enum/errorcode"
	"github.com/suyuan32/simple-admin-common/i18n"
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

func (l *InitDatabaseLogic) InitDatabase(req *types.InitReq) (resp *types.BaseMsgResp, err error) {

	if req.App {
		if err := l.svcCtx.DB.Schema.Create(l.ctx, schema.WithForeignKeys(false)); err != nil {
			logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
			return nil, errorx.NewCodeError(errorcode.Internal, err.Error())
		}
	}

	if req.CoreApi && l.svcCtx.Config.CoreRpc.Enabled {
		err = l.insertApiData()
		if err != nil {
			if status.Code(err) == codes.InvalidArgument {
				return nil, errorx.NewCodeError(errorcode.InvalidArgument,
					l.svcCtx.Trans.Trans(l.ctx, "init.alreadyInit"))
			}
			return nil, err
		}
	}

	if req.CoreMenu && l.svcCtx.Config.CoreRpc.Enabled {
		err = l.insertMenuData()
		if err != nil {
			return nil, err
		}
	}

	err = l.svcCtx.Casbin.LoadPolicy()
	if err != nil {
		logx.Errorw("failed to load Casbin Policy", logx.Field("detail", err))
		return nil, errorx.NewCodeInternalError(i18n.DatabaseError)
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.Success)}, nil
}

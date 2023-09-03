package cloudsdk

import (
	"github.com/suyuan32/simple-admin-file/ent"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploaderGroup struct {
	TencentCOS     map[string]*cos.Client
	ProviderIdData map[string]uint64
}

func NewUploaderGroup(db *ent.Client) *UploaderGroup {
	uploaderGroup := UploaderGroup{}
	uploaderGroup.TencentCOS = make(map[string]*cos.Client)
	uploaderGroup.ProviderIdData = make(map[string]uint64)

	err := uploaderGroup.NewTencentClient(db)
	logx.Errorw("failed to load tencent cos config from database, you may need to initialize database",
		logx.Field("detail", err))

	return &uploaderGroup
}

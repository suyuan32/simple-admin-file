package cloudsdk

import (
	"github.com/suyuan32/simple-admin-file/ent"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type UploaderGroup struct {
	TencentCOS     map[string]*cos.Client
	ProviderIdData map[string]uint64
}

func NewUploaderGroup(db *ent.Client) *UploaderGroup {
	uploaderGroup := UploaderGroup{}

	_ = uploaderGroup.NewTencentClient(db)
	//logx.Must(err)

	return &uploaderGroup
}

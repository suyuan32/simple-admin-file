package cloudsdk

import (
	"context"
	"fmt"
	"github.com/suyuan32/simple-admin-file/ent"
	"github.com/suyuan32/simple-admin-file/ent/storageprovider"
	"github.com/suyuan32/simple-admin-file/internal/utils/dberrorhandler"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"net/url"
	"time"
)

func (g *UploaderGroup) NewTencentClient(db *ent.Client) error {
	data, err := db.StorageProvider.Query().Where(storageprovider.StateEQ(true)).All(context.Background())
	if err != nil {
		return dberrorhandler.DefaultEntError(logx.WithContext(context.Background()), err, nil)
	}

	for _, v := range data {
		u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", v.Bucket, v.Region))
		b := &cos.BaseURL{BucketURL: u}
		c := cos.NewClient(b, &http.Client{
			Timeout: 100 * time.Second,
			Transport: &cos.AuthorizationTransport{
				SecretID:  v.SecretID,
				SecretKey: v.SecretKey,
			},
		})

		g.TencentCOS[v.Name] = c
		g.ProviderData[v.Name] = struct {
			Id     uint64
			Folder string
		}{Id: v.ID, Folder: v.Folder}

		if v.IsDefault == true {
			g.DefaultProvider = v.Name
		}
	}

	return nil
}

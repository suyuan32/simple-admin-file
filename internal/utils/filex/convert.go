package filex

import (
	"net/url"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/simple-admin-file/internal/enum/filetype"
)

// ConvertFileTypeToUint8 converts file type string to uint8.
func ConvertFileTypeToUint8(fileType string) uint8 {
	switch fileType {
	case "other":
		return filetype.Other
	case "image":
		return filetype.Image
	case "video":
		return filetype.Video
	case "audio":
		return filetype.Audio
	default:
		return filetype.Other
	}
}

func ConvertUrlStringToFileUUID(urlStr string) (string, error) {
	urlData, err := url.Parse(urlStr)
	if err != nil {
		logx.Error("failed to parse url", logx.Field("details", err), logx.Field("data", urlStr))
		return "", err
	}

	fileId := filepath.Base(urlData.Path)

	if len(fileId) >= 36 {
		fileId = fileId[:36]
	} else if len(fileId) < 36 {
		return "", errorx.NewApiBadRequestError("wrong file path")
	}

	return fileId, nil
}

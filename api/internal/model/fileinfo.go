package model

import "gorm.io/gorm"

type FileInfo struct {
	gorm.Model
	UUID     string `json:"UUID" gorm:"index;comment:UUID"`
	Name     string `json:"name" gorm:"comment:file name"`
	FileType string `json:"fileType" gorm:"comment:file type"`
	Size     int64  `json:"size" gorm:"comment:file size"`
	Path     string `json:"path" gorm:"comment:store path"`
	UserUUID string `json:"userUUID" gorm:"index;comment:owner uuid"`
	Md5      string `json:"md5" gorm:"index;comment:file md5 string data"`
	Status   bool   `json:"status" gorm:"comment:status false private true public"`
}

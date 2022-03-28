package oss

import (
	"fmt"
	"os"

	"github.com/gofrs/uuid"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"

	"xiaohuazhu/internal/config"
)

func PushObject(path string, prefix string) (string, error) {
	// path
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 文件名
	uu, _ := uuid.NewV4()
	fileName := fmt.Sprintf("%s_%s.jpg", prefix, uu.String())
	// 上传图片
	if _, err := config.AllConn.Oss.PutObject(
		"image",
		fileName,
		f,
		-1,
		minio.PutObjectOptions{
			UserMetadata: nil,
			Progress:     nil,
			// 文件类型
			ContentType:             "image/jpeg",
			ContentEncoding:         "",
			ContentDisposition:      "",
			ContentLanguage:         "",
			CacheControl:            "",
			ServerSideEncryption:    nil,
			NumThreads:              0,
			StorageClass:            "",
			WebsiteRedirectLocation: "",
		},
	); err != nil {
		logrus.Errorf("minio pushObject, err: %s", err.Error())
		return "", err
	}
	logrus.Infof("minio pushobject 成功, 名称: %s", fileName)

	return "image/" + fileName, nil
}

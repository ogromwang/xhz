package oss

import (
	"fmt"
	"os"
	"xiaohuazhu/internal/util"

	"github.com/gofrs/uuid"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"

	"xiaohuazhu/internal/config"
)

func GetUrlByProtocol(uri string) string {
	return config.AllConfig.Oss.Protocol + config.AllConfig.Oss.Endpoint + "/" + uri
}

func PushObject(path string, prefix string) (string, error) {
	// path
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return PushObjectByFile(f, prefix)
}

func PushObjectByFile(f *os.File, prefix string) (string, error) {
	// 文件名
	uu, _ := uuid.NewV4()
	fileName := fmt.Sprintf("%s_%s%s", prefix, uu.String(), util.GetFileExt(f.Name()))
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
	logrus.Infof("minio pushobject 成功, 名称: image/%s", fileName)

	return "image/" + fileName, nil
}

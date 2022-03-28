package oss

import (
	"fmt"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type OssClient struct {
	minioClient *minio.Client
}

func MinioClient() *OssClient {
	endpoint := "127.0.0.1:9000"

	// 初使化 minio client对象。
	client, err := minio.NewWithOptions(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("minio", "123456789", ""),
		Secure: false,
	})
	if err != nil {
		logrus.Errorf("初始化minio 失败, err: %s", err.Error())
	}

	logrus.Infof("初始化 minio client 成功 %v\n", client)
	return &OssClient{
		minioClient: client,
	}
}

func (c *OssClient) PushObject(path string, prefix string) (string, error) {
	// path
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 文件名
	fileName := fmt.Sprintf("%s_%d.jpg", prefix, time.Now().Unix())
	// 上传图片
	if _, err := c.minioClient.PutObject(
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
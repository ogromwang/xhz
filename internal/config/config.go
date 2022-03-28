package config

import (
	"flag"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"xiaohuazhu/internal/model"
)

var (
	AllConfig *model.Config
	AllConn   *model.Conn
	OrmDB     *gorm.DB
	path      = flag.String("f", "../config/", "配置文件的位置")
)

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}

	var err error
	AllConfig = new(model.Config)
	AllConn = new(model.Conn)

	// application
	if _, err = toml.DecodeFile(*path+"application.toml", &AllConfig); err != nil {
		panic("初始化配置失败... err: " + err.Error())
	}

	// db
	if _, err = toml.DecodeFile(*path+"db.toml", &AllConfig); err != nil {
		panic("初始化配置失败... err: " + err.Error())
	}

	initDB(err)
	initMinio(err)
}

func initDB(err error) {
	// postgresql 链接初始化
	AllConn.Db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: AllConfig.Db.Dns,
		// disables implicit prepared statement usage
		PreferSimpleProtocol: AllConfig.Db.PreferSimpleProtocol,
	}), &gorm.Config{})
	if err != nil {
		panic("初始化db失败... err: " + err.Error())
	}
	logrus.Infof("初始化 db 成功")
}

func initMinio(err error) {
	ossConfig := AllConfig.Oss

	// 初使化 minio client对象。
	AllConn.Oss, err = minio.NewWithOptions(ossConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(ossConfig.Id, ossConfig.Secret, ossConfig.Token),
		Secure: false,
	})
	if err != nil {
		logrus.Errorf("初始化minio 失败, err: %s", err.Error())
	}
	logrus.Infof("初始化 minio client 成功")

}

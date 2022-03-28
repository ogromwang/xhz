package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var path = flag.String("f", "../config/", "配置文件的位置")
var OrmDB *gorm.DB
var AllConfig *Config

type Config struct {
	Application Application
	Db          Db
}

type Application struct {
	Name string
	Port string
}

type Db struct {
	Dns                  string
	PreferSimpleProtocol bool
}

func init() {
	flag.Parse()

	var err error
	AllConfig = new(Config)

	// application
	if _, err = toml.DecodeFile(*path+"application.toml", &AllConfig); err != nil {
		panic("初始化配置失败... err: " + err.Error())
	}

	// db
	if _, err = toml.DecodeFile(*path+"db.toml", &AllConfig); err != nil {
		panic("初始化配置失败... err: " + err.Error())
	}

	// postgresql 链接初始化
	OrmDB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: AllConfig.Db.Dns,
		// disables implicit prepared statement usage
		PreferSimpleProtocol: AllConfig.Db.PreferSimpleProtocol,
	}), &gorm.Config{})
	if err != nil {
		panic("初始化db失败... err: " + err.Error())
	}
	logrus.Infof("初始化 db 成功")

}

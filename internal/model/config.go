package model

import (
	"github.com/minio/minio-go"
	"gorm.io/gorm"
)

// Config 配置初始化
type Config struct {
	Application *Application
	Db          *Db
	Oss         *Oss
}

type Application struct {
	Name string
	Port string
	Auth *Auth
}

type Auth struct {
	PasswordSalt  string
	JwtSigned     string
	JwtExpireHour int
}

type Db struct {
	Dns                  string
	PreferSimpleProtocol bool
}

type Oss struct {
	Endpoint string
	Protocol string
	Id       string
	Secret   string
	Token    string
}

// Conn 全局连接
type Conn struct {
	Db  *gorm.DB
	Oss *minio.Client
}

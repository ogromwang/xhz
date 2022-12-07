## xhz 一个基于 Flutter Go 的在线亲友记账共享 app
[![docs](https://img.shields.io/badge/docs-reference-green.svg)](https://github.com/ogromwang)
[![email](https://img.shields.io/badge/email-ogromwang@gmail.com-red.svg)](ogromwang@gmail.com)

### 项目介绍
后端使用 golang gin gorm 开发，移动端使用 Flutter 构建。
[前端项目地址](https://github.com/ogromwang/xhz_app)

## 技术选型
### 数据库
1. postgresql

### 对象存储
1. minio

### 鉴权
1. jwt

## 项目结构

```
.
├── Dockerfile
├── LICENSE
├── README.md
├── build.sh
├── cmd
│   └── main.go
├── config
│   ├── application.toml
│   └── db.toml
├── db
│   └── init.sql
├── go.mod
└── internal
    ├── config
    ├── dao
    ├── model
    ├── server
    ├── service
    └── util
```

## 从源码构建
 `config` 文件夹下

1. db 配置

 ```toml
[Db]
    Dns = "host=127.0.0.1 user=writer password=123 dbname=xiaohuazhu port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    PreferSimpleProtocol = true

[Oss]
    Endpoint = "127.0.0.1:9000"
    AccessPrefix = "http://127.0.0.1/"
    Id = "minio"
    Secret = "123456789"
    Token = ""
```

2. 应用配置

name、端口、默认 Icon 、JWT 配置
```toml
[Application]
    Name = "xiaohuazhu"
    Port = "8080"
    DefaultIcon = "image/picture_b27918ac-b14c-47b8-a93f-f07ccd122aa3.jpg"

[Application.Auth]
    PasswordSalt = "_salt"
    JwtSigned = "bf19d68e1caee5738c7c7194107ec897"
    JwtExpireHour = 168
```

## Dockerfile

```shell
# 第一阶段构建
FROM golang:1.17 AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.io
ENV GOARCH amd64

WORKDIR /opt/ogrom.com/xiaohuazhu
COPY go.mod .
COPY go.sum .
RUN  go mod tidy
COPY  . /opt/ogrom.com/xiaohuazhu
RUN go build -a -installsuffix cgo -o /opt/ogrom.com/xiaohuazhu/cmd/main ./cmd/

# 第二阶段构建
FROM scratch
WORKDIR /opt/ogrom.com/xiaohuazhu
COPY --from=builder /opt/ogrom.com/xiaohuazhu .
# 暴露的 port
EXPOSE 8000
# 启动
ENTRYPOINT ["./cmd/main", "-f", "./config/"]
```

## 二进制文件启动
```shell
nohup ./xhz -f ./config/ &
```
```text
[GIN-debug] Listening and serving HTTP on :8080
time="2022-04-25T18:40:18+08:00" level=info msg="初始化 db 成功"
time="2022-04-25T18:40:18+08:00" level=info msg="初始化 minio client 成功"
time="2022-04-25T18:40:18+08:00" level=info msg="程序启动中..."
```
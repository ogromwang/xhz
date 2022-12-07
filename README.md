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
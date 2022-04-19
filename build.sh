#!/bin/sh
echo building xiaohuazhu

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./cmd/main ./cmd/main.go

# 删除之前的镜像
docker rmi -f xiaohuazhu

docker build -t xiaohuazhu .
#    FROM scratch
#    # 使用基础镜像
#    # 环境变量
#    ENV GOPROXY https://goproxy.cn,direct
#
#    # work dir 在容器的的目录位置，没有则会新建
#    WORKDIR /opt/ogrom.com/xiaohuazhu
#    # 将当前 . 的内容 copy 到容器目标目录中
#    COPY . /opt/ogrom.com/xiaohuazhu
#
#    # 暴露的 port
#    EXPOSE 8000
#    # 容器启动时候执行的命名
#    ENTRYPOINT ["./cmd/main", "-f", "./config/"]



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

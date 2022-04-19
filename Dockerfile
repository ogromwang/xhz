FROM scratch
# 使用基础镜像
# 环境变量
ENV GOPROXY https://goproxy.cn,direct

# work dir 在容器的的目录位置，没有则会新建
WORKDIR /opt/ogrom.com/xiaohuazhu
# 将当前 . 的内容 copy 到容器目标目录中
COPY . /opt/ogrom.com/xiaohuazhu

# 暴露的 port
EXPOSE 8000
# 容器启动时候执行的命名
ENTRYPOINT ["./cmd/main", "-f", "./config/"]
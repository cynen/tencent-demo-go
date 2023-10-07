# 带自编译的Dockerfile
# alpine的源配置的是ustc.可以更换的.

# 如果多次在同一台服务器构建,注意构建命令,不要使用缓存.
# docker build --no-cache -t test .

# complie
FROM golang:1.20-alpine AS gobuilder
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
        && apk add git \
        && go env -w GOPROXY=https://goproxy.cn,direct \
        && go env -w GO111MODULE="auto" \
        && git clone https://gitee.com/cynen/tencent-demo-go.git \
        && cd tencent-demo-go/qywx-tyqw \
        && CGO_ENABLED=0 go build -o tyqw
# 代码仓github和gitee已经实现了同步,我们也可以配置 github的代码仓:
# https://github.com/cynen/tencent-demo-go.git

# build images

FROM alpine:3.18
ENV LANG C.UTF-8
ENV TZ=Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update && apk upgrade \
    && apk add ca-certificates\
    && update-ca-certificates \
    && apk --no-cache add openssl wget bash tzdata curl \
    && set -ex
# 从上一个镜像中copy到当前容器里.
COPY --from=gobuilder /go/tencent-demo-go/qywx-tyqw/tyqw /usr/bin/
# 将配置文件也copy进去.
COPY config.yml /etc/

ENTRYPOINT ["tyqw"]
CMD ["-c","/etc/config.yml"]

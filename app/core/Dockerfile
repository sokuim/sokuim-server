FROM golang:1.25-alpine AS Builder

WORKDIR /app

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o core ./app/core

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# etcd 地址构建参数：构建镜像时由 CI 从 GitHub Variables（vars.ETCD_HOST）注入，
# 本地 docker build 未传参时回退到默认地址
ARG ETCD_HOST=127.0.0.1:2379

# /app/core: core为上一步编译出来的
# /app/app/core/etc/core.build.yaml: core.build.yaml 为仓库中的文件
COPY --from=Builder /app/core .
COPY --from=Builder /app/app/core/etc/core.build.yaml ./etc/core.yaml

# 将 core.build.yaml 中的 ${ETCD_HOST} 占位符替换为实际 etcd 地址
RUN sed -i "s|\${ETCD_HOST}|${ETCD_HOST}|g" /app/etc/core.yaml

EXPOSE 8083

CMD ["./core", "-f", "etc/core.yaml"]
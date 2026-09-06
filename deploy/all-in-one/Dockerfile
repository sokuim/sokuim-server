# sokuim all-in-one 单镜像：一个镜像内打包 4 个 Go 服务（core/gateway/job/comet）+ etcd
# 本地构建：docker compose up -d --build（不再从 GHCR 拉取 4 个独立镜像）
#
# 结构：同一镜像在 compose 里被 5 个容器复用，各自用 command 指定跑哪个进程：
#   - etcd      -> etcd-entrypoint.sh（首次启动自动引导 root/admin 认证）
#   - core      -> ./core -f etc/core.yaml
#   - gateway   -> ./gateway -f etc/gateway.yaml
#   - job       -> ./job -f etc/job.yaml
#   - comet     -> ./comet -f etc/comet.yaml

# ============ 阶段 1：编译 4 个 Go 服务 ============
FROM golang:1.25-alpine AS Builder

WORKDIR /app

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# -p 1 串行编译 + GOGC=20 更激进 GC，压低峰值内存，避免低内存 Docker（如 1.94GiB）下编译 go-redis/v9 被 OOM kill
RUN GOGC=20 go build -p 1 -ldflags="-s -w" -o core ./app/core \
 && GOGC=20 go build -p 1 -ldflags="-s -w" -o gateway ./app/gateway \
 && GOGC=20 go build -p 1 -ldflags="-s -w" -o job ./app/job \
 && GOGC=20 go build -p 1 -ldflags="-s -w" -o comet ./app/comet

# ============ 阶段 2：提取 etcd 二进制 ============
FROM bitnamilegacy/etcd:3.6.4 AS Etcd

# ============ 阶段 3：运行时 ============
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl bash && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# grpc_health_probe：core 是 gRPC 服务，用它做健康检查（go-zero zrpc 已注册标准 health 服务）
# 用 uname -m 映射架构，避免依赖 builder 注入的 TARGETARCH（传统 builder 下为空导致 404）
RUN GRPC_HEALTH_PROBE_ARCH="$(uname -m | sed 's/aarch64/arm64/; s/x86_64/amd64/')" \
    && curl -fsSL -o /usr/local/bin/grpc_health_probe \
      "https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.24/grpc_health_probe-linux-${GRPC_HEALTH_PROBE_ARCH}" \
    && chmod +x /usr/local/bin/grpc_health_probe

WORKDIR /app

# etcd 二进制（静态链接，可独立运行于 alpine）
COPY --from=Etcd /opt/bitnami/etcd/bin/etcd /usr/local/bin/etcd
COPY --from=Etcd /opt/bitnami/etcd/bin/etcdctl /usr/local/bin/etcdctl

# etcd 首次启动认证引导脚本
COPY deploy/all-in-one/scripts/etcd-entrypoint.sh /usr/local/bin/etcd-entrypoint.sh
RUN chmod +x /usr/local/bin/etcd-entrypoint.sh

# 4 个 Go 服务二进制
COPY --from=Builder /app/core .
COPY --from=Builder /app/gateway .
COPY --from=Builder /app/job .
COPY --from=Builder /app/comet .

# 各服务配置（core.build.yaml 含 ${ETCD_HOST} 占位符，下面替换）
COPY --from=Builder /app/app/core/etc/core.build.yaml ./etc/core.yaml
COPY --from=Builder /app/app/gateway/etc/gateway.build.yaml ./etc/gateway.yaml
COPY --from=Builder /app/app/job/etc/job.build.yaml ./etc/job.yaml
COPY --from=Builder /app/app/comet/etc/comet.build.yaml ./etc/comet.yaml

# etcd 地址：all-in-one 编排内 etcd 服务名为 etcd，core 通过容器网络 DNS 访问
ARG ETCD_HOST=etcd:2379
RUN sed -i "s|\${ETCD_HOST}|${ETCD_HOST}|g" /app/etc/core.yaml

EXPOSE 8080 8081 8082 8083 2379 2380

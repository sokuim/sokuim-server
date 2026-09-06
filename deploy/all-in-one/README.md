# sokuim all-in-one 一键部署

> 本地编译 + 单镜像编排，一条命令把 sokuim-server 的 4 个服务（core / gateway / job / comet）+ etcd 全部拉起。

## 1. 这是什么

传统发布（`.github/workflows/release.yml`）是「CI 里 build → 推 GHCR → SSH 到服务器分别部署」，
每个服务一个独立镜像。本目录提供另一种**本地一键起**方式：

- **一个镜像**：一个 `sokuim-all-in-one:latest` 镜像内打包了 4 个 Go 服务二进制 + etcd 二进制。
- **本地编译**：`docker compose up --build` 时用本机 Docker 现场编译，不依赖 GHCR 上已推送的镜像。
- **5 个容器**：同一镜像被复用，各自用 `command` 指定跑哪个进程。

适合：本地开发联调、演示环境、单机快速验收。生产分步发布仍走 `release.yml` / `infra.yml`，互不影响。

## 2. 目录结构

```
deploy/all-in-one/
├── Dockerfile               # 多阶段构建：编译 4 服务 + 提取 etcd 二进制 -> 运行时
├── docker-compose.yml       # 编排：5 个容器共用同一镜像
├── deploy.sh                # 一键部署脚本（推荐入口，见 §4）
├── .env                     # 宿主机端口映射（按需覆盖）
├── scripts/
│   └── etcd-entrypoint.sh   # etcd 启动 + 首次认证引导脚本
└── README.md                # 本文档
```

## 3. 架构说明

### 3.1 单镜像，多容器

所有服务共用镜像 `sokuim-all-in-one:latest`。镜像里同时放了 4 个 Go 二进制和 etcd：

| 文件 | 路径 | 说明 |
|------|------|------|
| core / gateway / job / comet | `/app/` | 4 个 Go 服务二进制 |
| 各服务配置 | `/app/etc/*.yaml` | 由 `*.build.yaml` 复制而来 |
| etcd / etcdctl | `/usr/local/bin/` | 从 `bitnamilegacy/etcd:3.6.4` 提取 |
| grpc_health_probe | `/usr/local/bin/` | core（gRPC 服务）的健康检查探针 |
| etcd-entrypoint.sh | `/usr/local/bin/` | etcd 启动 + 认证引导 |

compose 里用 YAML anchor（`x-sokuim-image`）把同一份 `image + build` 复用给 5 个容器，
每个容器只靠 `command` 区分跑哪个进程：

```
etcd     -> etcd-entrypoint.sh
core     -> ./core     -f etc/core.yaml
gateway  -> ./gateway  -f etc/gateway.yaml
job      -> ./job      -f etc/job.yaml
comet    -> ./comet    -f etc/comet.yaml
```

### 3.2 容器与端口

| 容器 | 容器内端口 | 宿主机端口 | 协议 | 说明 |
|------|-----------|-----------|------|------|
| etcd   | 2379 / 2380 | 2379 / 2380 | gRPC | 服务注册中心，2380 为 peer 端口（单节点仍占用） |
| core   | 8083 | `${CORE_PORT:-8083}` | gRPC | 核心 RPC 服务，注册到 etcd |
| gateway | 8080 | `${GATEWAY_PORT:-8080}` | HTTP | 网关 |
| job    | 8081 | `${JOB_PORT:-8081}` | HTTP | 任务服务 |
| comet  | 8082 | `${COMET_PORT:-8082}` | HTTP | 长连接服务 |

宿主机端口可在 `.env` 覆盖（见 §5）。**容器内端口固定不可改**（写死在服务配置 `ListenOn` 里）。

### 3.3 网络与依赖

- 所有容器加入同一个 bridge 网络 `sokuim-net`，容器间通过服务名互访（如 `etcd:2379`）。
- 依赖链只有一条：`core` 依赖 `etcd`，且是 `condition: service_healthy`——即 **etcd 通过健康检查后 core 才启动**。
- gateway / job / comet 是纯 HTTP 服务，不依赖 etcd，各自独立健康检查。

### 3.4 etcd 认证

`core` 的配置写死了 `Etcd.User: root / Pass: admin`，go-zero 会带凭据访问 etcd，
因此 etcd **必须开启认证**。镜像内的 `etcd-entrypoint.sh` 负责：

1. **首次启动**（数据目录不存在）：后台以 `new` 状态启动 → 创建 root 用户（密码 admin）→ 授予 root 角色 → 开启认证 → 停掉 → 前台以 `existing` 状态常驻。
2. **再次启动**（数据目录已有 `member/`）：直接前台 `existing` 启动，认证保持不变（幂等）。

数据通过命名卷 `etcd-data` 持久化，容器重建后认证与注册信息不丢失。

## 4. 快速开始

### 4.1 前置条件

- Docker ≥ 20.x，且 `docker compose`（V2）可用。
- 建议给 Docker 分配 **≥ 2 GB 内存**。编译 `go-redis/v9`（go-zero 间接依赖）单包峰值约 1 GB，
  内存过小会 OOM。构建脚本已用 `-p 1` + `GOGC=20` 压低峰值，但 1.94 GB 以下仍可能失败。
- 宿主机端口 8086/8081/8082/8085/2379/2380 空闲（即 `.env` 里配置的端口；若被占用，见 §7 故障排查）。

### 4.2 一键启动（推荐）

```bash
cd deploy/all-in-one
./deploy.sh start
```

脚本会依次执行：停止并移除容器 → 删除旧镜像 → **全量重新编译镜像**（`--no-cache`）→
`docker compose up -d` 启动 → 等待全部容器 healthy → 自动验证（3 个 HTTP 服务 healthz + core 注册 etcd），
失败时非零退出并报错。

> 注意：`start` 每次都会全量重新编译（约几分钟），确保代码绝对最新。若只想快速重启容器（不重建镜像），
> 用 `./deploy.sh stop` 后手动 `docker compose up -d`，或直接 `docker compose restart`。

> 想手动分步执行：
> ```bash
> cd 项目根目录
> docker compose -f deploy/all-in-one/docker-compose.yml down
> docker rmi sokuim-all-in-one:latest
> docker build --no-cache -f deploy/all-in-one/Dockerfile -t sokuim-all-in-one:latest .
> cd deploy/all-in-one && docker compose up -d
> ```

### 4.3 脚本子命令

`deploy.sh` 提供以下子命令：

| 命令 | 作用 |
|------|------|
| `./deploy.sh start` | 停止 + 删旧镜像 + 全量重建 + 启动 + 等健康 + 验证（最常用） |
| `./deploy.sh restart` | 同 start（停止 + 删镜像 + 全量重建 + 启动） |
| `./deploy.sh stop` | 停止所有容器（保留容器与 etcd 数据） |
| `./deploy.sh status` | 查看容器状态 |
| `./deploy.sh logs [service]` | 查看日志（`-f` 跟随，可指定服务名） |
| `./deploy.sh verify` | 仅运行健康验证（不启动） |
| `./deploy.sh down` | 停止并移除容器（保留 etcd 数据卷） |
| `./deploy.sh clean` | 停止并移除容器 + 清空 etcd 数据 |

### 4.4 查看状态

```bash
./deploy.sh status          # 或 docker compose ps
```

5 个容器全部 `running (healthy)` 即部署成功。

### 4.5 验证

```bash
./deploy.sh verify
# 等价于手动执行：
# curl http://localhost:8086/common/healthz   # gateway（见 .env 的 GATEWAY_PORT）
# curl http://localhost:8081/common/healthz   # job
# curl http://localhost:8082/common/healthz   # comet
# 均返回 {"status":0,"success":true,...}
#
# docker exec all-in-one_etcd_1 \
#   etcdctl --endpoints=http://127.0.0.1:2379 --user root:admin \
#   get --prefix core.rpc
# 应看到 core.rpc/<id> 及其地址
```

### 4.6 停止 / 清理

```bash
./deploy.sh stop             # 停止（保留容器与数据）
./deploy.sh down             # 停止并移除容器（保留 etcd 数据卷）
./deploy.sh clean            # 停止并移除容器 + 清空 etcd 数据（下次启动重新引导认证）
```

## 5. 配置说明

### 5.1 宿主机端口（`.env`）

`.env` 只覆盖**宿主机侧端口**，容器内端口不变：

```ini
CORE_PORT=8085      # core   宿主机 8085 -> 容器 8083
GATEWAY_PORT=8086   # gateway 宿主机 8086 -> 容器 8080
JOB_PORT=8081       # job     宿主机 8081 -> 容器 8081
COMET_PORT=8082     # comet   宿主机 8082 -> 容器 8082
```

> 注：本仓库 `.env` 里 `CORE_PORT=8085`、`GATEWAY_PORT=8086`，是因为开发机 8083 被 `emqx` 容器占用、
> 8080 被其它容器占用。换一台端口空闲的机器可改回 `CORE_PORT=8083`、`GATEWAY_PORT=8080`。

### 5.2 etcd 环境变量

compose 里 etcd 服务的 `environment` 用的是 `SOKUIM_` 前缀：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `SOKUIM_ROOT_PASSWORD` | `admin` | etcd root 用户密码 |
| `SOKUIM_ADVERTISE_URL` | `http://etcd:2379` | core 访问 etcd 的地址，必须是容器网络内可达的服务名 |
| `SOKUIM_DATA_DIR` | `/var/lib/etcd` | etcd 数据目录（对应命名卷挂载点） |

> ⚠️ **不要改成 `ETCD_` 前缀**。etcd 官方二进制会把 `ETCD_*` 当成自己的 flag 环境变量
> （如 `ETCD_DATA_DIR` 与 `--data-dir` 冲突），导致启动直接 fatal。

### 5.3 镜像名与构建上下文

- 镜像名：`sokuim-all-in-one:latest`（compose 的 `image` 字段）。
- 构建上下文是**仓库根目录**（`../..`），`dockerfile` 指向 `deploy/all-in-one/Dockerfile`。
  因此必须在仓库根目录存在时才能构建（compose 已通过相对路径自动定位）。

### 5.4 core 的 etcd 地址

`core` 的 `etc/core.yaml` 在**镜像构建时**由 `Dockerfile` 用 `sed` 把 `${ETCD_HOST}` 占位符
替换为 `etcd:2379`（构建参数 `ARG ETCD_HOST=etcd:2379`）。这与 `release.yml` 用
`${{ vars.ETCD_HOST }}` 注入的方式对应——all-in-one 场景固定为容器网络内的 `etcd:2379`。

## 6. 镜像构建说明（Dockerfile）

三阶段构建：

1. **Builder**（`golang:1.25-alpine`）：`go mod download` → 编译 4 个服务。
   用 `GOGC=20` + `-p 1` 降低编译峰值内存（低内存 Docker 环境必需）。
2. **Etcd**（`bitnamilegacy/etcd:3.6.4`）：仅作为来源，供提取 etcd / etcdctl 二进制。
3. **运行时**（`alpine:3.19`）：装运行依赖 + `grpc_health_probe`，再从 Builder / Etcd 拷贝二进制与配置。

如需单独重新构建镜像（不改编排），加 `--no-cache` 强制全量编译，不加则走缓存：

```bash
cd 项目根目录
docker build --no-cache -f deploy/all-in-one/Dockerfile -t sokuim-all-in-one:latest .
```

## 7. 故障排查

| 现象 | 原因 | 处理 |
|------|------|------|
| 构建时 `signal: killed` | Docker 内存不足，编译大依赖被 OOM | 调大 Docker 内存至 ≥ 2 GB |
| 构建时 `failed to fetch anonymous token ... i/o timeout` | buildkit 走 IPv6 拉取 `auth.docker.io` 超时 | 用 `./deploy.sh start`（脚本内已用 legacy builder 规避）；或手动 `DOCKER_BUILDKIT=0 docker build ...` |
| `Bind for 0.0.0.0:8083 failed: port is already allocated` | 宿主机端口被占用 | 改 `.env` 的对应端口为其它空闲端口 |
| `etcd` 容器 unhealthy / `connection refused` | etcd 未正常启动 | 看日志 `docker logs all-in-one_etcd_1`；若报 `conflicting environment variable` 说明误用了 `ETCD_` 前缀环境变量 |
| core 一直 `created` 不启动 | `depends_on` 等 etcd healthy | etcd 健康后 core 会自动启动，先解决 etcd |
| core 报连不上 etcd / 认证失败 | etcd 认证未引导成功，或 `SOKUIM_ADVERTISE_URL` 不对 | 确认 `SOKUIM_ADVERTISE_URL=http://etcd:2379`，必要时 `./deploy.sh clean` 重新引导 |
| 想彻底重置环境 | — | `./deploy.sh clean`（清空 etcd 数据）后重新 `./deploy.sh start` |

### 查看日志

```bash
./deploy.sh logs -f etcd        # etcd 日志（看认证引导过程）
./deploy.sh logs -f core        # core 日志（看 RPC server 启动）
./deploy.sh logs gateway        # 网关日志（不带 -f 只输出最近日志）
```

## 8. 与其他部署方式的关系

| 方式 | 触发 | 用途 | 镜像来源 |
|------|------|------|----------|
| `release.yml` | push 到 main/master | 生产分步发布 | GHCR 构建推送 |
| `infra.yml` | 手动 | 基础设施（etcd/mysql/redis） | 固定镜像，纯部署 |
| **本目录** | 手动 `./deploy.sh start` | 本地一键联调/演示 | 本地编译 |

三者并存、互不影响。本目录未设置 `container_name`，避免与 release / infra 单独部署的同名容器冲突。

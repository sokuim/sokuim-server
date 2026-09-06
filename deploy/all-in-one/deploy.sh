#!/usr/bin/env bash
# sokuim all-in-one 一键部署脚本
#
# 用法：
#   ./deploy.sh start           停止+删旧镜像+全量重建+启动+等健康+验证（最常用）
#   ./deploy.sh restart         同 start：停止+删镜像+全量重建+启动（慢，确保代码最新）
#   ./deploy.sh stop            停止（保留容器与 etcd 数据）
#   ./deploy.sh status          查看容器状态
#   ./deploy.sh logs [service]  查看日志（-f 跟随，可指定服务名）
#   ./deploy.sh verify          仅运行健康验证（不启动）
#   ./deploy.sh down            停止并移除容器（保留 etcd 数据）
#   ./deploy.sh clean           停止并移除容器 + 清空 etcd 数据
#
# 前置条件：Docker >= 20.x 且 docker compose (V2) 可用；建议给 Docker 分配 >= 2GB 内存。

set -euo pipefail

# 切到脚本所在目录（compose 文件与 .env 同目录）
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

COMPOSE="docker compose"
SERVICES=(etcd core gateway job comet)
HEALTH_TIMEOUT=${HEALTH_TIMEOUT:-180}

# 使用 legacy builder（而非 buildkit）：部分网络环境下 buildkit 走 IPv6 拉取
# auth.docker.io 会超时，legacy builder 走本地缓存可正常构建。
export DOCKER_BUILDKIT=0
export COMPOSE_DOCKER_CLI_BUILD=0

# 读取 .env 里的宿主机端口（供验证/释放端口用），带默认值兜底
if [ -f .env ]; then
  set -a; source .env; set +a
fi
GATEWAY_PORT=${GATEWAY_PORT:-8080}
JOB_PORT=${JOB_PORT:-8081}
COMET_PORT=${COMET_PORT:-8082}

# ---------- 工具函数 ----------

info()  { echo -e "\033[36m[info]\033[0m $*"; }
ok()    { echo -e "\033[32m[ ok ]\033[0m $*"; }
fail()  { echo -e "\033[31m[fail]\033[0m $*" >&2; }

# 读取某个容器的健康状态：healthy / starting / unhealthy / missing
get_health() {
  local svc="$1"
  local cid
  cid=$($COMPOSE ps -q "$svc" 2>/dev/null || true)
  [ -z "$cid" ] && { echo "missing"; return; }
  docker inspect --format '{{.State.Health.Status}}' "$cid" 2>/dev/null || echo "missing"
}

# 等待所有容器 healthy
wait_healthy() {
  local timeout="$HEALTH_TIMEOUT"
  local elapsed=0 interval=3
  info "等待全部容器 healthy（最长 ${timeout}s）..."
  while [ "$elapsed" -lt "$timeout" ]; do
    local all=1
    for svc in "${SERVICES[@]}"; do
      local h
      h="$(get_health "$svc")"
      if [ "$h" != "healthy" ]; then
        all=0
        break
      fi
    done
    [ "$all" -eq 1 ] && { ok "全部容器已 healthy"; return 0; }
    sleep "$interval"
    elapsed=$((elapsed + interval))
  done
  fail "等待健康超时（${timeout}s）"
  $COMPOSE ps
  return 1
}

# 验证：3 个 HTTP 服务 healthz + core 注册到 etcd
# 关键点：curl 必须加 --noproxy '*' 强制绕过系统代理。本机开启了 127.0.0.1:7890 的
#         系统代理（Clash/V2Ray），curl 访问 localhost/127.0.0.1 时会走该代理，
#         代理对 localhost 请求间歇性返回 502，而浏览器对 localhost 绕过代理所以正常。
#         这就是「网页正常、脚本 curl 502」的根因。
# 每个 healthz 带重试（5 次，间隔 1s）+ 超时，覆盖启动抖动。
verify() {
  info "验证服务健康..."
  local svc port i up err
  for svc in gateway job comet; do
    case "$svc" in
      gateway) port="$GATEWAY_PORT" ;;
      job)     port="$JOB_PORT" ;;
      comet)   port="$COMET_PORT" ;;
    esac
    up=0
    for i in $(seq 1 5); do
      # -s 静默 body；-S 出错时把错误写到 stderr；--noproxy 绕过系统代理；--max-time 防止挂死
      if err="$(curl -sSf --noproxy '*' --max-time 3 "http://127.0.0.1:${port}/common/healthz" 2>&1)"; then
        up=1
        break
      fi
      sleep 1
    done
    if [ "$up" -eq 1 ]; then
      ok "$svc (:${port}) healthz 正常"
    else
      fail "$svc (:${port}) healthz 失败（重试 5 次仍不通）"
      [ -n "$err" ] && fail "  原因: $err"
      return 1
    fi
  done

  # core 注册检查也带重试：go-zero 周期续约 etcd 租约，偶发在两次续约间隙查不到注册记录
  up=0
  for i in $(seq 1 3); do
    if $COMPOSE exec -T etcd etcdctl \
        --endpoints=http://127.0.0.1:2379 --user root:admin \
        get --prefix core.rpc 2>/dev/null | grep -q "core.rpc"; then
      up=1
      break
    fi
    sleep 1
  done
  if [ "$up" -eq 1 ]; then
    ok "core 已注册到 etcd"
  else
    fail "core 未注册到 etcd"
    fail "  原因: etcd 中无 core.rpc 注册记录（core 可能启动失败或未连上 etcd）"
    return 1
  fi

  ok "验证通过"
}

# ---------- 子命令 ----------

# 停止容器并删除旧镜像（保留 etcd 数据卷）
stop_and_remove_image() {
  info "停止并移除容器（保留 etcd 数据卷）..."
  $COMPOSE down

  info "删除旧镜像 sokuim-all-in-one:latest ..."
  if docker image inspect sokuim-all-in-one:latest >/dev/null 2>&1; then
    docker rmi -f sokuim-all-in-one:latest
    ok "旧镜像已删除"
  else
    info "镜像不存在，跳过删除"
  fi
}

# 全量重建镜像（--no-cache：删镜像后 build 中间层缓存仍可能残留，须 --no-cache 才彻底）
build_image() {
  info "全量编译镜像（legacy builder + --no-cache）..."
  # 用 docker build（DOCKER_BUILDKIT=0）而非 compose 的 buildkit：
  # 部分网络环境下 buildkit 走 IPv6 拉取 auth.docker.io 会超时，legacy 走本地缓存可正常构建。
  docker build --no-cache \
    -f "$SCRIPT_DIR/Dockerfile" \
    -t sokuim-all-in-one:latest \
    "$SCRIPT_DIR/../.." \
  || { fail "镜像构建失败"; exit 1; }
}

# 完整部署流程：停止 + 删镜像 + 全量重建 + 启动 + 等健康 + 验证
deploy_full() {
  stop_and_remove_image
  build_image

  info "启动编排（docker compose up -d）..."
  $COMPOSE up -d

  wait_healthy
  verify
}

cmd_start() {
  deploy_full
}

cmd_stop() {
  info "停止所有容器..."
  $COMPOSE stop
}

cmd_restart() {
  deploy_full
}

cmd_status() {
  $COMPOSE ps
}

cmd_logs() {
  # 透传参数：logs -f / logs gateway 等
  $COMPOSE logs "$@"
}

cmd_down() {
  info "停止并移除容器（保留 etcd 数据卷）..."
  $COMPOSE down
}

cmd_clean() {
  info "停止并移除容器 + 清空 etcd 数据卷..."
  $COMPOSE down -v
}

usage() {
  sed -n '2,14p' "$0" | sed 's/^# \{0,1\}//'
}

# ---------- 入口 ----------

case "${1:-}" in
  start)    cmd_start ;;
  stop)     cmd_stop ;;
  restart)  cmd_restart ;;
  status)   cmd_status ;;
  logs)     shift; cmd_logs "$@" ;;
  verify)   verify ;;
  down)     cmd_down ;;
  clean)    cmd_clean ;;
  help|-h|--help|"") usage ;;
  *)        fail "未知命令: $1"; echo; usage; exit 1 ;;
esac

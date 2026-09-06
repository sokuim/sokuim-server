#!/bin/bash
# etcd 单节点 + 首次启动认证引导（root/admin）
# 与 bitnamilegacy/etcd 镜像行为对齐：首次启动时创建 root 用户、授予 root 角色、开启认证。
#
# 幂等：仅当数据目录不存在（首次启动）时引导认证；重启时数据目录已有 member，
#       直接以 existing 状态前台启动，认证保持不变。
#
# 环境变量（用 SOKUIM_ 前缀，避免与 etcd 官方二进制的 ETCD_* flag 环境变量冲突）：
#   SOKUIM_ROOT_PASSWORD   root 用户密码（默认 admin）
#   SOKUIM_ADVERTISE_URL   core 侧访问 etcd 的地址（默认 http://etcd:2379）
#                         必须是 compose 网络内其它容器可达的地址，否则 go-zero 的
#                         AutoSync 会同步到错误端点导致 core 连不上。
#   SOKUIM_DATA_DIR        etcd 数据目录（默认 /var/lib/etcd）

set -euo pipefail

DATA_DIR="${SOKUIM_DATA_DIR:-/var/lib/etcd}"
ROOT_PASSWORD="${SOKUIM_ROOT_PASSWORD:-admin}"
ADVERTISE_URL="${SOKUIM_ADVERTISE_URL:-http://etcd:2379}"

mkdir -p "$DATA_DIR"

# 等待 etcd 就绪（未开认证阶段，不带凭据）
wait_ready() {
  local i
  for i in $(seq 1 60); do
    if etcdctl --endpoints=http://127.0.0.1:2379 endpoint health >/dev/null 2>&1; then
      return 0
    fi
    sleep 1
  done
  echo "etcd failed to become ready within 60s" >&2
  return 1
}

# 首次启动：后台以 new 状态启动 -> 引导认证 -> 停止 -> 前台以 existing 状态常驻
if [ ! -d "$DATA_DIR/member" ]; then
  echo "[etcd] first boot: starting with initial-cluster-state=new"
  etcd --name etcd0 \
    --data-dir "$DATA_DIR" \
    --listen-client-urls http://0.0.0.0:2379 \
    --advertise-client-urls "$ADVERTISE_URL" \
    --listen-peer-urls http://0.0.0.0:2380 \
    --initial-advertise-peer-urls http://127.0.0.1:2380 \
    --initial-cluster etcd0=http://127.0.0.1:2380 \
    --initial-cluster-state new \
    --initial-cluster-token sokuim-etcd \
    >/dev/null 2>&1 &
  ETCD_PID=$!

  wait_ready

  echo "[etcd] bootstrapping authentication (root/${ROOT_PASSWORD})"
  echo "$ROOT_PASSWORD" | etcdctl --endpoints=http://127.0.0.1:2379 user add root --interactive=false
  etcdctl --endpoints=http://127.0.0.1:2379 user grant-role root root
  etcdctl --endpoints=http://127.0.0.1:2379 auth enable

  echo "[etcd] auth enabled, restarting in foreground"
  kill "$ETCD_PID" 2>/dev/null || true
  wait "$ETCD_PID" 2>/dev/null || true
fi

echo "[etcd] starting in foreground (initial-cluster-state=existing)"
exec etcd --name etcd0 \
  --data-dir "$DATA_DIR" \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls "$ADVERTISE_URL" \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://127.0.0.1:2380 \
  --initial-cluster etcd0=http://127.0.0.1:2380 \
  --initial-cluster-state existing \
  --initial-cluster-token sokuim-etcd

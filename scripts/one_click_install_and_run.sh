#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/deploy/docker-compose.yml"
LITE_COMPOSE_FILE="$ROOT_DIR/deploy/docker-compose.lite.yml"

MODE="full"
if [[ "${1:-}" == "--lite" ]]; then
  MODE="lite"
fi

log() {
  printf '\n[HyperCompute] %s\n' "$1"
}

warn() {
  printf '\n[HyperCompute][WARN] %s\n' "$1"
}

require_command() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    return 1
  fi
  return 0
}

run_as_root() {
  if [[ "${EUID:-$(id -u)}" -eq 0 ]]; then
    "$@"
  elif require_command sudo; then
    sudo "$@"
  else
    warn "当前用户不是 root，且未安装 sudo，无法自动安装 Docker。"
    return 1
  fi
}

install_docker_linux() {
  if require_command apt-get; then
    log "检测到 Debian/Ubuntu，尝试自动安装 Docker（需要 sudo 权限）"
    run_as_root apt-get update
    run_as_root apt-get install -y ca-certificates curl gnupg lsb-release
    run_as_root install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | run_as_root gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    run_as_root chmod a+r /etc/apt/keyrings/docker.gpg
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      run_as_root tee /etc/apt/sources.list.d/docker.list >/dev/null
    run_as_root apt-get update
    run_as_root apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    run_as_root usermod -aG docker "$USER" || true
    return 0
  fi

  if require_command yum; then
    log "检测到 RHEL/CentOS，尝试自动安装 Docker（需要 sudo 权限）"
    run_as_root yum install -y yum-utils
    run_as_root yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    run_as_root yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    run_as_root systemctl enable --now docker
    run_as_root usermod -aG docker "$USER" || true
    return 0
  fi

  return 1
}

install_docker_macos() {
  if require_command brew; then
    log "检测到 macOS + Homebrew，安装 docker 与 docker-compose"
    brew install --cask docker
    brew install docker-compose
    log "请手动启动 Docker Desktop 后，重新执行该脚本。"
    exit 0
  fi
  return 1
}

ensure_docker() {
  if require_command docker; then
    return 0
  fi

  log "未检测到 Docker，开始自动安装。"
  local os
  os="$(uname -s)"

  if [[ "$os" == "Linux" ]]; then
    install_docker_linux || {
      echo "自动安装 Docker 失败，请手动安装后重试：https://docs.docker.com/engine/install/"
      exit 1
    }
  elif [[ "$os" == "Darwin" ]]; then
    install_docker_macos || {
      echo "自动安装 Docker 失败，请先安装 Homebrew 或 Docker Desktop。"
      exit 1
    }
  else
    echo "暂不支持自动安装 Docker 的系统：$os"
    exit 1
  fi

  if ! require_command docker; then
    echo "Docker 安装步骤已执行，但仍未检测到 docker 命令。请重新打开终端后重试。"
    exit 1
  fi
}

ensure_compose() {
  if docker compose version >/dev/null 2>&1; then
    return 0
  fi

  if require_command docker-compose; then
    alias docker_compose='docker-compose'
    return 0
  fi

  echo "未检测到 docker compose，请确认 Docker Compose 已安装。"
  exit 1
}

start_docker_service_if_needed() {
  if docker info >/dev/null 2>&1; then
    return 0
  fi

  if require_command systemctl; then
    log "Docker 守护进程未运行，尝试启动。"
    run_as_root systemctl start docker || true
  fi

  if ! docker info >/dev/null 2>&1; then
    echo "Docker 服务未就绪。请启动 Docker 后重试。"
    exit 1
  fi
}

run_stack() {
  local compose="$COMPOSE_FILE"
  if [[ "$MODE" == "lite" ]]; then
    compose="$LITE_COMPOSE_FILE"
  fi

  log "执行项目初始化（依赖拉取、构建等）"
  bash "$ROOT_DIR/deploy/bootstrap.sh"

  log "启动系统服务（模式：$MODE）"
  docker compose -f "$compose" up -d --build

  log "等待数据库启动"
  bash "$ROOT_DIR/scripts/wait_for_postgres.sh"

  log "执行数据库迁移"
  bash "$ROOT_DIR/scripts/migrate_all.sh"

  log "执行初始化数据"
  bash "$ROOT_DIR/scripts/seed_data.sh"

  log "执行健康检查"
  bash "$ROOT_DIR/scripts/health_check.sh"

  log "✅ 环境安装与系统部署完成。"
  log "查看日志: make logs"
  log "停止系统: make down"
}

main() {
  log "开始一键安装并运行 HyperCompute"
  ensure_docker
  ensure_compose
  start_docker_service_if_needed
  run_stack
}

main "$@"

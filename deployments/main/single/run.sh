#!/bin/bash

# 网络名称
NETWORK_NAME="cowim-network"

# Compose 文件路径
COMPOSE_FILE1="./docker-compose.yaml"
COMPOSE_FILE2="../../base/docker-compose.yaml"

# 检测 docker-compose 或 podman-compose
# 如果同时安装，默认podman-compose
detect_compose_tool() {
  if command -v podman-compose &>/dev/null; then
    COMPOSE_TOOL="podman-compose"
  elif command -v docker-compose &>/dev/null; then
    COMPOSE_TOOL="docker-compose"
  else
    echo "Error: Neither docker-compose nor podman-compose is installed. Please install one of them."
    exit 1
  fi
  echo "Using $COMPOSE_TOOL"
}

# 创建共享网络
create_network() {
  echo "Creating network: $NETWORK_NAME"
  docker network create --driver bridge --subnet=10.0.0.0/24 "$NETWORK_NAME" || {
    echo "Network $NETWORK_NAME already exists."
  }
}

# 启动 Docker Compose 项目
start_compose() {
  echo "Starting Docker Compose projects..."
  $COMPOSE_TOOL -f "$COMPOSE_FILE1" up -d
  $COMPOSE_TOOL -f "$COMPOSE_FILE2" up -d
}

# 停止 Docker Compose 项目
stop_compose() {
  echo "Stopping Docker Compose projects..."
  $COMPOSE_TOOL -f "$COMPOSE_FILE1" down
  $COMPOSE_TOOL -f "$COMPOSE_FILE2" down
}

# 删除共享网络
delete_network() {
  echo "Deleting network: $NETWORK_NAME"
  docker network rm "$NETWORK_NAME" || {
    echo "Network $NETWORK_NAME does not exist or is in use."
  }
}

# 主流程
case "$1" in
  up)
    detect_compose_tool
    create_network
    start_compose
    ;;
  down)
    detect_compose_tool
    stop_compose
    delete_network
    ;;
  *)
    echo "Usage: \$0 {up|down}"
    ;;
esac

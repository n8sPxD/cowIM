#!/bin/bash

# 网络名称
NETWORK_NAME="cowim-network"

# Compose 文件路径
COMPOSE_FILE1="./docker-compose.yaml"
COMPOSE_FILE2="../../base/docker-compose.yaml"

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
  podman-compose -f "$COMPOSE_FILE1" up -d
  podman-compose -f "$COMPOSE_FILE2" up -d
}

# 停止 Docker Compose 项目
stop_compose() {
  echo "Stopping Docker Compose projects..."
  podman-compose -f "$COMPOSE_FILE1" down
  podman-compose -f "$COMPOSE_FILE2" down
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
    create_network
    start_compose
    ;;
  down)
    stop_compose
    delete_network
    ;;
  *)
    echo "Usage: \$0 {up|down}"
    ;;
esac

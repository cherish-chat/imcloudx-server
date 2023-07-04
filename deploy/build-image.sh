#!/bin/zsh

#services=(gateway app)
services=(gateway app)

IMAGE_REPO_NAME_PREFIX="registry.cn-shanghai.aliyuncs.com/xxim-dev/imcloudx-"

# 定义writeDockerFile方法
function writeDockerFile() {
service=$1
cat > Dockerfile << EOF
FROM debian:latest

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

# 创建 /app/etc文件夹
WORKDIR /app/etc

# 复制二进制文件
COPY $service /app/$service

# 工作目录切换到 /app
WORKDIR /app

# 映射 /app/etc/config.yaml 到宿主机
VOLUME /app/etc/config.yaml

# 启动服务
CMD ["/app/$service"]

EOF
}

# shellcheck disable=SC2034
SCRIPT_DIR=$(pwd)
TAG=$(date +%Y%m%d%H%M%S)
echo "docker TAG = ${TAG}"
cd ../ && echo "进入项目根目录"
ROOT_DIR=$(pwd)
echo "项目根目录为：${ROOT_DIR}"

# shellcheck disable=SC2128
for service in $services; do
  cd "${ROOT_DIR}/app/$service" && echo "进入$service目录"
  GOOS=linux GOARCH=amd64 go build -o "$SCRIPT_DIR/tmp/$service/$service" "./$service.go"
  cd "$SCRIPT_DIR/tmp/$service" && echo "进入$service build目录"
  writeDockerFile "$service"
  docker build --platform linux/x86_64 -t "${IMAGE_REPO_NAME_PREFIX}${service}:${TAG}" .
  docker tag "${IMAGE_REPO_NAME_PREFIX}${service}:${TAG}" "${IMAGE_REPO_NAME_PREFIX}${service}:latest"
  docker push "${IMAGE_REPO_NAME_PREFIX}${service}:${TAG}"
  docker push "${IMAGE_REPO_NAME_PREFIX}${service}:latest"
  cd "$SCRIPT_DIR" && echo "进入脚本目录"
done

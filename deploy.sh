#!/bin/bash
# =============================================================
# XATU 二手书交易平台 - Docker 一键部署脚本
# 使用方法：
#   1. 首次部署： bash deploy.sh
#   2. 更新代码： bash deploy.sh update
#   3. 查看状态： bash deploy.sh status
#   4. 停止服务： bash deploy.sh stop
# =============================================================

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &>/dev/null; then
        error "Docker 未安装！请先安装 Docker"
        echo ""
        echo "快速安装 Docker（Ubuntu）："
        echo "  curl -fsSL https://get.docker.com | sh"
        exit 1
    fi

    if ! command -v docker compose &>/dev/null; then
        error "Docker Compose 未安装！"
        echo ""
        echo "快速安装："
        echo "  apt install -y docker-compose-plugin"
        exit 1
    fi
    info "Docker 环境检查通过"
}

# 检查端口是否被占用
check_port() {
    local port=$1
    if ss -tlnp 2>/dev/null | grep -q ":$port "; then
        warn "端口 $port 已被占用，尝试释放..."
        fuser -k "$port/tcp" 2>/dev/null || true
        sleep 2
    fi
    info "端口 $port 可用"
}

# 首次部署
do_deploy() {
    echo "================================================"
    echo "   📚 XATU 二手书交易平台 - Docker 部署"
    echo "================================================"
    echo ""

    check_docker
    check_port 80

    # 确保脚本在项目根目录执行
    cd "$(dirname "$0")"

    # 拉取最新镜像并构建
    info "正在构建 Docker 镜像..."
    docker compose build --pull

    # 检查是否有旧容器在运行
    if docker compose ps --status running 2>/dev/null | grep -q "xatu"; then
        warn "检测到已有容器运行，正在重新创建..."
        docker compose down
    fi

    # 启动所有服务
    info "正在启动 MySQL + Redis + 应用..."
    docker compose up -d

    echo ""
    info "✅ 部署完成！服务已启动："
    echo ""
    echo "   🌐 访问地址：http://$(curl -s ifconfig.me 2>/dev/null || echo 'localhost')"
    echo "   📋 管理员账号：19992468036 / G20050111g"
    echo ""
    echo "   常用命令："
    echo "     docker compose logs app   查看应用日志"
    echo "     docker compose logs mysql 查看数据库日志"
    echo "     bash deploy.sh stop       停止服务"
    echo "     bash deploy.sh status     查看运行状态"
    echo ""
}

# 更新代码并重新部署
do_update() {
    info "正在更新代码..."
    cd "$(dirname "$0")"

    # 如果有 Git，拉取最新代码
    if git status &>/dev/null; then
        git pull
    else
        warn "非 Git 仓库，跳过代码拉取"
    fi

    info "重新构建并启动..."
    docker compose down
    docker compose build --no-cache
    docker compose up -d

    info "✅ 更新完成！"
    docker compose ps
}

# 查看状态
do_status() {
    cd "$(dirname "$0")"
    echo "================================================"
    echo "   📊 XATU 服务运行状态"
    echo "================================================"
    docker compose ps
    echo ""
    echo "--- 资源占用 ---"
    docker stats --no-stream $(docker compose ps -q 2>/dev/null) 2>/dev/null || echo "(暂无运行中的容器)"
}

# 停止服务
do_stop() {
    echo ""
    warn "正在停止所有服务..."
    cd "$(dirname "$0")"
    docker compose down
    info "✅ 已停止"
    echo ""
}

# ====== 主入口 ======
case "${1:-deploy}" in
    deploy)
        do_deploy
        ;;
    update)
        do_update
        ;;
    status)
        do_status
        ;;
    stop)
        do_stop
        ;;
    restart)
        do_stop
        do_deploy
        ;;
    *)
        echo "用法: bash deploy.sh [deploy|update|status|stop|restart]"
        echo ""
        echo "   deploy  首次部署（默认）"
        echo "   update  更新代码后重新部署"
        echo "   status  查看服务状态"
        echo "   stop    停止所有服务"
        echo "   restart 重启所有服务"
        ;;
esac

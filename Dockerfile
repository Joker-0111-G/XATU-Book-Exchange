# ========== 构建阶段 ==========
FROM golang:1.25-alpine AS builder

# 安装 git（go mod 需要）
RUN apk add --no-cache git

WORKDIR /build

# 先复制 go.mod/go.sum 充分利用 Docker 缓存
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制后端源码
COPY backend/ .

# 静态编译（alpine 用小即可，关键是要 CGO_ENABLED=0）
RUN CGO_ENABLED=0 GOOS=linux go build -o xatu-server .

# ========== 运行阶段 ==========
FROM alpine:3.21

# 安装 ca-certificates（HTTPS 调用需要）+ tzdata（时区）
RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/xatu-server .

# 复制配置文件
COPY backend/config/config.docker.yaml ./config/config.yaml

# 复制前端静态文件
COPY frontend/ /frontend/

# 创建上传目录
RUN mkdir -p ./uploads/images

# 暴露端口
EXPOSE 8080

# 启动
CMD ["./xatu-server"]

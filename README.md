<div align="center">

# 📚 XATU 二手书交易平台

**西安工业大学 · 校园二手书线上交易平台**

让知识循环，让校园更环保

[🌐 在线演示](http://47.96.255.177) · 
[📖 设计文档](DESIGN.md) ·
[🐳 Docker 部署](#docker-一键部署)

---

</div>

## ✨ 功能特性

### 🏪 用户端
| 功能 | 说明 |
|------|------|
| **用户注册/登录** | 手机号注册，JWT 认证，Token 自动刷新 |
| **首页浏览** | 轮播图、分类筛选、价格/新旧度筛选、搜索、分页 |
| **图书详情** | 多图轮播、卖家信息、状态展示 |
| **发布图书** | 选择分类、上传图片、填写描述，所见即所得 |
| **购买图书** | 填写联系方式，事务保护防并发抢购 |
| **收藏管理** | 添加/取消收藏，收藏列表查看 |
| **订单管理** | 买家订单（取消/完成）+ 卖家订单（确认） |
| **站内消息** | 会话列表 + 实时聊天界面 |
| **个人中心** | 查看/编辑资料、修改密码 |

### 🔧 管理端
| 功能 | 说明 |
|------|------|
| **数据统计** | 用户数、图书数、订单数、完成交易数 |
| **图书管理** | 查看/上架/下架所有图书 |
| **分类管理** | 添加/删除分类（含子分类检查） |
| **用户管理** | 启用/禁用用户 |

---

## 🏗️ 技术栈

| 层级 | 技术 | 用途 |
|------|------|------|
| **前端** | HTML / CSS / Vanilla JS | 纯前端单页应用，无框架依赖 |
| **后端** | **Go** + **Gin** | RESTful API 服务 |
| **ORM** | **GORM** | MySQL 数据访问 |
| **数据库** | **MySQL 8.0** | 持久化存储 |
| **缓存** | **Redis 7** | JWT 黑名单（登出）、接口限流 |
| **认证** | **JWT** (golang-jwt) | 无状态 Token 认证 |
| **部署** | **Docker** + **Docker Compose** | 一键容器化部署 |

---

## 🗂️ 项目结构

```
XATU-Book-Exchange/
├── backend/                    # Go 后端
│   ├── main.go                 # 程序入口
│   ├── config/
│   │   ├── config.go           # 配置加载
│   │   └── config.yaml         # 配置文件
│   ├── common/
│   │   ├── response.go         # 统一响应封装
│   │   ├── error_code.go       # 业务错误码
│   │   └── helper.go           # 辅助函数
│   ├── middleware/
│   │   ├── auth.go             # JWT 认证 + 管理员中间件
│   │   ├── cors.go             # 跨域处理
│   │   ├── logger.go           # 请求日志
│   │   └── rate_limit.go       # 内存限流
│   ├── handler/                # 控制器层
│   ├── service/                # 业务逻辑层
│   ├── repository/             # 数据访问层
│   ├── model/                  # 数据模型
│   ├── routes/
│   │   └── router.go           # 路由注册
│   ├── utils/
│   │   ├── jwt.go              # JWT 签发/解析/黑名单
│   │   ├── hash.go             # bcrypt 密码加密
│   │   └── validator.go        # 参数校验器
│   └── database/
│       ├── mysql.go            # MySQL 初始化
│       ├── redis.go            # Redis 初始化
│       └── seed.go             # 种子数据
├── frontend/
│   └── index.html              # 前端单页应用（纯 HTML）
├── docs/
│   ├── api.md                  # API 文档
│   └── schema.sql              # 数据库建表 SQL
├── Dockerfile                  # 多阶段构建镜像
├── docker-compose.yml          # Docker 编排
├── deploy.sh                   # 一键部署脚本
└── DESIGN.md                   # 详细设计文档
```

---

## 📦 本地开发

### 环境要求

- Go 1.25+
- MySQL 8.0+
- Redis 7+

### 启动步骤

```bash
# 1. 克隆仓库
git clone git@github.com:JokerGMJ/XATU-Book-Exchange.git
cd XATU-Book-Exchange

# 2. 创建数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS xatu_book_exchange DEFAULT CHARACTER SET utf8mb4;"

# 3. 修改配置（数据库密码等）
vim backend/config/config.yaml

# 4. 启动后端
cd backend
go run .
```

启动后访问 http://localhost:8080

> 首次启动时 `seed.go` 会自动创建管理员和分类数据。

---

## 🐳 Docker 一键部署

### 部署到服务器

```bash
# 1. 本地打包代码
tar czf xatu.tar.gz --exclude=.git .
scp xatu.tar.gz root@你的服务器IP:/opt/

# 2. SSH 到服务器
ssh root@你的服务器IP

# 3. 解压并部署
cd /opt && tar xzf xatu.tar.gz
cd XATU-Book-Exchange
bash deploy.sh
```

### 常用命令

```bash
bash deploy.sh          # 首次部署
bash deploy.sh status   # 查看运行状态
bash deploy.sh stop     # 停止服务
bash deploy.sh restart  # 重启服务
bash deploy.sh update   # 更新代码后重新部署

# 或者使用原生命令
docker compose logs app     # 查看应用日志
docker compose logs mysql   # 查看数据库日志
```

### 访问

部署完成后浏览器打开 `http://你的服务器IP`

---

## 🔑 默认管理员

| 角色 | 手机号 | 密码 | 说明 |
|------|--------|------|------|
| **管理员** | `1******` | `*******` | 拥有全部管理权限 |
| **普通用户** | 注册获取 | 注册设置 | 可发布图书、购买、聊天 |

---

## 📸 界面预览

### 首页
- 渐变轮播图，自动切换
- 实时统计（图书在售数、注册用户数、总订单）
- 一级分类快速筛选
- 图书网格展示，支持排序/筛选/搜索/分页

### 图书发布
- 下拉选择分类（层级化展示）
- 图片上传（JPG/PNG/WebP，≤5MB）
- 预览缩略图，可删除重传

### 站内信
- 会话列表 + 实时聊天
- 未读消息标记
- 支持图片/文字消息

---

## 🔐 安全特性

- ✅ JWT 双 Token（access_token 2h + refresh_token 7d）
- ✅ Redis 黑名单实现登出
- ✅ bcrypt 密码哈希（cost=10）
- ✅ 统一错误码 + 正确 HTTP 状态码（401/403/404/409）
- ✅ 上传文件类型/MIME 双重校验
- ✅ 防目录遍历路径检查
- ✅ 接口限流（内存令牌桶）
- ✅ 订单事务保护（防并发抢购）
- ✅ 删除分类前检查关联

---

## 📄 许可证

本项目仅供学习交流使用。

---

<div align="center">

**西安工业大学 · 二手书交易平台**

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev)
[![Gin](https://img.shields.io/badge/Gin-1.12-008ECF?logo=gin)](https://gin-gonic.com)
[![Docker](https://img.shields.io/badge/Docker-✓-2496ED?logo=docker)](https://docker.com)
[![License](https://img.shields.io/badge/License-MIT-green)]()

</div>

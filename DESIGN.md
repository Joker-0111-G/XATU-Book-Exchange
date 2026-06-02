# 西安工业大学 二手书交易平台 — 详细设计报告

> 项目：XATU-Book-Exchange  
> 技术栈：Go + Gin + MySQL + Redis + HTML  
> 撰写日期：2026-06-02

---

## 目录

1. [系统架构](#一系统架构)
2. [后端项目结构](#二后端项目结构)
3. [数据库设计](#三数据库设计-mysql)
4. [Redis 缓存策略](#四redis-缓存策略)
5. [API 路由设计](#五api-路由设计-gin)
6. [统一响应与错误码](#六统一响应格式)
7. [技术要点](#七技术要点)
8. [实施计划](#八实施计划分阶段)

---

## 一、系统架构

```
┌─────────────────────────────────────────────────────┐
│                  前端 (Vue 3)                        │
│           Element Plus / Vant UI                     │
├─────────────────────────────────────────────────────┤
│               Nginx (反向代理)                        │
├─────────────────────────────────────────────────────┤
│              后端 (Go + Gin)                         │
│   ├── 路由层 (routes)                                │
│   ├── 控制器层 (handlers)                            │
│   ├── 业务逻辑层 (services)                          │
│   ├── 数据访问层 (repositories)                      │
│   └── 中间件层 (middlewares)                         │
├──────────────────┬──────────────────────────────────┤
│   MySQL (持久化)   │   Redis (缓存/会话)               │
└──────────────────┴──────────────────────────────────┘
```

### 设计原则

| 原则 | 说明 |
|------|------|
| **分层架构** | handler → service → repository，每一层职责清晰，单向依赖 |
| **RESTful API** | 资源导向的 URL 设计，语义化 HTTP 方法 |
| **JWT 认证** | 无状态 Token，Redis 管理黑名单实现登出 |
| **统一错误响应** | 全局错误码 + 统一 JSON 格式，前端只需按 code 处理 |
| **软删除** | 所有核心表使用 `deleted_at` 进行逻辑删除 |

---

## 二、后端项目结构

```
XATU-Book-Exchange/
├── backend/
│   ├── main.go                      # 程序入口
│   ├── go.mod / go.sum
│   ├── config/
│   │   ├── config.go                # 配置结构体 + 加载逻辑
│   │   └── config.yaml              # 配置文件（MySQL、Redis、JWT密钥等）
│   ├── common/
│   │   ├── response.go              # 统一 JSON 响应封装
│   │   └── error_code.go            # 业务错误码定义
│   ├── middleware/
│   │   ├── auth.go                  # JWT 鉴权中间件
│   │   ├── cors.go                  # 跨域处理
│   │   ├── logger.go               # 请求日志
│   │   └── rate_limit.go           # 接口限流
│   ├── model/
│   │   ├── user.go
│   │   ├── book.go
│   │   ├── category.go
│   │   ├── order.go
│   │   ├── favorite.go
│   │   ├── message.go
│   │   └── banner.go
│   ├── handler/                     # 控制器：接收请求、参数校验、调用service
│   │   ├── user_handler.go
│   │   ├── book_handler.go
│   │   ├── category_handler.go
│   │   ├── order_handler.go
│   │   ├── favorite_handler.go
│   │   ├── message_handler.go
│   │   ├── upload_handler.go
│   │   └── admin_handler.go
│   ├── service/                      # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── book_service.go
│   │   ├── order_service.go
│   │   ├── favorite_service.go
│   │   └── message_service.go
│   ├── repository/                   # 数据访问层（gorm 查询封装）
│   │   ├── user_repo.go
│   │   ├── book_repo.go
│   │   ├── order_repo.go
│   │   ├── favorite_repo.go
│   │   └── message_repo.go
│   ├── database/
│   │   ├── mysql.go                  # 初始化 GORM + MySQL 连接
│   │   └── redis.go                  # 初始化 go-redis 连接
│   ├── routes/
│   │   └── router.go                 # 统一路由注册
│   └── utils/
│       ├── jwt.go                    # JWT 签发与解析
│       ├── hash.go                   # bcrypt 密码加密
│       ├── upload.go                 # 文件上传处理
│       └── validator.go             # 自定义参数校验器
├── frontend/                          # Vue 3 前端（后续开发）
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── docs/
    ├── api.md                        # API 文档
    └── schema.sql                    # 数据库建表脚本
```

---

## 三、数据库设计 (MySQL)

### 3.1 全局约定

| 项目 | 约定 |
|------|------|
| 引擎 | InnoDB |
| 字符集 | utf8mb4 |
| 排序规则 | utf8mb4_unicode_ci |
| 主键 | BIGINT UNSIGNED AUTO_INCREMENT |
| 时间字段 | DATETIME，精确到秒 |
| 软删除 | 统一使用 `deleted_at DATETIME NULLABLE` |

### 3.2 ER 关系

```
users ──1:N──> books      用户发布多本图书
users ──1:N──> orders (buyer)   用户购买订单
users ──1:N──> orders (seller)  用户出售订单
users ──1:N──> favorites        用户收藏
users ──1:N──> messages (from)  用户发送消息
users ──1:N──> messages (to)    用户接收消息
books ──N:1──> categories       图书属于一个分类
books ──1:N──> favorites        图书被收藏
books ──1:N──> orders           图书被下单
```

### 3.3 表结构

#### ① `users` — 用户表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK AI | 用户ID |
| phone | VARCHAR(20) | UNIQUE NOT NULL | **手机号（登录凭证）** |
| password_hash | VARCHAR(255) | NOT NULL | bcrypt 加密后的密码 |
| nickname | VARCHAR(50) | NOT NULL | 昵称 |
| avatar | VARCHAR(255) | DEFAULT '' | 头像 URL |
| email | VARCHAR(100) | DEFAULT '' | 邮箱 |
| major | VARCHAR(100) | DEFAULT '' | **专业** |
| wechat | VARCHAR(50) | DEFAULT '' | **微信号（交易联系方式）** |
| status | TINYINT | DEFAULT 1, INDEX | 0=禁用 1=正常 |
| is_admin | TINYINT | DEFAULT 0 | 0=普通用户 1=管理员 |
| created_at | DATETIME | NOT NULL | |
| updated_at | DATETIME | NOT NULL | |
| deleted_at | DATETIME | INDEX | 软删除 |

> **索引**：`phone`(UNIQUE), `status`

---

#### ② `categories` — 图书分类表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK AI | |
| name | VARCHAR(50) | NOT NULL | 分类名称 |
| parent_id | BIGINT UNSIGNED | DEFAULT 0, INDEX | 父分类ID（0=顶级） |
| icon | VARCHAR(255) | DEFAULT '' | 图标 URL |
| sort_order | INT | DEFAULT 0 | 排序编号 |
| created_at | DATETIME | NOT NULL | |

> **预设分类结构**（二级，`parent_id` 实现层级）：
>
> **一级：专业大类**
> | 分类 | 二级子类参考 |
> |------|-------------|
> | 计算机科学与工程学院 | 程序设计、数据结构、操作系统、计算机网络… |
> | 机电工程学院 | 机械制图、理论力学、材料力学… |
> | 电子信息工程学院 | 电路分析、模拟电子、数字电子… |
> | 经济管理学院 | 管理学、微观经济学、会计学… |
> | 建筑工程学院 | 土木工程材料、结构力学… |
> | 外国语学院 | 英美文学、翻译、语言学… |
> | 理学院 | 数学分析、高等代数、大学物理… |
> | 材料与化工学院 | 材料科学基础、物理化学… |
> | 艺术与传媒学院 | 设计基础、数字媒体… |
> | 马克思主义学院 | 思政类公共课 |
> | 公共课（跨专业） | 高等数学、大学英语、线性代数… |
>
> **一级：考研专区 ✦**
> | 二级子类 | 说明 |
> |---------|------|
> | 考研-数学 | 数一/数二/数三资料 |
> | 考研-英语 | 考研英语词汇/真题/作文 |
> | 考研-政治 | 肖秀荣、徐涛等 |
> | 考研-专业课 | 各校专业课真题/笔记 |
>
> **一级：其他**
> | 二级子类 | 说明 |
> |---------|------|
> | 考证教材 | 四六级、计算机二级、教资、CPA… |
> | 课外书 | 小说、科普、人文社科… |

---

#### ③ `books` — 图书表（核心表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK AI | |
| user_id | BIGINT UNSIGNED | FK NOT NULL INDEX | **卖家ID** |
| category_id | BIGINT UNSIGNED | FK NOT NULL INDEX | 分类ID |
| title | VARCHAR(200) | NOT NULL | **书名** |
| author | VARCHAR(100) | DEFAULT '' | 作者 |
| selling_price | DECIMAL(10,2) | NOT NULL | **售价** |
| `condition` | TINYINT | NOT NULL | **新旧程度（1-10）** |
| description | TEXT | | 描述/笔记情况 |
| images | JSON | | **图片URL数组** |
| status | TINYINT | DEFAULT 1, INDEX | 0=下架 1=**在售** 2=已售 |
| created_at | DATETIME | NOT NULL | |
| updated_at | DATETIME | NOT NULL | |
| deleted_at | DATETIME | INDEX | 软删除 |

> **索引**：`user_id`, `category_id`, `status`, `title`
>
> **外键**：`user_id` → users(id), `category_id` → categories(id)

---

#### ④ `favorites` — 收藏表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK AI | |
| user_id | BIGINT UNSIGNED | FK NOT NULL | 用户ID |
| book_id | BIGINT UNSIGNED | FK NOT NULL | 图书ID |
| created_at | DATETIME | NOT NULL | |

> **索引**：`(user_id, book_id)` **UNIQUE**（防止重复收藏）, `user_id`, `book_id`

---

#### ⑤ `orders` — 订单表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK AI | |
| order_no | VARCHAR(32) | **UNIQUE** NOT NULL | 订单号 |
| seller_id | BIGINT UNSIGNED | FK NOT NULL | 卖家ID |
| buyer_id | BIGINT UNSIGNED | FK NOT NULL | 买家ID |
| book_id | BIGINT UNSIGNED | FK NOT NULL | 图书ID |
| price | DECIMAL(10,2) | NOT NULL | 成交价 |
| status | TINYINT | DEFAULT 0, INDEX | **0=待确认 1=已确认 2=已完成 3=已取消** |
| contact_phone | VARCHAR(20) | DEFAULT '' | 联系电话 |
| contact_wechat | VARCHAR(50) | DEFAULT '' | 联系微信 |
| note | VARCHAR(500) | DEFAULT '' | 买家备注 |
| completed_at | DATETIME | NULLABLE | 完成时间 |
| created_at | DATETIME | NOT NULL | |
| updated_at | DATETIME | NOT NULL | |

> **索引**：`order_no`(UNIQUE), `seller_id`, `buyer_id`, `book_id`, `status`
>
> **订单号生成**：前缀 `XATU` + 时间戳(13位) + 4位随机数，共20位左右

---

#### ⑥ `messages` — 站内消息/聊天表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK AI | |
| from_user_id | BIGINT UNSIGNED | FK NOT NULL | 发送者 |
| to_user_id | BIGINT UNSIGNED | FK NOT NULL | 接收者 |
| book_id | BIGINT UNSIGNED | DEFAULT 0 | **关联图书（上下文）** |
| content | TEXT | NOT NULL | 消息内容 |
| content_type | TINYINT | DEFAULT 0 | 0=文字 1=图片 2=系统通知 |
| is_read | TINYINT | DEFAULT 0, INDEX | 0=未读 1=已读 |
| created_at | DATETIME | NOT NULL | |

> **索引**：`to_user_id`（查收件箱）, `(from_user_id, to_user_id)`, `is_read`

---

#### ⑦ `banners` — 首页轮播图表

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK AI | |
| title | VARCHAR(100) | DEFAULT '' | 标题 |
| image_url | VARCHAR(255) | NOT NULL | 图片URL |
| link_url | VARCHAR(255) | DEFAULT '' | 跳转链接 |
| sort_order | INT | DEFAULT 0 | 排序 |
| is_active | TINYINT | DEFAULT 1 | 0=隐藏 1=显示 |
| created_at | DATETIME | NOT NULL | |

---

## 四、Redis 缓存策略

### 4.1 缓存项

| Key 模式 | 数据类型 | TTL | 用途 |
|----------|---------|-----|------|
| `token:blacklist:{jti}` | string | JWT剩余有效期 | **JWT黑名单**，登出/封禁时加入 |
| `user:{id}:profile` | string (JSON) | 1h | 用户信息缓存 |
| `book:{id}:detail` | string (JSON) | 30min | **图书详情**热点缓存 |
| `book:hot:list` | string (JSON) | 5min | 首页热门图书列表 |
| `category:tree` | string (JSON) | 1h | 全部分类树 |
| `banner:list` | string (JSON) | 1h | 轮播图列表 |
| `rate:limit:{ip}:{api}` | string (计数器) | 1s~1min | **接口限流** |

### 4.2 缓存策略

| 策略 | 说明 |
|------|------|
| **被动过期** | 设置 TTL，到期自动淘汰，下次请求回源 DB 重建 |
| **主动淘汰** | 图书更新/删除时，立即删除 `book:{id}:detail` |
| **延迟双删** | 写操作：删缓存 → 更新 DB → 延迟再删（防并发脏读） |
| **缓存预热** | 启动时可加载分类树和轮播图到 Redis |

---

## 五、API 路由设计 (Gin)

### 5.1 路由分组

| 分组 | 路由前缀 | 中间件 | 说明 |
|------|---------|--------|------|
| Public | `/api/v1/public` | 无 | 无需登录的公开接口 |
| User | `/api/v1/user` | auth | 个人中心 |
| Book | `/api/v1/books` | auth | 图书管理（需登录） |
| Order | `/api/v1/orders` | auth | 订单管理 |
| Favorite | `/api/v1/favorites` | auth | 收藏管理 |
| Message | `/api/v1/messages` | auth | 站内消息 |
| Upload | `/api/v1/upload` | auth | 文件上传 |
| Admin | `/api/v1/admin` | auth + admin | 管理员接口 |

### 5.2 接口详细清单

#### 公开接口（无需登录）

| 方法 | 路径 | handler | 说明 |
|------|------|---------|------|
| POST | `/api/v1/public/register` | UserHandler.Register | 用户注册 |
| POST | `/api/v1/public/login` | UserHandler.Login | 用户登录，返回 JWT |
| POST | `/api/v1/public/refresh-token` | UserHandler.RefreshToken | 刷新 Token |
| GET | `/api/v1/public/books` | BookHandler.List | **图书列表（分页+筛选+排序）** |
| GET | `/api/v1/public/books/:id` | BookHandler.Get | 图书详情 |
| GET | `/api/v1/public/books/search` | BookHandler.Search | **搜索图书** |
| GET | `/api/v1/public/categories` | CategoryHandler.List | 分类列表 |
| GET | `/api/v1/public/banners` | BannerHandler.List | 轮播图 |

> **图书列表筛选参数**：`?category_id=&condition=&min_price=&max_price=&sort=created_at|price&order=asc|desc&page=1&page_size=20`

---

#### 用户接口

| 方法 | 路径 | handler | 说明 |
|------|------|---------|------|
| GET | `/api/v1/user/profile` | UserHandler.Profile | 获取个人信息 |
| PUT | `/api/v1/user/profile` | UserHandler.UpdateProfile | 修改个人信息 |
| PUT | `/api/v1/user/password` | UserHandler.ChangePassword | 修改密码 |
| GET | `/api/v1/user/books` | BookHandler.UserBooks | 我发布的图书 |
| GET | `/api/v1/user/orders` | OrderHandler.UserOrders | 我买的订单（买家视角） |
| GET | `/api/v1/user/sales` | OrderHandler.UserSales | 我的售出（卖家视角） |

---

#### 图书接口（需登录）

| 方法 | 路径 | handler | 说明 |
|------|------|---------|------|
| POST | `/api/v1/books` | BookHandler.Create | **发布图书** |
| PUT | `/api/v1/books/:id` | BookHandler.Update | 编辑图书 |
| PUT | `/api/v1/books/:id/status` | BookHandler.UpdateStatus | 上架/下架 |
| DELETE | `/api/v1/books/:id` | BookHandler.Delete | 删除图书（软删除） |

---

#### 收藏接口

| 方法 | 路径 | handler | 说明 |
|------|------|---------|------|
| GET | `/api/v1/favorites` | FavoriteHandler.List | 收藏列表 |
| POST | `/api/v1/favorites` | FavoriteHandler.Add | **添加收藏** |
| DELETE | `/api/v1/favorites/:bookId` | FavoriteHandler.Remove | 取消收藏 |
| GET | `/api/v1/favorites/check/:bookId` | FavoriteHandler.Check | 检查是否已收藏 |

---

#### 订单接口

| 方法 | 路径 | handler | 说明 |
|------|------|---------|------|
| POST | `/api/v1/orders` | OrderHandler.Create | **创建订单（购买）** |
| GET | `/api/v1/orders` | OrderHandler.List | 订单列表 |
| GET | `/api/v1/orders/sales` | OrderHandler.SalesList | 卖出的订单列表 |
| GET | `/api/v1/orders/:id` | OrderHandler.Get | 订单详情 |
| PUT | `/api/v1/orders/:id/confirm` | OrderHandler.Confirm | **卖家确认交易** |
| PUT | `/api/v1/orders/:id/complete` | OrderHandler.Complete | 确认完成 |
| PUT | `/api/v1/orders/:id/cancel` | OrderHandler.Cancel | 取消订单 |

---

#### 消息接口

| 方法 | 路径 | handler | 说明 |
|------|------|---------|------|
| GET | `/api/v1/messages/conversations` | MessageHandler.Conversations | **会话列表** |
| GET | `/api/v1/messages/conversations/:userId` | MessageHandler.History | 聊天记录 |
| POST | `/api/v1/messages` | MessageHandler.Send | **发送消息** |
| PUT | `/api/v1/messages/read` | MessageHandler.MarkRead | 批量标记已读 |
| GET | `/api/v1/messages/unread-count` | MessageHandler.UnreadCount | 未读数 |

---

#### 上传接口

| 方法 | 路径 | handler | 说明 |
|------|------|---------|------|
| POST | `/api/v1/upload/image` | UploadHandler.UploadImage | 上传单张图片 |
| DELETE | `/api/v1/upload/image` | UploadHandler.DeleteImage | 删除图片 |

---

#### 管理员接口

| 方法 | 路径 | handler | 说明 |
|------|------|---------|------|
| GET | `/api/v1/admin/users` | AdminHandler.Users | 用户管理列表 |
| PUT | `/api/v1/admin/users/:id/status` | AdminHandler.UpdateUserStatus | 启用/禁用用户 |
| GET | `/api/v1/admin/books` | AdminHandler.Books | 图书管理列表 |
| PUT | `/api/v1/admin/books/:id/status` | AdminHandler.UpdateBookStatus | 下架违规图书 |
| GET | `/api/v1/admin/orders` | AdminHandler.Orders | 全部订单 |
| GET | `/api/v1/admin/categories` | AdminHandler.Categories | 分类列表 |
| POST | `/api/v1/admin/categories` | AdminHandler.CreateCategory | 添加分类 |
| PUT | `/api/v1/admin/categories/:id` | AdminHandler.UpdateCategory | 修改分类 |
| DELETE | `/api/v1/admin/categories/:id` | AdminHandler.DeleteCategory | 删除分类 |
| POST | `/api/v1/admin/banners` | AdminHandler.CreateBanner | 添加轮播图 |
| PUT | `/api/v1/admin/banners/:id` | AdminHandler.UpdateBanner | 修改轮播图 |
| DELETE | `/api/v1/admin/banners/:id` | AdminHandler.DeleteBanner | 删除轮播图 |
| GET | `/api/v1/admin/statistics` | AdminHandler.Statistics | **数据统计** |

---

## 六、统一响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": { ... },
  "meta": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

### 错误响应

```json
{
  "code": 1001,
  "message": "参数错误",
  "data": null
}
```

### 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 0 | 200 | 成功 |
| 1001 | 400 | 参数错误 |
| 1002 | 401 | 未认证（Token失效/未登录） |
| 1003 | 403 | 权限不足 |
| 1004 | 404 | 资源不存在 |
| 1005 | 409 | 资源已存在（如重复注册） |
| 1006 | 409 | 操作不允许 |
| 2001 | 401 | 手机号或密码错误 |
| 2002 | 403 | 用户已被禁用 |
| 3001 | 409 | 图书已售出 |
| 3002 | 400 | 不能购买自己的图书 |
| 4001 | 409 | 订单状态不允许当前操作 |
| 5001 | 400 | 上传文件过大（限制5MB） |
| 5002 | 400 | 上传文件类型不支持 |
| 9999 | 500 | 系统内部错误 |

---

## 七、技术要点

### 7.1 JWT 认证流程

```
登录 ──> 校验账号密码 ──> 签发 access_token(2h) + refresh_token(7d)
                                    │
                   每次请求 Header: Authorization: Bearer <token>
                                    │
                              ┌─────┴─────┐
                              │  auth 中间件 │
                              └─────┬─────┘
                                    │
                         解析 JWT → 验证签名 → 检查 Redis 黑名单
                                    │
                             c.Set("user_id", uid)
                                    │
                               执行业务逻辑
```

### 7.2 密码安全
- 使用 `golang.org/x/crypto/bcrypt` 哈希存储
- 成本因子（cost）设为 10，兼顾安全与性能

### 7.3 图片存储
- **本地存储**：`uploads/images/{yyyy/mm/dd}/{uuid}.{ext}`
- **限制**：单张最大 5MB，支持 jpg/png/webp
- **访问**：开发阶段用 `r.Static()`，生产用 Nginx

### 7.4 搜索实现
- 使用 MySQL `LIKE` 模糊匹配，对 `title` 和 `author` 做简单搜索
- 配合 `category_id`、`condition` 等多维筛选
- 后续数据量增大后可升级为 Elasticsearch

### 7.5 订单与图书状态一致性

创建订单时必须保证 **订单 + 图书状态变更** 的原子性：

```go
tx := db.Begin()

// 1. 检查图书是否在售
book := getBook(tx, bookId)
if book.Status != 1 { tx.Rollback(); return error }

// 2. 创建订单
order := Order{SellerId: book.UserId, BuyerId: uid, ...}
tx.Create(&order)

// 3. 标记图书已售
tx.Model(&book).Update("status", 2)

// 4. 提交
tx.Commit()
```

### 7.6 核心依赖

| Go 依赖 | 用途 |
|---------|------|
| `github.com/gin-gonic/gin` | Web 框架 |
| `gorm.io/gorm` + `gorm.io/driver/mysql` | ORM |
| `github.com/redis/go-redis/v9` | Redis 客户端 |
| `github.com/golang-jwt/jwt/v5` | JWT Token |
| `golang.org/x/crypto/bcrypt` | 密码哈希 |
| `github.com/go-playground/validator/v10` | 参数校验 |
| `go.uber.org/zap` | 结构化日志 |
| `github.com/sony/sonyflake` | 分布式唯一ID |
| `github.com/swaggo/gin-swagger` | Swagger API 文档 |
| `github.com/spf13/viper` | 配置管理 |

---

## 八、实施计划（分阶段）

### 第一阶段：项目骨架搭建

| 任务 | 产出 |
|------|------|
| ① 初始化 Go Module、创建目录结构 | `go.mod` + 目录树 |
| ② 配置文件 (`config.yaml`) + 加载逻辑 | config 包 |
| ③ MySQL + Redis 初始化连接 | database 包 |
| ④ 统一响应 + 错误码定义 | common 包 |
| ⑤ 基础中间件（CORS、Logger） | middleware 包 |
| ⑥ JWT 工具 + 认证中间件 | utils/jwt.go + middleware/auth.go |
| ⑦ Gin 路由框架注册 | routes/router.go + main.go |
| ⑧ 验证：`go build` 通过，启动无报错 | |

### 第二阶段：核心业务功能

| 任务 | 产出 |
|------|------|
| ① 用户注册、登录、个人资料 | user_handler + service + repo |
| ② 图书发布、编辑、删除 | book_handler + service + repo |
| ③ 图书分类管理 | category_handler |
| ④ 图书列表（分页+筛选+排序） | BookHandler.List |
| ⑤ 收藏功能 | favorite_handler + service + repo |
| ⑥ 建表 SQL 脚本 | docs/schema.sql |

### 第三阶段：交易闭环

| 任务 | 产出 |
|------|------|
| ① 订单创建（含事务） | order_handler + service + repo |
| ② 订单确认/完成/取消 | 订单状态机流转 |
| ③ 站内消息/聊天 | message_handler + service + repo |
| ④ 图书状态自动关联（已售下架） | 订单创建时同步 |

### 第四阶段：管理完善

| 任务 | 产出 |
|------|------|
| ① 管理员后台接口 | admin_handler |
| ② 文件上传 | upload_handler + utils/upload.go |
| ③ Redis 缓存集成 | 各 service 层嵌入 |
| ④ 轮播图管理 | banner 相关接口 |
| ⑤ 数据统计接口 | 用户数、图书数、订单数 |

### 第五阶段：前端对接（后续）

| 任务 | 说明 |
|------|------|
| ① 初始化 Vue 3 + Vite 项目 | 前端骨架 |
| ② 用户端页面 | 登录、注册、首页、图书列表、详情、发布、订单、聊天 |
| ③ 管理端页面 | 用户管理、图书审核、数据看板 |
| ④ 联调测试 | 前后端对接 |

---

## 九、数据库初始化

GORM AutoMigrate 自动建表，只需手动创建数据库：

```sql
CREATE DATABASE IF NOT EXISTS xatu_book_exchange
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;
```

Go 代码中通过 GORM 自动迁移：

```go
db.AutoMigrate(
    &model.User{},
    &model.Category{},
    &model.Book{},
    &model.Favorite{},
    &model.Order{},
    &model.Message{},
    &model.Banner{},
)
```

---

## 十、config.yaml 配置模板

```yaml
server:
  port: 8080
  mode: debug          # debug / release / test

mysql:
  host: 127.0.0.1
  port: 3306
  user: root
  password: ""
  dbname: xatu_book_exchange
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-jwt-secret-key"
  access_expire: 2h      # access_token 过期时间
  refresh_expire: 168h   # refresh_token 过期时间（7天）

upload:
  max_size: 5            # 单文件最大 MB
  dir: "./uploads/images"
  allow_types:
    - ".jpg"
    - ".jpeg"
    - ".png"
    - ".webp"
```

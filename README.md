# 企业大模型公共服务平台 (AI Gateway Console)

轻量级企业大模型 API 管理平台（模拟版）。支持多模型统一管理、API 调用全链路日志、Token / 费用实时统计与 Dashboard 展示。

## 技术栈

| 层级 | 技术 |
|---|---|
| 后端 | Go + Gin + Gorm + MySQL + Redis |
| 前端 | Vue3 + Vite + TypeScript + Element Plus + ECharts + Pinia |
| 基础设施 | Docker Compose (MySQL 8.0 + Redis 7) |

## 快速启动

### Docker Compose（推荐）

```bash
docker-compose up -d
```

访问：
- 前端 Dashboard: http://localhost:3000
- 后端 API: http://localhost:8080

### 本地开发

```bash
# 启动基础设施
docker-compose up -d mysql redis

# 初始化数据库（首次）
docker exec -i ai_gateway_mysql mysql -u root -proot123 ai_gateway < server/migrations/001_init.sql
docker exec -i ai_gateway_mysql mysql -u root -proot123 ai_gateway < server/migrations/002_seed.sql

# 后端
cd server && go run cmd/server/main.go

# 前端
cd web && npm install && npm run dev
# → http://localhost:5173
```

## 种子数据

| 用途 | 值 |
|---|---|
| 测试 API Key | `sk-example-key-for-demo-20260805` |
| Admin 账号 | `admin` / `admin123` |
| Teacher 账号 | `teacher_wang` / `teacher123` |
| Student 账号 | `student_zhang` / `student123` |

## Demo 实现范围

### 实际实现

- POST /api/chat 完整闭环（API Key 校验 + 额度检查 + Mock 模型 + 费用计算 + 日志 + Redis 统计）
- 6 个统计 API（overview / model-usage / trend / logs / models / health）
- Dashboard 单页控制台（健康条 + 指标卡片 + 模型卡片 + 实时日志 + 趋势图 + 调用测试区）
- API Key bcrypt 安全存储 + 日额度校验
- 独立 goroutine 模拟器（3~6s 间隔，5 个模型差异化特征）
- 费用为核心指标（模型定价 → 实时计算 → Dashboard 展示）
- Docker Compose 一键启动

### 架构设计（暂不实现）

- CAS / OIDC 统一身份认证
- 消息队列（Kafka / RabbitMQ）
- 多租户隔离
- 完整 RBAC 管理后台
- 日志表分区归档
- 完整告警规则引擎
- 真实调用外部大模型 API

详见 [架构设计文档 V1.2](企业大模型公共服务平台系统架构设计文档%20V1.2.md) 第 1.3 节。

## 演示步骤

1. `docker-compose up -d` 启动全部服务
2. 浏览器打开 http://localhost:3000
3. 观察 Dashboard 实时数据（模拟器每 3~6s 生成一条）
4. 在调用测试区选择模型 → 输入消息 → 点击发送
5. 观察指标卡片数字变化 + 实时日志追加
6. 展示返回的 OpenAI 兼容格式（含 Token + 费用）
7. 讲解架构设计思路

## 接口文档

详见架构设计文档 [4.4 节](企业大模型公共服务平台系统架构设计文档%20V1.2.md)。

### 核心接口

| 接口 | 方法 | 认证 | 说明 |
|---|---|---|---|
| `/api/chat` | POST | API Key | 大模型调用 |
| `/api/stats/overview` | GET | 无 | 核心概览 |
| `/api/stats/model-usage` | GET | 无 | 模型用量分布 |
| `/api/stats/trend` | GET | 无 | 7 天趋势 |
| `/api/realtime/logs` | GET | 无 | 实时日志（分页） |
| `/api/models` | GET | 无 | 模型列表 + 今日指标 |
| `/api/health` | GET | 无 | 健康检查 |

## 项目结构

```
.
├── docker-compose.yml
├── server/
│   ├── cmd/server/main.go           # 入口
│   ├── internal/
│   │   ├── config/                  # 配置
│   │   ├── model/                   # GORM 数据模型
│   │   ├── handler/                 # HTTP 处理
│   │   ├── service/                 # 业务逻辑
│   │   ├── middleware/              # 中间件
│   │   ├── simulator/               # 模拟数据生成
│   │   ├── pkg/                     # 公共工具
│   │   └── router/                  # 路由注册
│   └── migrations/                  # SQL 脚本
├── web/
│   └── src/
│       ├── api/                     # API 封装
│       ├── components/              # Vue 组件
│       ├── stores/                  # Pinia Store
│       └── types/                   # TypeScript 类型
└── docs/                            # 设计文档
```

## 未来扩展方向

- 真实接入大模型 API 网关
- API Gateway: 鉴权、限流、路由、熔断
- 对接学校统一身份认证（CAS/OAuth2/OIDC）
- 消息队列 + WebSocket 实现真正实时推送
- 多租户隔离与配额管理
- 异常告警（Token 突增、失败率升高、模型离线）

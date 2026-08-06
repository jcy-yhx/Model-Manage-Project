# Changelog

## V1.2 (2026-08-07)

### 新增

- **POST /api/chat** 完整调用闭环（API Key bcrypt 校验 + 日额度检查 + Mock 模型 + OpenAI 兼容响应）
- **费用计算**：models 表新增 input_price / output_price，每次调用实时计算并写入 api_usage_logs.cost
- **模拟器**：后台 goroutine，3~6s 间隔，5 个模型差异化特征（Token 量、延迟、成功率）
- **6 个统计 API**：overview / model-usage / trend / realtime/logs / models / health
- **Redis 实时聚合**：当日 Token/请求/费用/延迟/成功率，48h TTL 自动过期
- **日汇总表** daily_usage_stats：为趋势查询提供轻量聚合
- **Vue3 Dashboard 单页控制台**：
  - 顶部健康状态条
  - 6 个核心指标卡片（今日费用突出显示）
  - 5 个模型状态卡片
  - 实时调用流表格（Element Plus）
  - 7 天趋势折线图（ECharts）
  - 模型使用占比环形图
  - API 调用测试区（现场演示闭环）
- **Docker Compose** 一键启动（MySQL 8.0 + Redis 7 + Go Server + Nginx + Vue）
- **API Key 安全**：bcrypt (cost=12) 哈希 + key_prefix 索引 + 日志脱敏

### 架构设计

- **统一身份认证方案**：JWT Demo → CAS/OAuth2/OIDC 演进路径 + RBAC 角色矩阵
- **日志表长期策略**：单表 → MySQL RANGE 分区 → 冷热分离 → ClickHouse/ES
- **实时监控演进**：前端轮询 → 消息队列 + WebSocket 推送
- **Gateway 分期**：Gin 中间件组合 → 独立 Kong/APISIX

### 工程化

- Go 三层架构：handler → service → model
- 统一 JSON 响应格式 `{"code":0,"message":"ok","data":{...}}`
- 优雅关闭（SIGINT/SIGTERM）
- 配置管理（YAML + 环境变量覆盖）
- 前端 API 层独立封装 + Pinia 状态管理 + TypeScript 类型定义
- 分阶段开发计划（阶段 0~7）

### 种子数据

- 3 个用户（admin / teacher / student）
- 1 个项目 + 1 个 API Key
- 5 个大模型（Qwen2.5-72B / DeepSeek-V3 / MiniMax-Text-01 / GLM-4 / Qwen2.5-Coder）
- 参考定价（元/1K tokens）

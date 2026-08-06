# Demo 检查清单

用于面试视频录制前的准备和现场演示指导。

## 一、启动前检查

- [ ] Docker Desktop 已启动
- [ ] 端口 3306、6391、8080、3000 未被占用
- [ ] 项目代码位于工作目录

## 二、启动步骤

```bash
# 1. 启动全部服务（首次启动会构建镜像，约 2 分钟）
docker-compose up -d

# 2. 等待所有容器健康
docker ps --filter "name=ai_gateway"
# 预期：4 个容器均为 healthy / Up

# 3. 初始化数据库（仅首次）
docker exec -i ai_gateway_mysql mysql -u root -proot123 ai_gateway < server/migrations/001_init.sql
docker exec -i ai_gateway_mysql mysql -u root -proot123 ai_gateway < server/migrations/002_seed.sql
```

## 三、演示操作步骤

### 环节 1：Dashboard 首页展示（约 30 秒）

1. 浏览器打开 http://localhost:3000
2. 展示内容：
   - [ ] 顶部健康条：模型在线数 / 总数、最后更新时间
   - [ ] 6 个指标卡片：今日费用 ¥xx.xx（蓝色边框突出）→ Token → 累计 → 请求数 → 成功率 → 延迟
   - [ ] 数字平滑变化（模拟器持续生成数据）

### 环节 2：模型卡片 + 占比图（约 20 秒）

3. 展示 5 张模型卡片
   - [ ] 每张卡片：模型名称 + 绿色状态灯 + 定价 + Token + 费用 + 延迟
4. 展示右侧环形图（模型用量占比）

### 环节 3：实时日志 + 趋势图（约 20 秒）

5. 展示实时调用流表格
   - [ ] 时间、项目名、模型名、Token、费用、状态（绿色 success 标签）、延迟
6. 展示 7 天趋势折线图

### 环节 4：API 调用闭环演示（约 1.5 分钟，核心高光时刻）

7. 滚动到底部「API 调用测试区」
8. 操作：
   - [ ] 下拉选择模型（保持默认 qwen2.5-72b）
   - [ ] API Key 已预填 `sk-example-key-for-demo-20260805`
   - [ ] 输入消息：「介绍一下量子计算」
   - [ ] 点击「发送请求」
9. 观察结果：
   - [ ] 右侧显示返回结果：模型名、Token 数、费用、回复内容（量子计算相关内容）
   - [ ] 顶部指标卡片「今日预计费用」数字变化
   - [ ] 实时日志中出现刚才的调用记录
   - [ ] 强调：**一次调用 → MySQL 记录 + Redis 统计 + Dashboard 刷新 = 完整闭环**

### 环节 5：异常场景演示（约 30 秒，可选）

10. 输入不存在的模型名 → 展示 404 错误
11. 修改 API Key 为无效值 → 展示 401 错误

## 四、架构讲解建议

```bash
# 同步打开另一个终端演示架构
docker ps                          # 展示 4 个容器
docker exec ai_gateway_mysql mysql -u root -proot123 ai_gateway \
  -e "SELECT COUNT(*) FROM api_usage_logs;"  # 展示累计调用量
docker exec ai_gateway_redis redis-cli KEYS "stats:*" | head  # 展示 Redis 实时指标
```

## 五、视频时间分配建议（5~10 分钟）

| 环节 | 时长 |
|---|---|
| 需求理解与项目定位 | 30s |
| 总体架构 + 技术选型 | 1min |
| 数据库设计与企业级抽象 | 1min |
| **Dashboard 实时效果** | 1.5min |
| **API 调用闭环演示** | 1.5min |
| API Key 安全 + 额度校验 | 1min |
| 统一身份认证设计方案 | 1min |
| 优化方向与未来扩展 | 1min |

## 六、常见问题准备

- Q: 「这些功能你全部实现了吗？」
- A: 不。项目有明确的 Demo 实现边界（见架构文档 1.3 节）。API 调用闭环、Dashboard、费用统计是全量实现的；CAS/OIDC、消息队列、告警引擎是架构设计方案，Demo 阶段不实现。这种边界控制体现了对项目周期的理解。

- Q: 「成本数据怎么算的？」
- A: models 表存储 input_price / output_price（元/1K tokens），每次调用实时计算 `(prompt/1000)×input + (completion/1000)×output`，写入 MySQL 日志 + Redis 实时指标，Dashboard 从 Redis 读取展示。

- Q: 「API Key 安全性？」
- A: bcrypt (cost=12) 哈希存储，仅首次创建时展示原始 Key；校验时先 key_prefix 缩小搜索范围，再 bcrypt 比对；日志中仅输出前缀，不记录原始 Key。

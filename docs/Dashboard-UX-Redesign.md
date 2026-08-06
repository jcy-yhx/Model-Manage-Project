# Dashboard UI/UX 重构方案

> 项目：企业大模型公共服务平台（AI Gateway Console）
>
> 版本：V1.3 Design Proposal
>
> 状态：等待评审确认

---

## 1. 当前页面问题分析

### 1.1 信息层级

| 问题 | 详情 |
|---|---|
| **顶层缺失** | 页面无 Header / 平台标识。用户打开后不知道"这是哪个系统"。健康条太轻——只一条浅蓝 bar，容易被跳过 |
| **同级堆砌** | 6 个指标卡片等宽等高，成本和 Token 视觉权重完全相同。没有主次 |
| **模型卡片信息过载** | 每张卡塞了名称、提供商 tag、输入/输出定价、Tokens、费用、延迟。体积小内容多，阅读负担高 |
| **调用测试区在下层** | TestPanel 放在页面最底部——这是面试演示的高光时刻，却被当"附赠功能"处理 |
| **实时日志可读性差** | Element Plus 默认表格样式，stripe 条纹 + 全边框，信息密度低。6 位小数费用列浪费宽度 |

### 1.2 视觉层级

| 问题 | 详情 |
|---|---|
| **无 Design System** | 每个组件各自定义 `border-radius: 10px` / `#fff` / `box-shadow`，没有 CSS 变量，未来改一个颜色要改 8 个文件 |
| **阴影泛滥** | 每个卡片都有 `box-shadow: 0 1px 4px rgba(0,0,0,.08)`。用阴影作为唯一分隔手段，而不用留白或背景色差 |
| **边框混乱** | 费用卡片用 `border: 2px solid #1677ff` 突出，其他卡片无边框。一致性弱 |
| **间距不统一** | `padding: 16px` vs `padding: 14px` vs `padding: 10px 20px` 混用，没有间距体系 |
| **Element Plus 默认感** | el-table 和 el-select 的默认样式直接使用，一看就是模板后台 |

### 1.3 页面节奏

```
当前：平铺 → 平铺 → 平铺
  [HealthBar]          ← 一条细带
  [6 Cards]            ← 6个等大卡片
  [ModelCards + Pie]   ← 左宽右窄，不均衡
  [TrendChart]         ← 孤立的图
  [RealtimeLog]        ← 大表格
  [TestPanel]          ← 收尾
```

问题：没有呼吸位。没有区域过渡。所有模块间距 `margin-top: 12px` 一刀切。

### 1.4 空间利用

| 问题 | 详情 |
|---|---|
| **右侧 280px 浪费** | PieChart 占 280px，模型卡片占剩余宽度。实际环形图只需要 180px，多出的 100px 是空白 |
| **全宽恒等 1400px** | 没有利用中等屏幕的灵活布局。笔记本上看 6 列指标卡片太挤 |
| **日志表格矮胖** | `max-height: 300px` 只能看 6~7 行，下面就是 TestPanel |

### 1.5 状态表达

| 问题 | 详情 |
|---|---|
| **健康灯太小** | `width: 8px` 的绿色圆点，需要仔细看才知道在线 |
| **趋势无异常标注** | 如果 Token 突增或成功率下降，折线图不会变红——没有阈值告警视觉呈现 |
| **日志状态标签** | `el-tag size="small"` 是最小尺寸，success 为绿色、fail 为红色，timeout 归为 danger——但 timeout 和 fail 对运维人员含义不同，应有独立颜色 |

### 1.6 企业产品感

当前页面给人的感觉：

> 一个用 Vue Admin 模板搭的学生大作业后台，把业务数据填充进去就完事了。

缺少以下企业产品特征：
- 品牌 Header
- 时间选择器 / 筛选器
- 数据对比（环比、同比）
- 异常状态的主动提示（而非被动展示）
- 操作入口（当前只有一个隐式的 TestPanel）

### 1.7 AI 平台属性

当前页面只是一个"数据看板"——Stats Dashboard。而非一个 **AI Gateway 管理控制台**。

AI 平台应有的属性：
- 模型即资产：一眼看到每个模型的健康度和调用量
- 调用即成本：费用不是"附加列"，是核心指标
- Gateway = 流量入口：应该有流控感的实时数据
- 测试即体验：Playground 不应该是"附赠功能"，而是平台入口的一部分

---

## 2. 产品定位

**产品定义：**

> **企业 AI 模型调用治理与监控控制台**（AI Gateway Console）

不是「数据大屏」，不是「Admin 后台」。是一线研发和运维团队用来**管理 API 调用、监控模型健康、追踪成本**的日常工具型控制台。

**用户打开页面后的关注路径：**

| 时刻 | 用户认知 |
|---|---|
| **0~3 秒** | 平台状态是否正常？（一个全局状态指示） |
| **3~10 秒** | 费用花了多少？模型都健康吗？ |
| **10~30 秒** | 哪个模型用量最大？有没有异常趋势？ |
| **30 秒+** | 我要试一下调用一个模型 |

> 首屏不需要展示一切。信息应该有节奏地递进，用户随滚动自然深化。

---

## 3. 用户认知路径

```
顶部全局状态条（0.5秒）
  "平台正常运行中  ·  5/5 模型在线  · 今日费用 ¥18.52"
  ↓
Hero 区域（3秒）  ← 首屏核心
  一块大面板：今日费用（超大数字）+ Token 用量 + 请求数 + 成功率 + 延迟
  对比昨日变化（绿色↑/红色↓ 箭头）
  ↓
模型健康面板（10秒）
  5个模型横向卡片（状态 + 调用量微柱 + 延迟 + 定价）
  ↓
双图区（20秒）
  左：7天 Token/费用趋势折线图
  右：当前模型用量占比环形图
  ↓
实时调用流表格（30秒）
  最近调用记录（可展开查看详情）
  ↓
API Playground（30秒+）
  调用测试区，内嵌在控制台底部的交互模块
```

---

## 4. 页面 Wireframe

```
┌──────────────────────────────────────────────────────────────────┐
│  HEADER                                                           │
│  AI Gateway Console    [在线 4/5] [今日 ¥18.52] [14:32:15]       │
├──────────────────────────────────────────────────────────────────┤
│                                                                    │
│  HERO METRICS                                                      │
│  ┌─────────────────────────────────────────────────────────┐      │
│  │  今日预计费用                 数据更新时间：14:32:15     │      │
│  │  ¥ 18.52    ↑12.5%                                      │      │
│  │                                                          │      │
│  │  ───────────────────────────────────────────────────────│      │
│  │                                                          │      │
│  │  1,245,670         3,842          98.3%         342 ms   │      │
│  │   今日 Token       今日请求       平均成功率     平均延迟  │      │
│  │   ↑8.3% vs 昨日    ↑5.1%          ↑0.2%         ↓12ms    │      │
│  └─────────────────────────────────────────────────────────┘      │
│                                                                    │
│  MODEL HEALTH                                                      │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐  │
│  │ Qwen     │ │ DeepSeek │ │ MiniMax  │ │ GLM-4    │ │ Qwen   │  │
│  │ 2.5-72B  │ │ V3       │ │ Text-01  │ │          │ │ Coder  │  │
│  │ ● online │ │ ● online │ │ ● online │ │ ● online │ │ ●online│  │
│  │ █████▌   │ │ ███▌     │ │ ██▌      │ │ █▌       │ │ ████▌  │  │
│  │ ¥3.82    │ │ ¥1.20    │ │ ¥0.48    │ │ ¥0.15    │ │ ¥2.88  │  │
│  │ 380ms   │ │ 220ms    │ │ 650ms    │ │ 120ms    │ │ 290ms  │  │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └────────┘  │
│                                                                    │
│  ANALYTICS                                                         │
│  ┌─────────────────────────────────┐ ┌──────────────────┐         │
│  │  调用趋势（近 7 天）            │ │  模型用量分布     │         │
│  │  📈 Token + 费用 双 Y 轴折线   │ │  🍩 环形图        │         │
│  │                                 │ │                  │         │
│  └─────────────────────────────────┘ └──────────────────┘         │
│                                                                    │
│  REALTIME LOGS                                                     │
│  ┌─────────────────────────────────────────────────────────┐      │
│  │  实时调用流                    [筛选: 全部模型 ▾]       │      │
│  │  ────────────────────────────────────────────────────── │      │
│  │  14:32:15 │ 智能问答 │ Qwen2.5 │ 235 │ ¥0.0028 │ ✅ │ 410ms│
│  │  14:32:12 │ 智能问答 │ DeepS-3 │ 189 │ ¥0.0015 │ ✅ │ 220ms│
│  │  14:32:08 │ 智能问答 │ MiniMax │ 420 │ ¥0.0312 │ ✅ │ 650ms│
│  │  14:32:03 │ 智能问答 │ Qwen2.5 │ 156 │ ¥0.0019 │ ❌ │ timeout│
│  │  ────────────────────────────────────────────────────── │      │
│  │                                       查看全部日志 →    │      │
│  └─────────────────────────────────────────────────────────┘      │
│                                                                    │
│  PLAYGROUND                                                        │
│  ┌─────────────────────────────────────────────────────────┐      │
│  │  API Playground                                          │      │
│  │  ┌──────────────────┐ ┌───────────────────────────────┐ │      │
│  │  │ 模型 [Qwen2.5 ▾] │ │                               │ │      │
│  │  │                   │ │  量子计算是一种利用量子力学   │ │      │
│  │  │ 输入消息...       │ │  原理进行信息处理的...        │ │      │
│  │  │                   │ │                               │ │      │
│  │  │ [发送]            │ │  qwen2.5-72b · 198 tokens    │ │      │
│  │  └──────────────────┘ │  ¥0.0020                      │ │      │
│  │                       └───────────────────────────────┘ │      │
│  └─────────────────────────────────────────────────────────┘      │
└──────────────────────────────────────────────────────────────────┘
```

**视觉优先级（从上到下递减）：**
1. Hero 指标面板 ← 首屏焦点
2. 模型健康卡片 ← 次焦点
3. 趋势图 + 用量分布 ← 分析层
4. 实时日志 ← 数据层
5. Playground ← 交互层

---

## 5. Design System

### 5.1 颜色规范

```
主色系统（蓝色）
  --color-primary:         #2563EB    主色（按钮、链接、费用数字）
  --color-primary-light:   #EFF6FF    主色浅底（Hero 面板背景）
  --color-primary-dark:    #1D4ED8    主色深（hover 态）

中性色
  --color-bg:              #F8FAFC    页面底色
  --color-surface:         #FFFFFF    卡片背景
  --color-border:          #E2E8F0    分隔线 / 细边框
  --color-border-light:    #F1F5F9    轻分隔

文本色
  --color-text-primary:    #0F172A    主文字（标题、指标数字）
  --color-text-secondary:  #475569    次文字（标签、说明）
  --color-text-tertiary:   #94A3B8    辅助文字（更新时间、单位）

语义色（状态）
  --color-success:         #16A34A    成功 / online
  --color-warning:         #EA580C    降级 / degraded
  --color-error:           #DC2626    失败 / offline / timeout
  --color-info:            #2563EB    信息

图表色（5 个模型）
  --chart-1:               #3B82F6    蓝色系渐变
  --chart-2:               #6366F1    靛蓝色
  --chart-3:               #8B5CF6    紫色
  --chart-4:               #06B6D4    青色
  --chart-5:               #10B981    绿色
```

### 5.2 字体规范

```
--font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif
--font-mono:   "SF Mono", "JetBrains Mono", "Fira Code", monospace

层级：
  Page Title    24px / 700  / --color-text-primary    页面大标题
  Section Title 15px / 600  / --color-text-primary    区域标题
  Metric Number 36px / 700  / --color-text-primary    Hero 指标数字（费用）
  Metric Number 28px / 700  / --color-text-primary    次指标数字（Token等）
  Card Title    13px / 600  / --color-text-secondary  卡片标题
  Body          14px / 400  / --color-text-primary    正文
  Caption       12px / 400  / --color-text-tertiary   辅助文字
  Mono Data     13px / 500  / --color-text-primary    表格数据（等宽数字）
```

### 5.3 间距规范（8px 体系）

```
--space-1:  4px     微间距
--space-2:  8px     组件内间距
--space-3:  12px    卡片内 padding 小
--space-4:  16px    卡片 padding  / 组件间标准间距
--space-5:  20px    区块内间距
--space-6:  24px    区块间间距
--space-8:  32px    大区块间距
--space-10: 40px    页面级间距
```

### 5.4 圆角规范

```
--radius-sm:  6px    输入框、标签、小按钮
--radius-md:  10px   卡片、面板
--radius-lg:  14px   Hero 面板、大容器
--radius-full: 9999px 状态指示灯、胶囊标签
```

### 5.5 阴影规范

```
--shadow-none:  none                   默认（无阴影，用留白和色差分隔）
--shadow-sm:    0 1px 2px rgba(0,0,0,0.04)   悬浮态轻阴影
--shadow-md:    0 4px 12px rgba(0,0,0,0.06)  下拉、弹窗
--shadow-lg:    0 8px 24px rgba(0,0,0,0.08)  模态框

原则：减少阴影使用。卡片之间优先用背景色差 + 留白区分，而不是阴影堆砌。
```

---

## 6. 组件拆分方案

```
Dashboard (App.vue)
├── AppHeader.vue          全局顶栏：品牌名 + 状态摘要 + 更新时间
├── HeroMetrics.vue        首屏 Hero：费用大数字 + 4 个次指标 + 环比箭头
├── ModelHealthPanel.vue   模型健康网格：5 个模型卡（状态灯 + 微柱图 + 关键指标）
├── AnalyticsRow.vue       分析区包裹
│   ├── TrendChart.vue     趋势折线图（Token + 费用 双 Y 轴，优化配色）
│   └── UsagePie.vue       环形图（精简尺寸）
├── RealtimeLogs.vue       日志表格（去边框、优化列宽）
└── Playground.vue         调用测试区（集成式设计）
```

**不再需要的组件：**
- `HealthBar.vue` → 融入 `AppHeader.vue`
- `MetricCards.vue` → 重构为 `HeroMetrics.vue`
- `ModelCards.vue` → 重构为 `ModelHealthPanel.vue`
- `TestPanel.vue` → 重构为 `Playground.vue`
- `PieChart.vue` → 重构为 `UsagePie.vue`
- `StatusDot.vue` → 保留并优化尺寸

### 组件职责

| 组件 | 数据来源 | 交互 |
|---|---|---|
| AppHeader | store.health, store.lastUpdate | 静态展示 |
| HeroMetrics | store.overview | 数字增长动画，hover 展示环比详情 |
| ModelHealthPanel | store.models | 卡片 hover 高亮 |
| TrendChart | store.trend | ECharts tooltip 交互 |
| UsagePie | store.modelUsage | ECharts 扇区高亮 |
| RealtimeLogs | store.logs | 自动滚动，行 hover |
| Playground | store.sendChat | 模型选择、输入消息、发送、看结果 |

---

## 7. 数据流不变

保持现有 Pinia Store 和 API 层完全不动。所有组件仅改变 UI 呈现，不修改任何：
- `web/src/api/*.ts`
- `web/src/stores/dashboard.ts`
- `web/src/types/index.ts`
- 后端任何文件

---

## 8. 实现阶段

### 阶段 6-1：Design System + 页面 Layout

**产出：**
- `web/src/styles/variables.css` — CSS 自定义属性（颜色/字体/间距/圆角/阴影）
- `web/src/App.vue` — 新 Layout 骨架（Header + 主体滚动区 + 组件占位）
- `web/src/components/layout/AppHeader.vue` — 全局顶栏

### 阶段 6-2：Hero 核心指标区域

**产出：**
- `web/src/components/dashboard/HeroMetrics.vue`

**特点：** 费用超大数字 + 环比箭头 + 4 个次指标。过渡动画平滑。

### 阶段 6-3：API Playground

**产出：**
- `web/src/components/dashboard/Playground.vue`

**原因：** 这是面试演示核心。放在阶段 6-2 后立即实现，确保闭环演示可用。

### 阶段 6-4：模型健康面板

**产出：**
- `web/src/components/dashboard/ModelHealthPanel.vue`
- `web/src/components/common/StatusDot.vue`（优化）

### 阶段 6-5：趋势分析 + 用量分布

**产出：**
- `web/src/components/dashboard/AnalyticsRow.vue`
- `web/src/components/charts/TrendChart.vue`（重写）
- `web/src/components/charts/UsagePie.vue`（重写）

### 阶段 6-6：实时日志

**产出：**
- `web/src/components/dashboard/RealtimeLogs.vue`（重写）

---

## 9. 不新增的内容（明确边界）

- 不新增任何后端 API
- 不修改 GORM Model 或数据库表
- 不新增页面路由
- 不引入新 npm 包（在 element-plus + echarts + pinia + vue-router + axios 范围内）
- 不修改 Pinia Store 逻辑

---

## 10. 预期对比

| 维度 | 当前 | 目标 |
|---|---|---|
| 平台识别 | 无品牌感知 | AppHeader 明确「AI Gateway Console」 |
| 首屏焦点 | 分散在 6 个等大卡片 | Hero 面板，费用大字 + 环比 |
| 模型展示 | 5 张密集卡片 | 5 张呼吸感卡片 + 微型用量条 |
| 图表体验 | ECharts 默认 | 企业配色、无多余图例元素 |
| Playground | 底部孤立区域 | 集成式 API 测试区 |
| 阴影 | 8 处 box-shadow | 精简到 0~2 处 |
| 间距 | 各处 12/14/16px 混用 | 8px 体系统一 |
| 颜色 | 硬编码 #1677ff / #52c41a 等 | CSS 变量 tokens |

---

> 方案结束。请评审后回复 **"同意 Design System，开始阶段 6-1"** 进入实现。

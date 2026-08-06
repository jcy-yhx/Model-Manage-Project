// === 统计概览 ===
export interface OverviewData {
  today: {
    total_tokens: number
    total_requests: number
    success_rate: number
    avg_latency_ms: number
    total_cost: number
  }
  total: {
    total_tokens: number
    total_requests: number
    total_cost: number
  }
  trend: {
    token_change_pct: number
    request_change_pct: number
  }
  model_count: {
    online: number
    total: number
  }
}

// === 模型用量 ===
export interface ModelUsageItem {
  model_name: string
  tokens: number
  percentage: number
  cost: number
}

// === 趋势 ===
export interface TrendPoint {
  date: string
  tokens: number
  requests: number
}

// === 日志 ===
export interface LogEntry {
  id: number
  model_name: string
  project_name: string
  total_tokens: number
  cost: number
  status: string
  latency_ms: number
  created_at: string
}

export interface Pagination {
  page: number
  page_size: number
  total: number
}

// === 模型 ===
export interface ModelInfo {
  id: number
  name: string
  display_name: string
  provider: string
  status: string
  input_price: number
  output_price: number
  description: string
  today_tokens: number
  today_requests: number
  avg_latency_ms: number
  today_cost: number
}

// === 健康 ===
export interface HealthData {
  db: string
  redis: string
  online_models: number
  total_models: number
  alerts: { level: string; model: string; message: string }[] | null
}

// === Chat ===
export interface ChatRequest {
  model: string
  messages: { role: string; content: string }[]
}

export interface ChatResponse {
  id: string
  object: string
  created: number
  model: string
  choices: {
    index: number
    message: { role: string; content: string }
    finish_reason: string
  }[]
  usage: {
    prompt_tokens: number
    completion_tokens: number
    total_tokens: number
    cost: number
  }
}

// === 统一响应 ===
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

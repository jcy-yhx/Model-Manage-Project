import http from './index'
import type {
  ApiResponse,
  OverviewData,
  ModelUsageItem,
  TrendPoint,
  LogEntry,
  Pagination,
  ModelInfo,
  HealthData,
  ChatRequest,
  ChatResponse,
} from '../types'

export async function fetchOverview() {
  const res = await http.get<ApiResponse<OverviewData>>('/stats/overview')
  return res.data.data
}

export async function fetchModelUsage() {
  const res = await http.get<ApiResponse<{ models: ModelUsageItem[]; total_tokens: number }>>('/stats/model-usage')
  return res.data.data
}

export async function fetchTrend(days = 7) {
  const res = await http.get<ApiResponse<{ points: TrendPoint[] }>>('/stats/trend', { params: { days } })
  return res.data.data.points
}

export async function fetchLogs(page = 1, pageSize = 20) {
  const res = await http.get<ApiResponse<{ logs: LogEntry[]; pagination: Pagination }>>('/realtime/logs', {
    params: { page, page_size: pageSize },
  })
  return res.data.data
}

export async function fetchModels() {
  const res = await http.get<ApiResponse<{ models: ModelInfo[] }>>('/models')
  return res.data.data.models
}

export async function fetchHealth() {
  const res = await http.get<ApiResponse<HealthData>>('/health')
  return res.data.data
}

export async function postChat(req: ChatRequest, apiKey: string) {
  const res = await http.post<ApiResponse<ChatResponse>>('/chat', req, {
    headers: { Authorization: `Bearer ${apiKey}` },
  })
  return res.data.data
}

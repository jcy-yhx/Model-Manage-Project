import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as api from '../api/stats'
import type {
  OverviewData, ModelUsageItem, TrendPoint, LogEntry, Pagination,
  ModelInfo, HealthData, ChatResponse,
} from '../types'

export const useDashboard = defineStore('dashboard', () => {
  // 各模块数据
  const overview = ref<OverviewData | null>(null)
  const modelUsage = ref<{ models: ModelUsageItem[]; total_tokens: number } | null>(null)
  const trend = ref<TrendPoint[]>([])
  const logs = ref<{ logs: LogEntry[]; pagination: Pagination } | null>(null)
  const models = ref<ModelInfo[]>([])
  const health = ref<HealthData | null>(null)
  const lastUpdate = ref('')

  // 调用测试
  const chatResult = ref<ChatResponse | null>(null)
  const chatError = ref('')
  const chatLoading = ref(false)

  // 轮询刷新
  let timer: ReturnType<typeof setInterval> | null = null

  async function refreshAll() {
    try {
      const [ov, mu, tr, lo, md, hl] = await Promise.all([
        api.fetchOverview(),
        api.fetchModelUsage(),
        api.fetchTrend(),
        api.fetchLogs(1, 10),
        api.fetchModels(),
        api.fetchHealth(),
      ])
      overview.value = ov
      modelUsage.value = mu
      trend.value = tr
      logs.value = lo
      models.value = md
      health.value = hl
      lastUpdate.value = new Date().toLocaleTimeString()
    } catch (e) {
      console.error('[Dashboard] refresh error:', e)
    }
  }

  function startPolling(intervalMs = 4000) {
    refreshAll()
    timer = setInterval(refreshAll, intervalMs)
  }

  function stopPolling() {
    if (timer) { clearInterval(timer); timer = null }
  }

  // 调用测试
  async function sendChat(model: string, message: string, apiKey: string) {
    chatLoading.value = true
    chatResult.value = null
    chatError.value = ''
    try {
      const res = await api.postChat(
        { model, messages: [{ role: 'user', content: message }] },
        apiKey
      )
      chatResult.value = res
      // 立即刷新 Dashboard 以显示新的用量
      await refreshAll()
    } catch (e: any) {
      chatError.value = e.response?.data?.message || e.message || '调用失败'
    } finally {
      chatLoading.value = false
    }
  }

  return {
    overview, modelUsage, trend, logs, models, health, lastUpdate,
    chatResult, chatError, chatLoading,
    refreshAll, startPolling, stopPolling, sendChat,
  }
})

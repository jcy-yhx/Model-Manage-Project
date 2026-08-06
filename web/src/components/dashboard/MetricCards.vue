<template>
  <div class="metric-cards">
    <div class="card cost">
      <div class="label">今日预计费用</div>
      <div class="value">¥{{ (store.overview?.today.total_cost ?? 0).toFixed(2) }}</div>
    </div>
    <div class="card">
      <div class="label">今日总 Token</div>
      <div class="value">{{ fmtNum(store.overview?.today.total_tokens ?? 0) }}</div>
    </div>
    <div class="card">
      <div class="label">累计总 Token</div>
      <div class="value">{{ fmtNum(store.overview?.total.total_tokens ?? 0) }}</div>
    </div>
    <div class="card">
      <div class="label">今日请求数</div>
      <div class="value">{{ fmtNum(store.overview?.today.total_requests ?? 0) }}</div>
    </div>
    <div class="card">
      <div class="label">平均成功率</div>
      <div class="value">{{ ((store.overview?.today.success_rate ?? 0) * 100).toFixed(1) }}%</div>
    </div>
    <div class="card">
      <div class="label">平均延迟</div>
      <div class="value">{{ store.overview?.today.avg_latency_ms ?? 0 }} ms</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDashboard } from '../../stores/dashboard'

const store = useDashboard()

function fmtNum(n: number): string {
  if (n >= 1e6) return (n / 1e6).toFixed(1) + 'M'
  if (n >= 1e3) return (n / 1e3).toFixed(1) + 'K'
  return String(n)
}
</script>

<style scoped>
.metric-cards { display: grid; grid-template-columns: repeat(6, 1fr); gap: 12px; }
.card { background: #fff; border-radius: 10px; padding: 16px; box-shadow: 0 1px 4px rgba(0,0,0,.08); text-align: center; }
.card.cost { border: 2px solid #1677ff; }
.label { font-size: 13px; color: #888; margin-bottom: 6px; }
.card.cost .label { color: #1677ff; font-weight: 600; }
.value { font-size: 22px; font-weight: 700; color: #222; transition: all .3s; }
.card.cost .value { color: #1677ff; font-size: 24px; }
</style>

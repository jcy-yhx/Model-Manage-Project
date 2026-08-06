<template>
  <div class="hero-panel">
    <div class="hero-top">
      <div class="hero-main">
        <span class="hero-label">今日预计费用</span>
        <div class="hero-value-row">
          <span class="hero-value">¥{{ costDisplay }}</span>
          <span v-if="costTrend !== 0" class="trend-badge" :class="costTrend > 0 ? 'up' : 'down'">
            <span class="trend-arrow">{{ costTrend > 0 ? '↑' : '↓' }}</span>
            {{ Math.abs(costTrend).toFixed(1) }}%
          </span>
        </div>
      </div>
      <div class="hero-time">
        数据更新于 {{ store.lastUpdate || '--' }}
      </div>
    </div>

    <div class="hero-divider"></div>

    <div class="hero-metrics-row">
      <div class="metric-item">
        <span class="metric-label">今日 Token</span>
        <span class="metric-value">{{ fmtLargeNum(today.tokens) }}</span>
        <span v-if="store.overview?.trend.token_change_pct !== 0" class="metric-trend" :class="store.overview?.trend.token_change_pct! > 0 ? 'up' : 'down'">
          {{ store.overview?.trend.token_change_pct! > 0 ? '↑' : '↓' }}{{ Math.abs(store.overview?.trend.token_change_pct!).toFixed(1) }}%
        </span>
      </div>
      <div class="metric-item">
        <span class="metric-label">今日请求</span>
        <span class="metric-value">{{ fmtLargeNum(today.requests) }}</span>
        <span v-if="store.overview?.trend.request_change_pct !== 0" class="metric-trend" :class="store.overview?.trend.request_change_pct! > 0 ? 'up' : 'down'">
          {{ store.overview?.trend.request_change_pct! > 0 ? '↑' : '↓' }}{{ Math.abs(store.overview?.trend.request_change_pct!).toFixed(1) }}%
        </span>
      </div>
      <div class="metric-item">
        <span class="metric-label">平均成功率</span>
        <span class="metric-value">{{ ((store.overview?.today.success_rate ?? 0) * 100).toFixed(1) }}%</span>
      </div>
      <div class="metric-item">
        <span class="metric-label">平均延迟</span>
        <span class="metric-value">{{ store.overview?.today.avg_latency_ms ?? 0 }} <span class="metric-unit">ms</span></span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useDashboard } from '../../stores/dashboard'

const store = useDashboard()

const today = computed(() => ({
  tokens: store.overview?.today.total_tokens ?? 0,
  requests: store.overview?.today.total_requests ?? 0,
}))

const costDisplay = computed(() => {
  const v = store.overview?.today.total_cost ?? 0
  return v.toFixed(2)
})

// 费用趋势用 Token 变化代替（费用和 Token 变化正相关）
const costTrend = computed(() => store.overview?.trend.token_change_pct ?? 0)

function fmtLargeNum(n: number): string {
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + 'M'
  if (n >= 1_000) return (n / 1_000).toFixed(1) + 'K'
  return String(n)
}
</script>

<style scoped>
.hero-panel {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-6) var(--space-8);
}

/* ── Top: Cost ── */
.hero-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
}
.hero-label {
  font-size: var(--text-card-title);
  font-weight: var(--weight-medium);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.hero-value-row {
  display: flex;
  align-items: baseline;
  gap: var(--space-3);
  margin-top: var(--space-1);
}
.hero-value {
  font-size: var(--text-hero);
  font-weight: var(--weight-bold);
  color: var(--color-text-primary);
  line-height: 1.1;
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
  transition: color var(--transition-fast);
}
.trend-badge {
  font-size: var(--text-body);
  font-weight: var(--weight-semibold);
  padding: 2px 10px;
  border-radius: var(--radius-full);
  display: inline-flex;
  align-items: center;
  gap: 2px;
}
.trend-badge.up {
  color: var(--color-success);
  background: var(--color-success-bg);
}
.trend-badge.down {
  color: var(--color-error);
  background: var(--color-error-bg);
}
.trend-arrow {
  font-size: 16px;
  line-height: 1;
}
.hero-time {
  font-size: var(--text-caption);
  color: var(--color-text-tertiary);
  white-space: nowrap;
}

/* ── Divider ── */
.hero-divider {
  height: 1px;
  background: var(--color-border);
  margin: var(--space-5) 0;
}

/* ── Bottom: 4 sub-metrics ── */
.hero-metrics-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-4);
}
.metric-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.metric-label {
  font-size: var(--text-caption);
  color: var(--color-text-tertiary);
  font-weight: var(--weight-normal);
}
.metric-value {
  font-size: var(--text-metric);
  font-weight: var(--weight-bold);
  color: var(--color-text-primary);
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
  line-height: 1.2;
}
.metric-unit {
  font-size: var(--text-body);
  font-weight: var(--weight-normal);
  color: var(--color-text-secondary);
}
.metric-trend {
  font-size: var(--text-caption);
  font-weight: var(--weight-medium);
  margin-top: 1px;
}
.metric-trend.up { color: var(--color-success); }
.metric-trend.down { color: var(--color-error); }
</style>

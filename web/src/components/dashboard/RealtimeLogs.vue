<template>
  <div class="rtl">
    <div class="rtl-header">
      <h3 class="rtl-title">实时调用流</h3>
      <div class="rtl-actions">
        <span class="rtl-count">{{ store.logs?.pagination?.total ?? 0 }} 条记录</span>
        <button class="rtl-filter" :class="{ active: statusFilter }" @click="toggleFilter" type="button">
          {{ statusFilter ? '筛选: ' + statusFilter : '全部状态 ▾' }}
        </button>
        <div v-if="showFilter" class="filter-dropdown">
          <button @click="setFilter('')" type="button">全部</button>
          <button @click="setFilter('success')" type="button">✅ success</button>
          <button @click="setFilter('fail')" type="button">❌ fail</button>
          <button @click="setFilter('timeout')" type="button">⏱ timeout</button>
        </div>
      </div>
    </div>

    <!-- 表头 -->
    <div class="rtl-thead">
      <span class="th-time">时间</span>
      <span class="th-model">模型</span>
      <span class="th-tokens">Token</span>
      <span class="th-cost">费用</span>
      <span class="th-status">状态</span>
      <span class="th-latency">延迟</span>
    </div>

    <!-- 表面（无传统 table 边框） -->
    <div class="rtl-body">
      <div
        v-for="log in store.logs?.logs ?? []"
        :key="log.id"
        class="rtl-row"
        :class="{ 'is-fail': log.status !== 'success' }"
      >
        <span class="td-time">{{ fmtTime(log.created_at) }}</span>
        <span class="td-model">
          <span class="model-badge" :style="{ background: modelColor(log.model_name) }"></span>
          {{ shortName(log.model_name) }}
        </span>
        <span class="td-tokens">{{ fmtNum(log.total_tokens) }}</span>
        <span class="td-cost">¥{{ log.cost.toFixed(4) }}</span>
        <span class="td-status">
          <span class="status-badge" :class="'badge-' + log.status">
            {{ statusLabel(log.status) }}
          </span>
        </span>
        <span class="td-latency">{{ log.latency_ms }}ms</span>
      </div>

      <div v-if="!store.logs?.logs?.length" class="rtl-empty">
        暂无调用记录
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useDashboard } from '../../stores/dashboard'

const CHART_COLORS: Record<string, string> = {
  'qwen2.5-72b': '#3B82F6',
  'qwen2.5-coder': '#6366F1',
  'deepseek-v3': '#06B6D4',
  'minimax-text-01': '#8B5CF6',
  'glm-4': '#10B981',
}

const store = useDashboard()
const showFilter = ref(false)
const statusFilter = ref('')
function toggleFilter() { showFilter.value = !showFilter.value }
function setFilter(v: string) { statusFilter.value = v; showFilter.value = false }

function fmtTime(t: string) {
  if (!t) return ''
  const d = new Date(t)
  return d.toLocaleTimeString('zh-CN', { hour12: false })
}

function fmtNum(n: number): string {
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + 'M'
  if (n >= 1_000) return (n / 1_000).toFixed(1) + 'K'
  return String(n)
}

function shortName(n: string): string {
  return n.replace('qwen2.5-', '').replace('minimax-', '').replace('deepseek-', '')
}

function modelColor(name: string): string {
  return CHART_COLORS[name] || '#94A3B8'
}

function statusLabel(s: string): string {
  switch (s) {
    case 'success': return '成功'
    case 'fail': return '失败'
    case 'timeout': return '超时'
    default: return s
  }
}
</script>

<style scoped>
.rtl {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

/* ── Header ── */
.rtl-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4) var(--space-6);
  border-bottom: 1px solid var(--color-border-light);
}
.rtl-title {
  font-size: var(--text-section);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
}
.rtl-actions {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  position: relative;
}
.rtl-count {
  font-size: var(--text-caption);
  color: var(--color-text-tertiary);
}
.rtl-filter {
  font-size: var(--text-caption);
  color: var(--color-text-secondary);
  background: var(--color-bg-alt);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: 4px 12px;
  cursor: pointer;
  transition: all var(--transition-fast);
}
.rtl-filter:hover, .rtl-filter.active { border-color: var(--color-primary); color: var(--color-primary); }

.filter-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 4px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-md);
  z-index: 10;
  display: flex;
  flex-direction: column;
  min-width: 120px;
  overflow: hidden;
}
.filter-dropdown button {
  padding: 6px 14px;
  font-size: var(--text-caption);
  color: var(--color-text-primary);
  background: none;
  border: none;
  text-align: left;
  cursor: pointer;
  transition: background var(--transition-fast);
}
.filter-dropdown button:hover { background: var(--color-bg-alt); }

/* ── Thead ── */
.rtl-thead {
  display: grid;
  grid-template-columns: 100px 1fr 100px 120px 80px 80px;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-6);
  background: var(--color-bg-alt);
  border-bottom: 1px solid var(--color-border);
  font-size: 11px;
  font-weight: var(--weight-semibold);
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.th-tokens, .th-cost, .th-status, .th-latency { text-align: right; }

/* ── Body (grid rows) ── */
.rtl-body {
  max-height: 360px;
  overflow-y: auto;
}
.rtl-row {
  display: grid;
  grid-template-columns: 100px 1fr 100px 120px 80px 80px;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-6);
  border-bottom: 1px solid var(--color-border-light);
  font-size: var(--text-caption);
  color: var(--color-text-primary);
  align-items: center;
  transition: background var(--transition-fast);
}
.rtl-row:hover { background: var(--color-bg); }
.rtl-row.is-fail { background: var(--color-error-bg); opacity: 0.85; }
.rtl-row.is-fail:hover { background: #FEE2E2; }

.td-time {
  color: var(--color-text-tertiary);
  font-family: var(--font-mono);
  font-size: 12px;
}
.td-model {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-weight: var(--weight-medium);
}
.model-badge {
  width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0;
}
.td-tokens {
  text-align: right;
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
}
.td-cost {
  text-align: right;
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
  font-weight: var(--weight-semibold);
  color: var(--color-primary);
}
.td-latency {
  text-align: right;
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
  color: var(--color-text-secondary);
}

/* Status badges */
.td-status { text-align: right; }
.status-badge {
  display: inline-block;
  padding: 1px 8px;
  border-radius: var(--radius-full);
  font-size: 11px;
  font-weight: var(--weight-medium);
}
.badge-success {
  color: var(--color-success);
  background: var(--color-success-bg);
}
.badge-fail {
  color: var(--color-error);
  background: var(--color-error-bg);
}
.badge-timeout {
  color: var(--color-warning);
  background: var(--color-warning-bg);
}

.rtl-empty {
  text-align: center;
  padding: var(--space-10);
  color: var(--color-text-tertiary);
  font-size: var(--text-body);
}
</style>

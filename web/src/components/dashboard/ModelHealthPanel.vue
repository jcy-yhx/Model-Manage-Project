<template>
  <div class="mhp">
    <div class="mhp-header">
      <h3 class="mhp-title">模型健康状态</h3>
      <div class="mhp-summary">
        <span class="summary-dot online"></span>
        {{ onlineCount }} 在线
        <span v-if="degradedCount" class="summary-dot degraded"></span>
        <span v-if="degradedCount">{{ degradedCount }} 降级</span>
      </div>
    </div>

    <div class="mhp-grid">
      <div
        v-for="m in store.models"
        :key="m.id"
        class="model-card"
        :class="{ 'is-degraded': m.status === 'degraded', 'is-offline': m.status === 'offline' }"
      >
        <!-- 顶部：状态 + 名称 + 提供商 -->
        <div class="mc-top">
          <StatusDot :status="m.status" />
          <div class="mc-identity">
            <span class="mc-name">{{ m.display_name || m.name }}</span>
            <span class="mc-provider">{{ m.provider }}</span>
          </div>
          <div class="mc-pricing">
            <span class="mc-price-tag">入 ¥{{ m.input_price }}</span>
            <span class="mc-price-tag">出 ¥{{ m.output_price }}</span>
          </div>
        </div>

        <!-- 中部：用量微型柱 -->
        <div class="mc-usage">
          <div class="mc-usage-bar">
            <div
              class="mc-usage-fill"
              :style="{ width: usagePercent(m) + '%' }"
            ></div>
          </div>
          <span class="mc-usage-tokens">{{ fmtCompactNum(m.today_tokens) }}</span>
        </div>

        <!-- 底部：关键指标 -->
        <div class="mc-bottom">
          <div class="mc-stat">
            <span class="mc-stat-label">费用</span>
            <span class="mc-stat-value cost">¥{{ m.today_cost.toFixed(4) }}</span>
          </div>
          <div class="mc-stat">
            <span class="mc-stat-label">延迟</span>
            <span class="mc-stat-value">{{ m.avg_latency_ms }}ms</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useDashboard } from '../../stores/dashboard'
import StatusDot from '../common/StatusDot.vue'
import type { ModelInfo } from '../../types'

const store = useDashboard()

const onlineCount = computed(() => store.models.filter((m) => m.status === 'online').length)
const degradedCount = computed(() => store.models.filter((m) => m.status !== 'online').length)

// 用量占总量的百分比（用于微柱宽度）
const maxTokens = computed(() => {
  const vals = store.models.map((m) => m.today_tokens)
  return Math.max(...vals, 1)
})

function usagePercent(m: ModelInfo): number {
  return Math.round((m.today_tokens / maxTokens.value) * 100)
}

function fmtCompactNum(n: number): string {
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + 'M'
  if (n >= 1_000) return (n / 1_000).toFixed(1) + 'K'
  return String(n)
}
</script>

<style scoped>
.mhp {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-5) var(--space-6);
}

/* ── Header ── */
.mhp-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-4);
}
.mhp-title {
  font-size: var(--text-section);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
}
.mhp-summary {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-caption);
  color: var(--color-text-secondary);
}
.summary-dot {
  width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0;
}
.summary-dot.online { background: var(--color-success); }
.summary-dot.degraded { background: var(--color-warning); }

/* ── Grid ── */
.mhp-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: var(--space-3);
}

/* ── Card ── */
.model-card {
  background: var(--color-bg);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  transition: border-color var(--transition-fast), background var(--transition-fast);
}
.model-card:hover {
  border-color: var(--color-primary-muted);
  background: var(--color-surface);
}

/* 异常模型卡片视觉区分 */
.model-card.is-degraded {
  border-left: 3px solid var(--color-warning);
  background: var(--color-warning-bg);
}
.model-card.is-offline {
  border-left: 3px solid var(--color-error);
  background: var(--color-error-bg);
  opacity: 0.75;
}

/* ── Top row ── */
.mc-top {
  display: flex;
  align-items: flex-start;
  gap: var(--space-2);
}
.mc-identity {
  flex: 1;
  min-width: 0;
}
.mc-name {
  display: block;
  font-size: var(--text-card-title);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.mc-provider {
  font-size: 11px;
  color: var(--color-text-tertiary);
}
.mc-pricing {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 1px;
}
.mc-price-tag {
  font-size: 10px;
  color: var(--color-text-tertiary);
  font-family: var(--font-mono);
  white-space: nowrap;
}

/* ── Usage bar ── */
.mc-usage {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}
.mc-usage-bar {
  flex: 1;
  height: 4px;
  background: var(--color-border);
  border-radius: var(--radius-full);
  overflow: hidden;
}
.mc-usage-fill {
  height: 100%;
  background: var(--color-primary);
  border-radius: var(--radius-full);
  transition: width var(--transition-slow);
  min-width: 2px;
}
.model-card.is-degraded .mc-usage-fill { background: var(--color-warning); }
.model-card.is-offline .mc-usage-fill { background: var(--color-error); }
.mc-usage-tokens {
  font-size: var(--text-caption);
  color: var(--color-text-secondary);
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

/* ── Bottom stats ── */
.mc-bottom {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-1);
}
.mc-stat {
  display: flex;
  flex-direction: column;
}
.mc-stat-label {
  font-size: 10px;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}
.mc-stat-value {
  font-size: var(--text-card-title);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
}
.mc-stat-value.cost { color: var(--color-primary); font-weight: var(--weight-bold); }
</style>

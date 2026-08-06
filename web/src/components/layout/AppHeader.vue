<template>
  <header class="app-header">
    <div class="header-left">
      <div class="brand-mark">AG</div>
      <div class="brand-text">
        <span class="brand-name">AI Gateway Console</span>
        <span class="brand-subtitle">企业大模型调用治理与监控</span>
      </div>
    </div>
    <div class="header-right">
      <div class="status-chip" :class="statusClass">
        <span class="status-pulse"></span>
        {{ store.health?.online_models ?? '-' }}/{{ store.health?.total_models ?? '-' }} 模型在线
      </div>
      <div class="cost-chip">今日 ¥{{ (store.overview?.today.total_cost ?? 0).toFixed(2) }}</div>
      <div class="time-chip">{{ store.lastUpdate || '--' }}</div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useDashboard } from '../../stores/dashboard'

const store = useDashboard()

const statusClass = computed(() => {
  const online = store.health?.online_models ?? 0
  const total = store.health?.total_models ?? 0
  if (total === 0) return 'status-unknown'
  if (online === total) return 'status-healthy'
  if (online === 0) return 'status-critical'
  return 'status-degraded'
})
</script>

<style scoped>
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--space-6);
  height: 56px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 0;
  z-index: 100;
}
.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}
.brand-mark {
  width: 32px;
  height: 32px;
  background: var(--color-primary);
  color: var(--color-text-inverse);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: var(--weight-bold);
  letter-spacing: 0.5px;
}
.brand-text {
  display: flex;
  flex-direction: column;
}
.brand-name {
  font-size: var(--text-section);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  line-height: 1.2;
}
.brand-subtitle {
  font-size: 11px;
  color: var(--color-text-tertiary);
  line-height: 1.2;
}
.header-right {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}
.status-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  font-size: var(--text-caption);
  font-weight: var(--weight-medium);
}
.status-pulse {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}
.status-healthy {
  background: var(--color-success-bg);
  color: var(--color-success);
}
.status-healthy .status-pulse { background: var(--color-success); }
.status-degraded {
  background: var(--color-warning-bg);
  color: var(--color-warning);
}
.status-degraded .status-pulse { background: var(--color-warning); }
.status-critical {
  background: var(--color-error-bg);
  color: var(--color-error);
}
.status-critical .status-pulse { background: var(--color-error); }
.status-unknown {
  background: var(--color-bg-alt);
  color: var(--color-text-tertiary);
}
.status-unknown .status-pulse { background: var(--color-text-tertiary); }

.cost-chip {
  padding: 4px 12px;
  border-radius: var(--radius-full);
  background: var(--color-primary-light);
  color: var(--color-primary);
  font-size: var(--text-caption);
  font-weight: var(--weight-semibold);
}
.time-chip {
  padding: 4px 12px;
  border-radius: var(--radius-full);
  background: var(--color-bg-alt);
  color: var(--color-text-tertiary);
  font-size: var(--text-caption);
}
</style>

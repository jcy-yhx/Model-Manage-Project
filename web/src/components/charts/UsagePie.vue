<template>
  <div class="usage-pie">
    <div class="up-header">
      <h3 class="up-title">模型用量分布</h3>
    </div>
    <v-chart :option="option" autoresize style="height: 220px" />
    <div class="up-legend">
      <div
        v-for="(m, i) in store.modelUsage?.models ?? []"
        :key="m.model_name"
        class="up-legend-item"
      >
        <span class="up-legend-dot" :style="{ background: CHART_COLORS[i] }"></span>
        <span class="up-legend-name">{{ shortName(m.model_name) }}</span>
        <span class="up-legend-pct">{{ (m.percentage * 100).toFixed(0) }}%</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { useDashboard } from '../../stores/dashboard'

const CHART_COLORS = ['#3B82F6', '#6366F1', '#8B5CF6', '#06B6D4', '#10B981']

const store = useDashboard()

const option = computed(() => {
  const data = (store.modelUsage?.models ?? []).map((m, i) => ({
    name: m.model_name,
    value: m.tokens,
    itemStyle: { color: CHART_COLORS[i] },
  }))

  return {
    tooltip: {
      trigger: 'item' as const,
      backgroundColor: '#fff',
      borderColor: '#E2E8F0',
      borderWidth: 1,
      textStyle: { color: '#0F172A', fontSize: 13 },
      formatter: (p: any) => {
        return `<b>${p.name}</b><br/>Token: ${p.value >= 1000 ? (p.value/1000).toFixed(1)+'K' : p.value}<br/>占比: ${p.percent}%`
      },
    },
    series: [{
      type: 'pie',
      radius: ['55%', '82%'],
      center: ['50%', '48%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderColor: '#fff',
        borderWidth: 2,
        borderRadius: 2,
      },
      label: { show: false },
      emphasis: {
        scale: true,
        scaleSize: 6,
        label: { show: true, fontSize: 14, fontWeight: 'bold' as const },
      },
      data,
    }],
  }
})

function shortName(name: string): string {
  return name.replace('qwen2.5-', '').replace('minimax-', '').replace('deepseek-', '')
}
</script>

<style scoped>
.usage-pie {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-5) var(--space-6);
}
.up-header {
  margin-bottom: var(--space-1);
}
.up-title {
  font-size: var(--text-section);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
}
.up-legend {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
  justify-content: center;
  margin-top: var(--space-2);
}
.up-legend-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: var(--text-caption);
  color: var(--color-text-secondary);
}
.up-legend-dot {
  width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0;
}
.up-legend-name {
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.up-legend-pct {
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
}
</style>

<template>
  <div class="trend-panel">
    <div class="trend-header">
      <h3 class="trend-title">调用趋势</h3>
      <span class="trend-subtitle">近 7 天</span>
    </div>
    <div class="chart-wrapper">
      <v-chart v-if="hasData" :option="option" :echarts="echarts" autoresize style="height: 260px; width: 100%" />
      <div v-else class="chart-empty">暂无趋势数据，模拟器运行一段时间后自动生成</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import * as echarts from 'echarts'
import { useDashboard } from '../../stores/dashboard'
import type { TrendPoint } from '../../types'

const store = useDashboard()

const hasData = computed(() => (store.trend ?? []).length > 0)

const option = computed(() => {
  const pts: TrendPoint[] = store.trend ?? []
  if (pts.length === 0) return {}

  return {
    tooltip: {
      trigger: 'axis' as const,
      backgroundColor: '#fff',
      borderColor: '#E2E8F0',
      borderWidth: 1,
      textStyle: { color: '#0F172A', fontSize: 13 },
    },
    legend: {
      data: ['Token', '费用估算'],
      bottom: 0,
      textStyle: { color: '#64748B', fontSize: 12 },
    },
    grid: {
      left: 52,
      right: 52,
      top: 16,
      bottom: 36,
    },
    xAxis: {
      type: 'category' as const,
      data: pts.map((p) => {
        const d = typeof p.date === 'string' ? p.date : String(p.date)
        return d.slice(5, 10)
      }),
      axisLine: { lineStyle: { color: '#E2E8F0' } },
      axisTick: { show: false },
      axisLabel: { color: '#94A3B8', fontSize: 12 },
    },
    yAxis: [
      {
        type: 'value' as const,
        name: 'Token',
        nameTextStyle: { color: '#94A3B8', fontSize: 11 },
        splitLine: { lineStyle: { color: '#F1F5F9', type: 'dashed' as const } },
        axisLabel: {
          color: '#94A3B8',
          fontSize: 11,
          formatter: (v: number) => (v >= 1000 ? (v / 1000).toFixed(0) + 'K' : String(v)),
        },
      },
    ],
    series: [
      {
        name: 'Token',
        type: 'line',
        data: pts.map((p) => p.tokens),
        smooth: pts.length > 2,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: { color: '#3B82F6', width: 2 },
        itemStyle: { color: '#3B82F6' },
        areaStyle: {
          color: {
            type: 'linear' as const,
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(59,130,246,0.12)' },
              { offset: 1, color: 'rgba(59,130,246,0.0)' },
            ],
          },
        },
      },
    ],
  }
})
</script>

<style scoped>
.trend-panel {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-5) var(--space-6);
}
.trend-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: var(--space-1);
}
.trend-title {
  font-size: var(--text-section);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
}
.trend-subtitle {
  font-size: var(--text-caption);
  color: var(--color-text-tertiary);
}
.chart-wrapper {
  min-height: 260px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.chart-empty {
  color: var(--color-text-tertiary);
  font-size: var(--text-body);
}
</style>

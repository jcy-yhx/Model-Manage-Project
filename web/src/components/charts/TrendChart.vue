<template>
  <div class="trend-panel">
    <div class="trend-header">
      <h3 class="trend-title">调用趋势</h3>
      <span class="trend-subtitle">近 7 天</span>
    </div>
    <v-chart :option="option" autoresize style="height: 260px" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { useDashboard } from '../../stores/dashboard'
import type { TrendPoint } from '../../types'

const store = useDashboard()

const option = computed(() => {
  const pts: TrendPoint[] = store.trend ?? []

  return {
    tooltip: {
      trigger: 'axis' as const,
      backgroundColor: '#fff',
      borderColor: '#E2E8F0',
      borderWidth: 1,
      textStyle: { color: '#0F172A', fontSize: 13 },
      // 自定义 tooltip 格式
      formatter: (params: any) => {
        if (!params?.length) return ''
        const date = pts[params[0]?.dataIndex]?.date?.slice(0, 10) ?? ''
        let html = `<div style="font-weight:600;margin-bottom:6px">${date}</div>`
        for (const p of params) {
          const val = p.value ?? 0
          const display = p.seriesName === 'Token'
            ? val >= 1e6 ? (val / 1e6).toFixed(1) + 'M' : val >= 1e3 ? (val / 1e3).toFixed(1) + 'K' : String(val)
            : '¥' + (val ?? 0).toFixed(2)
          html += `<div style="display:flex;align-items:center;gap:6px;margin:2px 0">
            <span style="display:inline-block;width:8px;height:8px;border-radius:50%;background:${p.color}"></span>
            ${p.seriesName}: <b>${display}</b>
          </div>`
        }
        return html
      },
    },
    legend: {
      data: ['Token', '费用'],
      bottom: 0,
      textStyle: { color: '#64748B', fontSize: 12 },
      itemWidth: 10,
      itemHeight: 10,
      itemGap: 20,
    },
    grid: {
      left: 52,
      right: 52,
      top: 16,
      bottom: 36,
    },
    xAxis: {
      type: 'category' as const,
      data: pts.map((p) => p.date.slice(5)),
      axisLine: { lineStyle: { color: '#E2E8F0' } },
      axisTick: { show: false },
      axisLabel: { color: '#94A3B8', fontSize: 12 },
    },
    yAxis: [
      {
        type: 'value' as const,
        name: 'Token',
        nameTextStyle: { color: '#94A3B8', fontSize: 11 },
        splitLine: { lineStyle: { color: '#F1F5F9', type: 'dashed' } },
        axisLabel: {
          color: '#94A3B8',
          fontSize: 11,
          formatter: (v: number) => v >= 1e3 ? (v / 1e3).toFixed(0) + 'K' : String(v),
        },
      },
      {
        type: 'value' as const,
        name: '费用 ¥',
        nameTextStyle: { color: '#94A3B8', fontSize: 11 },
        splitLine: { show: false },
        axisLabel: { color: '#94A3B8', fontSize: 11 },
      },
    ],
    series: [
      {
        name: 'Token',
        type: 'line',
        data: pts.map((p) => p.tokens),
        smooth: true,
        symbol: 'circle',
        symbolSize: 4,
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
      {
        name: '费用',
        type: 'line',
        yAxisIndex: 1,
        // 费用数据从 Redis cost key 读取——当前 trend API 不返回 cost，
        // 用 Token 估算费用做近似的视觉展示
        data: pts.map((p) => {
          // 简单估算：Token × 平均价格 ~0.008/K
          return Math.round(p.tokens * 0.008) / 1000
        }),
        smooth: true,
        symbol: 'circle',
        symbolSize: 4,
        lineStyle: { color: '#10B981', width: 2 },
        itemStyle: { color: '#10B981' },
        areaStyle: {
          color: {
            type: 'linear' as const,
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(16,185,129,0.10)' },
              { offset: 1, color: 'rgba(16,185,129,0.0)' },
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
</style>

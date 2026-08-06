<template>
  <div class="trend-chart">
    <h3>近 7 天趋势</h3>
    <v-chart :option="option" autoresize style="height: 240px" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { useDashboard } from '../../stores/dashboard'
import type { TrendPoint } from '../../types'

use([LineChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const store = useDashboard()

const option = computed(() => {
  const pts: TrendPoint[] = store.trend ?? []
  return {
    tooltip: { trigger: 'axis' as const },
    legend: { data: ['Token', '请求数'], bottom: 0 },
    grid: { left: 50, right: 20, top: 10, bottom: 30 },
    xAxis: { type: 'category' as const, data: pts.map((p) => p.date.slice(5)) },
    yAxis: [
      { type: 'value' as const, name: 'Token' },
      { type: 'value' as const, name: '请求数' },
    ],
    series: [
      { name: 'Token', type: 'line', data: pts.map((p) => p.tokens), smooth: true, itemStyle: { color: '#1677ff' } },
      { name: '请求数', type: 'line', yAxisIndex: 1, data: pts.map((p) => p.requests), smooth: true, itemStyle: { color: '#52c41a' } },
    ],
  }
})
</script>

<style scoped>
.trend-chart { background: #fff; border-radius: 10px; padding: 16px; box-shadow: 0 1px 4px rgba(0,0,0,.08); margin-top: 12px; }
.trend-chart h3 { margin: 0 0 6px; font-size: 15px; }
</style>

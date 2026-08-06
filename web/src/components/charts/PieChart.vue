<template>
  <div class="pie-chart">
    <h3>模型使用占比</h3>
    <v-chart :option="option" autoresize style="height: 220px" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { useDashboard } from '../../stores/dashboard'

const store = useDashboard()

const option = computed(() => ({
  tooltip: { trigger: 'item' as const },
  legend: { orient: 'vertical' as const, right: 10, top: 'center' },
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    data: (store.modelUsage?.models ?? []).map((m) => ({
      name: m.model_name,
      value: m.tokens,
    })),
    emphasis: { itemStyle: { shadowBlur: 10, shadowColor: 'rgba(0,0,0,.2)' } },
  }],
}))
</script>

<style scoped>
.pie-chart { background: #fff; border-radius: 10px; padding: 16px; box-shadow: 0 1px 4px rgba(0,0,0,.08); }
.pie-chart h3 { margin: 0 0 6px; font-size: 15px; }
</style>

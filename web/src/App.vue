<template>
  <div class="dashboard">
    <HealthBar />
    <MetricCards />
    <div class="middle-row">
      <div class="middle-left"><ModelCards /></div>
      <div class="middle-right"><PieChart /></div>
    </div>
    <TrendChart />
    <RealtimeLog />
    <TestPanel />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useDashboard } from './stores/dashboard'
import HealthBar from './components/layout/HealthBar.vue'
import MetricCards from './components/dashboard/MetricCards.vue'
import ModelCards from './components/dashboard/ModelCards.vue'
import RealtimeLog from './components/dashboard/RealtimeLog.vue'
import TrendChart from './components/dashboard/TrendChart.vue'
import TestPanel from './components/dashboard/TestPanel.vue'
import PieChart from './components/charts/PieChart.vue'

const store = useDashboard()

onMounted(() => store.startPolling(4000))
onUnmounted(() => store.stopPolling())
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f5f7fa; color: #333; }
.dashboard { max-width: 1400px; margin: 0 auto; padding: 16px; }
.middle-row { display: grid; grid-template-columns: 1fr 280px; gap: 12px; margin-top: 12px; }
</style>

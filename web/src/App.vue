<template>
  <div class="app-shell">
    <AppHeader />
    <main class="main-content">
      <!-- Hero 指标区域 -->
      <section class="section section-hero">
        <HeroMetrics />
      </section>

      <!-- 模型健康 -->
      <section class="section section-health">
        <ModelHealthPanel />
      </section>

      <!-- 用量分布 -->
      <section class="section section-pie">
        <PieChart />
      </section>

      <!-- 趋势分析 -->
      <section class="section section-analytics">
        <TrendChart />
      </section>

      <!-- 实时日志 -->
      <section class="section section-logs">
        <RealtimeLog />
      </section>

      <!-- API Playground -->
      <section class="section section-playground">
        <Playground />
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useDashboard } from './stores/dashboard'
import AppHeader from './components/layout/AppHeader.vue'
import HeroMetrics from './components/dashboard/HeroMetrics.vue'
import ModelHealthPanel from './components/dashboard/ModelHealthPanel.vue'
import RealtimeLog from './components/dashboard/RealtimeLog.vue'
import TrendChart from './components/dashboard/TrendChart.vue'
import Playground from './components/dashboard/Playground.vue'
import PieChart from './components/charts/PieChart.vue'

const store = useDashboard()

onMounted(() => store.startPolling(4000))
onUnmounted(() => store.stopPolling())
</script>

<style>
@import './styles/variables.css';

* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: var(--font-family);
  background: var(--color-bg);
  color: var(--color-text-primary);
  font-size: var(--text-body);
  line-height: 1.5;
  -webkit-font-smoothing: antialiased;
}

.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  max-width: 1280px;
  width: 100%;
  margin: 0 auto;
  padding: var(--space-6);
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.section { }
</style>

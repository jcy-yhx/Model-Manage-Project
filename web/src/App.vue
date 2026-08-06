<template>
  <div class="app-shell">
    <AppHeader />
    <main class="main-content">
      <!-- Hero 指标区域 (placeholder: 当前 MetricCards, 阶段 6-2 替换为 HeroMetrics) -->
      <section class="section section-hero">
        <MetricCards />
      </section>

      <!-- 模型健康 + 用量分布 -->
      <section class="section section-models">
        <ModelCards />
        <PieChart class="pie-side" />
      </section>

      <!-- 趋势分析 -->
      <section class="section section-analytics">
        <TrendChart />
      </section>

      <!-- 实时日志 -->
      <section class="section section-logs">
        <RealtimeLog />
      </section>

      <!-- API Playground (placeholder: 当前 TestPanel, 阶段 6-3 替换为 Playground) -->
      <section class="section section-playground">
        <TestPanel />
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useDashboard } from './stores/dashboard'
import AppHeader from './components/layout/AppHeader.vue'
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

.section-models {
  display: grid;
  grid-template-columns: 1fr 260px;
  gap: var(--space-5);
  align-items: start;
}

.pie-side { }
</style>

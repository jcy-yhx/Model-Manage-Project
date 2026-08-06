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

      <!-- 趋势分析 + 用量分布 -->
      <section class="section section-analytics">
        <div class="analytics-grid">
          <TrendChart class="analytics-main" />
          <UsagePie class="analytics-side" />
        </div>
      </section>

      <!-- 实时日志 -->
      <section class="section section-logs">
        <RealtimeLogs />
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
import RealtimeLogs from './components/dashboard/RealtimeLogs.vue'
import TrendChart from './components/charts/TrendChart.vue'
import Playground from './components/dashboard/Playground.vue'
import UsagePie from './components/charts/UsagePie.vue'

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

.analytics-grid {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: var(--space-5);
  align-items: stretch;
}
.analytics-main { }
.analytics-side { }
</style>

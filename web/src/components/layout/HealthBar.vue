<template>
  <div class="health-bar">
    <span class="status-indicator">
      <span class="dot online"></span>
      {{ store.health?.online_models ?? '-' }} / {{ store.health?.total_models ?? '-' }} 模型在线
    </span>
    <span v-if="alerts.length" class="alerts">
      <span v-for="a in alerts" :key="a.model" class="alert-tag">⚠ {{ a.model }} {{ a.message }}</span>
    </span>
    <span class="update-time">最后更新: {{ store.lastUpdate || '--' }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useDashboard } from '../../stores/dashboard'

const store = useDashboard()
const alerts = computed(() => store.health?.alerts ?? [])
</script>

<style scoped>
.health-bar {
  display: flex; align-items: center; gap: 16px;
  padding: 10px 20px; background: #f0f6ff; border-radius: 8px;
  font-size: 14px; color: #333;
}
.dot { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-right: 6px; }
.dot.online { background: #52c41a; }
.alert-tag { background: #fff3cd; color: #856404; padding: 2px 8px; border-radius: 4px; font-size: 12px; }
.update-time { margin-left: auto; color: #999; }
</style>

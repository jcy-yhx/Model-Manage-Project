<template>
  <div class="model-cards">
    <div v-for="m in store.models" :key="m.id" class="model-card">
      <div class="header">
        <StatusDot :status="m.status" />
        <span class="name">{{ m.display_name || m.name }}</span>
        <span class="provider">{{ m.provider }}</span>
      </div>
      <div class="prices">
        <span>入 ¥{{ m.input_price }}/1K</span>
        <span>出 ¥{{ m.output_price }}/1K</span>
      </div>
      <div class="stats">
        <div><span class="num">{{ fmtNum(m.today_tokens) }}</span><span class="unit">Tokens</span></div>
        <div><span class="num">¥{{ m.today_cost.toFixed(4) }}</span><span class="unit">费用</span></div>
        <div><span class="num">{{ m.avg_latency_ms }}ms</span><span class="unit">延迟</span></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDashboard } from '../../stores/dashboard'
import StatusDot from '../common/StatusDot.vue'

const store = useDashboard()

function fmtNum(n: number): string {
  if (n >= 1e3) return (n / 1e3).toFixed(1) + 'K'
  return String(n)
}
</script>

<style scoped>
.model-cards { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 12px; margin-top: 12px; }
.model-card { background: #fff; border-radius: 10px; padding: 14px; box-shadow: 0 1px 4px rgba(0,0,0,.08); }
.header { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.name { font-weight: 600; font-size: 15px; }
.provider { margin-left: auto; font-size: 11px; color: #999; padding: 1px 6px; background: #f5f5f5; border-radius: 4px; }
.prices { display: flex; gap: 12px; font-size: 11px; color: #888; margin-bottom: 10px; }
.stats { display: flex; gap: 10px; }
.stats div { flex: 1; text-align: center; }
.num { display: block; font-size: 16px; font-weight: 700; color: #333; }
.unit { font-size: 10px; color: #aaa; }
</style>

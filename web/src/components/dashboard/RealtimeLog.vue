<template>
  <div class="realtime-log">
    <h3>实时调用流</h3>
    <el-table :data="store.logs?.logs ?? []" size="small" max-height="300" stripe>
      <el-table-column prop="created_at" label="时间" width="160">
        <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column prop="project_name" label="项目" width="140" />
      <el-table-column prop="model_name" label="模型" width="140" />
      <el-table-column prop="total_tokens" label="Token" width="80" align="right" />
      <el-table-column prop="cost" label="费用" width="90" align="right">
        <template #default="{ row }">¥{{ row.cost.toFixed(6) }}</template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="latency_ms" label="延迟" width="70" align="right">
        <template #default="{ row }">{{ row.latency_ms }}ms</template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { useDashboard } from '../../stores/dashboard'

const store = useDashboard()

function fmtTime(t: string) {
  if (!t) return ''
  return new Date(t).toLocaleTimeString()
}
</script>

<style scoped>
.realtime-log { background: #fff; border-radius: 10px; padding: 16px; box-shadow: 0 1px 4px rgba(0,0,0,.08); margin-top: 12px; }
.realtime-log h3 { margin: 0 0 10px; font-size: 15px; }
</style>

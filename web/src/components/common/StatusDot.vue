<template>
  <span class="status-dot" :class="status"></span>
</template>

<script setup lang="ts">
defineProps<{ status: string }>()
</script>

<style scoped>
.status-dot { display: inline-block; width: 10px; height: 10px; border-radius: 50%; }
.status-dot.online { background: #52c41a; }
.status-dot.degraded { background: #faad14; }
.status-dot.offline { background: #ff4d4f; }
</style>

<template>
  <span class="status-dot" :class="status">
    <span class="dot-core"></span>
  </span>
</template>

<script setup lang="ts">
defineProps<{ status: string }>()
</script>

<style scoped>
.status-dot {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  flex-shrink: 0;
}
.dot-core {
  display: block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

/* online: green glow */
.status-dot.online { background: var(--color-success-bg); }
.status-dot.online .dot-core { background: var(--color-success); }

/* degraded: amber glow */
.status-dot.degraded { background: var(--color-warning-bg); }
.status-dot.degraded .dot-core { background: var(--color-warning); }

/* offline: red glow */
.status-dot.offline { background: var(--color-error-bg); }
.status-dot.offline .dot-core { background: var(--color-error); }
</style>

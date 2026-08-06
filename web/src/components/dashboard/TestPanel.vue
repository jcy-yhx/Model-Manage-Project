<template>
  <div class="test-panel">
    <h3>API 调用测试区</h3>
    <div class="test-body">
      <div class="test-inputs">
        <el-select v-model="selectedModel" placeholder="选择模型" size="default" style="width: 100%">
          <el-option v-for="m in store.models" :key="m.id" :label="m.display_name" :value="m.name" />
        </el-select>
        <el-input v-model="apiKey" placeholder="API Key" size="default" style="margin-top: 8px" />
        <el-input v-model="message" type="textarea" :rows="3" placeholder="输入消息..." size="default" style="margin-top: 8px" />
        <el-button type="primary" :loading="store.chatLoading" @click="send" style="margin-top: 8px; width: 100%">
          发送请求
        </el-button>
      </div>
      <div class="test-result">
        <div v-if="store.chatLoading">请求中...</div>
        <div v-else-if="store.chatError" class="error">{{ store.chatError }}</div>
        <div v-else-if="store.chatResult" class="result">
          <div class="result-meta">
            <span>{{ store.chatResult.model }}</span>
            <span>{{ store.chatResult.usage.total_tokens }} tokens</span>
            <span>¥{{ store.chatResult.usage.cost.toFixed(6) }}</span>
          </div>
          <div class="result-content">{{ store.chatResult.choices[0]?.message.content }}</div>
        </div>
        <div v-else class="placeholder">点击发送，体验一次真实 API 调用</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useDashboard } from '../../stores/dashboard'

const store = useDashboard()

const selectedModel = ref('qwen2.5-72b')
const apiKey = ref('sk-example-key-for-demo-20260805')
const message = ref('介绍一下量子计算')

function send() {
  store.sendChat(selectedModel.value, message.value, apiKey.value)
}
</script>

<style scoped>
.test-panel { background: #fff; border-radius: 10px; padding: 16px; box-shadow: 0 1px 4px rgba(0,0,0,.08); margin-top: 12px; }
.test-panel h3 { margin: 0 0 10px; font-size: 15px; }
.test-body { display: grid; grid-template-columns: 300px 1fr; gap: 16px; }
.test-result { background: #fafafa; border-radius: 8px; padding: 12px; min-height: 200px; font-size: 13px; }
.result-meta { display: flex; gap: 12px; margin-bottom: 8px; font-size: 12px; color: #1677ff; }
.result-content { white-space: pre-wrap; color: #333; line-height: 1.6; }
.error { color: #ff4d4f; }
.placeholder { color: #ccc; display: flex; align-items: center; justify-content: center; height: 180px; }
</style>

<template>
  <div class="playground">
    <div class="pg-header">
      <h3 class="pg-title">API Playground</h3>
      <span class="pg-desc">模拟大模型调用，实时查看 Token 消耗与费用</span>
    </div>

    <div class="pg-body">
      <!-- 左侧：输入区 -->
      <div class="pg-inputs">
        <div class="field">
          <label class="field-label">模型</label>
          <el-select
            v-model="selectedModel"
            placeholder="选择模型"
            class="pg-select"
            popper-class="pg-select-dropdown"
          >
            <el-option
              v-for="m in store.models"
              :key="m.id"
              :label="m.display_name"
              :value="m.name"
            >
              <div class="model-option">
                <span class="mo-name">{{ m.display_name }}</span>
                <span class="mo-provider">{{ m.provider }}</span>
                <span class="mo-price">¥{{ m.input_price }}/¥{{ m.output_price }}</span>
              </div>
            </el-option>
          </el-select>
        </div>

        <div class="field">
          <label class="field-label">消息</label>
          <textarea
            v-model="message"
            class="pg-textarea"
            placeholder="输入你想问的内容..."
            rows="4"
          ></textarea>
          <div class="quick-prompts">
            <button
              v-for="p in quickPrompts"
              :key="p"
              class="prompt-chip"
              @click="message = p"
              type="button"
            >
              {{ p }}
            </button>
          </div>
        </div>

        <div class="field field-key">
          <label class="field-label">API Key</label>
          <div class="key-row">
            <input
              :type="showKey ? 'text' : 'password'"
              v-model="apiKey"
              class="pg-input"
              placeholder="sk-..."
            />
            <button class="key-toggle" @click="showKey = !showKey" type="button">
              {{ showKey ? '隐藏' : '显示' }}
            </button>
          </div>
        </div>

        <button
          class="pg-send-btn"
          :class="{ loading: store.chatLoading }"
          :disabled="store.chatLoading || !message.trim()"
          @click="send"
        >
          <span v-if="store.chatLoading" class="btn-spinner"></span>
          <span v-else class="btn-icon">↵</span>
          {{ store.chatLoading ? '请求中...' : '发送请求' }}
        </button>
      </div>

      <!-- 右侧：结果区 -->
      <div class="pg-result" :class="{ 'has-content': store.chatResult, 'has-error': store.chatError }">
        <!-- 空状态 -->
        <div v-if="!store.chatResult && !store.chatError && !store.chatLoading" class="empty-state">
          <div class="empty-icon">⚡</div>
          <div class="empty-text">输入消息并点击发送</div>
          <div class="empty-hint">体验一次完整的 API 调用闭环</div>
        </div>

        <!-- 加载 -->
        <div v-if="store.chatLoading" class="loading-state">
          <div class="loading-dots"><span></span><span></span><span></span></div>
          <div class="loading-text">模型正在生成回复...</div>
        </div>

        <!-- 错误 -->
        <div v-if="store.chatError && !store.chatLoading" class="error-state">
          <div class="error-icon">!</div>
          <div class="error-text">{{ store.chatError }}</div>
        </div>

        <!-- 成功 -->
        <div v-if="store.chatResult && !store.chatLoading" class="success-state">
          <div class="result-meta-row">
            <div class="rm-item">
              <span class="rm-label">模型</span>
              <span class="rm-value">{{ store.chatResult.model }}</span>
            </div>
            <div class="rm-item">
              <span class="rm-label">Token 用量</span>
              <span class="rm-value">{{ store.chatResult.usage.total_tokens }}</span>
            </div>
            <div class="rm-item emphasize">
              <span class="rm-label">费用</span>
              <span class="rm-value">¥{{ store.chatResult.usage.cost.toFixed(6) }}</span>
            </div>
          </div>

          <div class="result-divider"></div>

          <div class="result-content">
            <div class="rc-header">
              <span class="rc-role">assistant</span>
              <span class="rc-reason">finish_reason: {{ store.chatResult.choices[0]?.finish_reason }}</span>
            </div>
            <div class="rc-body">{{ store.chatResult.choices[0]?.message.content }}</div>
          </div>

          <!-- Usage 详情 -->
          <details class="usage-detail">
            <summary>详细用量</summary>
            <div class="usage-grid">
              <div><span>prompt_tokens</span><code>{{ store.chatResult.usage.prompt_tokens }}</code></div>
              <div><span>completion_tokens</span><code>{{ store.chatResult.usage.completion_tokens }}</code></div>
              <div><span>total_tokens</span><code>{{ store.chatResult.usage.total_tokens }}</code></div>
              <div><span>cost</span><code>¥{{ store.chatResult.usage.cost.toFixed(6) }}</code></div>
            </div>
          </details>
        </div>
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
const message = ref('')
const showKey = ref(false)

const quickPrompts = [
  '介绍一下量子计算',
  '写一段 Python 快速排序代码',
  '什么是微服务架构？',
]

function send() {
  if (!message.value.trim()) return
  store.sendChat(selectedModel.value, message.value, apiKey.value)
}
</script>

<style scoped>
.playground {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

/* ── Header ── */
.pg-header {
  padding: var(--space-4) var(--space-6);
  border-bottom: 1px solid var(--color-border-light);
}
.pg-title {
  font-size: var(--text-section);
  font-weight: var(--weight-semibold);
  color: var(--color-text-primary);
  margin: 0 0 2px;
}
.pg-desc {
  font-size: var(--text-caption);
  color: var(--color-text-tertiary);
}

/* ── Body ── */
.pg-body {
  display: grid;
  grid-template-columns: 1fr 1fr;
  min-height: 360px;
}

/* ── Input side ── */
.pg-inputs {
  padding: var(--space-5) var(--space-6);
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  border-right: 1px solid var(--color-border-light);
}
.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.field-key { margin-top: auto; }
.field-label {
  font-size: var(--text-caption);
  font-weight: var(--weight-medium);
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}
.pg-select {
  width: 100%;
}
.pg-textarea {
  width: 100%;
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-family: var(--font-family);
  font-size: var(--text-body);
  color: var(--color-text-primary);
  resize: vertical;
  outline: none;
  transition: border-color var(--transition-fast);
}
.pg-textarea:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-muted);
}
.pg-textarea::placeholder { color: var(--color-text-tertiary); }

.quick-prompts {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-1);
}
.prompt-chip {
  font-size: 11px;
  color: var(--color-text-secondary);
  background: var(--color-bg-alt);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-full);
  padding: 3px 10px;
  cursor: pointer;
  transition: all var(--transition-fast);
}
.prompt-chip:hover {
  color: var(--color-primary);
  border-color: var(--color-primary);
  background: var(--color-primary-light);
}

.key-row {
  display: flex;
  gap: var(--space-2);
}
.pg-input {
  flex: 1;
  padding: 8px var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--color-text-primary);
  outline: none;
  transition: border-color var(--transition-fast);
}
.pg-input:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-muted);
}
.key-toggle {
  font-size: var(--text-caption);
  color: var(--color-text-secondary);
  background: none;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: 8px 12px;
  cursor: pointer;
  white-space: nowrap;
  transition: all var(--transition-fast);
}
.key-toggle:hover { border-color: var(--color-primary); color: var(--color-primary); }

.pg-send-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: 10px var(--space-4);
  background: var(--color-primary);
  color: var(--color-text-inverse);
  border: none;
  border-radius: var(--radius-sm);
  font-size: var(--text-body);
  font-weight: var(--weight-semibold);
  cursor: pointer;
  transition: all var(--transition-fast);
}
.pg-send-btn:hover:not(:disabled) { background: var(--color-primary-dark); }
.pg-send-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.pg-send-btn.loading { background: var(--color-primary); opacity: 0.85; }
.btn-icon { font-size: 18px; line-height: 1; }
.btn-spinner {
  width: 16px; height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* ── Result side ── */
.pg-result {
  padding: var(--space-5) var(--space-6);
  display: flex;
  flex-direction: column;
}
.pg-result.has-content { background: var(--color-bg); }
.pg-result.has-error { background: var(--color-error-bg); }

/* Empty */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
}
.empty-icon { font-size: 32px; }
.empty-text { font-size: var(--text-body); color: var(--color-text-secondary); }
.empty-hint { font-size: var(--text-caption); color: var(--color-text-tertiary); }

/* Loading */
.loading-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
}
.loading-dots span {
  display: inline-block; width: 6px; height: 6px; border-radius: 50%;
  background: var(--color-primary); margin: 0 3px;
  animation: bounce 0.6s infinite alternate;
}
.loading-dots span:nth-child(2) { animation-delay: 0.2s; }
.loading-dots span:nth-child(3) { animation-delay: 0.4s; }
@keyframes bounce { to { transform: translateY(-6px); opacity: 0.4; } }
.loading-text { font-size: var(--text-caption); color: var(--color-text-tertiary); }

/* Error */
.error-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
}
.error-icon {
  width: 40px; height: 40px; border-radius: 50%;
  background: var(--color-error); color: #fff;
  display: flex; align-items: center; justify-content: center;
  font-size: 20px; font-weight: var(--weight-bold);
}
.error-text { font-size: var(--text-body); color: var(--color-error); text-align: center; }

/* Success */
.success-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}
.result-meta-row {
  display: flex;
  gap: var(--space-5);
}
.rm-item { display: flex; flex-direction: column; gap: 2px; }
.rm-item.emphasize .rm-value {
  color: var(--color-primary);
  font-weight: var(--weight-bold);
}
.rm-label { font-size: 10px; color: var(--color-text-tertiary); text-transform: uppercase; letter-spacing: 0.3px; }
.rm-value { font-size: var(--text-card-title); font-weight: var(--weight-semibold); color: var(--color-text-primary); font-family: var(--font-mono); }

.result-divider {
  height: 1px;
  background: var(--color-border);
}

.result-content {
  background: var(--color-surface);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  overflow: hidden;
}
.rc-header {
  display: flex;
  justify-content: space-between;
  padding: var(--space-2) var(--space-3);
  background: var(--color-bg-alt);
  font-size: var(--text-caption);
}
.rc-role { color: var(--color-primary); font-weight: var(--weight-semibold); }
.rc-reason { color: var(--color-text-tertiary); }
.rc-body {
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-body);
  line-height: 1.7;
  white-space: pre-wrap;
  color: var(--color-text-primary);
  max-height: 200px;
  overflow-y: auto;
}

.usage-detail {
  font-size: var(--text-caption);
  color: var(--color-text-tertiary);
  cursor: pointer;
}
.usage-detail summary { margin-bottom: var(--space-2); }
.usage-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-1) var(--space-4);
}
.usage-grid div {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 2px 0;
}
.usage-grid code {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--color-text-primary);
  background: var(--color-bg-alt);
  padding: 1px 6px;
  border-radius: var(--radius-sm);
}

/* ── Model option in dropdown ── */
.model-option {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}
.mo-name { font-weight: var(--weight-semibold); }
.mo-provider { font-size: 11px; color: var(--color-text-tertiary); }
.mo-price { margin-left: auto; font-size: 11px; color: var(--color-primary); }
</style>

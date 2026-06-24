<template>
  <div class="json-view">
    <!-- 顶部工具栏 -->
    <header class="json-toolbar">
      <div class="toolbar-left">
        <button class="back-btn" @click="router.push('/')" title="返回工具箱">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"/></svg>
        </button>
        <h1 class="toolbar-title">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="4 7 4 4 20 4 20 7"/><line x1="9" y1="20" x2="15" y2="20"/><line x1="12" y1="4" x2="12" y2="20"/></svg>
          JSON 格式化
        </h1>
      </div>
      <div class="toolbar-center">
        <div class="btn-group">
          <button class="tool-btn primary" @click="formatJson(2)" :disabled="!inputJson.trim() || loading" title="缩进 2 空格">
            <svg v-if="loadingAction !== 'format2'" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="4 7 4 4 20 4 20 7"/><line x1="9" y1="20" x2="15" y2="20"/><line x1="12" y1="4" x2="12" y2="20"/></svg>
            <svg v-else class="spin" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
            格式化
          </button>
          <button class="tool-btn primary" @click="formatJson(4)" :disabled="!inputJson.trim() || loading" title="缩进 4 空格">
            <svg v-if="loadingAction === 'format4'" class="spin" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
            缩进 4
          </button>
          <button class="tool-btn" @click="minifyJson" :disabled="!inputJson.trim() || loading" title="压缩为单行">
            <svg v-if="loadingAction !== 'minify'" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="4 14 10 14 10 20"/><polyline points="20 10 14 10 14 4"/><line x1="14" y1="10" x2="21" y2="3"/><line x1="3" y1="21" x2="10" y2="14"/></svg>
            <svg v-else class="spin" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
            压缩
          </button>
          <div class="btn-separator"></div>
          <button class="tool-btn" @click="validateJson" :disabled="!inputJson.trim() || loading" title="校验 JSON 语法">
            <svg v-if="loadingAction !== 'validate'" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
            <svg v-else class="spin" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
            校验
          </button>
        </div>
      </div>
      <div class="toolbar-right">
        <span v-if="validateResult" class="validate-badge" :class="validateResult.valid ? 'ok' : 'err'">
          <svg v-if="validateResult.valid" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
          <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          {{ validateResult.valid ? '有效' : '无效' }}
        </span>
      </div>
    </header>

    <!-- 左右分栏主区域 -->
    <div class="json-split">
      <!-- 左侧：输入 -->
      <div class="split-pane input-pane">
        <div class="pane-header">
          <span class="pane-label">输入</span>
          <div class="pane-actions">
            <button class="pane-btn" @click="pasteFromClipboard" title="从剪贴板粘贴">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><rect x="8" y="2" width="8" height="4" rx="1" ry="1"/></svg>
              粘贴
            </button>
            <button class="pane-btn" @click="clearInput" title="清空输入">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              清空
            </button>
          </div>
        </div>
        <textarea
          class="json-editor"
          v-model="inputJson"
          placeholder='在此粘贴或输入 JSON…'
          spellcheck="false"
        ></textarea>
      </div>

      <!-- 分隔线 -->
      <div class="split-divider">
        <div class="divider-line"></div>
        <div class="divider-handle">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="8 4 4 8 8 12"/><polyline points="16 4 20 8 16 12"/></svg>
        </div>
        <div class="divider-line"></div>
      </div>

      <!-- 右侧：输出 -->
      <div class="split-pane output-pane">
        <div class="pane-header">
          <span class="pane-label">输出</span>
          <div class="pane-actions">
            <span v-if="outputInfo" class="size-badge">{{ outputInfo }}</span>
            <button class="pane-btn" @click="copyOutput" :disabled="!outputJson" title="复制结果">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
              复制
            </button>
          </div>
        </div>
        <textarea
          class="json-editor output"
          v-model="outputJson"
          readonly
          placeholder='格式化结果将显示在这里…'
          spellcheck="false"
        ></textarea>
      </div>
    </div>

    <!-- 校验结果浮动条 -->
    <Transition name="slide-up">
      <div v-if="validateResult" class="validate-bar" :class="validateResult.valid ? 'success' : 'error'">
        <svg v-if="validateResult.valid" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
        <span v-if="validateResult.valid">
          JSON 格式正确 · {{ validateResult.keys }} 个键 · 嵌套深度 {{ validateResult.depth }}
        </span>
        <span v-else class="error-msg">{{ validateResult.error }}</span>
        <button class="bar-close" @click="validateResult = null">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { jsonFormatterApi, type JsonValidateData } from '@/services/api/jsonFormatter'
import { getApiErrorMessage } from '@/services/api/errorHandling'
import { showToast } from '@/composables/useToast'

const inputJson = ref('')
const outputJson = ref('')
const validateResult = ref<JsonValidateData | null>(null)
const router = useRouter()

/** 当前正在执行的操作标识，用于区分各按钮的 loading 图标 */
const loadingAction = ref<'format2' | 'format4' | 'minify' | 'validate' | null>(null)
/** 任意操作运行中时为 true，用于禁用所有工具按钮防重复提交 */
const loading = computed(() => loadingAction.value !== null)

const outputInfo = computed(() => {
  if (!outputJson.value) return ''
  const size = new Blob([outputJson.value]).size
  if (size > 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${size} B`
})

async function formatJson(indent: number) {
  const action = indent === 4 ? 'format4' : 'format2'
  if (loading.value) return
  loadingAction.value = action
  try {
    const result = await jsonFormatterApi.format(inputJson.value, indent, false)
    outputJson.value = result.output
    validateResult.value = null
  } catch (e) {
    showToast(getApiErrorMessage(e, '格式化失败'), 'err')
  } finally {
    loadingAction.value = null
  }
}

async function minifyJson() {
  if (loading.value) return
  loadingAction.value = 'minify'
  try {
    const result = await jsonFormatterApi.format(inputJson.value, 2, true)
    outputJson.value = result.output
    validateResult.value = null
  } catch (e) {
    showToast(getApiErrorMessage(e, '压缩失败'), 'err')
  } finally {
    loadingAction.value = null
  }
}

async function validateJson() {
  if (loading.value) return
  loadingAction.value = 'validate'
  try {
    const result = await jsonFormatterApi.validate(inputJson.value)
    validateResult.value = result
    if (result.valid) {
      showToast('JSON 格式正确', 'ok')
    } else {
      showToast('JSON 格式错误', 'err')
    }
  } catch (e) {
    showToast(getApiErrorMessage(e, '校验失败'), 'err')
  } finally {
    loadingAction.value = null
  }
}

async function pasteFromClipboard() {
  try {
    const text = await navigator.clipboard.readText()
    inputJson.value = text
  } catch {
    showToast('无法访问剪贴板', 'err')
  }
}

function clearInput() {
  inputJson.value = ''
  outputJson.value = ''
  validateResult.value = null
}

async function copyOutput() {
  if (!outputJson.value) return
  try {
    await navigator.clipboard.writeText(outputJson.value)
    showToast('已复制到剪贴板', 'ok')
  } catch {
    showToast('复制失败', 'err')
  }
}
</script>

<style scoped>
.json-view {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 56px);
  overflow: hidden;
  animation: pageEnter 0.4s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes pageEnter {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* ===== 顶部工具栏 ===== */
.json-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 52px;
  border-bottom: 1px solid var(--border);
  background: var(--surface);
  flex-shrink: 0;
  gap: 16px;
}

.toolbar-left {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.back-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.15s;
}

.back-btn:hover {
  color: var(--accent);
  background: var(--accent-soft);
}

.toolbar-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 700;
  color: var(--text);
  white-space: nowrap;
}

.toolbar-title svg {
  color: var(--accent);
}

.toolbar-center {
  flex-shrink: 0;
}

.btn-group {
  display: flex;
  align-items: center;
  gap: 2px;
  background: var(--surface-alt);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 3px;
}

.tool-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 6px 14px;
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: calc(var(--radius) - 2px);
  cursor: pointer;
  transition: all 0.15s;
  white-space: nowrap;
}

.tool-btn:hover:not(:disabled) {
  color: var(--text);
  background: var(--surface);
}

.tool-btn.primary {
  color: #fff;
  background: var(--accent);
}

.tool-btn.primary:hover:not(:disabled) {
  background: var(--accent-hover);
  box-shadow: 0 2px 8px rgba(100, 116, 139, 0.18);
}

.tool-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.btn-separator {
  width: 1px;
  height: 20px;
  background: var(--border);
  margin: 0 4px;
}

.toolbar-right {
  flex-shrink: 0;
  min-width: 60px;
  display: flex;
  justify-content: flex-end;
}

.validate-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 10px;
  font-size: 11.5px;
  font-weight: 700;
  border-radius: 20px;
  letter-spacing: 0.03em;
}

.validate-badge.ok {
  color: #059669;
  background: #ecfdf5;
  border: 1px solid #a7f3d0;
}

.validate-badge.err {
  color: #dc2626;
  background: #fef2f2;
  border: 1px solid #fecaca;
}

/* ===== 左右分栏 ===== */
.json-split {
  display: flex;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.split-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

.pane-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  height: 38px;
  border-bottom: 1px solid var(--border);
  background: var(--surface-alt);
  flex-shrink: 0;
}

.pane-label {
  font-size: 11.5px;
  font-weight: 700;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.pane-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.size-badge {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-tertiary);
  font-variant-numeric: tabular-nums;
  padding: 2px 8px;
  background: var(--surface);
  border-radius: 10px;
}

.pane-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  font-size: 11.5px;
  font-weight: 500;
  color: var(--text-tertiary);
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.15s;
}

.pane-btn:hover {
  color: var(--accent);
  background: var(--accent-soft);
}

.pane-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.json-editor {
  flex: 1;
  width: 100%;
  padding: 16px 20px;
  border: none;
  outline: none;
  resize: none;
  font-family: 'Cascadia Code', 'Fira Code', 'JetBrains Mono', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.65;
  color: var(--text);
  background: var(--surface);
  tab-size: 2;
  min-height: 0;
}

.json-editor.output {
  background: #f8fafc;
}

.json-editor::placeholder {
  color: var(--text-tertiary);
  font-style: italic;
}

/* ===== 分隔线 ===== */
.split-divider {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 24px;
  flex-shrink: 0;
  gap: 6px;
}

.divider-line {
  flex: 1;
  width: 1px;
  background: var(--border);
}

.divider-handle {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: var(--text-tertiary);
  background: var(--surface);
  border: 1px solid var(--border);
  flex-shrink: 0;
  transition: all 0.15s;
}

.divider-handle:hover {
  color: var(--accent);
  border-color: var(--accent);
  background: var(--accent-soft);
}

.divider-handle svg {
  transform: rotate(90deg);
}

/* ===== 校验结果浮动条 ===== */
.validate-bar {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 18px;
  border-radius: var(--radius);
  font-size: 13px;
  font-weight: 500;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12), 0 2px 8px rgba(0, 0, 0, 0.06);
  z-index: 100;
  max-width: 80vw;
}

.validate-bar.success {
  background: #065f46;
  color: #d1fae5;
  border: 1px solid #059669;
}

.validate-bar.error {
  background: #991b1b;
  color: #fecaca;
  border: 1px solid #dc2626;
}

.error-msg {
  word-break: break-all;
  max-width: 60vw;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bar-close {
  background: transparent;
  border: none;
  cursor: pointer;
  color: inherit;
  opacity: 0.6;
  padding: 2px;
  margin-left: 4px;
  transition: opacity 0.15s;
  flex-shrink: 0;
}

.bar-close:hover {
  opacity: 1;
}

.slide-up-enter-active { transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1); }
.slide-up-leave-active { transition: all 0.2s ease; }
.slide-up-enter-from, .slide-up-leave-to { opacity: 0; transform: translateX(-50%) translateY(16px); }

/* loading spin 动画 */
@keyframes spin {
  to { transform: rotate(360deg); }
}
.spin {
  animation: spin 0.7s linear infinite;
  opacity: 0.8;
}
</style>

<template>
  <header class="app-header">
    <!-- 主行：Logo + 标题 + 文件名 + 主题/登录 -->
    <div class="header-row main-row">
      <div class="app-header-left">
        <router-link to="/" class="app-logo" title="返回首页">
          <svg viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
        </router-link>
        <span class="app-title">Markdown → HTML</span>
      </div>
      <div class="app-header-right">
        <div class="filename-field" :class="{ focused: filenameFocus }">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
            <polyline points="14 2 14 8 20 8" />
          </svg>
          <input
            :value="editorStore.filename"
            @input="onFilenameInput"
            @focus="filenameFocus = true"
            @blur="filenameFocus = false"
            placeholder="文件名"
          >
          <span class="ext">.html</span>
        </div>
        <span class="toolbar-sep"></span>
        <button class="toolbar-btn icon-only" @click="appStore.toggleTheme()" :title="appStore.theme === 'light' ? '切换暗色' : '切换亮色'">
          <svg v-if="appStore.theme === 'light'" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
          <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>
        </button>
        <span v-if="authStore.isLoggedIn" class="user-name" :title="authStore.username || ''">{{ authStore.username }}</span>
        <button v-if="authStore.isLoggedIn" class="toolbar-btn ghost" @click="onLogout" title="退出登录">退出</button>
        <router-link v-else to="/login" class="toolbar-btn primary" title="登录">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
          <span class="btn-label">登录</span>
        </router-link>
      </div>
    </div>

    <!-- 副行：编辑辅助 + 转换输出 + 历史记录 -->
    <div class="header-row sub-row">
      <div class="group">
        <span class="group-label">编辑</span>
        <div class="group-btns">
          <button class="toolbar-btn" @click="loadExample" title="加载示例">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="9" y1="13" x2="15" y2="13"/><line x1="9" y1="17" x2="15" y2="17"/></svg>
            <span class="btn-label">示例</span>
          </button>
          <button class="toolbar-btn" :class="{ active: editorStore.findBoxOpen }" @click="editorStore.findBoxOpen = !editorStore.findBoxOpen" title="查找替换 (Ctrl+F)">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
            <span class="btn-label">查找</span>
          </button>
          <button class="toolbar-btn" @click="emit('importFile')" title="导入 .md 文件">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
            <span class="btn-label">导入</span>
          </button>
        </div>
      </div>

      <div class="group" v-if="authStore.isLoggedIn">
        <span class="group-label">记录</span>
        <div class="group-btns">
          <router-link to="/history" class="toolbar-btn" title="历史记录">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            <span class="btn-label">历史</span>
          </router-link>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAppStore } from '@/stores/app'
import { useEditorStore } from '@/stores/editor'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const appStore = useAppStore()
const editorStore = useEditorStore()
const authStore = useAuthStore()
const router = useRouter()

const filenameFocus = ref(false)

function onLogout() {
  authStore.logout()
  router.push('/md2html')
}

const emit = defineEmits<{
  loadExample: []
  importFile: []
}>()

function loadExample() { emit('loadExample') }

function onFilenameInput(e: Event) {
  const val = (e.target as HTMLInputElement).value
  editorStore.setFilename(val)
}
</script>

<style scoped>
.app-header {
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  user-select: none;
}

/* === 主行 === */
.main-row {
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  gap: 12px;
}

.app-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.app-logo {
  width: 24px;
  height: 24px;
  background: var(--text);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  cursor: pointer;
  transition: transform 0.12s;
}

.app-logo:hover { transform: scale(1.05); }

.app-logo svg {
  width: 13px;
  height: 13px;
  stroke: var(--surface);
  fill: none;
  stroke-width: 2.5;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.app-title {
  font-size: 13.5px;
  font-weight: 600;
  letter-spacing: -0.01em;
  color: var(--text);
}

.app-header-right {
  display: flex;
  align-items: center;
  gap: 6px;
}

.filename-field {
  display: flex;
  align-items: center;
  gap: 6px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  height: 28px;
  padding: 0 10px;
  background: var(--bg);
  font-size: 12px;
  transition: border-color 0.12s, box-shadow 0.12s;
  color: var(--text-tertiary);
}

.filename-field.focused {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-soft);
}

.filename-field input {
  border: none;
  outline: none;
  background: transparent;
  color: var(--text);
  font-size: 12px;
  width: 130px;
  font-family: var(--font-sans);
}

.filename-field .ext {
  color: var(--text-tertiary);
  font-size: 11px;
  font-family: var(--font-mono);
}

.user-name {
  font-size: 12px;
  color: var(--text-secondary);
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 0 4px;
}

.toolbar-btn.ghost {
  background: transparent;
  border-color: var(--border);
  color: var(--text-secondary);
}

.toolbar-btn.ghost:hover {
  background: var(--surface-hover);
  color: var(--text);
}

/* === 副行 === */
.sub-row {
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 6px 16px 8px;
  border-top: 1px solid var(--border);
  background: var(--bg);
  flex-wrap: wrap;
}

.group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.group-label {
  font-size: 10.5px;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  user-select: none;
}

.group-btns {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* === 通用按钮 === */
.toolbar-btn {
  height: 26px;
  padding: 0 9px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-family: var(--font-sans);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  transition: all 0.12s ease;
  white-space: nowrap;
  text-decoration: none;
}

.toolbar-btn:hover {
  border-color: var(--border-strong);
  color: var(--text);
  background: var(--surface-hover);
}

.toolbar-btn:active { transform: translateY(0.5px); }

.toolbar-btn.primary {
  background: var(--text);
  color: var(--surface);
  border-color: var(--text);
}

.toolbar-btn.primary:hover {
  background: var(--text-secondary);
  border-color: var(--text-secondary);
  color: var(--surface);
}

.toolbar-btn.icon-only {
  width: 28px;
  padding: 0;
  justify-content: center;
  height: 28px;
}

.toolbar-btn.icon-only svg {
  width: 14px;
  height: 14px;
}

.toolbar-btn.active {
  background: var(--accent-soft);
  border-color: var(--accent);
  color: var(--accent);
}

.toolbar-sep {
  width: 1px;
  height: 18px;
  background: var(--border);
  margin: 0 2px;
}

@media (max-width: 900px) {
  .btn-label { display: none; }
  .filename-field input { width: 80px; }
  .group-label { display: none; }
  .sub-row { gap: 8px; }
}
</style>

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
        <button
          class="save-status"
          :class="editorStore.saveStatus"
          :title="saveStatusTitle"
          @click="emit('retry-save')"
        >
          <span class="status-dot"></span>
          <span>{{ saveStatusText }}</span>
        </button>
        <span class="toolbar-sep"></span>
        <router-link v-if="authStore.isLoggedIn" to="/history" class="toolbar-btn" title="历史记录">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
          <span class="btn-label">历史</span>
        </router-link>
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
  </header>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
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
  'retry-save': []
}>()

const saveStatusText = computed(() => {
  if (!authStore.isLoggedIn) return '仅本地'
  if (editorStore.saveStatus === 'saving') return '保存中'
  if (editorStore.saveStatus === 'saved') return '已云端保存'
  if (editorStore.saveStatus === 'error') return '保存失败'
  return '本地已保存'
})

const saveStatusTitle = computed(() => {
  if (!authStore.isLoggedIn) return '未登录，当前内容仅保存到本地浏览器'
  if (editorStore.saveStatus === 'error') return '自动保存失败，点击后等待下次内容变更或导出时重试'
  return saveStatusText.value
})

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
  box-shadow: var(--shadow-sm);
  position: relative;
  z-index: 20;
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
  width: 26px;
  height: 26px;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-hover) 100%);
  border-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  cursor: pointer;
  transition: transform 0.2s var(--ease-spring), box-shadow 0.2s;
  box-shadow: 0 1px 3px rgba(100, 116, 139, 0.15);
}

.app-logo:hover {
  transform: scale(1.08);
  box-shadow: 0 2px 8px rgba(100, 116, 139, 0.2);
}

.app-logo svg {
  width: 13px;
  height: 13px;
  stroke: #fff;
  fill: none;
  stroke-width: 2.5;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.app-title {
  font-size: 13.5px;
  font-weight: 700;
  letter-spacing: -0.02em;
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
  transition: border-color 0.2s var(--ease-out), box-shadow 0.2s var(--ease-out), background 0.2s;
  color: var(--text-tertiary);
}

.filename-field.focused {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-soft);
  background: var(--surface);
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

.save-status {
  height: 24px;
  padding: 0 8px;
  border: 1px solid transparent;
  background: transparent;
  color: var(--text-tertiary);
  border-radius: var(--radius-sm);
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-family: var(--font-sans);
  cursor: default;
  white-space: nowrap;
}

.save-status.error {
  cursor: pointer;
  color: var(--danger);
  background: rgba(180, 86, 86, 0.08);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--text-tertiary);
}

.save-status.saved .status-dot { background: var(--ok); }
.save-status.saving .status-dot {
  background: var(--accent);
  animation: savePulse 1s ease-in-out infinite;
}
.save-status.error .status-dot { background: var(--danger); }
.save-status.local .status-dot { background: var(--text-tertiary); }

@keyframes savePulse {
  0%, 100% { opacity: 0.45; transform: scale(0.85); }
  50% { opacity: 1; transform: scale(1); }
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
  transition: all 0.15s var(--ease-out);
  white-space: nowrap;
  text-decoration: none;
  position: relative;
}

.toolbar-btn:hover {
  border-color: var(--border-strong);
  color: var(--text);
  background: var(--surface-hover);
  box-shadow: var(--shadow-sm);
}

.toolbar-btn:active {
  transform: scale(0.97);
  transition-duration: 0.08s;
}

.toolbar-btn.primary {
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-hover) 100%);
  color: #fff;
  border-color: var(--accent);
  box-shadow: 0 1px 3px rgba(100, 116, 139, 0.15);
}

.toolbar-btn.primary:hover {
  box-shadow: 0 2px 8px rgba(100, 116, 139, 0.2);
  filter: brightness(1.05);
  color: #fff;
  border-color: var(--accent-hover);
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
  box-shadow: 0 0 0 2px var(--accent-soft);
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
  .save-status span:not(.status-dot) { display: none; }
  .save-status { padding: 0 4px; }
}
</style>

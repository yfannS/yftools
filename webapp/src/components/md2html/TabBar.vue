<template>
  <div class="tab-bar" v-if="editorStore.tabs.length > 0">
    <div class="tab-list" ref="tabListRef">
      <div
        v-for="tab in editorStore.tabs"
        :key="tab.id"
        class="tab-item"
        :class="{ active: tab.id === editorStore.activeTabId }"
        @click="editorStore.switchTab(tab.id)"
        @dblclick="startRename(tab.id)"
        @contextmenu.prevent="onContextMenu($event, tab.id)"
      >
        <svg class="tab-icon" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14 2 14 8 20 8"/>
        </svg>
        <span v-if="renamingTabId !== tab.id" class="tab-name" :title="tab.filename">{{ tab.filename }}</span>
        <input
          v-else
          class="tab-rename-input"
          v-model="renameValue"
          @blur="finishRename"
          @keydown.enter="finishRename"
          @keydown.escape="cancelRename"
          ref="renameInputRef"
        />
        <span class="tab-dot" v-if="!tab.renderedHtml && tab.content.trim()"></span>
        <button
          class="tab-close"
          @click.stop="editorStore.closeTab(tab.id)"
          v-if="editorStore.tabs.length > 1"
          title="关闭"
        >&times;</button>
      </div>
    </div>
    <button class="tab-new" @click="editorStore.createTab()" title="新建标签">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useEditorStore } from '@/stores/editor'
import { showConfirm } from '@/composables/useConfirm'

const editorStore = useEditorStore()

const renamingTabId = ref<string | null>(null)
const renameValue = ref('')
const renameInputRef = ref<HTMLInputElement>()

function startRename(tabId: string) {
  const tab = editorStore.tabs.find(t => t.id === tabId)
  if (!tab) return
  renamingTabId.value = tabId
  renameValue.value = tab.filename
  nextTick(() => {
    renameInputRef.value?.focus()
    renameInputRef.value?.select()
  })
}

function finishRename() {
  if (renamingTabId.value && renameValue.value.trim()) {
    const savedId = editorStore.activeTabId
    editorStore.switchTab(renamingTabId.value)
    editorStore.setFilename(renameValue.value.trim())
    if (savedId !== renamingTabId.value) {
      editorStore.switchTab(savedId)
    }
  }
  renamingTabId.value = null
}

function cancelRename() {
  renamingTabId.value = null
}

async function onContextMenu(_e: MouseEvent, tabId: string) {
  const ok = await showConfirm({ message: '关闭其他标签？', confirmText: '关闭' })
  if (ok) {
    editorStore.closeOtherTabs(tabId)
  }
}
</script>

<style scoped>
.tab-bar {
  display: flex;
  align-items: center;
  background: var(--surface-raised);
  border-bottom: 1px solid var(--border);
  height: 34px;
  flex-shrink: 0;
  user-select: none;
}

.tab-list {
  display: flex;
  align-items: center;
  flex: 1;
  overflow-x: auto;
  overflow-y: hidden;
  gap: 0;
  scrollbar-width: none;
}

.tab-list::-webkit-scrollbar { display: none; }

.tab-item {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 0 14px;
  height: 34px;
  font-size: 12px;
  color: var(--text-tertiary);
  cursor: pointer;
  border-right: 1px solid var(--border);
  white-space: nowrap;
  min-width: 0;
  max-width: 160px;
  transition: background 0.15s var(--ease-out), color 0.15s var(--ease-out);
  position: relative;
}

.tab-item:hover {
  background: var(--surface);
  color: var(--text-secondary);
}

.tab-item.active {
  background: var(--surface);
  color: var(--text);
  border-bottom: 2px solid var(--accent);
  box-shadow: inset 0 -2px 0 var(--accent);
}

.tab-icon {
  flex-shrink: 0;
  opacity: 0.4;
  transition: opacity 0.15s, color 0.15s;
}

.tab-item.active .tab-icon {
  opacity: 1;
  color: var(--accent);
}

.tab-name {
  overflow: hidden;
  text-overflow: ellipsis;
}

.tab-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--accent);
  flex-shrink: 0;
  opacity: 0.6;
  animation: tabDotPulse 2s ease-in-out infinite;
}

@keyframes tabDotPulse {
  0%, 100% { opacity: 0.4; }
  50% { opacity: 0.8; }
}

.tab-rename-input {
  border: none;
  outline: none;
  background: var(--bg);
  color: var(--text);
  font-size: 12px;
  font-family: inherit;
  width: 80px;
  padding: 0 2px;
  border-bottom: 1px solid var(--accent);
}

.tab-close {
  display: none;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border: none;
  background: transparent;
  color: var(--text-tertiary);
  font-size: 13px;
  line-height: 1;
  cursor: pointer;
  border-radius: 3px;
  flex-shrink: 0;
  padding: 0;
  margin-left: 2px;
  transition: background 0.12s, color 0.12s;
}

.tab-item:hover .tab-close {
  display: inline-flex;
}

.tab-close:hover {
  background: var(--danger);
  color: #fff;
}

.tab-new {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: none;
  background: transparent;
  color: var(--text-tertiary);
  cursor: pointer;
  flex-shrink: 0;
  transition: color 0.15s var(--ease-out), background 0.15s var(--ease-out);
}

.tab-new:hover {
  color: var(--accent);
  background: var(--accent-soft);
}
</style>

<template>
  <div class="outline-trigger" @mouseenter="open" @mouseleave="onLeave" ref="triggerRef">
    <button class="outline-fab" title="目录导航" aria-label="目录导航">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="8" y1="6" x2="21" y2="6" />
        <line x1="8" y1="12" x2="21" y2="12" />
        <line x1="8" y1="18" x2="21" y2="18" />
        <line x1="3" y1="6" x2="3.01" y2="6" />
        <line x1="3" y1="12" x2="3.01" y2="12" />
        <line x1="3" y1="18" x2="3.01" y2="18" />
      </svg>
    </button>
    <Teleport to="body">
      <div
        v-show="visible"
        class="outline-dropdown"
        :style="dropdownStyle"
        @mouseenter="cancelLeave"
        @mouseleave="onLeave"
        ref="dropdownRef"
      >
        <div class="outline-dropdown-header">
          <span class="outline-dropdown-title">目录</span>
          <span class="outline-dropdown-count">{{ outline.length }}</span>
        </div>
        <div class="outline-dropdown-body">
          <template v-if="!outline.length">
            <p class="empty-hint">暂无标题 · 使用 # 添加</p>
          </template>
          <ul v-else class="outline-list">
            <li
              v-for="item in outline"
              :key="item.id"
              :class="['outline-item', `lv-${item.level}`, { active: item.id === activeId }]"
              :style="{ paddingLeft: 12 + (item.level - 1) * 12 + 'px' }"
              @click="onJump(item.id)"
              :title="item.text"
            >
              <span class="dot" :class="`dot-${Math.min(item.level, 4)}`"></span>
              <span class="text">{{ item.text || '(无标题)' }}</span>
            </li>
          </ul>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { usePreviewStore } from '@/stores/preview'
import type { OutlineItem } from '@/stores/preview'

const previewStore = usePreviewStore()
const outline = computed<OutlineItem[]>(() => previewStore.outline)
const activeId = ref<string>('')

const visible = ref(false)

const triggerRef = ref<HTMLElement>()
const dropdownRef = ref<HTMLElement>()

const dropdownStyle = ref<Record<string, string>>({})
let leaveTimer: ReturnType<typeof setTimeout> | null = null

function open() {
  if (leaveTimer) { clearTimeout(leaveTimer); leaveTimer = null }
  if (triggerRef.value) {
    const rect = triggerRef.value.getBoundingClientRect()
    dropdownStyle.value = {
      top: (rect.bottom + 6) + 'px',
      right: (window.innerWidth - rect.right) + 'px',
    }
  }
  visible.value = true
}

function onLeave() {
  leaveTimer = setTimeout(() => { visible.value = false }, 200)
}

function cancelLeave() {
  if (leaveTimer) { clearTimeout(leaveTimer); leaveTimer = null }
}

/** 获取预览滚动容器 */
function getPreviewPane(): HTMLElement | null {
  return document.querySelector('.preview-pane.active')
}

function onJump(id: string) {
  const el = document.getElementById(id)
  if (!el) return

  // 先尝试在预览面板的滚动容器内滚动
  const pane = getPreviewPane()
  if (pane) {
    // 计算 heading 相对于 pane 的偏移量
    const paneRect = pane.getBoundingClientRect()
    const elRect = el.getBoundingClientRect()
    const offset = elRect.top - paneRect.top + pane.scrollTop - 16
    pane.scrollTo({ top: offset, behavior: 'smooth' })
  } else {
    // 降级到 scrollIntoView
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }

  activeId.value = id
  visible.value = false
}

function onPreviewScroll() {
  if (!outline.value.length) return
  const pane = getPreviewPane()
  if (!pane) return
  const top = pane.scrollTop + 40
  let current = outline.value[0]?.id || ''
  for (const item of outline.value) {
    const el = document.getElementById(item.id)
    if (el) {
      // 计算元素在 pane 内的 offsetTop
      let offsetTop = 0
      let node: HTMLElement | null = el
      while (node && node !== pane) {
        offsetTop += node.offsetTop
        node = node.offsetParent as HTMLElement | null
      }
      if (offsetTop <= top) current = item.id
    }
  }
  activeId.value = current
}

onMounted(() => {
  const pane = getPreviewPane()
  pane?.addEventListener('scroll', onPreviewScroll, { passive: true })
})

onUnmounted(() => {
  const pane = getPreviewPane()
  pane?.removeEventListener('scroll', onPreviewScroll)
  if (leaveTimer) clearTimeout(leaveTimer)
})
</script>

<style scoped>
.outline-trigger {
  position: relative;
  flex-shrink: 0;
}

.outline-fab {
  width: 30px;
  height: 30px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text-tertiary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
  padding: 0;
}

.outline-fab:hover {
  color: var(--accent);
  border-color: var(--accent);
  background: var(--accent-soft);
}

/* 弹窗 —— 通过 Teleport 挂载到 body，使用 fixed 定位 */
.outline-dropdown {
  position: fixed;
  width: 260px;
  max-height: 480px;
  background: var(--surface);
  border: 1px solid var(--border-strong);
  border-radius: var(--radius);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15), 0 2px 8px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  z-index: 9999;
}

.outline-dropdown-header {
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.outline-dropdown-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.outline-dropdown-count {
  font-size: 11px;
  color: var(--text-tertiary);
  font-family: var(--font-mono);
}

.outline-dropdown-body {
  flex: 1;
  overflow-y: auto;
  padding: 6px 0;
}

.empty-hint {
  font-size: 12px;
  color: var(--text-tertiary);
  text-align: center;
  padding: 24px 16px;
  margin: 0;
}

.outline-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.outline-item {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 28px;
  padding-right: 12px;
  cursor: pointer;
  color: var(--text-secondary);
  font-size: 12.5px;
  line-height: 1;
  border-left: 2px solid transparent;
  transition: background 0.1s, color 0.1s, border-color 0.1s;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.outline-item:hover {
  background: var(--surface-hover);
  color: var(--text);
}

.outline-item.active {
  background: var(--accent-soft);
  color: var(--accent);
  border-left-color: var(--accent);
  font-weight: 500;
}

.outline-item .dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  flex-shrink: 0;
  background: var(--border-strong);
}

.outline-item .dot-1 { background: #2563eb; }
.outline-item .dot-2 { background: #0891b2; }
.outline-item .dot-3 { background: #7c3aed; }
.outline-item .dot-4 { background: #db2777; }

.outline-item.active .dot { background: currentColor; }

.outline-item .text {
  overflow: hidden;
  text-overflow: ellipsis;
}

@media (max-width: 900px) {
  .outline-trigger { display: none; }
}
</style>

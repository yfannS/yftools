<template>
  <div class="history-view">
    <div class="history-header">
      <div class="header-left">
        <h2>历史记录</h2>
        <span class="header-count" v-if="authStore.isLoggedIn && historyStore.items.length > 0">
          {{ historyStore.items.length }} 条记录
        </span>
      </div>
      <router-link to="/md2html" class="back-link">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"/><polyline points="12 19 5 12 12 5"/></svg>
        <span>返回编辑器</span>
      </router-link>
    </div>

    <!-- 未登录 -->
    <Transition name="fade-scale" appear>
      <div v-if="!authStore.isLoggedIn" class="empty-hint">
        <div class="empty-icon pulse">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
        </div>
        <p class="empty-title">需要登录</p>
        <p class="empty-desc">请先<router-link to="/login">登录</router-link>后查看历史记录</p>
      </div>
    </Transition>

    <!-- 加载中 -->
    <Transition name="fade-scale" appear>
      <div v-if="authStore.isLoggedIn && historyStore.loading" class="empty-hint">
        <div class="spinner-ring">
          <div class="ring"></div>
          <div class="ring"></div>
          <div class="ring"></div>
        </div>
        <p class="empty-title">加载中</p>
        <p class="empty-desc shimmer">正在获取您的历史记录...</p>
      </div>
    </Transition>

    <!-- 空列表 -->
    <Transition name="fade-scale" appear>
      <div v-if="authStore.isLoggedIn && !historyStore.loading && !historyStore.items.length" class="empty-hint">
        <div class="empty-icon breathe">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
        </div>
        <p class="empty-title">暂无历史记录</p>
        <p class="empty-desc">编辑 Markdown 文档后，内容会自动保存到历史记录</p>
      </div>
    </Transition>

    <!-- 列表 -->
    <div v-if="authStore.isLoggedIn && !historyStore.loading && historyStore.items.length" class="history-list-wrapper">
      <TransitionGroup name="list-stagger" tag="div" class="history-list" appear>
        <div
          v-for="(item, index) in historyStore.items"
          :key="item.id"
          class="history-item"
          :style="{ '--stagger-index': index }"
          @click="onItemClick(item)"
        >
          <div class="item-accent-bar"></div>
          <div class="item-main">
            <div class="item-row">
              <div class="item-title">{{ item.title || '未命名文档' }}</div>
              <div class="item-meta">
                <span class="meta-chars">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
                  {{ item.char_count }} 字
                </span>
                <span class="meta-theme" v-if="item.theme && item.theme !== 'default'">{{ item.theme }}</span>
                <span class="meta-dot"></span>
                <span class="meta-time">{{ formatTime(item.updated_at || item.created_at) }}</span>
              </div>
            </div>
          </div>
          <button class="delete-btn" @click.stop="onDelete(item.id)" title="删除">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
          </button>
        </div>
      </TransitionGroup>

      <Transition name="fade-up" appear>
        <div class="pagination" v-if="historyStore.total > historyStore.pageSize">
          <button class="page-btn" :disabled="historyStore.page <= 1" @click="historyStore.fetchHistory(historyStore.page - 1)">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"/></svg>
            上一页
          </button>
          <span class="page-info">{{ historyStore.page }} / {{ totalPages }}</span>
          <button class="page-btn" :disabled="historyStore.page >= totalPages" @click="historyStore.fetchHistory(historyStore.page + 1)">
            下一页
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"/></svg>
          </button>
        </div>
      </Transition>
    </div>

    <!-- 详情加载遮罩 -->
    <Transition name="overlay">
      <div class="detail-loading-overlay" v-if="historyStore.detailLoading">
        <div class="overlay-content">
          <div class="spinner-ring">
            <div class="ring"></div>
            <div class="ring"></div>
            <div class="ring"></div>
          </div>
          <p>正在加载文档...</p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useHistoryStore, type HistoryListItem } from '@/stores/history'
import { useEditorStore } from '@/stores/editor'

const authStore = useAuthStore()
const historyStore = useHistoryStore()
const editorStore = useEditorStore()
const router = useRouter()

const totalPages = computed(() => Math.ceil(historyStore.total / historyStore.pageSize))

function formatTime(t: string) {
  if (!t) return ''
  const d = new Date(t)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const mins = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins} 分钟前`
  if (hours < 24) return `${hours} 小时前`
  if (days < 7) return `${days} 天前`
  return d.toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

async function onItemClick(item: HistoryListItem) {
  const detail = await historyStore.getDetail(item.id)
  if (detail) {
    editorStore.setContent(detail.markdown)
    router.push('/md2html')
  }
}

async function onDelete(id: number) {
  if (!confirm('确定删除这条记录？')) return
  await historyStore.deleteItem(id)
}

onMounted(() => {
  if (authStore.isLoggedIn) {
    historyStore.fetchHistory(1)
  }
})
</script>

<style scoped>
/* ===== 页面容器 ===== */
.history-view {
  padding: 32px 24px;
  max-width: 960px;
  margin: 0 auto;
  position: relative;
  animation: pageEnter 0.5s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes pageEnter {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* ===== 头部 ===== */
.history-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 28px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border);
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.history-header h2 {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--text);
}

.header-count {
  font-size: 13px;
  color: var(--text-tertiary);
  font-weight: 400;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-tertiary);
  text-decoration: none;
  padding: 6px 12px;
  border-radius: var(--radius-sm);
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.back-link:hover {
  color: var(--accent);
  background: var(--accent-soft);
}

.back-link svg {
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.back-link:hover svg {
  transform: translateX(-2px);
}

/* ===== 空状态 ===== */
.empty-hint {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100px 0;
  text-align: center;
}

.empty-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: var(--surface-raised);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
  color: var(--text-tertiary);
}

.empty-icon svg {
  width: 28px;
  height: 28px;
}

.empty-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text);
  margin-bottom: 6px;
}

.empty-desc {
  font-size: 13px;
  color: var(--text-tertiary);
  line-height: 1.6;
}

.empty-desc a {
  color: var(--accent);
  font-weight: 500;
  transition: color 0.15s;
}

.empty-desc a:hover {
  color: var(--accent-hover);
}

/* 脉冲动画 */
.pulse {
  animation: pulseAnim 2s ease-in-out infinite;
}

@keyframes pulseAnim {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.7; transform: scale(0.95); }
}

/* 呼吸动画 */
.breathe {
  animation: breatheAnim 3s ease-in-out infinite;
}

@keyframes breatheAnim {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}

/* 闪烁文字 */
.shimmer {
  background: linear-gradient(90deg, var(--text-tertiary) 25%, var(--text-secondary) 50%, var(--text-tertiary) 75%);
  background-size: 200% 100%;
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: shimmerMove 2s linear infinite;
}

@keyframes shimmerMove {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* ===== 加载动画（三环） ===== */
.spinner-ring {
  position: relative;
  width: 40px;
  height: 40px;
  margin-bottom: 16px;
}

.ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 2px solid transparent;
  border-top-color: var(--accent);
  animation: ringSpin 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
}

.ring:nth-child(1) { animation-delay: -0.45s; border-top-color: var(--accent); }
.ring:nth-child(2) { animation-delay: -0.3s; border-top-color: var(--accent-hover); opacity: 0.6; }
.ring:nth-child(3) { animation-delay: -0.15s; border-top-color: var(--accent); opacity: 0.3; }

@keyframes ringSpin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* ===== 列表容器 ===== */
.history-list-wrapper {
  animation: pageEnter 0.5s cubic-bezier(0.16, 1, 0.3, 1) both;
}

/* ===== 列表 ===== */
.history-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* ===== 历史记录卡片 ===== */
.history-item {
  display: flex;
  align-items: stretch;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--surface);
  overflow: hidden;
  cursor: pointer;
  position: relative;
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  animation: cardEnter 0.45s cubic-bezier(0.16, 1, 0.3, 1) both;
  animation-delay: calc(var(--stagger-index, 0) * 50ms);
}

@keyframes cardEnter {
  from {
    opacity: 0;
    transform: translateY(12px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.history-item::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: var(--radius);
  opacity: 0;
  transition: opacity 0.25s;
  pointer-events: none;
  box-shadow: 0 4px 20px rgba(100, 116, 139, 0.06), 0 1px 3px rgba(0, 0, 0, 0.04);
}

.history-item:hover {
  border-color: var(--accent);
  transform: translateY(-2px);
}

.history-item:hover::before {
  opacity: 1;
}

.history-item:active {
  transform: translateY(0) scale(0.995);
  transition-duration: 0.1s;
}

/* 左侧装饰条 */
.item-accent-bar {
  width: 3px;
  background: var(--accent);
  opacity: 0;
  transition: opacity 0.25s;
  flex-shrink: 0;
}

.history-item:hover .item-accent-bar {
  opacity: 1;
}

.item-main {
  flex: 1;
  padding: 16px 20px;
  min-width: 0;
  display: flex;
  align-items: center;
}

.item-row {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 16px;
  min-width: 0;
}

.item-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
  transition: color 0.2s;
  flex-shrink: 1;
  min-width: 0;
}

.history-item:hover .item-title {
  color: var(--accent);
}

.item-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  color: var(--text-tertiary);
  font-variant-numeric: tabular-nums;
  flex-shrink: 0;
  margin-left: auto;
}

.meta-chars {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-family: var(--font-mono);
  font-size: 11px;
}

.meta-chars svg {
  opacity: 0.5;
}

.meta-theme {
  padding: 1px 8px;
  background: var(--surface-raised);
  border-radius: 4px;
  font-size: 11px;
  line-height: 20px;
  font-weight: 500;
  color: var(--text-secondary);
}

.meta-dot {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--border-strong);
}

.meta-time {
  font-size: 11px;
}

/* ===== 删除按钮 ===== */
.delete-btn {
  width: 48px;
  border: none;
  border-left: 1px solid var(--border);
  background: transparent;
  color: var(--text-tertiary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  opacity: 0;
  transform: translateX(8px);
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.history-item:hover .delete-btn {
  opacity: 1;
  transform: translateX(0);
}

.delete-btn:hover {
  background: rgba(180, 86, 86, 0.08);
  color: var(--danger);
  border-left-color: rgba(180, 86, 86, 0.15);
}

.delete-btn:active {
  transform: scale(0.92);
}

:root[data-theme="dark"] .delete-btn:hover {
  background: rgba(180, 86, 86, 0.12);
  color: #f87171;
}

/* ===== 分页 ===== */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 28px;
  padding-top: 20px;
  border-top: 1px solid var(--border);
}

.page-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 34px;
  padding: 0 16px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.page-btn:hover:not(:disabled) {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--accent-soft);
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(100, 116, 139, 0.1);
}

.page-btn:active:not(:disabled) {
  transform: translateY(0) scale(0.97);
}

.page-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.page-info {
  font-size: 13px;
  color: var(--text-tertiary);
  font-variant-numeric: tabular-nums;
  font-weight: 500;
  min-width: 60px;
  text-align: center;
}

/* ===== 详情加载遮罩 ===== */
.detail-loading-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.25);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 100;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

.overlay-content {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 40px 48px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  color: var(--text);
  font-size: 14px;
  font-weight: 500;
}

.overlay-content .spinner-ring {
  margin-bottom: 0;
}

/* ===== Vue Transition 动画 ===== */

/* fade-scale: 空状态等 */
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(10px);
}

/* fade-up: 分页等 */
.fade-up-enter-active,
.fade-up-leave-active {
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.fade-up-enter-from,
.fade-up-leave-to {
  opacity: 0;
  transform: translateY(12px);
}

/* list-stagger: 列表 stagger 动画 */
.list-stagger-enter-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.list-stagger-leave-active {
  transition: all 0.25s cubic-bezier(0.7, 0, 0.84, 0);
}

.list-stagger-enter-from,
.list-stagger-leave-to {
  opacity: 0;
  transform: translateX(-16px) scale(0.98);
}

.list-stagger-move {
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

/* overlay 遮罩过渡 */
.overlay-enter-active,
.overlay-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}

.overlay-enter-active .overlay-content,
.overlay-leave-active .overlay-content {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.overlay-enter-from .overlay-content,
.overlay-leave-to .overlay-content {
  opacity: 0;
  transform: scale(0.92) translateY(8px);
}
</style>

<template>
  <div class="history-view">
    <div class="history-header">
      <h2>历史记录</h2>
      <router-link to="/md2html" class="back-link">← 返回编辑器</router-link>
    </div>

    <!-- 未登录 -->
    <div v-if="!authStore.isLoggedIn" class="empty-hint">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
      <p>请先<router-link to="/login">登录</router-link>后查看历史记录</p>
    </div>

    <!-- 加载中 -->
    <div v-else-if="historyStore.loading" class="empty-hint">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 空列表 -->
    <div v-else-if="!historyStore.items.length" class="empty-hint">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
      <p>暂无历史记录</p>
      <p class="sub">编辑 Markdown 文档后，内容会自动保存到历史记录</p>
    </div>

    <!-- 列表 -->
    <template v-else>
      <div class="history-list">
        <div
          v-for="item in historyStore.items"
          :key="item.id"
          class="history-item"
          @click="onItemClick(item)"
        >
          <div class="item-main">
            <div class="item-title">{{ item.title || '未命名文档' }}</div>
            <div class="item-meta">
              <span class="meta-chars">{{ item.char_count }} 字</span>
              <span class="meta-theme" v-if="item.theme && item.theme !== 'default'">{{ item.theme }}</span>
              <span class="meta-time">{{ formatTime(item.updated_at || item.created_at) }}</span>
            </div>
          </div>
          <button class="delete-btn" @click.stop="onDelete(item.id)" title="删除">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
          </button>
        </div>
      </div>

      <div class="pagination" v-if="historyStore.total > historyStore.pageSize">
        <button class="page-btn" :disabled="historyStore.page <= 1" @click="historyStore.fetchHistory(historyStore.page - 1)">上一页</button>
        <span class="page-info">{{ historyStore.page }} / {{ totalPages }}</span>
        <button class="page-btn" :disabled="historyStore.page >= totalPages" @click="historyStore.fetchHistory(historyStore.page + 1)">下一页</button>
      </div>
    </template>

    <!-- 详情加载遮罩 -->
    <div class="detail-loading-overlay" v-if="historyStore.detailLoading">
      <div class="spinner"></div>
      <p>加载文档内容...</p>
    </div>
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
  // 点击后请求详情接口获取完整 markdown
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
.history-view {
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
  position: relative;
}

.history-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.history-header h2 {
  font-size: 20px;
  font-weight: 600;
  letter-spacing: -0.01em;
}

.back-link {
  font-size: 13px;
  color: var(--accent);
  text-decoration: none;
  transition: opacity 0.12s;
}

.back-link:hover { opacity: 0.8; }

.empty-hint {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;
  color: var(--text-tertiary);
  font-size: 14px;
  text-align: center;
}

.empty-hint svg {
  width: 40px;
  height: 40px;
  stroke: var(--border-strong);
  margin-bottom: 12px;
}

.empty-hint a { color: var(--accent); }
.empty-hint .sub { font-size: 12px; color: var(--text-tertiary); margin-top: 4px; }

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  margin-bottom: 8px;
}

@keyframes spin { to { transform: rotate(360deg); } }

.history-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.history-item {
  display: flex;
  align-items: center;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--surface);
  overflow: hidden;
  transition: border-color 0.12s, box-shadow 0.12s;
  cursor: pointer;
}

.history-item:hover {
  border-color: var(--accent);
  box-shadow: 0 1px 4px rgba(37, 99, 235, 0.06);
}

.item-main {
  flex: 1;
  padding: 14px 16px;
  min-width: 0;
}

.item-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 6px;
  line-height: 1.4;
}

.item-meta {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: var(--text-tertiary);
  font-variant-numeric: tabular-nums;
}

.meta-chars {
  font-family: var(--font-mono);
}

.meta-theme {
  padding: 0 6px;
  background: var(--surface-raised);
  border-radius: 3px;
  font-size: 10px;
  line-height: 18px;
}

.delete-btn {
  width: 40px;
  height: 100%;
  min-height: 48px;
  border: none;
  border-left: 1px solid var(--border);
  background: transparent;
  color: var(--text-tertiary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.12s;
}

.delete-btn:hover {
  background: #fee2e2;
  color: #dc2626;
  border-left-color: #fecaca;
}

:root[data-theme="dark"] .delete-btn:hover {
  background: rgba(220, 38, 38, 0.15);
  color: #f87171;
  border-left-color: rgba(220, 38, 38, 0.3);
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-top: 20px;
}

.page-btn {
  height: 30px;
  padding: 0 14px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.12s;
}

.page-btn:hover:not(:disabled) {
  border-color: var(--accent);
  color: var(--accent);
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-info {
  font-size: 12px;
  color: var(--text-tertiary);
  font-variant-numeric: tabular-nums;
}

.detail-loading-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 100;
  color: white;
  font-size: 14px;
  backdrop-filter: blur(4px);
}
</style>

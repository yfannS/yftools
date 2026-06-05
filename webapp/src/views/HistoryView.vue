<template>
  <div class="history-view">
    <div class="history-header">
      <h2>历史记录</h2>
      <router-link to="/md2html" class="back-link">← 返回编辑器</router-link>
    </div>

    <div v-if="!authStore.isLoggedIn" class="empty-hint">
      <p>请先<router-link to="/login">登录</router-link>后查看历史记录</p>
    </div>

    <div v-else-if="historyStore.loading" class="empty-hint">
      <p>加载中...</p>
    </div>

    <div v-else-if="!historyStore.items.length" class="empty-hint">
      <p>暂无历史记录</p>
    </div>

    <template v-else>
      <div class="history-list">
        <div v-for="item in historyStore.items" :key="item.id" class="history-item">
          <div class="item-main" @click="loadItem(item)">
            <div class="item-preview">{{ item.markdown.slice(0, 120) }}{{ item.markdown.length > 120 ? '...' : '' }}</div>
            <div class="item-meta">
              <span class="item-time">{{ formatTime(item.created_at) }}</span>
              <span class="item-theme">{{ item.theme || 'default' }}</span>
            </div>
          </div>
          <button class="delete-btn" @click.stop="onDelete(item.id)" title="删除">✕</button>
        </div>
      </div>

      <div class="pagination" v-if="historyStore.total > historyStore.pageSize">
        <button class="page-btn" :disabled="historyStore.page <= 1" @click="historyStore.fetchHistory(historyStore.page - 1)">上一页</button>
        <span class="page-info">{{ historyStore.page }} / {{ totalPages }}</span>
        <button class="page-btn" :disabled="historyStore.page >= totalPages" @click="historyStore.fetchHistory(historyStore.page + 1)">下一页</button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useHistoryStore, type HistoryItem } from '@/stores/history'
import { useEditorStore } from '@/stores/editor'

const authStore = useAuthStore()
const historyStore = useHistoryStore()
const editorStore = useEditorStore()
const router = useRouter()

const totalPages = computed(() => Math.ceil(historyStore.total / historyStore.pageSize))

function formatTime(t: string) {
  if (!t) return ''
  const d = new Date(t)
  return d.toLocaleString('zh-CN')
}

function loadItem(item: HistoryItem) {
  editorStore.setContent(item.markdown)
  router.push('/md2html')
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
}

.back-link {
  font-size: 13px;
  color: var(--accent);
}

.empty-hint {
  text-align: center;
  padding: 60px 0;
  color: var(--text-tertiary);
  font-size: 14px;
}

.empty-hint a {
  color: var(--accent);
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.history-item {
  display: flex;
  align-items: stretch;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--surface);
  overflow: hidden;
  transition: border-color 0.12s;
}

.history-item:hover {
  border-color: var(--accent);
}

.item-main {
  flex: 1;
  padding: 12px 16px;
  cursor: pointer;
  min-width: 0;
}

.item-preview {
  font-size: 13px;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 6px;
  font-family: var(--font-mono);
  line-height: 1.5;
}

.item-meta {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: var(--text-tertiary);
}

.delete-btn {
  width: 40px;
  border: none;
  background: transparent;
  color: var(--text-tertiary);
  cursor: pointer;
  font-size: 14px;
  transition: all 0.12s;
  flex-shrink: 0;
}

.delete-btn:hover {
  background: var(--danger);
  color: white;
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
}
</style>

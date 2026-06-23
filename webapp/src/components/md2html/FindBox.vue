<template>
  <div class="find-box" v-if="editorStore.findBoxOpen">
    <div class="find-row">
      <input
        class="find-input"
        v-model="editorStore.findQuery"
        placeholder="查找内容"
        ref="findInputRef"
        @keydown.enter="onFindNext"
      />
      <button class="toolbar-btn" @click="onFindNext">下一个</button>
    </div>
    <div class="find-row">
      <input
        class="find-input"
        v-model="editorStore.replaceQuery"
        placeholder="替换为"
      />
      <button class="toolbar-btn" @click="onReplaceOne">替换</button>
      <button class="toolbar-btn" @click="onReplaceAll">全部</button>
    </div>
    <div class="find-row" style="justify-content:flex-end">
      <button class="toolbar-btn" @click="editorStore.findBoxOpen = false">关闭</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, inject } from 'vue'
import { useEditorStore } from '@/stores/editor'
import type { FindReplaceAPI } from '@/composables/useFindReplace'

const editorStore = useEditorStore()
const findInputRef = ref<HTMLInputElement>()

// 通过 inject 获取父组件提供的 findReplace API
const findReplace = inject<FindReplaceAPI>('findReplace')

onMounted(() => {
  nextTick(() => findInputRef.value?.focus())
})

function onFindNext() {
  findReplace?.findNext()
}

function onReplaceOne() {
  findReplace?.replaceOne()
}

function onReplaceAll() {
  findReplace?.replaceAll()
}
</script>

<style scoped>
.find-box {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 30;
  background: var(--surface);
  border: 1px solid var(--border-strong);
  border-radius: var(--radius);
  box-shadow: 0 8px 30px rgba(0,0,0,0.12);
  padding: 10px;
  width: 280px;
  display: flex;
  gap: 8px;
  flex-direction: column;
}

.find-row { display: flex; gap: 6px; }
.find-row input {
  flex: 1;
  min-width: 0;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 5px 8px;
  background: var(--editor-bg);
  color: var(--text);
  outline: none;
  font-size: 12px;
  font-family: var(--font-sans);
}
.find-row input:focus { border-color: var(--accent); }
.find-row .toolbar-btn { height: 26px; padding: 0 8px; font-size: 11px; }
</style>

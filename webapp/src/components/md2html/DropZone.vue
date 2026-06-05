<template>
  <div class="drop-zone" :class="{ active: isDragging }">
    释放文件以导入 Markdown
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { SUPPORTED_FILE_TYPES } from '@/utils/constants'

const emit = defineEmits<{ 'file-dropped': [file: File] }>()
const isDragging = ref(false)
let dragCounter = 0

function handleFile(file: File) {
  const ext = '.' + file.name.split('.').pop()?.toLowerCase()
  if (SUPPORTED_FILE_TYPES.includes(ext)) {
    emit('file-dropped', file)
  }
}

function onDragEnter(e: DragEvent) {
  e.preventDefault()
  dragCounter++
  isDragging.value = true
}

function onDragOver(e: DragEvent) {
  e.preventDefault()
}

function onDragLeave(e: DragEvent) {
  e.preventDefault()
  dragCounter--
  if (dragCounter <= 0) {
    dragCounter = 0
    isDragging.value = false
  }
}

function onDrop(e: DragEvent) {
  e.preventDefault()
  dragCounter = 0
  isDragging.value = false
  const file = e.dataTransfer?.files?.[0]
  if (file) handleFile(file)
}

onMounted(() => {
  document.addEventListener('dragenter', onDragEnter)
  document.addEventListener('dragover', onDragOver)
  document.addEventListener('dragleave', onDragLeave)
  document.addEventListener('drop', onDrop)
})

onUnmounted(() => {
  document.removeEventListener('dragenter', onDragEnter)
  document.removeEventListener('dragover', onDragOver)
  document.removeEventListener('dragleave', onDragLeave)
  document.removeEventListener('drop', onDrop)
})
</script>

<style scoped>
.drop-zone {
  position: fixed;
  inset: 6px;
  border: 2px dashed var(--accent);
  border-radius: var(--radius);
  background: var(--accent-soft);
  display: none;
  align-items: center;
  justify-content: center;
  color: var(--accent);
  font-size: 15px;
  font-weight: 500;
  z-index: 150;
  pointer-events: none;
  letter-spacing: -0.01em;
}

.drop-zone.active { display: flex; }
</style>

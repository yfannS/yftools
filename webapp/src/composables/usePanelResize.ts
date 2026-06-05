import { ref, onUnmounted } from 'vue'

export function usePanelResize() {
  const splitRatio = ref(parseFloat(localStorage.getItem('md2html-split-ratio') || '0.5'))
  const isResizing = ref(false)

  // Clamp ratio
  splitRatio.value = Math.min(0.8, Math.max(0.2, splitRatio.value))

  function applyRatio(ratio: number, mainEl: HTMLElement) {
    mainEl.style.gridTemplateColumns = `${ratio * 100}% 4px ${(1 - ratio) * 100}%`
  }

  function onResizeStart(e: MouseEvent, mainEl: HTMLElement) {
    e.preventDefault()
    isResizing.value = true
    document.body.classList.add('resize-active')

    const startX = e.clientX
    const startRatio = splitRatio.value
    const totalWidth = mainEl.offsetWidth

    function onMouseMove(e: MouseEvent) {
      if (!isResizing.value) return
      const offsetX = e.clientX - mainEl.getBoundingClientRect().left
      let ratio = offsetX / totalWidth
      ratio = Math.min(0.8, Math.max(0.2, ratio))
      splitRatio.value = ratio
      applyRatio(ratio, mainEl)
    }

    function onMouseUp() {
      isResizing.value = false
      document.body.classList.remove('resize-active')
      localStorage.setItem('md2html-split-ratio', splitRatio.value.toFixed(4))
      document.removeEventListener('mousemove', onMouseMove)
      document.removeEventListener('mouseup', onMouseUp)
    }

    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
  }

  function initRatio(mainEl: HTMLElement) {
    applyRatio(splitRatio.value, mainEl)
  }

  return {
    splitRatio,
    isResizing,
    onResizeStart,
    initRatio
  }
}

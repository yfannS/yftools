export function useToast() {
  function toast(msg: string, type?: string) {
    const el = document.createElement('div')
    el.className = 'toast'
    el.textContent = msg
    if (type === 'err') el.style.background = '#7f1d1d'
    const container = document.getElementById('toast-container')
    if (container) {
      container.appendChild(el)
      setTimeout(() => el.remove(), 2100)
    }
  }

  return { toast }
}

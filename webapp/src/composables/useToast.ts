const TOAST_CONTAINER_ID = 'toast-container'
const TOAST_STYLE_ID = 'toast-global-style'

export type ToastType = 'err' | 'ok' | 'info'

function ensureToastStyle() {
  if (document.getElementById(TOAST_STYLE_ID)) return

  const style = document.createElement('style')
  style.id = TOAST_STYLE_ID
  style.textContent = `
    .toast-container {
      position: fixed;
      top: 24px;
      right: 24px;
      display: flex;
      flex-direction: column;
      gap: 10px;
      z-index: 3000;
      pointer-events: none;
      max-width: min(360px, calc(100vw - 32px));
    }

    .toast {
      min-height: 44px;
      padding: 11px 14px;
      border-radius: 10px;
      color: #f8fafc;
      background: rgba(15, 23, 42, 0.94);
      box-shadow: 0 12px 24px rgba(15, 23, 42, 0.18);
      border: 1px solid rgba(148, 163, 184, 0.18);
      font-size: 13px;
      line-height: 1.45;
      display: flex;
      align-items: center;
      opacity: 0;
      transform: translateY(-6px);
      animation: toastIn 0.25s ease forwards, toastOut 0.2s 1.9s ease forwards;
      pointer-events: auto;
      word-break: break-word;
    }

    .toast.toast-err {
      background: rgba(127, 29, 29, 0.96);
      border-color: rgba(248, 113, 113, 0.28);
    }

    .toast.toast-ok {
      background: rgba(22, 101, 52, 0.96);
      border-color: rgba(74, 222, 128, 0.22);
    }

    @keyframes toastIn { to { opacity: 1; transform: translateY(0); } }
    @keyframes toastOut { to { opacity: 0; transform: translateY(-4px); } }

    @media (max-width: 640px) {
      .toast-container {
        top: 14px;
        left: 14px;
        right: 14px;
        max-width: none;
      }
    }
  `

  document.head.appendChild(style)
}

function ensureToastContainer() {
  let container = document.getElementById(TOAST_CONTAINER_ID)
  if (container) return container

  ensureToastStyle()

  container = document.createElement('div')
  container.id = TOAST_CONTAINER_ID
  container.className = 'toast-container'
  document.body.appendChild(container)
  return container
}

export function showToast(msg: string, type: ToastType = 'info') {
  if (typeof document === 'undefined') return

  const el = document.createElement('div')
  el.className = `toast toast-${type}`
  el.textContent = msg

  const container = ensureToastContainer()
  container.appendChild(el)
  setTimeout(() => el.remove(), 2200)
}

export function useToast() {
  function toast(msg: string, type?: string) {
    showToast(msg, type === 'err' ? 'err' : type === 'ok' ? 'ok' : 'info')
  }

  return { toast }
}

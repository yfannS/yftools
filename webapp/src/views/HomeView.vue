<template>
  <div class="home-view">
    <div class="hero">
      <h1>工具箱</h1>
      <p>在线工具集合，提升你的工作效率</p>
    </div>
    <div class="tools-grid">
      <router-link to="/md2html" class="tool-card">
        <div class="tool-icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
        </div>
        <div class="tool-info">
          <h3>Markdown → HTML</h3>
          <p>实时编辑预览，代码高亮，流程图渲染，导出独立 HTML</p>
        </div>
      </router-link>
      <router-link to="/json-formatter" class="tool-card">
        <div class="tool-icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><polyline points="4 7 4 4 20 4 20 7"/><line x1="9" y1="20" x2="15" y2="20"/><line x1="12" y1="4" x2="12" y2="20"/></svg>
        </div>
        <div class="tool-info">
          <h3>JSON 格式化</h3>
          <p>格式化、压缩、校验 JSON，支持自定义缩进</p>
        </div>
      </router-link>
    </div>
  </div>
</template>

<style scoped>
.home-view {
  max-width: 800px;
  margin: 0 auto;
  padding: 60px 24px;
  animation: homeFadeIn 0.5s var(--ease-out);
}

@keyframes homeFadeIn {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}

.hero {
  text-align: center;
  margin-bottom: 48px;
}

.hero h1 {
  font-size: 32px;
  font-weight: 800;
  letter-spacing: -0.03em;
  color: var(--text);
  margin-bottom: 8px;
}

.hero p {
  font-size: 16px;
  color: var(--text-secondary);
}

.tools-grid {
  display: grid;
  gap: 16px;
}

.tool-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  text-decoration: none;
  transition: all 0.2s var(--ease-out);
  box-shadow: var(--shadow-sm);
}

.tool-card:hover {
  border-color: var(--accent);
  box-shadow: var(--shadow-md), 0 0 0 3px var(--accent-soft);
  transform: translateY(-2px);
}

.tool-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--accent-soft);
  border-radius: var(--radius);
  color: var(--accent);
  flex-shrink: 0;
  transition: transform 0.2s var(--ease-spring);
}

.tool-card:hover .tool-icon {
  transform: scale(1.08);
}

.tool-info h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text);
  margin-bottom: 4px;
}

.tool-info p {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
}
</style>

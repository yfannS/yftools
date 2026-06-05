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
      <!-- 未来工具卡片在此添加 -->
    </div>
  </div>
</template>

<style scoped>
.home-view {
  max-width: 800px;
  margin: 0 auto;
  padding: 60px 24px;
}

.hero {
  text-align: center;
  margin-bottom: 48px;
}

.hero h1 {
  font-size: 32px;
  font-weight: 700;
  letter-spacing: -0.02em;
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
  transition: all 0.15s ease;
}

.tool-card:hover {
  border-color: var(--accent);
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.08);
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

<template>
  <div class="not-found">
    <div class="code">404</div>
    <p class="desc">页面不存在</p>
    <router-link to="/" class="back">返回首页</router-link>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  gap: 12px;
}

.code {
  font-size: 64px;
  font-weight: 500;
  color: var(--text-tertiary);
  letter-spacing: -0.04em;
  line-height: 1;
}

.desc {
  font-size: 15px;
  color: var(--text-secondary);
  margin: 0;
}

.back {
  margin-top: 8px;
  font-size: 13px;
  color: var(--accent);
  text-decoration: none;
  padding: 8px 16px;
  border: 0.5px solid var(--border);
  border-radius: var(--radius-sm);
  transition: background 0.2s;
}

.back:hover {
  background: var(--surface-raised);
}
</style>

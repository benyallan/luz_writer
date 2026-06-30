<script setup>
import { useWorkspaceStore } from '@/stores/workspace'
import FileTree from '@/components/sidebar/FileTree.vue'

const workspace = useWorkspaceStore()
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <span class="sidebar-title">
        {{ workspace.rootName ? workspace.rootName.toUpperCase() : 'EXPLORADOR' }}
      </span>
    </div>

    <div class="sidebar-body">
      <FileTree v-if="workspace.rootPath" />
      <div v-else class="sidebar-empty">
        <p class="empty-hint">Nenhuma pasta aberta</p>
        <button class="open-btn" @click="workspace.openFolder">
          Abrir pasta
        </button>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  background: var(--color-sidebar);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  height: 100%;
}

.sidebar-header {
  padding: 9px 12px 6px;
  flex-shrink: 0;
}

.sidebar-title {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: var(--color-text);
  user-select: none;
}

.sidebar-body {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

.sidebar-empty {
  padding: 24px 16px;
  text-align: center;
}

.empty-hint {
  color: var(--color-text-muted);
  font-size: 12px;
  margin-bottom: 12px;
}

.open-btn {
  padding: 5px 14px;
  background: var(--color-accent);
  color: #fff;
  border-radius: 2px;
  font-size: 12px;
  transition: background 0.1s;
}

.open-btn:hover {
  background: #1177bb;
}
</style>

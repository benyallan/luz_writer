<script setup>
import MenuBar from '@/components/layout/MenuBar.vue'
import Sidebar from '@/components/layout/Sidebar.vue'
import NewProjectDialog from '@/components/dialogs/NewProjectDialog.vue'
</script>

<template>
  <div class="app-shell">
    <NewProjectDialog />
    <MenuBar class="app-menubar" />
    <Sidebar class="app-sidebar" />
    <main class="app-editor">
      <div class="editor-welcome">
        <p>Abra uma pasta para começar</p>
      </div>
    </main>
  </div>
</template>

<style scoped>
.app-shell {
  display: grid;
  grid-template-rows: var(--menubar-height) 1fr;
  grid-template-columns: var(--sidebar-width) 1fr;
  height: 100vh;
  overflow: hidden;
}

.app-menubar {
  grid-column: 1 / -1;
  grid-row: 1;
}

.app-sidebar {
  grid-column: 1;
  grid-row: 2;
  overflow: hidden;
}

.app-editor {
  grid-column: 2;
  grid-row: 2;
  background: var(--color-bg);
  overflow: auto;
  display: flex;
  align-items: center;
  justify-content: center;
}

.editor-welcome {
  color: var(--color-text-muted);
  text-align: center;
  pointer-events: none;
  user-select: none;
}
</style>

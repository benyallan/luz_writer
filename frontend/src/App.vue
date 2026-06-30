<script setup>
import MenuBar from '@/components/layout/MenuBar.vue'
import Sidebar from '@/components/layout/Sidebar.vue'
import NewProjectDialog from '@/components/dialogs/NewProjectDialog.vue'
import FileEditor from '@/components/editor/FileEditor.vue'
import { useEditorStore } from '@/stores/editor'

const editor = useEditorStore()
</script>

<template>
  <div class="app-shell">
    <NewProjectDialog />
    <MenuBar class="app-menubar" />
    <Sidebar class="app-sidebar" />
    <main class="app-editor">
      <FileEditor v-if="editor.filePath" />
      <div v-else class="editor-welcome">
        <p>Selecione um arquivo para editar</p>
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
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.editor-welcome {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-muted);
  user-select: none;
  pointer-events: none;
}
</style>

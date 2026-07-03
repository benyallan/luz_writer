<script setup lang="ts">
import {ref} from 'vue'
import ExplorerGroup from './ExplorerGroup.vue'
import NewDocumentDialog from './NewDocumentDialog.vue'
import SettingsPanel from './SettingsPanel.vue'
import BuildTargetsPanel from './BuildTargetsPanel.vue'
import ExtensionsPanel from './ExtensionsPanel.vue'
import {useWorkspaceStore} from '../stores/workspace'
import {useEditorStore} from '../stores/editor'
import {useUiStore} from '../stores/ui'

const workspace = useWorkspaceStore()
const editor = useEditorStore()
const ui = useUiStore()
const newDocOpen = ref(false)

type GroupKey = 'pre' | 'cap' | 'pos'

function handleReorder(groupKey: GroupKey, newIds: string[]) {
  const pre = groupKey === 'pre' ? newIds : workspace.preTextuais.map(c => c.id)
  const cap = groupKey === 'cap' ? newIds : workspace.capitulos.map(c => c.id)
  const pos = groupKey === 'pos' ? newIds : workspace.posTextuais.map(c => c.id)
  workspace.reorderChapters([...pre, ...cap, ...pos])
}

function handleRename(id: string, title: string) {
  workspace.renameChapter(id, title)
}

function handleDelete(id: string) {
  if (editor.activeChapterId === id) editor.closeChapter()
  workspace.deleteChapter(id)
}

function handleOpen(id: string) {
  editor.openChapter(id)
}
</script>

<template>
  <aside class="sidebar">
    <template v-if="ui.activePanel === 'explorer'">
      <div class="sidebar__header">
        <span class="sidebar__title">Explorer</span>
        <button class="sidebar__new-button" type="button" title="Novo documento" @click="newDocOpen = true">+</button>
      </div>

      <p v-if="!workspace.chapters.length" class="sidebar__empty">Nenhum documento ainda.</p>

      <ExplorerGroup
        label="Pré-textuais"
        :items="workspace.preTextuais"
        :active-id="editor.activeChapterId"
        @reorder="ids => handleReorder('pre', ids)"
        @rename="handleRename"
        @delete="handleDelete"
        @open="handleOpen"
      />
      <ExplorerGroup
        label="Capítulos"
        :items="workspace.capitulos"
        :active-id="editor.activeChapterId"
        @reorder="ids => handleReorder('cap', ids)"
        @rename="handleRename"
        @delete="handleDelete"
        @open="handleOpen"
      />
      <ExplorerGroup
        label="Pós-textuais"
        :items="workspace.posTextuais"
        :active-id="editor.activeChapterId"
        @reorder="ids => handleReorder('pos', ids)"
        @rename="handleRename"
        @delete="handleDelete"
        @open="handleOpen"
      />

      <NewDocumentDialog v-model:open="newDocOpen" />
    </template>

    <BuildTargetsPanel v-else-if="ui.activePanel === 'buildTargets'" />

    <ExtensionsPanel v-else-if="ui.activePanel === 'extensions'" />

    <SettingsPanel v-else-if="ui.activePanel === 'settings'" />
  </aside>
</template>

<style scoped>
.sidebar {
  width: 260px;
  flex: 0 0 auto;
  background: var(--luz-bg-sidebar);
  border-right: 1px solid var(--luz-border);
  overflow-y: auto;
}

.sidebar__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px 0;
}

.sidebar__title {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--luz-fg-muted);
}

.sidebar__new-button {
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg-muted);
  border-radius: 4px;
  width: 20px;
  height: 20px;
  line-height: 1;
  cursor: pointer;
}

.sidebar__new-button:hover {
  background: var(--luz-bg-hover);
  color: var(--luz-fg);
}

.sidebar__empty {
  padding: 12px;
  margin: 0;
  font-size: 0.8rem;
  color: var(--luz-fg-muted);
}
</style>

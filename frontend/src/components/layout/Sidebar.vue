<script setup>
import { ref, nextTick } from 'vue'
import { FilePlus, FolderPlus, Folder, FileText } from '@lucide/vue'
import { CreateFile, CreateDirectory } from '@wails/go/main/App'
import { useWorkspaceStore } from '@/stores/workspace'
import FileTree from '@/components/sidebar/FileTree.vue'

const workspace = useWorkspaceStore()

const rootCreatingType = ref(null)   // 'file' | 'folder' | null
const rootNewItemName = ref('')
const rootCreatingInput = ref(null)

async function startRootCreating(type) {
  rootCreatingType.value = type
  rootNewItemName.value = ''
  await nextTick()
  rootCreatingInput.value?.focus()
}

function cancelRootCreate() {
  rootCreatingType.value = null
  rootNewItemName.value = ''
}

async function confirmRootCreate() {
  const name = rootNewItemName.value.trim()
  if (!name) { cancelRootCreate(); return }

  const type = rootCreatingType.value
  rootCreatingType.value = null
  rootNewItemName.value = ''

  try {
    if (type === 'file') {
      await CreateFile(workspace.rootPath, name)
    } else {
      await CreateDirectory(workspace.rootPath, name)
    }
    await workspace.refreshRoot()
  } catch (e) {
    console.error(e)
  }
}
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <span class="sidebar-title">
        {{ workspace.rootName ? workspace.rootName.toUpperCase() : 'EXPLORADOR' }}
      </span>
      <div v-if="workspace.rootPath" class="header-actions">
        <button class="header-btn" title="Novo arquivo" @click="startRootCreating('file')">
          <FilePlus :size="15" />
        </button>
        <button class="header-btn" title="Nova pasta" @click="startRootCreating('folder')">
          <FolderPlus :size="15" />
        </button>
      </div>
    </div>

    <div class="sidebar-body">
      <template v-if="workspace.rootPath">
        <!-- Root-level inline creating input -->
        <div v-if="rootCreatingType" class="tree-row root-creating-row">
          <Folder v-if="rootCreatingType === 'folder'" :size="16" class="tree-icon icon-folder" />
          <FileText v-else :size="16" class="tree-icon icon-file" />
          <input
            ref="rootCreatingInput"
            v-model="rootNewItemName"
            class="creating-input"
            type="text"
            @keydown.enter.prevent="confirmRootCreate"
            @keydown.escape="cancelRootCreate"
            @blur="cancelRootCreate"
          />
        </div>

        <FileTree />
      </template>

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
  display: flex;
  align-items: center;
  padding: 9px 6px 6px 12px;
  flex-shrink: 0;
  min-height: 30px;
}

.sidebar-title {
  flex: 1;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: var(--color-text);
  user-select: none;
}

.header-actions {
  display: flex;
  gap: 2px;
}

.header-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 3px;
  color: var(--color-text-muted);
  transition: background 0.1s, color 0.1s;
}

.header-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--color-text);
}

.sidebar-body {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

/* Root-level creating row (no indent, follows depth 0 file style) */
.root-creating-row {
  display: flex;
  align-items: center;
  gap: 4px;
  height: 22px;
  padding-left: 22px;
  padding-right: 4px;
  cursor: default;
}

.tree-icon { flex-shrink: 0; }
.icon-folder { color: var(--color-folder); }
.icon-file   { color: var(--color-file); }

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

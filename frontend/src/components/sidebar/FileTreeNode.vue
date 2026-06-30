<script setup>
import { ref, nextTick } from 'vue'
import {
  ChevronRight, ChevronDown,
  Folder, FolderOpen, FileText,
  FilePlus, FolderPlus,
} from '@lucide/vue'
import { CreateFile, CreateDirectory } from '@wails/go/main/App'
import { useWorkspaceStore } from '@/stores/workspace'

defineOptions({ name: 'FileTreeNode' })

const props = defineProps({
  node: { type: Object, required: true },
  depth: { type: Number, default: 0 },
})

const workspace = useWorkspaceStore()
const isOpen = ref(false)
const children = ref([])
const isLoaded = ref(false)
const creatingType = ref(null)   // 'file' | 'folder' | null
const newItemName = ref('')
const creatingInput = ref(null)

async function toggle() {
  if (!isLoaded.value) {
    children.value = await workspace.loadChildren(props.node.path)
    isLoaded.value = true
  }
  isOpen.value = !isOpen.value
}

function select() {
  workspace.setActiveFile(props.node.path)
}

async function startCreating(type) {
  if (!isLoaded.value) {
    children.value = await workspace.loadChildren(props.node.path)
    isLoaded.value = true
  }
  isOpen.value = true
  creatingType.value = type
  newItemName.value = ''
  await nextTick()
  creatingInput.value?.focus()
}

function cancelCreate() {
  creatingType.value = null
  newItemName.value = ''
}

async function confirmCreate() {
  const name = newItemName.value.trim()
  if (!name) { cancelCreate(); return }

  const type = creatingType.value
  creatingType.value = null   // clear before await so blur doesn't double-fire
  newItemName.value = ''

  try {
    if (type === 'file') {
      await CreateFile(props.node.path, name)
    } else {
      await CreateDirectory(props.node.path, name)
    }
    children.value = await workspace.loadChildren(props.node.path)
  } catch (e) {
    console.error(e)
  }
}
</script>

<template>
  <div class="tree-node">
    <!-- Folder row -->
    <div
      v-if="node.isDir"
      class="tree-row"
      :class="{ 'is-active': workspace.activeFilePath === node.path }"
      :style="{ paddingLeft: `${depth * 20 + 4}px` }"
      role="treeitem"
      :aria-expanded="isOpen"
      @click="toggle"
    >
      <span class="tree-chevron">
        <ChevronDown v-if="isOpen" :size="14" />
        <ChevronRight v-else :size="14" />
      </span>
      <FolderOpen v-if="isOpen" :size="16" class="tree-icon icon-folder" />
      <Folder v-else :size="16" class="tree-icon icon-folder" />
      <span class="tree-label">{{ node.name }}</span>

      <span class="row-actions">
        <button class="action-btn" title="Novo arquivo" @click.stop="startCreating('file')">
          <FilePlus :size="14" />
        </button>
        <button class="action-btn" title="Nova pasta" @click.stop="startCreating('folder')">
          <FolderPlus :size="14" />
        </button>
      </span>
    </div>

    <!-- File row -->
    <div
      v-else
      class="tree-row"
      :class="{ 'is-active': workspace.activeFilePath === node.path }"
      :style="{ paddingLeft: `${depth * 20 + 22}px` }"
      role="treeitem"
      @click="select"
    >
      <FileText :size="16" class="tree-icon icon-file" />
      <span class="tree-label">{{ node.name }}</span>
    </div>

    <!-- Children -->
    <template v-if="isOpen && node.isDir">
      <!-- Inline creating input -->
      <div
        v-if="creatingType"
        class="tree-row creating-row"
        :style="{ paddingLeft: `${(depth + 1) * 20 + 22}px` }"
      >
        <Folder v-if="creatingType === 'folder'" :size="16" class="tree-icon icon-folder" />
        <FileText v-else :size="16" class="tree-icon icon-file" />
        <input
          ref="creatingInput"
          v-model="newItemName"
          class="creating-input"
          type="text"
          @keydown.enter.prevent="confirmCreate"
          @keydown.escape="cancelCreate"
          @blur="cancelCreate"
        />
      </div>

      <FileTreeNode
        v-for="child in children"
        :key="child.path"
        :node="child"
        :depth="depth + 1"
      />
    </template>
  </div>
</template>

<style scoped>
.tree-row {
  display: flex;
  align-items: center;
  gap: 4px;
  height: 22px;
  cursor: pointer;
  padding-right: 4px;
}

.tree-row:hover {
  background: var(--color-hover);
}

.tree-row.is-active {
  background: var(--color-active);
  color: var(--color-active-text);
}

.tree-row.creating-row {
  cursor: default;
  background: transparent;
}

.tree-chevron {
  width: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: var(--color-text-muted);
}

.tree-icon {
  flex-shrink: 0;
}

.icon-folder { color: var(--color-folder); }
.icon-file   { color: var(--color-file); }

.tree-label {
  flex: 1;
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

.row-actions {
  display: flex;
  gap: 2px;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.1s;
}

.tree-row:hover .row-actions {
  opacity: 1;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 3px;
  color: var(--color-text-muted);
  transition: background 0.1s, color 0.1s;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--color-text);
}
</style>

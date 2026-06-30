<script setup>
import { ref } from 'vue'
import { ChevronRight, ChevronDown, Folder, FolderOpen, FileText } from '@lucide/vue'
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
  border-radius: 0;
  padding-right: 8px;
  white-space: nowrap;
  overflow: hidden;
}

.tree-row:hover {
  background: var(--color-hover);
}

.tree-row.is-active {
  background: var(--color-active);
  color: var(--color-active-text);
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

.icon-folder {
  color: var(--color-folder);
}

.icon-file {
  color: var(--color-file);
}

.tree-label {
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>

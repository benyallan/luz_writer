<script setup>
import { useWorkspaceStore } from '@/stores/workspace'
import FileTreeNode from './FileTreeNode.vue'

const workspace = useWorkspaceStore()
</script>

<template>
  <div class="file-tree" role="tree">
    <FileTreeNode
      v-for="node in workspace.rootNodes"
      :key="node.path"
      :node="node"
      :depth="0"
    />
  </div>
</template>

<style scoped>
.file-tree {
  padding: 2px 0;
}
</style>

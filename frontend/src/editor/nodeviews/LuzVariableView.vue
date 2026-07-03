<script setup lang="ts">
import {computed} from 'vue'
import {nodeViewProps, NodeViewWrapper} from '@tiptap/vue-3'
import {useWorkspaceStore} from '../../stores/workspace'

const props = defineProps(nodeViewProps)
const workspace = useWorkspaceStore()

const name = computed(() => props.node.attrs.name as string)
const variable = computed(() => workspace.project?.variables.find(v => v.name === name.value))
const broken = computed(() => !variable.value)
</script>

<template>
  <NodeViewWrapper as="span" class="luz-variable" :class="{'luz-variable--broken': broken}" :title="name">
    {{ broken ? name : variable!.value }}
  </NodeViewWrapper>
</template>

<style scoped>
.luz-variable {
  display: inline-block;
  background: var(--luz-bg-hover);
  border-radius: 10px;
  padding: 0 8px;
  font-size: 0.9em;
}

.luz-variable--broken {
  background: rgba(192, 57, 43, 0.15);
  color: #c0392b;
}
</style>

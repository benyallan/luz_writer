<script setup lang="ts">
import {ref} from 'vue'
import {ContextMenuRoot, ContextMenuTrigger, ContextMenuPortal, ContextMenuContent, ContextMenuItem} from 'reka-ui'
import type {model} from '../../wailsjs/go/models'
import {useUiStore} from '../stores/ui'

const props = defineProps<{
  label: string
  items: model.ChapterMeta[]
  activeId?: string | null
}>()

const emit = defineEmits<{
  reorder: [ids: string[]]
  rename: [id: string, title: string]
  delete: [id: string]
  open: [id: string]
}>()

const ui = useUiStore()

const draggingId = ref<string | null>(null)
const renamingId = ref<string | null>(null)
const renameValue = ref('')

function onDragStart(id: string) {
  draggingId.value = id
}

function onDrop(targetId: string) {
  if (!draggingId.value || draggingId.value === targetId) return
  const ids = props.items.map(i => i.id)
  const from = ids.indexOf(draggingId.value)
  const to = ids.indexOf(targetId)
  if (from === -1 || to === -1) return
  ids.splice(to, 0, ids.splice(from, 1)[0])
  emit('reorder', ids)
  draggingId.value = null
}

function startRename(item: model.ChapterMeta) {
  renamingId.value = item.id
  renameValue.value = item.title
}

function commitRename(id: string) {
  if (renamingId.value === id && renameValue.value.trim()) {
    emit('rename', id, renameValue.value.trim())
  }
  renamingId.value = null
}
</script>

<template>
  <div class="explorer-group" v-if="items.length">
    <p class="explorer-group__label">{{ label }}</p>
    <ul class="explorer-group__list">
      <ContextMenuRoot v-for="item in items" :key="item.id">
        <ContextMenuTrigger as-child>
          <li
            class="explorer-item"
            :class="{'explorer-item--active': item.id === activeId}"
            draggable="true"
            @dragstart="onDragStart(item.id)"
            @dragover.prevent
            @drop="onDrop(item.id)"
          >
            <input
              v-if="renamingId === item.id"
              class="explorer-item__rename-input"
              type="text"
              v-model="renameValue"
              autofocus
              @keyup.enter="commitRename(item.id)"
              @blur="commitRename(item.id)"
            />
            <span v-else class="explorer-item__title" @click="emit('open', item.id)" @dblclick="startRename(item)">{{ item.title }}</span>
            <button
              v-if="item.hasOverrides"
              class="explorer-item__overrides-badge"
              type="button"
              title="Tem configurações de página — clique para editar"
              @click="ui.openPageOverrides(item.id)"
            >
              ⚙
            </button>
            <span class="explorer-item__words">{{ item.wordCount }}</span>
            <button class="explorer-item__delete" type="button" title="Excluir" @click="emit('delete', item.id)">×</button>
          </li>
        </ContextMenuTrigger>
        <ContextMenuPortal>
          <ContextMenuContent class="toolbar-menu">
            <ContextMenuItem class="toolbar-menu__item" @select="ui.openPageOverrides(item.id)">
              Configurações desta página...
            </ContextMenuItem>
            <ContextMenuItem class="toolbar-menu__item" @select="startRename(item)"> Renomear </ContextMenuItem>
            <ContextMenuItem class="toolbar-menu__item" @select="emit('delete', item.id)"> Excluir </ContextMenuItem>
          </ContextMenuContent>
        </ContextMenuPortal>
      </ContextMenuRoot>
    </ul>
  </div>
</template>

<style scoped>
.explorer-group {
  margin-bottom: 8px;
}

.explorer-group__label {
  margin: 0;
  padding: 8px 12px 2px;
  font-size: 0.68rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--luz-fg-muted);
}

.explorer-group__list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.explorer-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  cursor: grab;
  font-size: 0.85rem;
}

.explorer-item:hover {
  background: var(--luz-bg-hover);
}

.explorer-item--active {
  background: var(--luz-bg-hover);
  font-weight: 600;
}

.explorer-item__title {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.explorer-item__rename-input {
  flex: 1;
  min-width: 0;
  background: var(--luz-bg-editor);
  border: 1px solid var(--luz-border);
  border-radius: 4px;
  color: var(--luz-fg);
  padding: 1px 4px;
}

.explorer-item__overrides-badge {
  border: none;
  background: transparent;
  color: var(--luz-fg-muted);
  cursor: pointer;
  font-size: 0.75rem;
  padding: 0;
  line-height: 1;
}

.explorer-item__overrides-badge:hover {
  color: var(--luz-fg);
}

.explorer-item__words {
  font-size: 0.68rem;
  color: var(--luz-fg-muted);
}

.explorer-item__delete {
  border: none;
  background: transparent;
  color: var(--luz-fg-muted);
  cursor: pointer;
  font-size: 0.9rem;
  line-height: 1;
  padding: 0 2px;
  visibility: hidden;
}

.explorer-item:hover .explorer-item__delete {
  visibility: visible;
}

.explorer-item__delete:hover {
  color: #c0392b;
}

/* .toolbar-menu/.toolbar-menu__item vivem em style.css (global) — ver
   comentário lá: conteúdo teleportado pelo Reka UI não é alcançado por
   <style scoped>. */
</style>

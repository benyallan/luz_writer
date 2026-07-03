<script setup lang="ts">
import {ref} from 'vue'
import {DialogRoot, DialogPortal, DialogOverlay, DialogContent, DialogTitle} from 'reka-ui'
import {usePluginsStore} from '../stores/plugins'
import type {model} from '../../wailsjs/go/models'

const open = defineModel<boolean>('open', {default: false})
const pluginsStore = usePluginsStore()

const newName = ref('')

function slugify(s: string) {
  return s
    .toLowerCase()
    .normalize('NFD')
    .replace(/[̀-ͯ]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

async function addStyle() {
  const name = newName.value.trim()
  if (!name) return
  const id = slugify(name)
  if (pluginsStore.styles.some(s => s.id === id)) return
  const style = {id, name, italic: false, bold: false, smallCaps: false, color: undefined} as unknown as model.CustomStyle
  await pluginsStore.saveStyles([...pluginsStore.styles, style])
  newName.value = ''
}

async function updateStyle(id: string, patch: Partial<model.CustomStyle>) {
  const next = pluginsStore.styles.map(s => (s.id === id ? {...s, ...patch} : s))
  await pluginsStore.saveStyles(next as model.CustomStyle[])
}

async function removeStyle(id: string) {
  await pluginsStore.saveStyles(pluginsStore.styles.filter(s => s.id !== id))
}

function checked(e: Event) {
  return (e.target as HTMLInputElement).checked
}
</script>

<template>
  <DialogRoot v-model:open="open">
    <DialogPortal>
      <DialogOverlay class="dialog-overlay" />
      <DialogContent class="dialog-content dialog-content--wide">
        <DialogTitle class="dialog-title">Estilos Personalizados</DialogTitle>

        <ul class="styles-list">
          <li v-for="s in pluginsStore.styles" :key="s.id" class="style-row">
            <span class="style-row__name">{{ s.name }}</span>
            <label class="style-row__opt"><input type="checkbox" :checked="s.italic" @change="updateStyle(s.id, {italic: checked($event)})" /> Itálico</label>
            <label class="style-row__opt"><input type="checkbox" :checked="s.bold" @change="updateStyle(s.id, {bold: checked($event)})" /> Negrito</label>
            <label class="style-row__opt"
              ><input type="checkbox" :checked="s.smallCaps" @change="updateStyle(s.id, {smallCaps: checked($event)})" /> Versalete</label
            >
            <input
              type="color"
              class="style-row__color"
              :value="s.color || '#000000'"
              @change="updateStyle(s.id, {color: ($event.target as HTMLInputElement).value})"
            />
            <button class="style-row__delete" type="button" title="Excluir" @click="removeStyle(s.id)">×</button>
          </li>
        </ul>
        <p v-if="!pluginsStore.styles.length" class="styles-list__empty">Nenhum estilo criado ainda.</p>

        <div class="new-style-row">
          <input class="form-input" type="text" v-model="newName" placeholder="Nome do novo estilo" @keyup.enter="addStyle" />
          <button class="welcome__button" type="button" :disabled="!newName.trim()" @click="addStyle">Adicionar</button>
        </div>

        <div class="dialog-actions">
          <button class="welcome__button welcome__button--primary" type="button" @click="open = false">Fechar</button>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
}

.dialog-content {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: var(--luz-bg-editor);
  border: 1px solid var(--luz-border);
  border-radius: 8px;
  padding: 20px;
}

.dialog-content--wide {
  width: 420px;
}

.dialog-title {
  margin: 0 0 16px;
  font-size: 1.1rem;
}

.styles-list {
  list-style: none;
  margin: 0 0 12px;
  padding: 0;
  max-height: 260px;
  overflow-y: auto;
}

.styles-list__empty {
  color: var(--luz-fg-muted);
  font-size: 0.82rem;
}

.style-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 0;
  border-bottom: 1px solid var(--luz-border);
  font-size: 0.78rem;
}

.style-row__name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.style-row__opt {
  display: flex;
  align-items: center;
  gap: 3px;
  white-space: nowrap;
}

.style-row__color {
  width: 24px;
  height: 24px;
  padding: 0;
  border: none;
  background: none;
}

.style-row__delete {
  border: none;
  background: transparent;
  color: var(--luz-fg-muted);
  cursor: pointer;
}

.style-row__delete:hover {
  color: #c0392b;
}

.new-style-row {
  display: flex;
  gap: 8px;
}

.form-input {
  padding: 6px 8px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  flex: 1;
  min-width: 0;
}

.dialog-actions {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.welcome__button {
  padding: 6px 14px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  cursor: pointer;
  font-size: 0.85rem;
}

.welcome__button:hover {
  background: var(--luz-bg-hover);
}

.welcome__button--primary {
  background: var(--luz-fg);
  color: var(--luz-bg-editor);
  border-color: var(--luz-fg);
}

.welcome__button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>

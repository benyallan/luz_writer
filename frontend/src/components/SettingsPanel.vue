<script setup lang="ts">
import {ref} from 'vue'
import {useWorkspaceStore} from '../stores/workspace'

const workspace = useWorkspaceStore()

const newName = ref('')
const newValue = ref('')

function slugify(s: string) {
  return s
    .toLowerCase()
    .normalize('NFD')
    .replace(/[̀-ͯ]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

async function addVariable() {
  const name = slugify(newName.value)
  if (!name || !workspace.project) return
  if (workspace.project.variables.some(v => v.name === name)) return
  const project = {...workspace.project, variables: [...workspace.project.variables, {name, value: newValue.value}]}
  await workspace.saveProject(project)
  newName.value = ''
  newValue.value = ''
}

async function updateVariableValue(name: string, value: string) {
  if (!workspace.project) return
  const project = {
    ...workspace.project,
    variables: workspace.project.variables.map(v => (v.name === name ? {...v, value} : v)),
  }
  await workspace.saveProject(project)
}

async function removeVariable(name: string) {
  if (!workspace.project) return
  const project = {...workspace.project, variables: workspace.project.variables.filter(v => v.name !== name)}
  await workspace.saveProject(project)
}
</script>

<template>
  <div class="settings-panel">
    <p class="settings-panel__label">Variáveis do Projeto</p>

    <ul class="settings-panel__list">
      <li v-for="v in workspace.project?.variables ?? []" :key="v.name" class="variable-row">
        <span class="variable-row__name">{{ v.name }}</span>
        <input
          class="variable-row__value"
          type="text"
          :value="v.value"
          @change="updateVariableValue(v.name, ($event.target as HTMLInputElement).value)"
        />
        <button class="variable-row__delete" type="button" title="Excluir" @click="removeVariable(v.name)">×</button>
      </li>
    </ul>

    <div class="settings-panel__new">
      <input class="settings-panel__input" type="text" v-model="newName" placeholder="nome-da-variavel" />
      <input class="settings-panel__input" type="text" v-model="newValue" placeholder="valor" />
      <button class="settings-panel__add" type="button" :disabled="!newName.trim()" @click="addVariable">Adicionar</button>
    </div>
  </div>
</template>

<style scoped>
.settings-panel {
  padding: 8px 12px;
}

.settings-panel__label {
  margin: 0 0 8px;
  font-size: 0.68rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--luz-fg-muted);
}

.settings-panel__list {
  list-style: none;
  margin: 0 0 12px;
  padding: 0;
}

.variable-row {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 4px;
  font-size: 0.82rem;
}

.variable-row__name {
  flex: 0 0 auto;
  max-width: 40%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--luz-fg-muted);
}

.variable-row__value {
  flex: 1;
  min-width: 0;
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
}

.variable-row__delete {
  border: none;
  background: transparent;
  color: var(--luz-fg-muted);
  cursor: pointer;
}

.variable-row__delete:hover {
  color: #c0392b;
}

.settings-panel__new {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.settings-panel__input {
  padding: 4px 6px;
  border-radius: 4px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  font-size: 0.8rem;
}

.settings-panel__add {
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  cursor: pointer;
  font-size: 0.8rem;
}

.settings-panel__add:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.settings-panel__add:hover:not(:disabled) {
  background: var(--luz-bg-hover);
}
</style>

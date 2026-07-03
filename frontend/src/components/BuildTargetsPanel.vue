<script setup lang="ts">
import {ref, computed} from 'vue'
import {
  AccordionRoot,
  AccordionItem,
  AccordionTrigger,
  AccordionContent,
  DialogRoot,
  DialogPortal,
  DialogOverlay,
  DialogContent,
  DialogTitle,
  SelectRoot,
  SelectTrigger,
  SelectValue,
  SelectIcon,
  SelectPortal,
  SelectContent,
  SelectViewport,
  SelectItem,
  SelectItemText,
} from 'reka-ui'
import SchemaForm from './SchemaForm.vue'
import {usePluginsStore} from '../stores/plugins'
import {useWorkspaceStore} from '../stores/workspace'
import type {model} from '../../wailsjs/go/models'

const pluginsStore = usePluginsStore()
const workspace = useWorkspaceStore()

const selectedId = ref<string | null>(null)
const selected = computed(() => pluginsStore.targets.find(t => t.id === selectedId.value) ?? null)

// Abas fixas (núcleo) + opcionais habilitados com formulário não-vazio
// (seção 8.0 — as abas em Build Targets só existem para plugins com config).
const modules = computed(() => {
  const core = pluginsStore.available.filter(p => p.core)
  const optional = pluginsStore.available.filter(p => !p.core && p.enabled)
  return [...core, ...optional].filter(p => p.schema.fields.length > 0)
})

function pluginConfigOf(target: model.Target, name: string): Record<string, any> {
  const raw = (target.pluginConfig as unknown as Record<string, unknown>)?.[name]
  if (!raw) return {}
  return typeof raw === 'string' ? JSON.parse(raw) : (raw as Record<string, any>)
}

async function updateModuleConfig(name: string, value: Record<string, any>) {
  if (!selected.value) return
  const pluginConfig = {...(selected.value.pluginConfig as any), [name]: value}
  await pluginsStore.saveTarget({...selected.value, pluginConfig} as model.Target)
}

async function updateField(field: 'name' | 'fontSize' | 'documentClass' | 'kind', value: string) {
  if (!selected.value) return
  await pluginsStore.saveTarget({...selected.value, [field]: value} as model.Target)
}

async function makeActive(id: string) {
  await pluginsStore.setActiveTarget(id)
}

async function removeTarget(id: string) {
  await pluginsStore.deleteTarget(id)
  if (selectedId.value === id) selectedId.value = null
}

const BLANK_PRESET = '__blank__'

const newDialogOpen = ref(false)
const newName = ref('')
const newPresetId = ref(BLANK_PRESET)

async function createTarget() {
  const name = newName.value.trim()
  if (!name) return
  const preset = pluginsStore.presets.find(p => p.id === newPresetId.value)
  const base: model.Target = preset
    ? ({...preset, id: '', name} as model.Target)
    : ({id: '', name, kind: 'print', documentClass: 'book', fontSize: '11pt', includeToc: true, pluginConfig: {}} as unknown as model.Target)

  await pluginsStore.saveTarget(base)
  newDialogOpen.value = false
  newName.value = ''
  newPresetId.value = BLANK_PRESET
  const created = pluginsStore.targets.find(t => t.name === name)
  if (created) selectedId.value = created.id
}
</script>

<template>
  <div class="build-targets">
    <div class="build-targets__header">
      <span class="build-targets__title">Build Targets</span>
      <button class="sidebar__new-button" type="button" title="Novo target" @click="newDialogOpen = true">+</button>
    </div>

    <p v-if="!pluginsStore.targets.length" class="sidebar__empty">Nenhum target ainda.</p>

    <ul class="target-list">
      <li
        v-for="t in pluginsStore.targets"
        :key="t.id"
        class="target-row"
        :class="{'target-row--selected': t.id === selectedId}"
      >
        <span class="target-row__name" @click="selectedId = t.id">{{ t.name }}</span>
        <span v-if="workspace.project?.activeTarget === t.id" class="target-row__badge" title="Target ativo">●</span>
        <button v-else class="target-row__activate" type="button" title="Tornar ativo" @click="makeActive(t.id)">usar</button>
        <button class="target-row__delete" type="button" title="Excluir" @click="removeTarget(t.id)">×</button>
      </li>
    </ul>

    <div v-if="selected" class="target-editor">
      <label class="target-editor__field">
        Nome
        <input type="text" :value="selected.name" @change="updateField('name', ($event.target as HTMLInputElement).value)" />
      </label>
      <label class="target-editor__field">
        Classe do documento
        <select :value="selected.documentClass" @change="updateField('documentClass', ($event.target as HTMLSelectElement).value)">
          <option value="book">Livro (book)</option>
          <option value="report">Relatório (report)</option>
          <option value="article">Artigo (article)</option>
        </select>
      </label>
      <label class="target-editor__field">
        Tamanho da fonte
        <input type="text" :value="selected.fontSize" @change="updateField('fontSize', ($event.target as HTMLInputElement).value)" />
      </label>
      <label class="target-editor__field">
        Tipo (kind)
        <select :value="selected.kind" @change="updateField('kind', ($event.target as HTMLSelectElement).value)">
          <option value="print">Impresso (print)</option>
          <option value="ebook">e-Book</option>
          <option value="article">Artigo</option>
        </select>
      </label>

      <AccordionRoot type="multiple" class="target-editor__modules">
        <AccordionItem v-for="m in modules" :key="m.name" :value="m.name" class="accordion-item">
          <AccordionTrigger class="accordion-trigger">{{ m.displayName }}</AccordionTrigger>
          <AccordionContent class="accordion-content">
            <SchemaForm
              :schema="m.schema"
              :model-value="pluginConfigOf(selected, m.name)"
              @update:model-value="v => updateModuleConfig(m.name, v)"
            />
          </AccordionContent>
        </AccordionItem>
      </AccordionRoot>
    </div>

    <DialogRoot v-model:open="newDialogOpen">
      <DialogPortal>
        <DialogOverlay class="dialog-overlay" />
        <DialogContent class="dialog-content">
          <DialogTitle class="dialog-title">Novo Target</DialogTitle>

          <div class="form-field">
            <label>Nome</label>
            <input class="form-input" type="text" v-model="newName" placeholder="Meu Target" />
          </div>

          <div class="form-field">
            <label>Começar de</label>
            <SelectRoot v-model="newPresetId">
              <SelectTrigger class="select-trigger">
                <SelectValue placeholder="Vazio" />
                <SelectIcon>▾</SelectIcon>
              </SelectTrigger>
              <SelectPortal>
                <SelectContent class="select-content">
                  <SelectViewport>
                    <SelectItem :value="BLANK_PRESET" class="select-item"><SelectItemText>Vazio</SelectItemText></SelectItem>
                    <SelectItem v-for="p in pluginsStore.presets" :key="p.id" :value="p.id" class="select-item">
                      <SelectItemText>{{ p.name }}</SelectItemText>
                    </SelectItem>
                  </SelectViewport>
                </SelectContent>
              </SelectPortal>
            </SelectRoot>
          </div>

          <div class="dialog-actions">
            <button class="welcome__button" type="button" @click="newDialogOpen = false">Cancelar</button>
            <button class="welcome__button welcome__button--primary" type="button" :disabled="!newName.trim()" @click="createTarget">
              Criar
            </button>
          </div>
        </DialogContent>
      </DialogPortal>
    </DialogRoot>
  </div>
</template>

<style scoped>
.build-targets {
  padding-bottom: 12px;
}

.build-targets__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px 0;
}

.build-targets__title {
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

.sidebar__empty {
  padding: 12px;
  margin: 0;
  font-size: 0.8rem;
  color: var(--luz-fg-muted);
}

.target-list {
  list-style: none;
  margin: 4px 0 0;
  padding: 0;
}

.target-row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  font-size: 0.85rem;
}

.target-row:hover,
.target-row--selected {
  background: var(--luz-bg-hover);
}

.target-row__name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
}

.target-row__badge {
  color: #2e8b57;
}

.target-row__activate {
  border: none;
  background: transparent;
  color: var(--luz-fg-muted);
  font-size: 0.68rem;
  cursor: pointer;
  text-transform: uppercase;
}

.target-row__activate:hover {
  color: var(--luz-fg);
}

.target-row__delete {
  border: none;
  background: transparent;
  color: var(--luz-fg-muted);
  cursor: pointer;
}

.target-row__delete:hover {
  color: #c0392b;
}

.target-editor {
  padding: 10px 12px 0;
  border-top: 1px solid var(--luz-border);
  margin-top: 8px;
}

.target-editor__field {
  display: flex;
  flex-direction: column;
  gap: 3px;
  font-size: 0.78rem;
  color: var(--luz-fg-muted);
  margin-bottom: 8px;
}

.target-editor__field input,
.target-editor__field select {
  padding: 4px 6px;
  border-radius: 4px;
  border: 1px solid var(--luz-border);
  color: var(--luz-fg);
  font-size: 0.85rem;
}

.target-editor__field input {
  background: transparent;
}

/* Native <select> não pode ter fundo transparente: no WebKitGTK (Wails no
   Linux) a lista de opções herda o background do próprio elemento, então
   "transparent" a torna ilegível contra o que estiver atrás da janela. */
.target-editor__field select {
  background: var(--luz-bg-editor);
}

.target-editor__modules {
  margin-top: 8px;
}

.accordion-item {
  border-top: 1px solid var(--luz-border);
}

.accordion-trigger {
  width: 100%;
  text-align: left;
  padding: 6px 0;
  border: none;
  background: transparent;
  color: var(--luz-fg);
  font-size: 0.82rem;
  cursor: pointer;
}

.accordion-content {
  padding: 4px 0 10px;
}

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
  width: 320px;
  background: var(--luz-bg-editor);
  border: 1px solid var(--luz-border);
  border-radius: 8px;
  padding: 20px;
}

.dialog-title {
  margin: 0 0 16px;
  font-size: 1.1rem;
}

.dialog-actions {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.form-field {
  margin-bottom: 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  text-align: left;
}

.form-field label {
  font-size: 0.78rem;
  color: var(--luz-fg-muted);
}

.form-input {
  padding: 6px 8px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  width: 100%;
}

.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  border-radius: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  cursor: pointer;
}

/* .select-content/.select-item vivem em style.css (global) — ver comentário
   lá: conteúdo teleportado pelo Reka UI não é alcançado por <style scoped>. */

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

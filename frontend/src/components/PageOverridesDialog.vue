<script setup lang="ts">
import {ref, computed, watch} from 'vue'
import {DialogRoot, DialogPortal, DialogOverlay, DialogContent, DialogTitle} from 'reka-ui'
import SchemaForm from './SchemaForm.vue'
import {usePluginsStore} from '../stores/plugins'
import {useWorkspaceStore} from '../stores/workspace'
import * as App from '../../wailsjs/go/main/App'

// Painel "Configurações desta página" (seções 5.6 e 8.4): SchemaForm em modo
// override, um módulo por vez, com toggle por campo e valor herdado do
// target esmaecido quando desligado.
const props = defineProps<{documentId: string | null}>()
const open = defineModel<boolean>('open', {default: false})

const pluginsStore = usePluginsStore()
const workspace = useWorkspaceStore()

const overrides = ref<Record<string, any>>({})

const activeTarget = computed(() => pluginsStore.targets.find(t => t.id === workspace.project?.activeTarget) ?? null)

// Só módulos ativos (núcleo + opcionais habilitados) com DocumentScope==true
// aceitam sobrescrita por página.
const scopedModules = computed(() => {
  const core = pluginsStore.available.filter(p => p.core && p.documentScope)
  const optional = pluginsStore.available.filter(p => !p.core && p.enabled && p.documentScope)
  return [...core, ...optional]
})

function targetConfigOf(name: string): Record<string, any> {
  const raw = (activeTarget.value?.pluginConfig as unknown as Record<string, unknown>)?.[name]
  if (!raw) return {}
  return typeof raw === 'string' ? JSON.parse(raw) : (raw as Record<string, any>)
}

async function load() {
  if (!props.documentId) return
  const json = await App.GetDocumentOverrides(props.documentId)
  overrides.value = JSON.parse(json)
}

watch(open, v => {
  if (v) load()
})

async function save() {
  if (!props.documentId) return
  await App.SaveDocumentOverrides(props.documentId, JSON.stringify(overrides.value))
  await workspace.refreshChapters()
}

async function updateModuleOverride(name: string, value: Record<string, any>) {
  const next = {...overrides.value}
  if (Object.keys(value).length === 0) {
    delete next[name]
  } else {
    next[name] = value
  }
  overrides.value = next
  await save()
}

async function clearAll() {
  overrides.value = {}
  await save()
}
</script>

<template>
  <DialogRoot v-model:open="open">
    <DialogPortal>
      <DialogOverlay class="dialog-overlay" />
      <DialogContent class="dialog-content dialog-content--wide">
        <DialogTitle class="dialog-title">Configurações desta Página</DialogTitle>

        <div v-for="m in scopedModules" :key="m.name" class="page-overrides__module">
          <p class="page-overrides__module-title">{{ m.displayName }}</p>
          <SchemaForm
            :schema="m.schema"
            :model-value="overrides[m.name] || {}"
            override-mode
            :inherited-values="targetConfigOf(m.name)"
            :inherited-label="activeTarget ? `de “${activeTarget.name}”` : 'do target'"
            @update:model-value="v => updateModuleOverride(m.name, v)"
          />
        </div>
        <p v-if="!scopedModules.length" class="page-overrides__empty">
          Nenhum módulo ativo aceita configuração por página.
        </p>

        <div class="dialog-actions">
          <button class="welcome__button" type="button" @click="clearAll">Limpar sobrescritas</button>
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
  max-height: 80vh;
  overflow-y: auto;
}

.dialog-content--wide {
  width: 380px;
}

.dialog-title {
  margin: 0 0 16px;
  font-size: 1.1rem;
}

.page-overrides__module {
  border-top: 1px solid var(--luz-border);
  padding-top: 10px;
  margin-top: 10px;
}

.page-overrides__module-title {
  margin: 0 0 8px;
  font-size: 0.78rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--luz-fg-muted);
}

.page-overrides__empty {
  color: var(--luz-fg-muted);
  font-size: 0.82rem;
}

.dialog-actions {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
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
</style>

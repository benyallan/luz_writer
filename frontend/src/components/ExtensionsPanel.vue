<script setup lang="ts">
import {ref, computed} from 'vue'
import {usePluginsStore} from '../stores/plugins'
import StylesManagerDialog from './StylesManagerDialog.vue'

const pluginsStore = usePluginsStore()
const stylesDialogOpen = ref(false)

const core = computed(() => pluginsStore.available.filter(p => p.core))
const optional = computed(() => pluginsStore.available.filter(p => !p.core))

async function toggle(name: string, enabled: boolean) {
  await pluginsStore.setPluginEnabled(name, enabled)
}
</script>

<template>
  <div class="extensions-panel">
    <p class="extensions-panel__label">Essenciais</p>
    <ul class="extensions-panel__list">
      <li v-for="p in core" :key="p.name" class="plugin-row">
        <div class="plugin-row__main">
          <span class="plugin-row__name">{{ p.displayName }}</span>
          <input type="checkbox" checked disabled title="Módulo do núcleo — sempre ativo" />
        </div>
        <p class="plugin-row__description">{{ p.description }}</p>
      </li>
    </ul>

    <p class="extensions-panel__label">Plugins</p>
    <ul class="extensions-panel__list">
      <li v-for="p in optional" :key="p.name" class="plugin-row">
        <div class="plugin-row__main">
          <span class="plugin-row__name">{{ p.displayName }}</span>
          <input type="checkbox" :checked="p.enabled" @change="toggle(p.name, ($event.target as HTMLInputElement).checked)" />
        </div>
        <p class="plugin-row__description">{{ p.description }}</p>
        <button
          v-if="p.name === 'customStyles' && p.enabled"
          class="plugin-row__manage"
          type="button"
          @click="stylesDialogOpen = true"
        >
          Gerenciar estilos...
        </button>
      </li>
    </ul>

    <StylesManagerDialog v-model:open="stylesDialogOpen" />
  </div>
</template>

<style scoped>
.extensions-panel {
  padding-bottom: 12px;
}

.extensions-panel__label {
  margin: 12px 0 4px;
  padding: 0 12px;
  font-size: 0.68rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--luz-fg-muted);
}

.extensions-panel__list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.plugin-row {
  padding: 6px 12px;
  border-bottom: 1px solid var(--luz-border);
}

.plugin-row__main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.plugin-row__name {
  font-size: 0.85rem;
}

.plugin-row__description {
  margin: 2px 0 0;
  font-size: 0.72rem;
  color: var(--luz-fg-muted);
}

.plugin-row__manage {
  margin-top: 6px;
  border: 1px solid var(--luz-border);
  background: transparent;
  color: var(--luz-fg);
  border-radius: 4px;
  padding: 3px 8px;
  font-size: 0.72rem;
  cursor: pointer;
}

.plugin-row__manage:hover {
  background: var(--luz-bg-hover);
}
</style>

<script setup lang="ts">
import {useUiStore} from '../stores/ui'
import type {ActivityPanel} from '../stores/ui'

const ui = useUiStore()

const items: {panel: ActivityPanel; label: string}[] = [
  {panel: 'explorer', label: 'Explorer'},
  {panel: 'buildTargets', label: 'Build Targets'},
  {panel: 'extensions', label: 'Extensions'},
  {panel: 'settings', label: 'Settings'},
]
</script>

<template>
  <nav class="activity-bar">
    <button
      v-for="item in items"
      :key="item.panel"
      class="activity-bar__item"
      :class="{'activity-bar__item--active': ui.activePanel === item.panel}"
      type="button"
      :title="item.label"
      @click="ui.activePanel = item.panel"
    >
      {{ item.label[0] }}
    </button>
  </nav>
</template>

<style scoped>
.activity-bar {
  width: 48px;
  flex: 0 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding-top: 8px;
  background: var(--luz-bg-activity-bar);
  border-right: 1px solid var(--luz-border);
}

.activity-bar__item {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--luz-fg-muted);
  cursor: pointer;
  font-size: 0.8rem;
}

.activity-bar__item:hover {
  background: var(--luz-bg-hover);
  color: var(--luz-fg);
}

.activity-bar__item--active {
  background: var(--luz-bg-hover);
  color: var(--luz-fg);
}
</style>

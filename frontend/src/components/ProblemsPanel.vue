<script setup lang="ts">
import {ref, computed} from 'vue'
import {CollapsibleRoot, CollapsibleTrigger, CollapsibleContent} from 'reka-ui'
import {useBuildStore} from '../stores/build'
import {useProblemsStore} from '../stores/problems'
import {useEditorStore} from '../stores/editor'
import {useUiStore} from '../stores/ui'
import type {model} from '../../wailsjs/go/models'

// Problems Panel (Etapa 5): lista os Problems do Rule Engine, revalidado
// automaticamente pelo backend a cada save relevante, mais o LogTail de uma
// falha real de compilação (quando o Tectonic falha por algo que o Rule
// Engine não cobre).
const open = ref(false)
const buildStore = useBuildStore()
const problemsStore = useProblemsStore()
const editorStore = useEditorStore()
const ui = useUiStore()

const problems = computed(() => problemsStore.items)
const compileFailed = computed(() => buildStore.lastResult !== null && !buildStore.lastResult.success)

function navigateTo(p: model.Problem) {
  const sep = p.source.indexOf(':')
  const kind = sep === -1 ? p.source : p.source.slice(0, sep)
  const id = sep === -1 ? '' : p.source.slice(sep + 1)

  switch (kind) {
    case 'chapter':
      ui.activePanel = 'explorer'
      if (id) editorStore.openChapter(id)
      break
    case 'override':
      if (id) ui.openPageOverrides(id)
      break
    case 'target':
      ui.activePanel = 'buildTargets'
      break
    case 'plugin':
    case 'styles':
      ui.activePanel = 'extensions'
      break
    default:
      break
  }
}
</script>

<template>
  <CollapsibleRoot v-model:open="open" class="problems-panel">
    <CollapsibleTrigger class="problems-panel__trigger">
      <span>Problems</span>
      <span class="problems-panel__count" :class="{'problems-panel__count--error': problemsStore.hasErrors}">
        {{ problems.length }}
      </span>
    </CollapsibleTrigger>
    <CollapsibleContent class="problems-panel__content">
      <p v-if="!problems.length && !compileFailed" class="problems-panel__empty">Nenhum problema encontrado.</p>

      <div v-if="compileFailed" class="problems-panel__failure">
        <p class="problems-panel__failure-title">A compilação falhou.</p>
        <pre class="problems-panel__log">{{ buildStore.lastResult?.logTail }}</pre>
      </div>

      <ul v-if="problems.length" class="problems-panel__list">
        <li
          v-for="(p, i) in problems"
          :key="i"
          class="problems-panel__item"
          :class="{'problems-panel__item--clickable': p.source.includes(':')}"
          @click="navigateTo(p)"
        >
          <span :class="`problems-panel__severity problems-panel__severity--${p.severity}`">{{ p.severity }}</span>
          <span v-if="p.code" class="problems-panel__code">{{ p.code }}</span>
          {{ p.message }}
        </li>
      </ul>
    </CollapsibleContent>
  </CollapsibleRoot>
</template>

<style scoped>
.problems-panel {
  flex: 0 0 auto;
  border-top: 1px solid var(--luz-border);
  background: var(--luz-bg-sidebar);
}

.problems-panel__trigger {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 12px;
  border: none;
  background: transparent;
  color: var(--luz-fg-muted);
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  cursor: pointer;
}

.problems-panel__count {
  background: var(--luz-bg-hover);
  border-radius: 8px;
  padding: 0 6px;
  font-size: 0.68rem;
}

.problems-panel__count--error {
  background: #c0392b;
  color: white;
}

.problems-panel__content {
  max-height: 220px;
  overflow-y: auto;
  padding: 4px 12px 12px;
}

.problems-panel__empty {
  margin: 0;
  color: var(--luz-fg-muted);
  font-size: 0.8rem;
}

.problems-panel__failure-title {
  margin: 4px 0;
  color: #c0392b;
  font-size: 0.82rem;
  font-weight: 600;
}

.problems-panel__log {
  background: var(--luz-bg-editor);
  border: 1px solid var(--luz-border);
  border-radius: 4px;
  padding: 8px;
  font-size: 0.72rem;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 140px;
  overflow-y: auto;
}

.problems-panel__list {
  list-style: none;
  margin: 4px 0 0;
  padding: 0;
}

.problems-panel__item {
  font-size: 0.8rem;
  padding: 3px 0;
  display: flex;
  gap: 6px;
  align-items: baseline;
}

.problems-panel__item--clickable {
  cursor: pointer;
}

.problems-panel__item--clickable:hover {
  background: var(--luz-bg-hover);
}

.problems-panel__severity {
  font-size: 0.65rem;
  text-transform: uppercase;
  color: var(--luz-fg-muted);
  flex: 0 0 auto;
}

.problems-panel__severity--error {
  color: #c0392b;
}

.problems-panel__severity--warning {
  color: #b8860b;
}

.problems-panel__code {
  font-size: 0.68rem;
  color: var(--luz-fg-muted);
  flex: 0 0 auto;
}
</style>

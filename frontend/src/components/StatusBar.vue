<script setup lang="ts">
import {computed} from 'vue'
import {DropdownMenuRoot, DropdownMenuTrigger, DropdownMenuPortal, DropdownMenuContent, DropdownMenuItem} from 'reka-ui'
import {useEditorStore} from '../stores/editor'
import {useBuildStore} from '../stores/build'
import {usePluginsStore} from '../stores/plugins'
import {useWorkspaceStore} from '../stores/workspace'
import {useProblemsStore} from '../stores/problems'

const editorStore = useEditorStore()
const buildStore = useBuildStore()
const pluginsStore = usePluginsStore()
const workspace = useWorkspaceStore()
const problemsStore = useProblemsStore()

const activeTargetName = computed(() => {
  const t = pluginsStore.targets.find(t => t.id === workspace.project?.activeTarget)
  return t?.name ?? 'Nenhum target ativo'
})

const saveLabel = computed(() => {
  switch (editorStore.saveStatus) {
    case 'saving':
      return 'Salvando...'
    case 'unsaved':
      return 'Alterações não salvas'
    default:
      return 'Salvo'
  }
})

const exportLabel = computed(() => {
  if (!buildStore.compiling) return 'Exportar PDF'
  const stageLabels: Record<string, string> = {
    validating: 'Validando...',
    generating: 'Gerando...',
    compiling: 'Compilando...',
    done: 'Concluído',
  }
  const label = buildStore.stage ? stageLabels[buildStore.stage] ?? buildStore.stage : 'Exportando...'
  return `${label} ${buildStore.percent}%`
})

const problemsLabel = computed(() => {
  if (problemsStore.hasErrors) return `✗ ${problemsStore.errorCount} problema(s)`
  if (problemsStore.items.length) return `⚠ ${problemsStore.items.length} aviso(s)`
  return '✓ Sem problemas'
})

// Seção 10: ausência do tectonic nunca trava o app, só desabilita a
// exportação — sinalizado aqui de forma persistente (não some sozinho).
const tectonicMissing = computed(() => buildStore.tectonicAvailable === false)
</script>

<template>
  <footer class="status-bar">
    <DropdownMenuRoot v-if="pluginsStore.targets.length">
      <DropdownMenuTrigger class="status-bar__target" title="Trocar target ativo">
        {{ activeTargetName }}
      </DropdownMenuTrigger>
      <DropdownMenuPortal>
        <DropdownMenuContent class="status-bar__target-menu" side="top">
          <DropdownMenuItem
            v-for="t in pluginsStore.targets"
            :key="t.id"
            class="status-bar__target-menu-item"
            @select="pluginsStore.setActiveTarget(t.id)"
          >
            {{ t.name }}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenuPortal>
    </DropdownMenuRoot>
    <span v-else class="status-bar__item">Nenhum target ativo</span>

    <span class="status-bar__spacer" />
    <span v-if="editorStore.activeChapterMeta" class="status-bar__item">
      {{ editorStore.activeChapterMeta.wordCount }} palavras
    </span>
    <span v-if="editorStore.activeChapterId" class="status-bar__item">{{ saveLabel }}</span>
    <span v-if="tectonicMissing" class="status-bar__item status-bar__item--error" title="Instale o Tectonic e reinicie o Luz Writer para habilitar a exportação">
      ✗ Tectonic não encontrado — exportação desabilitada
    </span>
    <button
      class="status-bar__export"
      type="button"
      :disabled="buildStore.compiling || problemsStore.hasErrors || tectonicMissing"
      :title="problemsStore.hasErrors ? 'Corrija os erros no painel Problems antes de exportar' : tectonicMissing ? 'Tectonic não encontrado' : 'Ctrl/Cmd+E'"
      @click="buildStore.compile()"
    >
      {{ exportLabel }}
    </button>
    <span class="status-bar__item" :class="{'status-bar__item--error': problemsStore.hasErrors}">{{ problemsLabel }}</span>
  </footer>
</template>

<style scoped>
.status-bar {
  height: 24px;
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 0 12px;
  background: var(--luz-bg-status-bar);
  border-top: 1px solid var(--luz-border);
  font-size: 0.72rem;
  color: var(--luz-fg-muted);
}

.status-bar__spacer {
  flex: 1;
}

.status-bar__target {
  border: none;
  background: transparent;
  color: var(--luz-fg-muted);
  cursor: pointer;
  font-size: 0.72rem;
  padding: 0;
}

.status-bar__target:hover {
  color: var(--luz-fg);
  text-decoration: underline;
}

/* :global() porque o DropdownMenuContent do Reka UI é teleportado — <style
   scoped> não alcança conteúdo dentro do Teleport (o atributo de escopo do
   Vue cai no wrapper de posicionamento, não nestes elementos). */
:global(.status-bar__target-menu) {
  background: var(--luz-bg-editor);
  border: 1px solid var(--luz-border);
  border-radius: 6px;
  padding: 4px;
  min-width: 160px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

:global(.status-bar__target-menu-item) {
  padding: 6px 8px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8rem;
  outline: none;
}

:global(.status-bar__target-menu-item[data-highlighted]) {
  background: var(--luz-bg-hover);
}

.status-bar__export {
  border: none;
  background: transparent;
  color: var(--luz-fg-muted);
  cursor: pointer;
  font-size: 0.72rem;
  padding: 0 6px;
}

.status-bar__export:hover:not(:disabled) {
  color: var(--luz-fg);
  text-decoration: underline;
}

.status-bar__export:disabled {
  cursor: default;
}

.status-bar__item--error {
  color: #c0392b;
}
</style>

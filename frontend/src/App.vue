<script setup lang="ts">
import {computed, onMounted, onBeforeUnmount, watch} from 'vue'
import ActivityBar from './components/ActivityBar.vue'
import Sidebar from './components/Sidebar.vue'
import EditorArea from './components/EditorArea.vue'
import StatusBar from './components/StatusBar.vue'
import ProblemsPanel from './components/ProblemsPanel.vue'
import WelcomeScreen from './components/WelcomeScreen.vue'
import PageOverridesDialog from './components/PageOverridesDialog.vue'
import ToastHost from './components/ToastHost.vue'
import {useWorkspaceStore} from './stores/workspace'
import {usePluginsStore} from './stores/plugins'
import {useUiStore} from './stores/ui'
import {useProblemsStore} from './stores/problems'
import {useEditorStore} from './stores/editor'
import {useBuildStore} from './stores/build'
import {useToastStore} from './stores/toast'

const workspace = useWorkspaceStore()
const pluginsStore = usePluginsStore()
const ui = useUiStore()
const problemsStore = useProblemsStore()
const editorStore = useEditorStore()
const buildStore = useBuildStore()
const toastStore = useToastStore()

// Rede de segurança para chamadas ao backend sem tratamento local de erro
// (ex.: salvar um target com dado inválido a partir de um diálogo) — em vez
// de a falha ser engolida em silêncio, vira uma notificação visível.
function handleRejection(e: PromiseRejectionEvent) {
  const reason = e.reason
  const message = typeof reason === 'string' ? reason : (reason?.message ?? String(reason))
  toastStore.push(message)
  e.preventDefault()
}

onMounted(() => {
  workspace.restoreLastWorkspace()
})

// Atalhos de teclado (seção 12, Etapa 6): Ctrl/Cmd+S salva o capítulo aberto
// imediatamente (sem esperar o debounce do autosave); Ctrl/Cmd+E exporta,
// respeitando o mesmo bloqueio por erros do botão "Exportar PDF".
function handleKeydown(e: KeyboardEvent) {
  if (!(e.ctrlKey || e.metaKey)) return
  if (e.key.toLowerCase() === 's') {
    e.preventDefault()
    if (editorStore.activeChapterId) editorStore.flushPendingSave()
  } else if (e.key.toLowerCase() === 'e') {
    e.preventDefault()
    if (!buildStore.compiling && !problemsStore.hasErrors) buildStore.compile()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
  window.addEventListener('unhandledrejection', handleRejection)
})
onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('unhandledrejection', handleRejection)
})

watch(
  () => workspace.isOpen,
  isOpen => {
    if (isOpen) {
      pluginsStore.refreshAll()
      problemsStore.refresh()
    }
  },
)

// workspace.error é tratado (não vira unhandledrejection) e só é exibido
// inline enquanto a WelcomeScreen está montada — depois de abrir o
// workspace, precisa de outro jeito de chegar até o usuário.
watch(
  () => workspace.error,
  message => {
    if (message && workspace.isOpen) toastStore.push(message)
  },
)

const pageOverridesOpen = computed({
  get: () => ui.pageOverridesDocId !== null,
  set: (v: boolean) => {
    if (!v) ui.pageOverridesDocId = null
  },
})
</script>

<template>
  <div class="app-shell">
    <template v-if="workspace.isOpen">
      <div class="app-shell__main">
        <ActivityBar />
        <Sidebar />
        <EditorArea />
      </div>
      <ProblemsPanel />
      <StatusBar />
      <PageOverridesDialog :document-id="ui.pageOverridesDocId" v-model:open="pageOverridesOpen" />
    </template>
    <WelcomeScreen v-else />
    <ToastHost />
  </div>
</template>

<style scoped>
.app-shell {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.app-shell__main {
  flex: 1 1 auto;
  min-height: 0;
  display: flex;
}
</style>

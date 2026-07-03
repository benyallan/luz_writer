import {ref, computed} from 'vue'
import {defineStore} from 'pinia'
import {useWorkspaceStore} from './workspace'
import * as App from '../../wailsjs/go/main/App'

export type SaveStatus = 'saved' | 'saving' | 'unsaved'

export const useEditorStore = defineStore('editor', () => {
  const workspace = useWorkspaceStore()

  const activeChapterId = ref<string | null>(null)
  const initialContent = ref<any>(null)
  const saveStatus = ref<SaveStatus>('saved')

  const activeChapterMeta = computed(() => workspace.chapters.find(c => c.id === activeChapterId.value) ?? null)

  async function openChapter(id: string) {
    const json = await App.LoadChapter(id)
    activeChapterId.value = id
    initialContent.value = JSON.parse(json)
    saveStatus.value = 'saved'
  }

  function closeChapter() {
    activeChapterId.value = null
    initialContent.value = null
  }

  let saveTimer: ReturnType<typeof setTimeout> | null = null
  let pendingGetContent: (() => unknown) | null = null

  // Autosave com debounce ~800ms (Etapa 2). getContent é chamado só quando o
  // timer dispara, para sempre gravar o conteúdo mais recente do editor.
  function scheduleSave(getContent: () => unknown) {
    if (!activeChapterId.value) return
    saveStatus.value = 'unsaved'
    pendingGetContent = getContent
    if (saveTimer) clearTimeout(saveTimer)
    saveTimer = setTimeout(() => {
      flushPendingSave()
    }, 800)
  }

  // Grava imediatamente qualquer alteração ainda no debounce — usado antes de
  // Exportar PDF, para a compilação nunca ler uma versão desatualizada do
  // capítulo em disco (seção 10 lê direto de capitulos/<id>.json).
  async function flushPendingSave() {
    if (saveTimer) {
      clearTimeout(saveTimer)
      saveTimer = null
    }
    const id = activeChapterId.value
    const getContent = pendingGetContent
    if (!id || !getContent) return
    pendingGetContent = null
    saveStatus.value = 'saving'
    await App.SaveChapter(id, JSON.stringify(getContent()))
    saveStatus.value = 'saved'
    await workspace.refreshChapters()
  }

  return {
    activeChapterId,
    initialContent,
    saveStatus,
    activeChapterMeta,
    openChapter,
    closeChapter,
    scheduleSave,
    flushPendingSave,
  }
})

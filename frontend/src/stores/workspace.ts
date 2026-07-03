import {computed, ref} from 'vue'
import {defineStore} from 'pinia'
import * as App from '../../wailsjs/go/main/App'
import {model} from '../../wailsjs/go/models'

const LAST_WORKSPACE_KEY = 'luz-writer:lastWorkspacePath'

export const PRE_TEXTUAL_ROLES = ['dedication', 'epigraph', 'acknowledgments', 'preface']
export const POST_TEXTUAL_ROLES = ['aboutAuthor', 'appendix']

export const useWorkspaceStore = defineStore('workspace', () => {
  const path = ref<string | null>(null)
  const project = ref<model.Project | null>(null)
  const chapters = ref<model.ChapterMeta[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isOpen = computed(() => project.value !== null)

  const preTextuais = computed(() => chapters.value.filter(c => PRE_TEXTUAL_ROLES.includes(c.role)))
  const capitulos = computed(() => chapters.value.filter(c => c.role === 'chapter'))
  const posTextuais = computed(() => chapters.value.filter(c => POST_TEXTUAL_ROLES.includes(c.role)))

  // Recarrega tanto a lista de capítulos quanto o project.json: CreateChapter/
  // DeleteChapter/ReorderChapters mudam chapterOrder no disco por baixo dos
  // pés do estado em memória, e um SaveProject subsequente (ex.: editar uma
  // variável em Settings) não pode sobrescrever isso com um chapterOrder
  // desatualizado.
  async function refreshChapters() {
    chapters.value = await App.ListChapters()
    if (project.value) {
      project.value = await App.GetProject()
    }
  }

  function applyOpenedWorkspace(info: model.WorkspaceInfo) {
    path.value = info.path
    project.value = info.project
    localStorage.setItem(LAST_WORKSPACE_KEY, info.path)
  }

  async function run<T>(fn: () => Promise<T>): Promise<T | undefined> {
    loading.value = true
    error.value = null
    try {
      return await fn()
    } catch (e) {
      error.value = String(e)
      return undefined
    } finally {
      loading.value = false
    }
  }

  async function createWorkspace(newPath: string, title: string, author: string, language: string) {
    return run(async () => {
      const info = await App.CreateWorkspace(newPath, title, author, language)
      applyOpenedWorkspace(info)
      await refreshChapters()
      return info
    })
  }

  async function openWorkspace(existingPath: string) {
    return run(async () => {
      const info = await App.OpenWorkspace(existingPath)
      applyOpenedWorkspace(info)
      await refreshChapters()
      return info
    })
  }

  async function openWorkspaceDialog() {
    return run(async () => {
      const info = await App.OpenWorkspaceDialog()
      applyOpenedWorkspace(info)
      await refreshChapters()
      return info
    })
  }

  async function pickDirectory() {
    return run(() => App.PickDirectory())
  }

  // Tenta restaurar o workspace aberto na última sessão (seção "Aceite" da
  // Etapa 1: reabrir o app restaura o estado). Silencioso em caso de falha
  // (ex.: pasta movida/apagada) — apenas limpa o registro e volta à tela
  // inicial de criar/abrir projeto.
  async function restoreLastWorkspace() {
    const lastPath = localStorage.getItem(LAST_WORKSPACE_KEY)
    if (!lastPath) return
    try {
      const info = await App.OpenWorkspace(lastPath)
      applyOpenedWorkspace(info)
      await refreshChapters()
    } catch {
      localStorage.removeItem(LAST_WORKSPACE_KEY)
    }
  }

  async function saveProject(p: model.Project) {
    return run(async () => {
      await App.SaveProject(p)
      project.value = p
    })
  }

  async function createChapter(title: string, role: string) {
    return run(async () => {
      const meta = await App.CreateChapter(title, role)
      await refreshChapters()
      return meta
    })
  }

  async function deleteChapter(id: string) {
    return run(async () => {
      await App.DeleteChapter(id)
      await refreshChapters()
    })
  }

  async function reorderChapters(order: string[]) {
    return run(async () => {
      await App.ReorderChapters(order)
      await refreshChapters()
    })
  }

  // Renomear é feito reescrevendo o texto do primeiro nó do documento — não
  // existe um método de ponte dedicado (o título é derivado do conteúdo).
  async function renameChapter(id: string, newTitle: string) {
    return run(async () => {
      const content = {
        type: 'doc',
        content: [
          {type: 'paragraph', content: [{type: 'text', text: newTitle}]},
        ],
      }
      await App.SaveChapter(id, JSON.stringify(content))
      await refreshChapters()
    })
  }

  return {
    path,
    project,
    chapters,
    loading,
    error,
    isOpen,
    preTextuais,
    capitulos,
    posTextuais,
    createWorkspace,
    openWorkspace,
    openWorkspaceDialog,
    pickDirectory,
    restoreLastWorkspace,
    refreshChapters,
    saveProject,
    createChapter,
    deleteChapter,
    reorderChapters,
    renameChapter,
  }
})

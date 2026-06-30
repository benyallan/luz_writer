import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { OpenFolder, ReadDirectory } from '@wails/go/main/App'

export const useWorkspaceStore = defineStore('workspace', () => {
  const rootPath = ref(null)
  const rootNodes = ref([])
  const activeFilePath = ref(null)
  const isNewProjectDialogOpen = ref(false)

  const rootName = computed(() =>
    rootPath.value ? rootPath.value.split('/').pop() : null
  )

  async function openFolder() {
    const path = await OpenFolder()
    if (!path) return
    await openAt(path)
  }

  async function openAt(path) {
    rootPath.value = path
    rootNodes.value = await ReadDirectory(path)
  }

  async function loadChildren(dirPath) {
    return await ReadDirectory(dirPath)
  }

  function setActiveFile(path) {
    activeFilePath.value = path
  }

  async function refreshRoot() {
    if (rootPath.value) {
      rootNodes.value = await ReadDirectory(rootPath.value)
    }
  }

  function newProject() {
    isNewProjectDialogOpen.value = true
  }

  return {
    rootPath,
    rootName,
    rootNodes,
    activeFilePath,
    isNewProjectDialogOpen,
    openFolder,
    openAt,
    refreshRoot,
    loadChildren,
    setActiveFile,
    newProject,
  }
})

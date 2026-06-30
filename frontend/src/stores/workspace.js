import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { OpenFolder, ReadDirectory } from '@wails/go/main/App'

export const useWorkspaceStore = defineStore('workspace', () => {
  const rootPath = ref(null)
  const rootNodes = ref([])
  const activeFilePath = ref(null)

  const rootName = computed(() =>
    rootPath.value ? rootPath.value.split('/').pop() : null
  )

  async function openFolder() {
    const path = await OpenFolder()
    if (!path) return
    rootPath.value = path
    rootNodes.value = await ReadDirectory(path)
  }

  async function loadChildren(dirPath) {
    return await ReadDirectory(dirPath)
  }

  function setActiveFile(path) {
    activeFilePath.value = path
  }

  function newProject() {
    // TODO: open new project dialog
  }

  return {
    rootPath,
    rootName,
    rootNodes,
    activeFilePath,
    openFolder,
    loadChildren,
    setActiveFile,
    newProject,
  }
})

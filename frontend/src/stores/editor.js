import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { ReadFile, SaveFile } from '@wails/go/main/App'

export const useEditorStore = defineStore('editor', () => {
  const filePath = ref(null)
  const content = ref('')
  const savedContent = ref('')

  const fileName = computed(() =>
    filePath.value ? filePath.value.split('/').pop() : null
  )
  const isDirty = computed(() => content.value !== savedContent.value)

  async function openFile(path) {
    const text = await ReadFile(path)
    filePath.value = path
    content.value = text
    savedContent.value = text
  }

  async function save() {
    if (!filePath.value || !isDirty.value) return
    await SaveFile(filePath.value, content.value)
    savedContent.value = content.value
  }

  function close() {
    filePath.value = null
    content.value = ''
    savedContent.value = ''
  }

  return { filePath, content, savedContent, fileName, isDirty, openFile, save, close }
})

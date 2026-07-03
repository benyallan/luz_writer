import {ref} from 'vue'
import {defineStore} from 'pinia'
import * as App from '../../wailsjs/go/main/App'
import {EventsOn, BrowserOpenURL} from '../../wailsjs/runtime/runtime'
import type {model} from '../../wailsjs/go/models'
import {useEditorStore} from './editor'

export const useBuildStore = defineStore('build', () => {
  const stage = ref<string | null>(null)
  const percent = ref(0)
  const compiling = ref(false)
  const lastResult = ref<model.BuildResult | null>(null)

  EventsOn('luz:build:progress', (p: {stage: string; percent: number}) => {
    stage.value = p.stage
    percent.value = p.percent
  })

  EventsOn('luz:build:done', (result: model.BuildResult) => {
    lastResult.value = result
  })

  async function compile() {
    compiling.value = true
    stage.value = 'validating'
    percent.value = 0
    try {
      // Garante que nenhuma alteração recente fique presa no debounce do
      // autosave antes de ler os capítulos do disco.
      await useEditorStore().flushPendingSave()
      const result = await App.Compile()
      lastResult.value = result
      if (result.success) {
        openContainingFolder(result.outputPath)
      }
      return result
    } finally {
      compiling.value = false
    }
  }

  function openContainingFolder(filePath: string) {
    const dir = filePath.replace(/[\\/][^\\/]*$/, '')
    BrowserOpenURL('file://' + dir)
  }

  return {stage, percent, compiling, lastResult, compile}
})

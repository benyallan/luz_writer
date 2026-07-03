import {ref} from 'vue'
import {defineStore} from 'pinia'

export type ActivityPanel = 'explorer' | 'buildTargets' | 'extensions' | 'settings'

// Estado mínimo de qual painel a Activity Bar tem selecionado (seção 2), e
// qual documento tem o painel "Configurações desta página" (seções 5.6/8.4)
// aberto — compartilhado entre o Explorer e o editor, que podem abri-lo.
export const useUiStore = defineStore('ui', () => {
  const activePanel = ref<ActivityPanel>('explorer')
  const pageOverridesDocId = ref<string | null>(null)

  function openPageOverrides(id: string) {
    pageOverridesDocId.value = id
  }

  return {activePanel, pageOverridesDocId, openPageOverrides}
})

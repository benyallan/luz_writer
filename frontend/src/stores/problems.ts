import {ref, computed} from 'vue'
import {defineStore} from 'pinia'
import {EventsOn} from '../../wailsjs/runtime/runtime'
import * as App from '../../wailsjs/go/main/App'
import type {model} from '../../wailsjs/go/models'

// Estado do Rule Engine (seção 9): revalidado automaticamente pelo backend a
// cada save relevante (projeto/capítulo/target/plugins/overrides), que emite
// luz:problems — este store só escuta e expõe.
export const useProblemsStore = defineStore('problems', () => {
  const items = ref<model.Problem[]>([])

  EventsOn('luz:problems', (problems: model.Problem[]) => {
    items.value = problems ?? []
  })

  const errorCount = computed(() => items.value.filter(p => p.severity === 'error').length)
  const hasErrors = computed(() => errorCount.value > 0)

  async function refresh() {
    items.value = (await App.Validate()) ?? []
  }

  return {items, errorCount, hasErrors, refresh}
})

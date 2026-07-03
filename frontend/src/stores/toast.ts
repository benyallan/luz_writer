import {ref} from 'vue'
import {defineStore} from 'pinia'

export interface Toast {
  id: number
  message: string
  kind: 'error' | 'info'
}

let nextId = 1

// Estados de erro amigáveis (seção 12, Etapa 6): rede de segurança global
// para chamadas ao backend que falham sem tratamento local (ex.: salvar um
// target com dado inválido) — em vez de um erro engolido em silêncio, o
// usuário vê uma notificação com a mensagem real do Go.
export const useToastStore = defineStore('toast', () => {
  const items = ref<Toast[]>([])

  function push(message: string, kind: Toast['kind'] = 'error') {
    const id = nextId++
    items.value.push({id, message, kind})
    setTimeout(() => dismiss(id), 6000)
  }

  function dismiss(id: number) {
    items.value = items.value.filter(t => t.id !== id)
  }

  return {items, push, dismiss}
})

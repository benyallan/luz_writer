import {Mark, mergeAttributes} from '@tiptap/core'

// Visual de intenção canônico da seção 6.0: fundo suave + aspas geradas por
// CSS (pseudo-elementos) — nenhum caractere de aspas é digitado/armazenado.
export const LuzInlineQuote = Mark.create({
  name: 'luzInlineQuote',

  parseHTML() {
    return [{tag: 'span[data-luz-mark="luzInlineQuote"]'}]
  },

  renderHTML({HTMLAttributes}) {
    return ['span', mergeAttributes(HTMLAttributes, {'data-luz-mark': 'luzInlineQuote', class: 'luz-inline-quote'}), 0]
  },
})

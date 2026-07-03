import {Mark, mergeAttributes} from '@tiptap/core'

// Idiomas do MVP (seção 6, tabela de marks) — mapeados para nomes babel no Go.
export const LUZ_LANG_OPTIONS = [
  {value: 'en', label: 'Inglês'},
  {value: 'fr', label: 'Francês'},
  {value: 'de', label: 'Alemão'},
  {value: 'es', label: 'Espanhol'},
  {value: 'it', label: 'Italiano'},
  {value: 'la', label: 'Latim'},
  {value: 'pt-BR', label: 'Português (Brasil)'},
]

export const LuzLang = Mark.create({
  name: 'luzLang',

  addAttributes() {
    return {
      lang: {
        default: null,
        parseHTML: element => element.getAttribute('data-lang'),
        renderHTML: attributes => ({'data-lang': attributes.lang}),
      },
    }
  },

  parseHTML() {
    return [{tag: 'span[data-luz-mark="luzLang"]'}]
  },

  renderHTML({HTMLAttributes}) {
    return ['span', mergeAttributes(HTMLAttributes, {'data-luz-mark': 'luzLang', class: 'luz-lang'}), 0]
  },
})

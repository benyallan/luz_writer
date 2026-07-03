import {Mark, mergeAttributes} from '@tiptap/core'

// Visual de intenção simplificado (seção 6.0): sublinhado pontilhado sutil
// para distinguir de formatação manual. A composição literal do estilo
// (itálico/negrito/versalete/cor) é aplicada no momento da exportação pelo
// \newcommand gerado no Go (internal/plugins.CustomStyleCommand) — o editor
// não recompõe visualmente o estilo em tempo real.
export const LuzCustomStyle = Mark.create({
  name: 'luzCustomStyle',

  addAttributes() {
    return {
      styleId: {
        default: null,
        parseHTML: element => element.getAttribute('data-style-id'),
        renderHTML: attributes => ({'data-style-id': attributes.styleId}),
      },
    }
  },

  parseHTML() {
    return [{tag: 'span[data-luz-mark="luzCustomStyle"]'}]
  },

  renderHTML({HTMLAttributes}) {
    return ['span', mergeAttributes(HTMLAttributes, {'data-luz-mark': 'luzCustomStyle', class: 'luz-custom-style'}), 0]
  },
})
